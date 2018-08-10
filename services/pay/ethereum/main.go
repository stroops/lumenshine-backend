package ethereum

import (
	"math/big"

	"github.com/stellar/go/support/errors"
	"github.com/tyler-smith/go-bip32"
)

var (
	ten      = big.NewInt(10)
	eighteen = big.NewInt(18)
	// weiInEth = 10^18
	weiInEth = new(big.Rat).SetInt(new(big.Int).Exp(ten, eighteen, nil))
)

//AddressGenerator generator for new address
type AddressGenerator struct {
	masterPublicKey *bip32.Key
}

//EthToWei converts ethereum to wei
func EthToWei(eth string) (*big.Int, error) {

	valueRat := new(big.Rat)
	_, ok := valueRat.SetString(eth)
	if !ok {
		return nil, errors.New("Could not convert to *big.Rat")
	}

	// Calculate value in Wei
	valueRat.Mul(valueRat, weiInEth)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("Invalid precision, is value smaller than 1 Wei?")
	}

	return valueRat.Num(), nil
}
