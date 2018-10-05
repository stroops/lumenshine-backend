package crypto

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

//PrivateECDSAKey generates the private key for a given accountkey
func PrivateECDSAKey(accountKey *bip32.Key) (string, error) {
	decoded := base58.Decode(accountKey.B58Serialize())
	privateKey := decoded[46:78]

	// Hex private key to ECDSA private key
	privateKeyECDSA, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	// ECDSA private key to hex private key
	privateKey = crypto.FromECDSA(privateKeyECDSA)
	fmt.Println("2- ", hexutil.Encode(privateKey))
	return hexutil.Encode(privateKey), nil
}
