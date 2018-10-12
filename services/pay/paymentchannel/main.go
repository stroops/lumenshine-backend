package paymentchannel

import (
	"icop/admin/config"
	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/sirupsen/logrus"
)

//Listener is the listener, that chacks the specific pamynet-channel for incoming payments
type Listener struct {
	db     *db.DB
	log    *logrus.Entry
	client interface{}
	cnf    *config.Config
}

//Channel reepresents one PaymentChannel
type Channel interface {
	GetDefaultCurrency() string
	GeneratePaymentAddress() (Address string, PrivateKey string, err error)
	RefundPayment(db *db.DB, Order *m.UserOrder, TxHash string, Amount *big.Int) error
	NewListener(DB *db.DB, cnf *config.Config) (*Listener, error)
	Start() error
}
