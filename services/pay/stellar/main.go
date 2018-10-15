package stellar

import (
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
	"github.com/sirupsen/logrus"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/log"
)

//Ensure, that we implement all methods
var _ paymentchannel.Channel = (*Channel)(nil)

//Channel reprensents the connection to the eth blochain
type Channel struct {
	db     *db.DB
	log    *logrus.Entry
	cnf    *config.Config
	client *horizon.Client
}

//NewStellarChannel connects the stellar-client
func NewStellarChannel(DB *db.DB, cnf *config.Config) *Channel {
	stl := new(Channel)
	stl.db = DB
	stl.cnf = cnf
	stl.log = helpers.GetDefaultLog("Stellar-Listener", "")

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

	stl.log.Info("Stellar-Channel created")

	return stl
}

//TransferAmount transfers the given amount to the given address in the btc network
//also adds the transaction logs
func (l *Channel) TransferAmount(Order *m.UserOrder, TxHash string, Amount *big.Int, fromAddress string, PaymentType string, BTCOutIndex int) error {
	return nil
}

//Start the listener for the eth blockchain
func (l *Channel) Start() error {
	l.log.Info("StellarListener starting")

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
