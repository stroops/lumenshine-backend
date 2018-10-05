package bitcoin

import (
	paycrypto "github.com/Soneso/lumenshine-backend/services/pay/crypto"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stellar/go/support/errors"
	"github.com/tyler-smith/go-bip32"
)

//NewAddressGenerator returns a new Bitcoin address generator
func NewAddressGenerator(masterPublicKeyString string, chainParams *chaincfg.Params) (*AddressGenerator, error) {
	deserializedMasterPublicKey, err := bip32.B58Deserialize(masterPublicKeyString)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing master public key")
	}

	if deserializedMasterPublicKey.IsPrivate {
		return nil, errors.New("Key is not a master public key")
	}

	return &AddressGenerator{deserializedMasterPublicKey, chainParams}, nil
}

//Generate generates a new bitcoin address
func (g *AddressGenerator) Generate(index uint32) (address string, privateKey string, err error) {
	if g.masterPublicKey == nil {
		return "", "", errors.New("No master public key set")
	}
	accountKey, err := g.masterPublicKey.NewChildKey(index)
	if err != nil {
		return "", "", err
	}

	privateKey, err = paycrypto.PrivateECDSAKey(accountKey)
	if err != nil {
		return "", "", err
	}

	addr, err := btcutil.NewAddressPubKey(accountKey.Key, g.chainParams)
	if err != nil {
		return "", "", err
	}

	address = addr.AddressPubKeyHash().EncodeAddress()
	//return address, string(addr.PubKey().SerializeUncompressed()), privateKey, nil
	return address, privateKey, nil
}
