package stellar

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/channel"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	h "github.com/Soneso/lumenshine-backend/services/pay/horizon"
	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Soneso/lumenshine-backend/services/pay/constants"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/log"
	"github.com/volatiletech/null"
)

//Ensure, that we implement all methods
var _ paymentchannel.Channel = (*Channel)(nil)

//Channel reprensents the connection to the eth blochain
type Channel struct {
	dbh *DBH
	db  *db.DB

	log        *logrus.Entry
	cnf        *config.Config
	client     *horizon.Client
	ChannelMgr *channel.Manager
}

//DBH connection to horizon db
type DBH struct {
	*sql.DB
}

//NewStellarChannel connects the stellar-client
func NewStellarChannel(DB *db.DB, cnf *config.Config, cm *channel.Manager) *Channel {
	stl := new(Channel)
	stl.db = DB
	stl.cnf = cnf
	stl.log = helpers.GetDefaultLog("Stellar-Listener", "")
	stl.ChannelMgr = cm

	stl.client = &horizon.Client{
		URL: cnf.Stellar.Horizon,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	root, err := stl.client.Root()
	if err != nil {
		log.WithField("err", err).Error("Error loading Horizon root")
		os.Exit(-1)
	}

	if root.NetworkPassphrase != cnf.Stellar.NetworkPassphrase {
		log.Errorf("Invalid network passphrase (have=%s, want=%s)", root.NetworkPassphrase, cnf.Stellar.NetworkPassphrase)
		os.Exit(-1)
	}

	//connect the horizon database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.StellarHorizenDB.DBHost, cnf.StellarHorizenDB.DBPort, cnf.StellarHorizenDB.DBUser, cnf.StellarHorizenDB.DBPassword, cnf.StellarHorizenDB.DBName)

	dbh, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Errorf("Failed to connect to horizon-db: %v", err)
		os.Exit(-1)
	}

	err = dbh.Ping()
	if err != nil {
		log.Errorf("Failed to ping horizon-database: %v", err)
		os.Exit(-1)
	}

	stl.dbh = &DBH{dbh}

	stl.log.Info("Stellar-Channel created")

	return stl
}

//TransferAmount transfers the given amount to the given address in the stellar network
//also adds the transaction logs
func (l *Channel) TransferAmount(Order *m.UserOrder, inTx *m.IncomingTransaction, Amount *big.Int, fromAddress string, TransationStatus string, BTCOutIndex int) error {

	phase, err := Order.IcoPhase().One(l.db)
	if err != nil {
		return errors.Wrap(err, "Could not read order-phase")
	}

	ch, err := l.ChannelMgr.GetChannel(phase.DistPK, phase.DistPresignerSeed, phase.DistPostsignerSeed)
	if err != nil {
		return errors.Wrap(err, "Could not get free channel")
	}
	defer l.ChannelMgr.ReleaseChannel(ch.PK)

	nc := db.NewNativeCalculator(constants.StellarDecimalPlaces) // stellar uses 8 decimals
	opDenom, err := nc.DenomFromString(l.cnf.StellarOperationFeeDenom)
	if err != nil {
		return errors.Wrap(err, "Could not read stellar base-base")
	}

	realAmount := Amount.Sub(Amount, opDenom)

	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: ch.PK},
		build.Network{Passphrase: l.cnf.Stellar.NetworkPassphrase},
		build.AutoSequence{SequenceProvider: l.client},
		build.Payment(
			build.SourceAccount{AddressOrSeed: Order.PaymentAddress},
			build.Destination{AddressOrSeed: fromAddress},
			build.NativeAmount{Amount: nc.ToNativ(realAmount)},
		),
	)
	if err != nil {
		return errors.Wrap(err, "could not create transfer transaction")
	}

	txe, err := tx.Sign(ch.Seed, Order.PaymentSeed)
	txeStr, err := txe.Base64()
	if err != nil {
		return errors.Wrap(err, "error signing transfer transaction")
	}

	hash, result, err := h.SubmitTransaction(Order, txeStr, l.log, false)
	//create out-transaction log
	var oTx m.OutgoingTransaction
	oTx.IncomingTransactionID = null.IntFrom(inTx.ID)
	oTx.OrderID = Order.ID
	oTx.Status = TransationStatus
	oTx.PaymentNetwork = m.PaymentNetworkStellar
	oTx.SenderAddress = Order.PaymentAddress
	oTx.ReceivingAddress = fromAddress
	oTx.TransactionString = txeStr
	oTx.TransactionHash = hash
	oTx.TransactionError = result
	oTx.PaymentNetworkAmountDenomination = Amount.String()
	oTx.ExecuteStatus = (err == nil)

	if err := oTx.Insert(l.db, boil.Infer()); err != nil {
		//just log the error
		l.log.WithError(err).WithField("order_id", Order.ID).Error("error saving outgoing transaction")
	}

	return err
}

//Start the listener for the eth blockchain
func (l *Channel) Start() error {
	l.log.Info("StellarListener starting")

	legderNumber, err := l.db.GetStellarLedgerToProcess()
	if err != nil {
		err = errors.Wrap(err, "Error getting stellar legderID to process from DB")
		l.log.Error(err)
		return err
	}

	go l.processLedgers(legderNumber)
	return nil
}

//GeneratePaymentAddress generates a ethereum address and seed
func (l *Channel) GeneratePaymentAddress() (Address string, Seed string, err error) {

	a, err := keypair.Random()
	if err != nil {
		return "", "", err
	}
	return a.Address(), a.Seed(), nil
}

//Name returns the name of the channel
func (l *Channel) Name() string {
	return m.PaymentNetworkStellar
}
