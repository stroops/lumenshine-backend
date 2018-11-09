package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Soneso/lumenshine-backend/db/querying"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/services/sse/db"
	"github.com/Soneso/lumenshine-backend/services/sse/environment"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"database/sql"
)

var (
	//last processed transaction id
	startupMaxLedgerID int64
	errorNoLedger      = errors.New("Ledger does not exists")
)

//StellarProcessor processes the stellar data
type StellarProcessor struct {
	env *environment.Environment
	db  *db.DB
	log *logrus.Entry
}

//NewStellarProcessor returns a new StellarProcessor
func NewStellarProcessor(log *logrus.Entry) *StellarProcessor {
	// start from actual opID
	c := new(StellarProcessor)
	c.db = environment.Env.DBH
	c.env = environment.Env
	c.log = log

	//we need to read the current latest OrderID
	//we need this information, in order to know, which operations must be signaled and which not.
	//this is configured via the with_resume flag in the DB
	var tmpID nextID

	err := queries.Raw("select max("+m.HistoryLedgerColumns.Sequence+") as value from "+m.TableNames.HistoryLedgers).Bind(nil, c.db, &tmpID)
	if err != nil {
		c.log.WithError(err).Error("Error reading intial index")
		startupMaxLedgerID = 0
	}
	startupMaxLedgerID = tmpID.NextID

	return c
}

type nextID struct {
	NextID int64 `boil:"value" json:"next_id"`
}

func (s *StellarProcessor) getNextID() (int64, error) {
	sqlStr := querying.GetSQLKeyString(`update @sse_index set value=value+1 where @name=$1 and @value>0 returning value`,
		map[string]string{
			"@sse_index": m.TableNames.SseIndex,
			"@name":      m.SseIndexColumns.Name,
			"@value":     m.SseIndexColumns.Value,
		})

	var ni nextID

	//read next operation_id to process
	err := queries.Raw(sqlStr, m.SseIndexNamesLastLedgerID).Bind(nil, s.db, &ni)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.WithError(err).Error("Error selecting next ledgerid")
			return 0, err
		}
		return 0, nil
	}

	return ni.NextID, nil

}

func (s *StellarProcessor) getLedger(ledgerID int64) (*m.HistoryLedger, error) {
	l, err := m.HistoryLedgers(qm.Where(m.HistoryLedgerColumns.Sequence+"=?", fmt.Sprintf("%d", ledgerID))).One(s.db)
	if err == sql.ErrNoRows {
		return nil, errorNoLedger
	}
	return l, err
}

//StartProcessing starts processing the stellar data
//should be started as a goroutine
//Set the value in the db for the last index to the one before the one you want to process. the value will be incremented after the first run
func (s *StellarProcessor) StartProcessing() {
	//read netx legderID to process
	//we wait until the value is set != 0
	var nextID int64
	for {
		nextID, _ = s.getNextID()
		if nextID != 0 {
			break
		}
		s.log.Info("No index set. Waiting...")
		time.Sleep(2 * time.Second)
	}

	for {
		legder, err := s.getLedger(nextID)
		if err != nil {
			if err == errorNoLedger {
				if nextID < startupMaxLedgerID {
					//if we did not find an ledger and are below the current startupMaxID,
					//there is a gap in the horizon db. therefore we process the next ledger ID
					nextID, err = s.getNextID()
					if err != nil {
						s.log.WithError(err).Error("Error retreiving next id")
					}
				} else {
					//we are behind the startupMaxID, so there where no new legderIDs in horizon. wait some time
					s.log.WithField("ledger_seq", nextID).Info("Waiting for ledger")
					time.Sleep(5 * time.Second) // ledger closes roughly every 5 seconds
				}
			} else {
				s.log.WithError(err).WithField("ledger_seq", nextID).Error("Error retreiving ledger")
			}
			//retry either old ledger, or new one if below startupMaxID
			continue
		}

		err = s.processLedger(legder)
		if err != nil {
			s.log.WithError(err).WithField("ledger_seq", nextID).Error("Error processing ledger")
		}

		//read next operation_id to process
		nextID, err = s.getNextID()
		if err != nil {
			s.log.WithError(err).Error("Error selecting next id")
			time.Sleep(2 * time.Second)
		}
	}
}

//{ CREATE_ACCOUNT = 0, PAYMENT = 1, PATH_PAYMENT = 2, MANAGE_OFFER = 3, CREATE_PASSIVE_OFFER = 4,
//	SET_OPTIONS = 5, CHANGE_TRUST = 6, ALLOW_TRUST = 7, ACCOUNT_MERGE = 8, INFLATION = 9, MANAGE_DATA = 10, BUMP_SEQUENCE = 11 }
func (s *StellarProcessor) processLedger(l *m.HistoryLedger) error {
	s.log.WithField("ledger_seq", l.Sequence).Info("processing ledger")

	sqlStr := `
	SELECT sse_config.*, history_transactions.id as transaction_id, history_operations.id as operation_id, history_ledgers.id as ledger_id
	FROM sse_config
	  INNER JOIN history_operations on cast(details->>'to' as character varying)=stellar_account
	  inner join history_transactions on history_operations.transaction_id = history_transactions.id
	  inner join history_ledgers on history_transactions.ledger_sequence = history_ledgers.sequence
	WHERE
	  (history_ledgers.sequence=$1) and
	  (2<<type&operation_types=operation_types)
	  AND (case when type=1 or type=2 then cast(details->>'to' as character varying)=stellar_account else true end)
	  AND (case when type=0 then cast(details->>'account' as character varying)=stellar_account else true end)
	  AND (case when with_resume=false then history_ledgers.sequence>=$2 else true end)`

	var orders []sseConfigData
	err := queries.Raw(sqlStr, l.Sequence, startupMaxLedgerID).Bind(nil, s.db, &orders)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // nothing found take next ledger
		}
		return err
	}

	for _, sC := range orders {
		var sD m.SseDatum
		sD.SseConfigID = sC.SseConfigID
		sD.SourceReceiver = sC.SourceReceiver
		sD.Status = m.SseDataStatusNew
		sD.StellarAccount = sC.StellarAccount
		sD.OperationType = sC.OperationType
		sD.OperationData = sC.OperationData
		sD.TransactionID = sC.TransactionID
		sD.OperationID = sC.OperationID
		sD.LedgerID = sC.LedgerID
		err := sD.Insert(s.db, boil.Infer())
		if err != nil {
			s.log.WithError(err).Error("Error inserting sse-data")
		}
	}

	return nil
}

type sseConfigData struct {
	SseConfigID    int       `boil:"id"`
	SourceReceiver string    `boil:"source_receiver"`
	StellarAccount string    `boil:"stellar_account"`
	OperationType  int       `boil:"type"`
	OperationData  null.JSON `boil:"details"`
	TransactionID  int64     `boil:"transaction_id"`
	OperationID    int64     `boil:"operation_id"`
	LedgerID       int64     `boil:"ledger_id"`
}
