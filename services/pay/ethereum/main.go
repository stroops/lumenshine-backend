package ethereum

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"

	m "github.com/Soneso/lumenshine-backend/services/db/models"

	"encoding/hex"

	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/errors"

	"github.com/ethereum/go-ethereum/crypto"
)

//Ensure, that we implement all methods
var _ paymentchannel.Channel = (*Channel)(nil)

//Channel reprensents the connection to the eth blochain
type Channel struct {
	db     *db.DB
	log    *logrus.Entry
	cnf    *config.Config
	client *ethclient.Client
}

//NewEthereumChannel connects the eth-client
func NewEthereumChannel(DB *db.DB, cnf *config.Config) *Channel {
	var err error
	eth := new(Channel)
	eth.db = DB
	eth.cnf = cnf
	eth.log = helpers.GetDefaultLog("Ethereum-Listener", "")

	ethereumClient, err := ethclient.Dial("http://" + eth.cnf.Ethereum.RPCServer)
	if err != nil {
		eth.log.WithField("err", err).Error("Error connecting to geth")
		os.Exit(-1)
	}
	eth.client = ethereumClient

	eth.log.Info("Ethereum-Channel created")

	return eth
}

//TransferAmount transfers the given amount to the given address in the btc network
//also adds the transaction logs
func (l *Channel) TransferAmount(Order *m.UserOrder, inTx *m.IncomingTransaction, Amount *big.Int, fromAddress string, PaymentType string, BTCOutIndex int) error {
	return nil
}

//Start the listener for the eth blockchain
func (l *Channel) Start() error {
	l.log.Info("EthereumListener starting")

	blockNumber, err := l.db.GetEthereumBlockToProcess()
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum block to process from DB")
		l.log.Error(err)
		return err
	}

	// Check if connected to correct network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	id, err := l.client.NetworkID(ctx)
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum network ID")
		l.log.Error(err)
		return err
	}

	if id.String() != l.cnf.Ethereum.NetworkID {
		return errors.Errorf("Invalid network ID (have=%s, want=%s)", id.String(), l.cnf.Ethereum.NetworkID)
	}

	go l.processBlocks(blockNumber)
	return nil
}

//GeneratePaymentAddress generates a ethereum address and seed
func (l *Channel) GeneratePaymentAddress() (Address string, Seed string, err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	Address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	// 0x8ee3333cDE801ceE9471ADf23370c48b011f82a6

	Seed = hex.EncodeToString(key.D.Bytes())
	// 05b14254a1d0c77a49eae3bdf080f926a2df17d8e2ebdf7af941ea001481e57f
	return Address, Seed, nil
}

//Name returns the name of the channel
func (l *Channel) Name() string {
	return m.PaymentNetworkEthereum
}
