package modext

import (
	"fmt"
	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

//UserOrderDenomination returns the value of the string from the order
func UserOrderDenomination(o *m.UserOrder) *big.Int {
	i := new(big.Int)
	i, ok := i.SetString(o.ExchangeCurrencyDenominationAmount, 0)
	if !ok {
		i = i.SetInt64(0)
	}
	return i
}

//UserOrderFiatPaymentUsage creates the usage string for an order
func UserOrderFiatPaymentUsage(usageString string, o *m.UserOrder) string {
	return fmt.Sprintf(usageString, fmt.Sprintf("%d", o.ID))
}
