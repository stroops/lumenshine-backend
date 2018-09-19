h1 Problem with ethereum

In order to build the ethereum client, we need to install some header giles manualy. this is a bug in the geth dependencies.
To do so make:
_go get github.com/ethereum/go-ethereum_

then move into the root directory of the project and do:

_/usr/bin/cp -rf "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"_

h1 dependecy error

If you get 
'# github.com/Soneso/lumenshine-backend/vendor/github.com/stellar/go/support/log
../../vendor/github.com/stellar/go/support/log/entry.go:122:14: cannot use hook.Entries (type []logrus.Entry) as type []*logrus.Entry in return argument
'
on go build,
please look at:
https://github.com/stellar/go/issues/565#issuecomment-409611689
this fixes the issue, till stellar updates.
Do the changes in the vendor folder...