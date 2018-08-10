package modext

import (
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"math/big"
)

//UserOrderDenomination returns the value of the string from the order
func UserOrderDenomination(o *m.UserOrder) *big.Int {
	i := new(big.Int)
	i, ok := i.SetString(o.ChainAmountDenom, 0)
	if !ok {
		i = i.SetInt64(0)
	}
	return i
}
