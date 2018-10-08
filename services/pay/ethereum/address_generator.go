package ethereum

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

//NewAddressGenerator returns a new ethereum address generator
func NewAddressGenerator(masterPublicKeyString string) (*AddressGenerator, error) {
	/*deserializedMasterPublicKey, err := bip32.B58Deserialize(masterPublicKeyString)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing master public key")
	}

	if deserializedMasterPublicKey.IsPrivate {
		return nil, errors.New("Key is not a master public key")
	}

	return &AddressGenerator{deserializedMasterPublicKey}, nil
	*/
	return &AddressGenerator{}, nil
}

//Generate generates a new ethereum address
func (g *AddressGenerator) Generate(index uint32) (address string, privateKey string, err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	// 0x8ee3333cDE801ceE9471ADf23370c48b011f82a6

	privateKey = hex.EncodeToString(key.D.Bytes())
	// 05b14254a1d0c77a49eae3bdf080f926a2df17d8e2ebdf7af941ea001481e57f

	return address, privateKey, nil

	/*if g.masterPublicKey == nil {
		return "", "", errors.New("No master public key set")
	}

	accountKey, err := g.masterPublicKey.NewChildKey(index)
	if err != nil {
		return "", "", errors.Wrap(err, "Error creating new child key")
	}

	privateKey, err = paycrypto.PrivateECDSAKey(accountKey)
	if err != nil {
		return "", "", err
	}

	uncompressed := secp256k1.UncompressPubkey(accountKey.Key)
	uncompressed = uncompressed[1:]
	keccak := crypto.Keccak256(uncompressed)
	address := ethereumCommon.BytesToAddress(keccak[12:]).Hex() // Encode lower 160 bits/20 bytes

	return address, privateKey, nil*/
}
