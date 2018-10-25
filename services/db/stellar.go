package main

import (
	"time"

	"github.com/Soneso/lumenshine-backend/pb"

	hm "github.com/Soneso/lumenshine-backend/db/horizon/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

//CheckWallet checks whether the given name and fedname are ok
func (s *server) GetStellarTransactions(ctx context.Context, r *pb.GetStellarTransactionsRequest) (*pb.StellarTransactionResponse, error) {
	timeFrom := time.Unix(r.StartTimestamp, 0)
	timeTo := time.Unix(r.EndTimestamp, 0)

	txs, err := hm.HistoryTransactions(qm.Where(
		hm.HistoryTransactionColumns.Account+"=? and "+hm.HistoryTransactionColumns.CreatedAt+">=? and "+hm.HistoryTransactionColumns.CreatedAt+"<=?",
		r.StellarAccountPk, timeFrom, timeTo,
	)).All(hdb)

	if err != nil {
		return nil, err
	}

	ops, err := hm.HistoryOperations(
		qm.Where(hm.HistoryOperationColumns.TransactionID+" in (select "+hm.HistoryTransactionColumns.ID+" from "+hm.TableNames.HistoryTransactions+" where "+hm.HistoryTransactionColumns.Account+"=?)", r.StellarAccountPk),
	).All(hdb)

	if err != nil {
		return nil, err
	}

	ret := new(pb.StellarTransactionResponse)
	ret.Transactions = make([]*pb.StellarTransaction, len(txs))

	for i, tx := range txs {
		ret.Transactions[i] = &pb.StellarTransaction{}
		ret.Transactions[i].TransactionHash = tx.TransactionHash
		ret.Transactions[i].FeePaid = int64(tx.FeePaid)
		ret.Transactions[i].OperationCount = int64(tx.OperationCount)
		ret.Transactions[i].CreatedAt = tx.CreatedAt.Time.Unix()
		//ret.Transactions[i].TxResult = tx.TXResult
		ret.Transactions[i].MemoType = tx.MemoType
		ret.Transactions[i].Memo = tx.Memo.String
		ret.Transactions[i].Operations = getOperations(ops, tx.ID)
	}
	return ret, nil
}

func getOperations(ops hm.HistoryOperationSlice, txID int64) []*pb.StellarOperations {
	var ret []*pb.StellarOperations
	for _, op := range ops {
		if op.TransactionID == txID {
			ret = append(ret, &pb.StellarOperations{
				ApplicationOrder: int64(op.ApplicationOrder),
				Type:             int64(op.Type),
				Details:          string(op.Details.JSON),
				SourceAccount:    op.SourceAccount,
			})
		}
	}
	return ret
}
