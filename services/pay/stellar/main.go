package stellar

import (
	"errors"
	"math/big"
)

//AddressGenerator generates a new address
type AddressGenerator struct{}

// Status describes status of account processing
type Status string

var (
	ten         = big.NewInt(10)
	seven       = big.NewInt(7)
	stroopInXLM = new(big.Rat).SetInt(new(big.Int).Exp(ten, seven, nil))
)

//XLMToStroops converts a xlm string to bigInt
func XLMToStroops(xlm string) (*big.Int, error) {
	valueRat := new(big.Rat)

	_, ok := valueRat.SetString(xlm)
	if !ok {
		return nil, errors.New("Could not convert to *big.Rat")
	}

	// Calculate value in stroop
	valueRat.Mul(valueRat, stroopInXLM)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("Invalid precision, is value smaller than 1 Stroop?")
	}

	return valueRat.Num(), nil
}

//StroopToCoin converts a stroop string to bigInt
func StroopToCoin(stroop int64) (*big.Int, error) {
	valueRat := new(big.Rat)
	valueRat.SetInt64(stroop)

	// Calculate value in stroop
	valueRat.Quo(valueRat, stroopInXLM)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("Invalid precision, is value smaller than 1 Coin?")
	}
	return valueRat.Num(), nil
}
