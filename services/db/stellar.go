package main

import (
	"time"

	"github.com/Soneso/lumenshine-backend/pb"

	hm "github.com/Soneso/lumenshine-backend/db/horizon/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

type stellarOperation struct {
	TxTransactionHash string      `boil:"tx_transaction_hash" json:"tx_transaction_hash"`
	TxCreatedAt       null.Time   `boil:"tx_created_at" json:"tx_created_at"`
	TxMemoType        string      `boil:"tx_memo_type" json:"tx_memo_type"`
	TxMemo            null.String `boil:"tx_memo" json:"tx_memo"`
	TxOperationCount  int         `boil:"tx_operation_count" json:"tx_operation_count"`
	TxFeePaid         int         `boil:"tx_fee_paid" json:"tx_fee_paid"`
	TxAccount         string      `boil:"tx_account" json:"tx_account"`

	OpID               int64     `boil:"op_id" json:"op_id"`
	OpApplicationOrder int       `boil:"op_application_order" json:"op_application_order"`
	OpType             int       `boil:"op_type" json:"op_type"`
	OpDtails           null.JSON `boil:"op_details" json:"op_details"`
}

//CheckWallet checks whether the given name and fedname are ok
func (s *server) GetStellarTransactions(ctx context.Context, r *pb.GetStellarTransactionsRequest) (*pb.StellarOperations, error) {
	timeFrom := time.Unix(r.StartTimestamp, 0)
	timeTo := time.Unix(r.EndTimestamp, 0)

	cT := hm.HistoryTransactionColumns
	cO := hm.HistoryOperationColumns
	cOP := hm.HistoryOperationParticipantColumns
	cA := hm.HistoryAccountColumns
	cL := hm.HistoryLedgerColumns
	tN := hm.TableNames

	boil.DebugMode = true

	q := hm.NewQuery(
		qm.From(tN.HistoryOperations+" t1"),

		qm.Select("t4."+cT.TransactionHash+" as tx_transaction_hash"),
		qm.Select("t5."+cL.ClosedAt+" as tx_created_at"),
		qm.Select("t4."+cT.MemoType+" as tx_memo_type"),
		qm.Select("t4."+cT.Memo+" as tx_memo"),
		qm.Select("t4."+cT.OperationCount+" as tx_operation_count"),
		qm.Select("t4."+cT.FeePaid+" as tx_fee_paid"),
		qm.Select("t4."+cT.Account+" as tx_account"),

		qm.Select("t1."+cO.ID+" as op_id"),
		qm.Select("t1."+cO.ApplicationOrder+" as op_application_order"),
		qm.Select("t1."+cO.Type+" as op_type"),
		qm.Select("t1."+cO.Details+" as op_details"),

		qm.InnerJoin(tN.HistoryOperationParticipants+" t2 on t1.id=t2."+cOP.HistoryOperationID),
		qm.InnerJoin(tN.HistoryAccounts+" t3 on t3.id=t2."+cOP.HistoryAccountID),
		qm.InnerJoin(tN.HistoryTransactions+" t4 on t4.id=t1."+cO.TransactionID),
		qm.InnerJoin(tN.HistoryLedgers+" t5 on t5.sequence=t4."+cT.LedgerSequence),

		qm.Where(
			"t3."+cA.Address+"=? and t5."+cL.ClosedAt+">=? and t5."+cL.ClosedAt+"<=?",
			r.StellarAccountPk, timeFrom, timeTo,
		),
	)

	var ops []stellarOperation
	err := q.Bind(nil, hdb, &ops)
	if err != nil {
		return nil, err
	}

	ret := new(pb.StellarOperations)
	ret.Operations = make([]*pb.StellarOperation, len(ops))

	for i, op := range ops {
		ret.Operations[i] = &pb.StellarOperation{
			TxTransactionHash: op.TxTransactionHash,
			TxCreatedAt:       op.TxCreatedAt.Time.Unix(),
			TxMemoType:        op.TxMemoType,
			TxMemo:            op.TxMemo.String,
			TxOperationCount:  int64(op.TxOperationCount),
			TxFeePaid:         int64(op.TxFeePaid),
			TxAccount:         op.TxAccount,

			OpId:               op.OpID,
			OpApplicationOrder: int64(op.OpApplicationOrder),
			OpType:             int64(op.OpType),
			OpDetails:          string(op.OpDtails.JSON),
		}
	}
	return ret, nil
}
