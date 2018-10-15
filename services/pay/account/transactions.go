package account

import (
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

//buildTransaction creates a stellar transaction
//we use AutoSequence on the source, which is the user account mostly
//signer is the seed of the source
func (c *Configurator) buildTransaction(source string, signers []string, sequence uint64, mutators ...build.TransactionMutator) (string, error) {
	muts := []build.TransactionMutator{
		build.SourceAccount{AddressOrSeed: source},
		build.Network{Passphrase: c.cnf.Stellar.NetworkPassphrase},
	}

	if sequence != 0 {
		muts = append(muts, build.Sequence{Sequence: sequence})
	} else {
		muts = append(muts, build.AutoSequence{SequenceProvider: c.Horizon})
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

//submitTransaction submits the transaction to the network and returns the hash if successfull
func (c *Configurator) submitTransaction(o *m.UserOrder, transaction string) (string, error) {
	localLog := c.log.WithFields(logrus.Fields{
		"order_id": o.ID,
	})
	localLog.WithField("transaction", transaction).Info("Submitting transaction")

	var ol m.OrderTransactionLog
	ol.TX = transaction
	ol.OrderID = o.ID

	resp, err := c.Horizon.SubmitTransaction(transaction)
	ol.TXHash = resp.Hash

	if err != nil {
		fields := logrus.Fields{"err": err}
		if err, ok := err.(*horizon.Error); ok {
			ol.ResultCode = null.StringFrom(string(err.Problem.Extras["result_codes"]))
			fields["result_codes"] = string(err.Problem.Extras["result_codes"])

			ol.ResultXDR = null.StringFrom(string(err.Problem.Extras["result_xdr"]))
			fields["result_xdr"] = string(err.Problem.Extras["result_xdr"])
		}

		ol.ErrorText = null.StringFrom(err.Error())
		ol.Status = false
		err = ol.Insert(c.DB, boil.Infer())
		if err != nil {
			localLog.Error("could not save transaction log")
		}
		localLog.WithFields(fields).Error("Error submitting transaction")
		return "", err
	}

	ol.Status = true

	err = ol.Insert(c.DB, boil.Infer())
	if err != nil {
		localLog.Error("could not save transaction log")
	}
	localLog.Info("Transaction successfully submitted")

	return resp.Hash, nil
}

// decodeFromBase64 decodes the transaction from a base64 string into a TransactionEnvelopeBuilder
func (c *Configurator) decodeFromBase64(encodedXdr string) (*build.TransactionEnvelopeBuilder, error) {

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
	n := build.Network{Passphrase: c.cnf.Stellar.NetworkPassphrase}
	err = txEnvelopeBuilder.MutateTX(n)
	if err != nil {
		return nil, err
	}

	return &txEnvelopeBuilder, nil
}
