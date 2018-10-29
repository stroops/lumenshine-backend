package main

import (
	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/gin-gonic/gin"
)

// StellarTransactionRequest request-data
// swagger:parameters StellarTransactionRequest GetStellarTransactions
type StellarTransactionRequest struct {
	// Filter for stellar account
	// required: true
	StellarAccountPK string `json:"stellar_account_pk" form:"stellar_account_pk" query:"stellar_account_pk" validate:"required"`

	// TX-Starting timestamp without time zone, eg. 2006-01-02 15:04:05
	// required: true
	StartTimestamp string `json:"start_timestamp" form:"start_timestamp" query:"start_timestamp" validate:"required"`

	// TX-Endtime timestamp without time zone, eg. 2006-01-02 15:04:05
	// required: true
	EndTimestamp string `json:"end_timestamp" form:"end_timestamp" query:"end_timestamp" validate:"required"`
}

type stellarOperations struct {
	ApplicationOrder int64  `json:"application_order"`
	Type             int64  `json:"type"`
	Details          string `json:"details"`
	SourceAccount    string `json:"source_account"`
}

// StellarTransactionResponse transactions for one user account, including all operations
// swagger:model
type StellarTransactionResponse struct {
	TransactionHash string    `json:"transaction_hash"`
	FeePaid         int64     `json:"fee_paid"`
	OperationCount  int64     `json:"operation_count"`
	CreatedAt       time.Time `json:"created_at"`
	//TXResult        string               `json:"tx_result"`
	MemoType   string               `json:"memo_type"`
	Memo       string               `json:"memo"`
	Operations []*stellarOperations `json:"operations"`
}

// GetStellarTransactions returns all transactions including all operations for one account
// swagger:route GET /portal/user/dashboard/get_stellar_transactions stellar GetStellarTransactions
//
// Lists all transactions including all operations for one account
//
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []StellarTransactionResponse
func GetStellarTransactions(uc *mw.IcopContext, c *gin.Context) {
	var l StellarTransactionRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", l.StartTimestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("start_timestamp", cerr.InvalidArgument, "start_timestamp wrong format", ""))
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", l.EndTimestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("start_timestamp", cerr.InvalidArgument, "start_timestamp wrong format", ""))
		return
	}

	txs, err := dbClient.GetStellarTransactions(c, &pb.GetStellarTransactionsRequest{
		Base:             NewBaseRequest(uc),
		StellarAccountPk: l.StellarAccountPK,
		StartTimestamp:   startTime.Unix(),
		EndTimestamp:     endTime.Unix(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	var ret []*StellarTransactionResponse
	for _, tx := range txs.Transactions {
		ret = append(ret, &StellarTransactionResponse{
			TransactionHash: tx.TransactionHash,
			FeePaid:         tx.FeePaid,
			OperationCount:  tx.OperationCount,
			CreatedAt:       time.Unix(tx.CreatedAt, 0),
			//TXResult:        tx.TxResult,
			MemoType:   tx.MemoType,
			Memo:       tx.Memo,
			Operations: getOperations(tx.Operations),
		})
	}

	c.JSON(http.StatusOK, ret)
}

func getOperations(ops []*pb.StellarOperations) []*stellarOperations {
	var ret []*stellarOperations
	for _, op := range ops {
		ret = append(ret, &stellarOperations{
			ApplicationOrder: op.ApplicationOrder,
			Type:             op.Type,
			Details:          op.Details,
			SourceAccount:    op.SourceAccount,
		})
	}
	return ret
}
