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

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries"
)

const (
	ethereumAddressIndexKey = "eth_address_index"
	ethereumLastBlockKey    = "eth_last_block"

	bitcoinAddressIndexKey = "btc_address_index"
	bitcoinLastBlockKey    = "btc_last_block"
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

//GetNextChainAddressIndex returns the next chain index
func (db *DB) GetNextChainAddressIndex(chain string) (uint32, error) {
	key := ""
	var index uint32
	if chain == m.PaymentNetworkBitcoin {
		key = "btc_address_index"
	} else if chain == m.PaymentNetworkEthereum {
		key = "eth_address_index"
	}

	if key != "" {
		//get and update
		var v m.KeyValueStore
		sql := querying.GetSQLKeyString(`update @key_value_store set @int_value = @int_value+1 where @key = 
			(select @key from @key_value_store where @key=$1 limit 1 for update) returning @int_value`,
			map[string]string{
				"@key_value_store": m.TableNames.KeyValueStore,
				"@int_value":       m.KeyValueStoreColumns.IntValue,
				"@key":             m.KeyValueStoreColumns.Key,
			})

		err := queries.Raw(sql, key).Bind(nil, db, &v)
		if err != nil {
			return 0, err
		}
		index = uint32(v.IntValue)
	}

	return index, nil
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

func (db *DB) getBlockToProcess(key string) (uint64, error) {
	kv, err := m.KeyValueStores(qm.Where(m.KeyValueStoreColumns.Key+"=?", key)).One(db)
	if err != nil {
		return 0, errors.Wrap(err, "Error getting `"+key+"` from DB")
	}

	block, err := strconv.ParseUint(kv.STRValue, 10, 64)
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

	lastBlock, err := strconv.ParseUint(kv.STRValue, 10, 64)
	if err != nil {
		return err
	}

	if block > lastBlock {
		kv.STRValue = fmt.Sprintf("%d", block)
		_, err := kv.Update(tx, boil.Whitelist(m.KeyValueStoreColumns.STRValue))
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

//GetOrderForAddress reads the user orders for open payments for the specified address and chain
//The method updates the order to refelect, that it has been processed, if the order was waiting_for_payment
//It will set the order in status OrderStatusPaymentReceived
//If no open order was found(OrderStatusWaitingForPayment), the function will return either nil or filter for any other order with the given address
//The function will be called for EVERY payment transaction in the external PaymentNetworks
//paymentUsage is either an empty string or, for stellar the MEMO with the orderID
func (db *DB) GetOrderForAddress(paymentChannel string, address string, paymentUsage string) (*m.UserOrder, error) {
	userOrder := new(m.UserOrder)

	var sqlStr string
	if paymentUsage == m.PaymentNetworkStellar {
		//payment must be the order ID
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
	err := queries.Raw(sqlStr, m.OrderStatusPaymentReceived, paymentChannel, address, m.OrderStatusWaitingForPayment).Bind(nil, db, userOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			//if we did not find an open order, we need to search for all orders on the user, in order to get the multiple transaction triggered
			userOrder, err = m.UserOrders(qm.Where(m.UserOrderColumns.PaymentNetwork+"=? and "+m.UserOrderColumns.PaymentAddress+"=?", paymentChannel, address)).One(db)
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

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

//AddNewTransaction adds a transaction to the database and returns true, if it was allready present
func (db DB) AddNewTransaction(log *logrus.Entry, paymentChannel string, txHash string, toAddress string, orderID int, denomAmount *big.Int) (isDuplicate bool, err error) {

	d := new(m.ProcessedTransaction)
	d.PaymentNetwork = paymentChannel
	d.ReceivingAddress = toAddress
	d.TransactionID = txHash
	d.OrderID = orderID
	d.Status = m.TransactionStatusNew
	d.PaymentNetworkAmountDenomination = denomAmount.String()

	err = d.Insert(db, boil.Infer())
	isDuplicate = isDuplicateError(err)
	if isDuplicate {
		//add the transaction to the multiple table for manual handling
		b := new(m.MultipleTransaction)
		b.PaymentNetwork = paymentChannel
		b.ReceivingAddress = toAddress
		b.TransactionID = txHash
		b.OrderID = orderID
		b.PaymentNetworkAmountDenom = denomAmount.String()
		errB := b.Insert(db, boil.Infer())
		if errB != nil {
			//we don't handle this error, just log it
			log.WithError(err).WithFields(logrus.Fields{"order_id": orderID, "transaction_id": txHash}).Error("Error saving multiple transaction")
		}
		return true, nil
	}

	if err != nil {
		return isDuplicate, err
	}

	//update the order to hold the transaction-ref
	order, err := m.FindUserOrder(db, orderID, "id")
	if err != nil {
		return isDuplicate, err
	}
	order.ProcessedTransactionID = null.IntFrom(d.ID)
	_, err = order.Update(db, boil.Whitelist(m.UserOrderColumns.ProcessedTransactionID, m.UserKycDocumentColumns.UpdatedAt))
	if err != nil {
		return isDuplicate, err
	}

	return isDuplicate, db.handleNewTransaction(log, d, denomAmount)
}

//handleNewTransaction checks the transaction data and updates the user_profile to reflect the payment
//the order must be in status OrderStatusPaymentReceived and will be set to status OrderStatusWaitingUserTX
func (db DB) handleNewTransaction(log *logrus.Entry, tx *m.ProcessedTransaction, denomAmount *big.Int) (err error) {

	order := new(m.UserOrder)

	sqlStr := querying.GetSQLKeyString(`update @user_order set @order_status=$1, @updated_at=current_timestamp where id =
			(select id from @user_order where id=$2 and @order_status=$3 limit 1 for update) returning
			*`,
		map[string]string{
			"@user_order":   m.TableNames.UserOrder,
			"@order_status": m.UserOrderColumns.OrderStatus,
			"@updated_at":   m.UserOrderColumns.UpdatedAt,
		})

	err = queries.Raw(sqlStr, m.OrderStatusWaitingForPayment, tx.OrderID, m.OrderStatusPaymentReceived).Bind(nil, db, order)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"order_id": tx.OrderID, "transaction_id": tx.TransactionID}).Error("Error selecting phasedata")
		return err
	}

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
	} else {
		//amount payed is exactly the amount bought. we can check/update, if there are coins left
		ph := new(m.IcoPhase)

		sqlStr = querying.GetSQLKeyString(`update @ico_phase set @tokens_left=@tokens_left-$1, @updated_at=current_timestamp where id=$2 and @ico_phase_status=$3 and 
		  start_time<=current_timestamp and end_time>=current_timestamp returning *`,
			map[string]string{
				"@ico_phase":   m.TableNames.IcoPhase,
				"@tokens_left": m.IcoPhaseColumns.TokensLeft,
				"@updated_at":  m.IcoPhaseColumns.UpdatedAt,
			})

		err = queries.Raw(sqlStr, order.TokenAmount, order.IcoPhaseID, m.IcoPhaseStatusActive).Bind(nil, db, ph)
		if err != nil {
			//something is not ok
			//either amount was to small, or phase is already gone... we will read the data againe and check
			if err != sql.ErrNoRows {
				// log error
				log.WithError(err).WithFields(logrus.Fields{"order_id": order.ID, "transaction_id": tx.TransactionID}).Error("Error selecting phasedata")
				return err
			}
			ph, err = m.IcoPhases(qm.Where(m.IcoPhaseColumns.IcoPhaseStatus+"=?", m.IcoPhaseStatusActive)).One(db)
			if err != nil {
				return err
			}
			if ph.TokensLeft < order.TokenAmount {
				order.OrderStatus = m.OrderStatusNoCoinsLeft
			} else {
				order.OrderStatus = m.OrderStatusPhaseExpired
			}
			_, err = order.Update(db, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.UpdatedAt))
			if err != nil {
				return err
			}
		}

		//everything seems ok -> update the order
		order.OrderStatus = m.OrderStatusWaitingUserTransaction
		_, err = order.Update(db, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.UpdatedAt))
		if err != nil {
			return err
		}
	}

	//check all user orders and if one is payed, set flag, if not, remove flag
	user, err := m.UserProfiles(qm.Where("id=?", order.UserID)).One(db)
	if err != nil {
		return err
	}

	cnt, err := m.UserOrders(qm.Where("user_id=? and order_status=?", order.UserID, m.OrderStatusWaitingUserTransaction)).Count(db)
	if cnt > 0 {
		user.PaymentState = m.PaymentStateOpen
	} else {
		user.PaymentState = m.PaymentStateClose
	}
	_, err = user.Update(db, boil.Whitelist(m.UserProfileColumns.PaymentState, m.UserProfileColumns.UpdatedAt))
	if err != nil {
		return err
	}

	return nil
}
