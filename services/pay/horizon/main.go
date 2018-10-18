package horizon

import (
	"fmt"
	"net/http"
	"time"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	dbc "github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
	"github.com/volatiletech/sqlboiler/boil"
)

var (
	//Horizon server
	Horizon *horizon.Client
	db      *dbc.DB
	cnf     *config.Config
)

//InitHorizon initializes the horizon configuration
func InitHorizon(DB *dbc.DB, Cnf *config.Config) error {
	db = DB
	cnf = Cnf

	if Horizon == nil {
		Horizon = &horizon.Client{
			URL: cnf.Stellar.Horizon,
			HTTP: &http.Client{
				Timeout: 60 * time.Second,
			},
		}

		root, err := Horizon.Root()
		if err != nil {
			return err
		}

		if root.NetworkPassphrase != cnf.Stellar.NetworkPassphrase {
			return fmt.Errorf("Invalid network passphrase (have=%s, want=%s)", root.NetworkPassphrase, cnf.Stellar.NetworkPassphrase)
		}
	}

	return nil
}

//BuildTransaction creates a stellar transaction
//we use AutoSequence on the source, which is the user account mostly
//signer is the seed of the source
func BuildTransaction(source string, signers []string, sequence uint64, mutators ...build.TransactionMutator) (string, error) {
	muts := []build.TransactionMutator{
		build.SourceAccount{AddressOrSeed: source},
		build.Network{Passphrase: cnf.Stellar.NetworkPassphrase},
	}

	if sequence != 0 {
		muts = append(muts, build.Sequence{Sequence: sequence})
	} else {
		muts = append(muts, build.AutoSequence{SequenceProvider: Horizon})
	}

	muts = append(muts, mutators...)
	tx, err := build.Transaction(muts...)
	if err != nil {
		return "", err
	}

	var txe build.TransactionEnvelopeBuilder
	txe.Init()
	err = tx.MutateTransactionEnvelope(&txe)
	if err != nil {
		return "", err
	}

	for _, s := range signers {
		s := build.Sign{Seed: s}
		err = s.MutateTransactionEnvelope(&txe)
		if err != nil {
			return "", err
		}
	}

	return txe.Base64()
}

//SubmitTransaction submits the transaction to the network and returns the hash if successfull
//if saveOrderTxLog is set, we will save a log in OrderTransactions, if not, we will return the result as a string
func SubmitTransaction(o *m.UserOrder, transaction string, log *logrus.Entry, saveOrderTxLog bool) (hash string, resultString string, err error) {
	localLog := log.WithFields(logrus.Fields{
		"order_id": o.ID,
	})
	localLog.WithField("transaction", transaction).Info("Submitting transaction")

	var ol m.OrderTransactionLog
	ol.OrderID = o.ID
	ol.TX = transaction
	ol.Status = true //default

	resp, err := Horizon.SubmitTransaction(transaction)
	ol.TXHash = resp.Hash

	if err != nil {
		fields := logrus.Fields{"err": err}
		if err, ok := err.(*horizon.Error); ok {
			ol.ResultCode = string(err.Problem.Extras["result_codes"])
			fields["result_codes"] = ol.ResultCode

			ol.ResultXDR = string(err.Problem.Extras["result_xdr"])
			fields["result_xdr"] = ol.ResultXDR
		}

		ol.ErrorText = err.Error()
		ol.Status = false

		localLog.WithFields(fields).Error("Error submitting transaction")
	}

	if saveOrderTxLog {
		if err := ol.Insert(db, boil.Infer()); err != nil {
			localLog.Error("could not save transaction log")
			return ol.TXHash, "", err
		}
		return ol.TXHash, "", err
	}

	//if the log was not saved, we will return the messages as string
	return ol.TXHash, ol.ResultCode + " - " + ol.ResultXDR + " - " + ol.ErrorText, err
}

// DecodeFromBase64 decodes the transaction from a base64 string into a TransactionEnvelopeBuilder
func DecodeFromBase64(encodedXdr string) (*build.TransactionEnvelopeBuilder, error) {
	// Unmarshall from base64 encoded XDR format
	var decoded xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(encodedXdr, &decoded)
	if err != nil {
		return nil, err
	}

	// convert to TransactionEnvelopeBuilder
	txEnvelopeBuilder := build.TransactionEnvelopeBuilder{E: &decoded}
	txEnvelopeBuilder.Init()

	//the passphrase needs to be added
	n := build.Network{Passphrase: cnf.Stellar.NetworkPassphrase}
	err = txEnvelopeBuilder.MutateTX(n)
	if err != nil {
		return nil, err
	}

	return &txEnvelopeBuilder, nil
}
