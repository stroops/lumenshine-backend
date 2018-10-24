package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/Soneso/lumenshine-backend/services/pay/config"

	"github.com/Soneso/lumenshine-backend/db/querying"
	m "github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq" //needed for SQL access
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
	"github.com/volatiletech/sqlboiler/queries"
)

const (
	ethereumLastBlockKey = "eth_last_block"
	bitcoinLastBlockKey  = "btc_last_block"
	stellarLastLedgerKey = "xlm_last_ledger_id"
)

//DB general DB struct
type DB struct {
	*sql.DB
}

//CreateNewDB creates a new DB connection
func CreateNewDB(cnf *config.Config) (*DB, error) {
	var err error

	//connect the customer database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.CustomerDB.DBHost, cnf.CustomerDB.DBPort, cnf.CustomerDB.DBUser, cnf.CustomerDB.DBPassword, cnf.CustomerDB.DBName)

	DBC, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to customer-db: %v", err)
	}

	err = DBC.Ping()
	if err != nil {
		log.Fatalf("Failed to ping customer-database: %v", err)
	}

	boil.SetDB(DBC)

	return &DB{DBC}, nil
}

//GetEthereumBlockToProcess gets the last processed eth block
func (db *DB) GetEthereumBlockToProcess() (uint64, error) {
	return db.getBlockToProcess(ethereumLastBlockKey)
}

//SaveLastProcessedEthereumBlock saves the last processed eth block
func (db *DB) SaveLastProcessedEthereumBlock(block uint64) error {
	return db.saveLastProcessedBlock(ethereumLastBlockKey, block)
}

//GetBitcoinBlockToProcess saves the last processed btc block
func (db *DB) GetBitcoinBlockToProcess() (uint64, error) {
	return db.getBlockToProcess(bitcoinLastBlockKey)
}

//SaveLastProcessedBitcoinBlock saves the last processed btc block
func (db *DB) SaveLastProcessedBitcoinBlock(block uint64) error {
	return db.saveLastProcessedBlock(bitcoinLastBlockKey, block)
}

//GetStellarLedgerToProcess returns the new ledgerid to process
func (db *DB) GetStellarLedgerToProcess() (int, error) {
	id, err := db.getBlockToProcess(stellarLastLedgerKey)
	return int(id), err
}

//SaveLastProcessedStellarLedger saves the last processed stellar ledger
func (db *DB) SaveLastProcessedStellarLedger(ledgerID int) error {
	return db.saveLastProcessedBlock(stellarLastLedgerKey, uint64(ledgerID))
}

func (db *DB) getBlockToProcess(key string) (uint64, error) {
	kv, err := m.KeyValueStores(qm.Where(m.KeyValueStoreColumns.Key+"=?", key)).One(db)
	if err != nil {
		return 0, errors.Wrap(err, "Error getting `"+key+"` from DB")
	}

	block, err := strconv.ParseUint(kv.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	// If set, `block` is the last processed block so we need to start processing from the next one.
	if block > 0 {
		block++
	}
	return block, nil
}

func (db *DB) saveLastProcessedBlock(key string, block uint64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	kv, err := m.KeyValueStores(qm.Where(m.KeyValueStoreColumns.Key+"=?", key)).One(tx)
	if err != nil {
		return err
	}

	lastBlock, err := strconv.ParseUint(kv.Value, 10, 64)
	if err != nil {
		return err
	}

	if block > lastBlock {
		kv.Value = fmt.Sprintf("%d", block)
		_, err := kv.Update(tx, boil.Whitelist(m.KeyValueStoreColumns.Value))
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

//GetOrderForAddress reads the user orders for open payments for the specified address and chain
//If no open order was found(OrderStatusWaitingForPayment), the function will return either nil or filter for any other order with the given address
//The function will be called for EVERY payment transaction in the external PaymentNetworks
//paymentUsage is either an empty string or, for stellar, the MEMO with the orderID. If no memo/or wrong was given, the order will not be processed
func (db *DB) GetOrderForAddress(l paymentchannel.Channel, address string, paymentUsage string) (*m.UserOrder, error) {
	var sqlStr string
	if l.Name() == m.PaymentNetworkStellar {
		//usage must be the order ID
		paymentUsage = strings.Trim(paymentUsage, " \n\t")
		if _, err := strconv.ParseInt(paymentUsage, 10, 64); err != nil {
			return nil, fmt.Errorf("Could not convert paymentUsage '%s' to id", paymentUsage)
		}

		sqlStr = querying.GetSQLKeyString(`update @user_order set @order_status=$1, @updated_at=current_timestamp where id =
			(select id from @user_order where @payment_network=$2 and @payment_address=$3 and @order_status=$4 and id=@id limit 1 for update) returning
			*`,
			map[string]string{
				"@user_order":      m.TableNames.UserOrder,
				"@order_status":    m.UserOrderColumns.OrderStatus,
				"@updated_at":      m.UserOrderColumns.UpdatedAt,
				"@payment_network": m.UserOrderColumns.PaymentNetwork,
				"@payment_address": m.UserOrderColumns.PaymentAddress,
				"@id":              paymentUsage,
			})
	} else {
		sqlStr = querying.GetSQLKeyString(`update @user_order set @order_status=$1, @updated_at=current_timestamp where id =
		(select id from @user_order where @payment_network=$2 and @payment_address=$3 and @order_status=$4 limit 1 for update) returning
		*`,
			map[string]string{
				"@user_order":      m.TableNames.UserOrder,
				"@order_status":    m.UserOrderColumns.OrderStatus,
				"@updated_at":      m.UserOrderColumns.UpdatedAt,
				"@payment_network": m.UserOrderColumns.PaymentNetwork,
				"@payment_address": m.UserOrderColumns.PaymentAddress,
			})
	}

	//set order to payment recived
	userOrder := new(m.UserOrder)
	err := queries.Raw(sqlStr, m.OrderStatusPaymentReceived, l.Name(), address, m.OrderStatusWaitingForPayment).Bind(nil, db, userOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			//TODO: do not trigger multiple transactions

			//if we did not find an open order, we need to search for all orders on the user, in order to get the multiple transaction triggered
			userOrder, err = m.UserOrders(
				qm.Where(m.UserOrderColumns.PaymentNetwork+"=? and "+m.UserOrderColumns.PaymentAddress+"=?", l.Name(), address),
				qm.OrderBy("id desc"),
			).One(db)

			if err != nil {
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return userOrder, nil
}

//AddNewTransaction adds a transaction from the payment network to the database and returns true, if it was allready present
//The function implies, that the order is in status OrderStatusPaymentReceived
func (db DB) AddNewTransaction(log *logrus.Entry, paymentChannel paymentchannel.Channel, txHash string,
	toAddress string, fromAddress string, order *m.UserOrder, denomAmount *big.Int, BTCOutIndex int) (isDuplicate bool, err error) {

	iTx := new(m.IncomingTransaction)
	iTx.PaymentNetwork = paymentChannel.Name()
	iTx.ReceivingAddress = toAddress
	iTx.SenderAddress = fromAddress
	iTx.TransactionHash = txHash
	iTx.BTCSRCOutIndex = BTCOutIndex
	iTx.OrderID = order.ID
	iTx.Status = m.TransactionStatusNew
	iTx.PaymentNetworkAmountDenomination = denomAmount.String()

	if order.OrderStatus != m.OrderStatusPaymentReceived {
		//TODO check with multiple transactions removing
		//this means, that we updated the order in some other process or the user send multiple payments, thus it's a additional payment --> refund
		iTx.Status = m.TransactionStatusRefund
		err = iTx.Insert(db, boil.Infer())
		if err != nil {
			return true, err
		}

		err = paymentChannel.TransferAmount(order, iTx, denomAmount, fromAddress, m.TransactionStatusRefund, BTCOutIndex)
		if err != nil {
			//we don't handle this error, just log it TransferAmount will set the error text in the order
			log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_id": txHash}).Error("Error refunding")
		}
		return true, nil
	}

	err = iTx.Insert(db, boil.Infer())
	if err != nil {
		return false, err
	}

	return false, db.handleNewTransaction(log, paymentChannel, order, iTx, denomAmount)
}

//handleNewTransaction checks the transaction data and updates the user_profile to reflect the payment
//the order must be in status OrderStatusWaitingForPayment and will be set to status OrderStatusWaitingUserTX
func (db DB) handleNewTransaction(log *logrus.Entry, paymentChannel paymentchannel.Channel, order *m.UserOrder, iTx *m.IncomingTransaction, denomAmount *big.Int) (err error) {
	//check order amount
	oa := new(big.Int)
	oa.SetString(order.ExchangeCurrencyDenominationAmount, 0)

	cmp := oa.Cmp(denomAmount)
	if cmp == -1 || cmp == 1 {
		if cmp == -1 {
			//order amount < denomAmount
			order.OrderStatus = m.OrderStatusOverPay
		} else if cmp == 1 {
			order.OrderStatus = m.OrderStatusUnderPay
		}

		_, err = order.Update(db, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.UpdatedAt))
		if err != nil {
			return err
		}

		err = paymentChannel.TransferAmount(order, iTx, denomAmount, iTx.SenderAddress, m.TransactionStatusRefund, iTx.BTCSRCOutIndex)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error refunding wrong amount")
		}
		return nil
	}

	//amount payed is exactly the amount bought. we can check/update, if there are coins left
	ph := new(m.IcoPhase)

	sqlStr := querying.GetSQLKeyString(`update @ico_phase set @tokens_left=@tokens_left-$1, @updated_at=current_timestamp where id=$2 and @ico_phase_status=$3 and
		  start_time<=current_timestamp and end_time>=current_timestamp and @tokens_left>=$4 returning *`,
		map[string]string{
			"@ico_phase":   m.TableNames.IcoPhase,
			"@tokens_left": m.IcoPhaseColumns.TokensLeft,
			"@updated_at":  m.IcoPhaseColumns.UpdatedAt,
		})

	err = queries.Raw(sqlStr, order.TokenAmount, order.IcoPhaseID, m.IcoPhaseStatusActive, order.TokenAmount).Bind(nil, db, ph)
	if err != nil {
		//something is not ok
		//either left token-amount is to small, or phase is already gone... we will read the data againe and check
		if err != sql.ErrNoRows {
			// log error
			log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error selecting phasedata")
			iTx.Status = m.TransactionStatusError
			if _, err := iTx.Update(db, boil.Whitelist(m.IncomingTransactionColumns.Status, m.IncomingTransactionColumns.UpdatedAt)); err != nil {
				return err
			}
			return err
		}

		ph, err = m.IcoPhases(qm.Where("id=?", order.IcoPhaseID)).One(db)
		if err != nil {
			iTx.Status = m.TransactionStatusError
			if _, err := iTx.Update(db, boil.Whitelist(m.IncomingTransactionColumns.Status, m.IncomingTransactionColumns.UpdatedAt)); err != nil {
				return err
			}
			return err
		}

		if ph.TokensLeft < order.TokenAmount {
			order.OrderStatus = m.OrderStatusNoCoinsLeft
		}
		if ph.IcoPhaseStatus != m.IcoPhaseStatusActive {
			order.OrderStatus = m.OrderStatusPhaseExpired
		}
		if _, err := order.Update(db, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.UpdatedAt)); err != nil {
			return err
		}

		if err := paymentChannel.TransferAmount(order, iTx, denomAmount, iTx.SenderAddress, m.TransactionStatusRefund, iTx.BTCSRCOutIndex); err != nil {
			log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error refunding wrong phase_status or tokenamount")
			return err
		}

		return nil
	}

	//everything seems ok -> update the order but first re-check current status
	err = order.Reload(db)
	if err != nil {
		iTx.Status = m.TransactionStatusError
		if _, err := iTx.Update(db, boil.Whitelist(m.IncomingTransactionColumns.Status, m.IncomingTransactionColumns.UpdatedAt)); err != nil {
			return err
		}
		return err
	}

	if order.OrderStatus != m.OrderStatusPaymentReceived {
		//order changed meanwhile, eg from second process, we refund and exit
		if err := paymentChannel.TransferAmount(order, iTx, denomAmount, iTx.SenderAddress, m.TransactionStatusRefund, iTx.BTCSRCOutIndex); err != nil {
			log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error refunding wrong order status")
			return err
		}
		return nil
	}

	order.OrderStatus = m.OrderStatusWaitingUserTransaction
	_, err = order.Update(db, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.UpdatedAt))
	if err != nil {
		iTx.Status = m.TransactionStatusError
		if _, err := iTx.Update(db, boil.Whitelist(m.IncomingTransactionColumns.Status, m.IncomingTransactionColumns.UpdatedAt)); err != nil {
			return err
		}
		return err
	}

	//TODO: move amount to payout-account in payment network

	//check all user orders and if one is payed, set flag, if not, remove flag
	user, err := m.UserProfiles(qm.Where("id=?", order.UserID)).One(db)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error selecting user profile-payment-status")
		return nil
	}

	cnt, err := m.UserOrders(qm.Where("user_id=? and order_status=?", order.UserID, m.OrderStatusWaitingUserTransaction)).Count(db)
	if cnt > 0 {
		user.PaymentState = m.PaymentStateOpen
	} else {
		user.PaymentState = m.PaymentStateClose
	}
	_, err = user.Update(db, boil.Whitelist(m.UserProfileColumns.PaymentState, m.UserProfileColumns.UpdatedAt))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_hash": iTx.TransactionHash}).Error("Error updating user profile-payment-status")
	}

	return nil
}
