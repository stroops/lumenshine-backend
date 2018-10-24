package paymentchannel

import (
	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

const ()

//Channel reepresents one PaymentChannel
type Channel interface {
	//GetDefaultCurrency() string
	GeneratePaymentAddress() (Address string, Seed string, err error)
	TransferAmount(Order *m.UserOrder, inTx *m.IncomingTransaction, Amount *big.Int, destAddress string, TransationStatus string, BTCOutIndex int) error
	Start() error
	Name() string
}
