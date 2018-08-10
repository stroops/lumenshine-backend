package bitcoin

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stellar/go/support/errors"
	"github.com/tyler-smith/go-bip32"
)

var (
	eight = big.NewInt(8)
	ten   = big.NewInt(10)
	// satInBtc = 10^8
	satInBtc = new(big.Rat).SetInt(new(big.Int).Exp(ten, eight, nil))
)

// for test
//https://testnet.manu.backend.hamburg/faucet
//pub: mpJ1f6Dwr3m9dBAWaFCKpHu7up3FKEjAss
//secrete: cPx61d2RohCdm8DuMbc7jMjkfHPcfJXKhSYwUAJdU3UCiXJhjVJm
//Wallet used: https://copay.io/#download

//AddressGenerator generator for bitcoin addresses
type AddressGenerator struct {
	masterPublicKey *bip32.Key
	chainParams     *chaincfg.Params
}

//BtcToSat converts bitcoin to satoshi
func BtcToSat(btc string) (*big.Int, error) {
	valueRat := new(big.Rat)
	_, ok := valueRat.SetString(btc)
	if !ok {
		return nil, errors.New("Could not convert to *big.Rat")
	}

	// Calculate value in satoshi
	valueRat.Mul(valueRat, satInBtc)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("Invalid precision, is value smaller than 1 satoshi?")
	}

	return valueRat.Num(), nil
}
