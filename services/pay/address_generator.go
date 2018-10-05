package main

import (
	"errors"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/bitcoin"
	"github.com/Soneso/lumenshine-backend/services/pay/ethereum"
	"github.com/btcsuite/btcd/chaincfg"
)

var (
	//generators holds a list of all generators. they are initialized the first time tehy are used
	//they key is composed from the payment-network + masterkey
	generators map[string]interface{}
)

//GeneratePaymentAddress generates a new address in the payment network
//checks if the generator allredy exists and if not, creates one
func (s *server) GeneratePaymentAddress(paymentNetwork string, masterKey string) (address string, privatekey string, index uint32, err error) {
	key := paymentNetwork + "-" + masterKey
	if generators == nil {
		generators = make(map[string]interface{})
	}
	g, exist := generators[key]

	if paymentNetwork == m.PaymentNetworkEthereum || paymentNetwork == m.PaymentNetworkBitcoin {
		if index, err = s.Env.DBC.GetNextChainAddressIndex(paymentNetwork); err != nil {
			return
		}
	}

	switch paymentNetwork {
	case m.PaymentNetworkEthereum:
		if !exist {
			if g, err = ethereum.NewAddressGenerator(masterKey); err != nil {
				return
			}

			generators[key] = g
		}
		address, privatekey, err = g.(*ethereum.AddressGenerator).Generate(index)
	case m.PaymentNetworkBitcoin:
		if !exist {
			var bitcoinChainParams *chaincfg.Params
			if s.Env.Config.Bitcoin.Testnet {
				bitcoinChainParams = &chaincfg.TestNet3Params
			} else {
				bitcoinChainParams = &chaincfg.MainNetParams
			}
			if g, err = bitcoin.NewAddressGenerator(masterKey, bitcoinChainParams); err != nil {
				return
			}

			generators[key] = g
		}
		address, privatekey, err = g.(*bitcoin.AddressGenerator).Generate(index)
	default:
		err = errors.New("Not a valid payment network")
	}
	return
}
