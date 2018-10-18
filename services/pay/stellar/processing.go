package stellar

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"database/sql"

	mh "github.com/Soneso/lumenshine-backend/db/horizon/models"
)

func (l *Channel) processLedgers(ledgerID int) {
	if ledgerID == 0 {
		l.log.Info("Starting from the latest ledger")
	} else {
		l.log.Infof("Starting from ledgerID %d", ledgerID)
	}

	// Time when last new ledger has been seen
	lastLedgerSeen := time.Now()
	noLedgerkWarningLogged := false

	for {
		ledger, err := l.getLedger(ledgerID)
		if err != nil {
			l.log.WithFields(logrus.Fields{"err": err, "ledgerSeq": ledgerID}).Error("Error getting ledger")
			time.Sleep(1 * time.Second)
			continue
		}

		// Block doesn't exist yet
		if ledger == nil {
			if time.Since(lastLedgerSeen) > 3*time.Minute && !noLedgerkWarningLogged {
				l.log.Warn("No new block in more than 3 minutes")
				noLedgerkWarningLogged = true
			}

			time.Sleep(5 * time.Second) //ledger closes roughly every 5 seconds
			continue
		}

		// Reset counter when new block appears
		lastLedgerSeen = time.Now()
		noLedgerkWarningLogged = false

		if ledger.Sequence == 0 {
			l.log.Error("Stellar node is not synced yet. Unable to process ledger. Sleeping 30 seconds")
			time.Sleep(30 * time.Second)
			continue
		}

		err = l.processLedger(ledger)
		if err != nil {
			l.log.WithError(err).WithFields(logrus.Fields{"legerSeq": ledger.Sequence}).Error("Error processing ledger")
			time.Sleep(1 * time.Second)
			continue
		}

		// Persist block number
		err = l.db.SaveLastProcessedStellarLedger(ledgerID)
		if err != nil {
			l.log.WithError(err).Error("Error saving last processed block")
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		ledgerID++
	}
}

// getLedger returns (nil, nil) if ledger has not been found (not exists yet)
func (l *Channel) getLedger(LedgerNumber int) (*mh.HistoryLedger, error) {
	ledger, err := mh.HistoryLedgers(qm.Where(mh.HistoryLedgerColumns.Sequence+"=?", LedgerNumber)).One(l.dbh)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return ledger, nil
}

//Operation one operation in the horizon db
type Operation struct {
	To          string `json:"to"`
	From        string `json:"from"`
	Amount      string `json:"amount"`
	AssetCode   string `json:"asset_code"`
	AssetType   string `json:"asset_type"`
	AssetIssuer string `json:"asset_issuer"`
}

func (l *Channel) processLedger(ledger *mh.HistoryLedger) error {
	transactions, err := mh.HistoryTransactions(qm.Where(mh.HistoryTransactionColumns.LedgerSequence+"=?", ledger.Sequence)).All(l.dbh)
	if err != nil {
		return err
	}

	localLog := l.log.WithFields(logrus.Fields{
		"ledgerSeq":    ledger.Sequence,
		"legderTime":   ledger.ClosedAt,
		"transactions": len(transactions),
	})
	localLog.Info("Processing ledger")

	for _, transaction := range transactions {
		//get all operations for transaction
		operations, err := mh.HistoryOperations(qm.Where(mh.HistoryOperationColumns.TransactionID+"=?", transaction.ID)).All(l.dbh)
		if err != nil {
			if err != sql.ErrNoRows {
				l.log.WithError(err).WithField("tx_id", transaction.ID).Error("Error fetching operations")
			}
			continue
		}

		for _, operation := range operations {
			if operation.Type == 1 {
				//we handle only payment operations
				if !operation.Details.IsZero() {
					var op Operation
					if err := json.Unmarshal(operation.Details.JSON, &op); err != nil {
						l.log.WithError(err).WithField("op_id", operation.ID).Error("Error unmarshaling operation")
					}

					//we handle only nativ asset
					if op.AssetType == "native" {
						if err := l.processTransaction(ledger, transaction, &op); err != nil {
							l.log.WithError(err).WithField("op_id", operation.ID).Error("Error processing operation")
						}
					}
				}
			}
		}
	}

	localLog.Info("Processed ledger")

	return nil
}

func (l *Channel) processTransaction(ledger *mh.HistoryLedger, transaction *mh.HistoryTransaction, operation *Operation) error {

	localLog := l.log.WithFields(logrus.Fields{"transaction": transaction.TransactionHash, "rail": "stellar"})
	localLog.Debug("Processing transaction")

	//get the order from the database
	memo := ""
	if !transaction.Memo.IsZero() {
		memo = transaction.Memo.String
	}

	order, err := l.db.GetOrderForAddress(l, operation.To, memo)
	if err != nil {
		return errors.Wrap(err, "Error getting association")
	}

	if order == nil {
		localLog.Debug("Associated address not found, skipping")
		return nil
	}

	ec, err := l.db.GetExchangeCurrecnyByID(order.ExchangeCurrencyID, localLog)

	//_, aec, err := l.db.GetActiveExchangeCurrecnyByID(order.ExchangeCurrencyID, order.IcoPhaseID, localLog)
	if err != nil {
		return err
	}

	valueStoops, err := ec.DenomFromNativ(operation.Amount)
	if err != nil {
		return err
	}

	// Add transaction as processing.
	isDuplicate, err := l.db.AddNewTransaction(l.log, l, transaction.TransactionHash, operation.To, operation.From, order, valueStoops, 0)
	if err != nil {
		return err
	}

	if isDuplicate {
		localLog.Debug("Transaction already processed, skipping")
		return nil
	}

	return nil
}
