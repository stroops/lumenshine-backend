package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/errors"
)

func (l *Channel) processBlocks(blockNumber uint64) {
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
		err = l.db.SaveLastProcessedEthereumBlock(blockNumber)
		if err != nil {
			l.log.WithField("err", err).Error("Error saving last processed block")
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		blockNumber = block.NumberU64() + 1
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Channel) getBlock(blockNumber uint64) (*types.Block, error) {
	var blockNumberInt *big.Int
	if blockNumber > 0 {
		blockNumberInt = big.NewInt(int64(blockNumber))
	}

	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, err := l.client.BlockByNumber(ctx, blockNumberInt)
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

func (l *Channel) processBlock(block *types.Block) error {
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
			"",
		)
		if err != nil {
			return errors.Wrap(err, "Error processing transaction")
		}
	}

	localLog.Info("Processed block")

	return nil
}

func (l *Channel) processTransaction(hash string, valueWei *big.Int, toAddress string, fromAddress string) error {
	localLog := l.log.WithFields(logrus.Fields{"transaction": hash, "rail": "ethereum"})
	localLog.Debug("Processing transaction")

	//get the order from the database
	order, err := l.db.GetOrderForAddress(l, toAddress, "")
	if err != nil {
		return errors.Wrap(err, "Error getting association")
	}

	if order == nil {
		localLog.Debug("Associated address not found, skipping")
		return nil
	}

	// Add transaction as processing.
	isDuplicate, err := l.db.AddNewTransaction(l.log, l, hash, toAddress, fromAddress, order, valueWei, 0)
	if err != nil {
		return err
	}

	if isDuplicate {
		localLog.Debug("Transaction already processed, skipping")
		return nil
	}

	return nil
}
