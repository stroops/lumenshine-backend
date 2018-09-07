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

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/errors"
)

type Listener struct {
	DB     *db.DB
	log    *logrus.Entry
	Client *ethclient.Client
	cnf    *config.Config
}

//NewListener createa a new listener and connects the eth-client
func NewListener(DB *db.DB, cnf *config.Config) *Listener {
	var err error

	l := new(Listener)
	l.DB = DB
	l.cnf = cnf
	l.log = helpers.GetDefaultLog("Ethereum-Listener", "")

	ethereumClient, err := ethclient.Dial("http://" + l.cnf.Ethereum.RPCServer)
	if err != nil {
		l.log.WithField("err", err).Error("Error connecting to geth")
		os.Exit(-1)
	}
	l.Client = ethereumClient

	l.cnf.Ethereum.MinimumWeiValueEth, err = EthToWei(cnf.Ethereum.MinimumWeiValueEthStr)
	if err != nil {
		l.log.Error("Invalid minimum accepted Ethereum transaction value")
		os.Exit(-1)
	}

	if l.cnf.Ethereum.MinimumWeiValueEth.Cmp(new(big.Int)) == 0 {
		l.log.Error("Minimum accepted Ethereum transaction value must be larger than 0")
		os.Exit(-1)
	}

	l.log.Info("Ethereum-Listener created")
	return l
}

func (l *Listener) Start() error {
	l.log.Info("EthereumListener starting")

	blockNumber, err := l.DB.GetEthereumBlockToProcess()
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum block to process from DB")
		l.log.Error(err)
		return err
	}

	// Check if connected to correct network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	id, err := l.Client.NetworkID(ctx)
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

func (l *Listener) processBlocks(blockNumber uint64) {
	if blockNumber == 0 {
		l.log.Info("Starting from the latest block")
	} else {
		l.log.Infof("Starting from block %d", blockNumber)
	}

	// Time when last new block has been seen
	lastBlockSeen := time.Now()
	noBlockWarningLogged := false

	for {
		block, err := l.getBlock(blockNumber)
		if err != nil {
			l.log.WithFields(logrus.Fields{"err": err, "blockNumber": blockNumber}).Error("Error getting block")
			time.Sleep(1 * time.Second)
			continue
		}

		// Block doesn't exist yet
		if block == nil {
			if time.Since(lastBlockSeen) > 3*time.Minute && !noBlockWarningLogged {
				l.log.Warn("No new block in more than 3 minutes")
				noBlockWarningLogged = true
			}

			time.Sleep(1 * time.Second)
			continue
		}

		// Reset counter when new block appears
		lastBlockSeen = time.Now()
		noBlockWarningLogged = false

		if block.NumberU64() == 0 {
			l.log.Error("Etheruem node is not synced yet. Unable to process blocks. Sleeping 30 seconds")
			time.Sleep(30 * time.Second)
			continue
		}

		err = l.processBlock(block)
		if err != nil {
			l.log.WithFields(logrus.Fields{"err": err, "blockNumber": block.NumberU64()}).Error("Error processing block")
			time.Sleep(1 * time.Second)
			continue
		}

		// Persist block number
		err = l.DB.SaveLastProcessedEthereumBlock(blockNumber)
		if err != nil {
			l.log.WithField("err", err).Error("Error saving last processed block")
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		blockNumber = block.NumberU64() + 1
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Listener) getBlock(blockNumber uint64) (*types.Block, error) {
	var blockNumberInt *big.Int
	if blockNumber > 0 {
		blockNumberInt = big.NewInt(int64(blockNumber))
	}

	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, err := l.Client.BlockByNumber(ctx, blockNumberInt)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		err = errors.Wrap(err, "Error getting block from geth")
		l.log.WithField("block", blockNumberInt.String()).Error(err)
		return nil, err
	}

	return block, nil
}

func (l *Listener) processBlock(block *types.Block) error {
	transactions := block.Transactions()
	blockTime := time.Unix(block.Time().Int64(), 0)

	localLog := l.log.WithFields(logrus.Fields{
		"blockNumber":  block.NumberU64(),
		"blockTime":    blockTime,
		"transactions": len(transactions),
	})
	localLog.Info("Processing block")

	for _, transaction := range transactions {
		to := transaction.To()
		if to == nil {
			// Contract creation
			continue
		}

		err := l.processTransaction(
			transaction.Hash().Hex(),
			transaction.Value(),
			to.Hex(),
		)
		if err != nil {
			return errors.Wrap(err, "Error processing transaction")
		}
	}

	localLog.Info("Processed block")

	return nil
}

func (l *Listener) processTransaction(hash string, valueWei *big.Int, toAddress string) error {
	localLog := l.log.WithFields(logrus.Fields{"transaction": hash, "rail": "ethereum"})
	localLog.Debug("Processing transaction")

	// Let's check if tx is valid first.

	// Check if value is above minimum required
	if valueWei.Cmp(l.cnf.Ethereum.MinimumWeiValueEth) < 0 {
		localLog.Debug("Value is below minimum required amount, skipping")
		return nil
	}

	//get the order from the database
	order, err := l.DB.GetOpenOrderForAddress(m.BlockChainEthereum, toAddress)
	if err != nil {
		return errors.Wrap(err, "Error getting association")
	}

	if order == nil {
		localLog.Debug("Associated address not found, skipping")
		return nil
	}

	// Add transaction as processing.
	processed, err := l.DB.AddNewTransaction(l.log, m.BlockChainEthereum, hash, toAddress, order.ID, valueWei)
	if err != nil {
		return err
	}

	if processed {
		localLog.Debug("Transaction already processed, skipping")
		return nil
	}

	return nil
}
