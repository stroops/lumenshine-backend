package modext

import (
	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

//UserOrderDenomination returns the value of the string from the order
func UserOrderDenomination(o *m.UserOrder) *big.Int {
	return new(big.Int).SetInt64(o.CurrencyDenomAmount)
}
