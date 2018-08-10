h1 Problem with ethereum

In order to build the ethereum client, we need to install some header giles manualy. this is a bug in the geth dependencies.
To do so make:
_go get github.com/ethereum/go-ethereum_

then move into the root directory of the project and do:

_cp -r "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"_