package paymentchannel

import (
	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

const (
	//TransferRefund used when refunding
	TransferRefund = "refund"

	//TransferCashout used when payment processed to end
	TransferCashout = "cashout"
)

//Channel reepresents one PaymentChannel
type Channel interface {
	//GetDefaultCurrency() string
	GeneratePaymentAddress() (Address string, Seed string, err error)
	TransferAmount(Order *m.UserOrder, TxHash string, Amount *big.Int, destAddress string, PaymentType string, BTCOutIndex int) error
	Start() error
	Name() string
}
