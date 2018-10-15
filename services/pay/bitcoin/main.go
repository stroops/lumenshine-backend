package bitcoin

import (
	"math/big"
	"os"

	"github.com/Soneso/lumenshine-backend/helpers"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
)

// for test
//https://testnet.manu.backend.hamburg/faucet
//pub: mpJ1f6Dwr3m9dBAWaFCKpHu7up3FKEjAss
//secrete: cPx61d2RohCdm8DuMbc7jMjkfHPcfJXKhSYwUAJdU3UCiXJhjVJm
//Wallet used: https://copay.io/#download

//Ensure, that we implement all methods
var _ paymentchannel.Channel = (*Channel)(nil)

//Channel reprensents the connection to the eth blochain
type Channel struct {
	db          *db.DB
	log         *logrus.Entry
	client      *rpcclient.Client
	cnf         *config.Config
	chainParams *chaincfg.Params
}

//NewBitcoinChannel connects the btc-client
func NewBitcoinChannel(DB *db.DB, cnf *config.Config) *Channel {
	var err error
	btc := new(Channel)
	btc.db = DB
	btc.cnf = cnf
	btc.log = helpers.GetDefaultLog("Bitcoin-Listener", "")

	connConfig := &rpcclient.ConnConfig{
		Host:         cnf.Bitcoin.RPCServer,
		User:         cnf.Bitcoin.RPCUser,
		Pass:         cnf.Bitcoin.RPCPass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	bitcoinClient, err := rpcclient.New(connConfig, nil)
	if err != nil {
		log.WithField("err", err).Error("Error connecting to bitcoin-core")
		os.Exit(-1)
	}
	btc.client = bitcoinClient

	log.Info("Bitcoin-Listener created")
	return btc
}

//TransferAmount transfers the given amount to the given address in the btc network
//also adds the transaction logs
func (l *Channel) TransferAmount(Order *m.UserOrder, TxHash string, Amount *big.Int, fromAddress string, PaymentType string, BTCOutIndex int) error {
	return nil
}

//Start starts the bitcoin listener
func (l *Channel) Start() error {
	l.log = helpers.GetDefaultLog("0", "BitcoinListener")
	l.log.Info("BitcoinListener starting")

	genesisBlockHash, err := l.client.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "Error getting genesis block")
	}

	if l.cnf.Bitcoin.Testnet {
		l.chainParams = &chaincfg.TestNet3Params
	} else {
		l.chainParams = &chaincfg.MainNetParams
	}

	if !genesisBlockHash.IsEqual(l.chainParams.GenesisHash) {
		return errors.New("Invalid genesis hash")
	}

	blockNumber, err := l.db.GetBitcoinBlockToProcess()
	if err != nil {
		err = errors.Wrap(err, "Error getting bitcoin block to process from DB")
		l.log.Error(err)
		return err
	}

	if blockNumber == 0 {
		blockNumberTmp, err := l.client.GetBlockCount()
		if err != nil {
			err = errors.Wrap(err, "Error getting the block count from bitcoin-core")
			l.log.Error(err)
			return err
		}
		blockNumber = uint64(blockNumberTmp)
	}

	go l.processBlocks(blockNumber)
	return nil
}

//GeneratePaymentAddress generates a ethereum address and seed
func (l *Channel) GeneratePaymentAddress() (Address string, Seed string, err error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", err
	}
	wifGen, err := btcutil.NewWIF(secret, l.chainParams, false)
	if err != nil {
		return "", "", err
	}

	addressGen, err := btcutil.NewAddressPubKey(wifGen.PrivKey.PubKey().SerializeUncompressed(), l.chainParams)
	if err != nil {
		return "", "", err
	}

	Address = addressGen.EncodeAddress()
	Seed = wifGen.String()

	return Address, Seed, nil
}

//Name returns the name of the channel
func (l *Channel) Name() string {
	return m.PaymentNetworkBitcoin
}
