package account

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/pkg/errors"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/support/log"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Configurator is responsible for configuring new Stellar accounts that participate in ICO.
type Configurator struct {
	DB      *db.DB
	log     *logrus.Entry
	cnf     *config.Config
	Horizon *horizon.Client
}

//NewAccountConfigurator initialises the account configurator
func NewAccountConfigurator(DB *db.DB, cnf *config.Config) *Configurator {
	var err error

	l := new(Configurator)
	l.DB = DB
	l.cnf = cnf
	l.log = helpers.GetDefaultLog("Account-Configurator", "")

	l.Horizon = &horizon.Client{
		URL: cnf.Stellar.Horizon,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	root, err := l.Horizon.Root()
	if err != nil {
		log.WithField("err", err).Error("Error loading Horizon root")
		os.Exit(-1)
	}

	if root.NetworkPassphrase != cnf.Stellar.NetworkPassphrase {
		log.Errorf("Invalid network passphrase (have=%s, want=%s)", root.NetworkPassphrase, cnf.Stellar.NetworkPassphrase)
		os.Exit(-1)
	}

	l.log.Info("Accountconfigurator created")
	return l
}

func (c *Configurator) hasTrustline(acc *horizon.Account, assetCode string, assetIssuerPK string) bool {
	//check if trustline present
	for _, b := range acc.Balances {
		if b.Asset.Code == assetCode && b.Asset.Issuer == assetIssuerPK {
			return true
		}
	}
	return false
}

//GetPaymentTransaction creates the transaction for a valid payment
//returns the transaction or an error code
//the transaction is signed with the pre-signer and must be signed againe with the postsigner on execute
func (c *Configurator) GetPaymentTransaction(o *m.UserOrder) (string, int64, error) {
	var acc horizon.Account
	var err error
	var exists bool

	_, aec, err := c.DB.GetActiveExchangeCurrecnyByID(o.ExchangeCurrencyID, o.IcoPhaseID, c.log)
	if err != nil {
		return "", 0, err
	}

	phase := aec.R.IcoPhase
	ico := phase.R.Ico

	acc, exists, err = c.GetAccount(o.StellarUserPublicKey)
	if err != nil {
		return "", 0, errors.Wrap(err, "Could not read stellar account")
	}

	if !exists {
		return "", cerr.StellarAccountNotExists, nil
	}

	hasTrustline := c.hasTrustline(&acc, ico.AssetCode, ico.IssuerPK)

	muts := []build.TransactionMutator{
		build.SourceAccount{AddressOrSeed: o.StellarUserPublicKey},
		build.Network{Passphrase: c.cnf.Stellar.NetworkPassphrase},
		build.AutoSequence{SequenceProvider: c.Horizon},
	}

	if !hasTrustline {
		//add base-fee-amount
		muts = append(muts,
			build.Payment(
				build.SourceAccount{AddressOrSeed: phase.DistPK},
				build.Destination{AddressOrSeed: o.StellarUserPublicKey},
				build.NativeAmount{Amount: c.cnf.StellarBaseFeeTrustline},
			),
		)

		//add the trustline
		muts = append(muts,
			build.Trust(
				ico.AssetCode,
				ico.IssuerPK,
				build.SourceAccount{AddressOrSeed: o.StellarUserPublicKey}),
		)
	} else {
		muts = append(muts,
			build.Payment(
				build.SourceAccount{AddressOrSeed: phase.DistPK},
				build.Destination{AddressOrSeed: o.StellarUserPublicKey},
				build.NativeAmount{Amount: c.cnf.StellarBaseFee},
			),
		)
	}

	//also add the coin transfer to the transaction
	muts = append(muts, build.Payment(
		build.SourceAccount{AddressOrSeed: phase.DistPK},
		build.Destination{AddressOrSeed: o.StellarUserPublicKey},
		build.CreditAmount{
			Code:   ico.AssetCode,
			Issuer: ico.IssuerPK,
			Amount: fmt.Sprintf("%d", o.TokenAmount),
		},
	))

	tx, err := build.Transaction(muts...)
	if err != nil {
		return "", 0, errors.Wrap(err, "Could not create transaction")
	}

	//sign with pre-signer
	txe, err := tx.Sign(phase.DistPresignerSeed)
	if err != nil {
		return "", 0, errors.Wrap(err, "Could not pre sign transaction")
	}

	txes, err := txe.Base64()
	if err != nil {
		return "", 0, errors.Wrap(err, "Could not bas64 encode TransactionEnvelopeBuilder")
	}

	return txes, 0, nil
}

//ExecuteTransaction checks the transaction, signs it and executes it
func (c *Configurator) ExecuteTransaction(o *m.UserOrder, tx string) error {

	_, aec, err := c.DB.GetActiveExchangeCurrecnyByID(o.ExchangeCurrencyID, o.IcoPhaseID, c.log)
	if err != nil {
		c.log.WithError(err).WithFields(logrus.Fields{"order_id": o.ID}).Error("Could not get active ExchangeCurrency")
		return err
	}

	txe, err := c.decodeFromBase64(tx)
	if err != nil {
		c.log.WithError(err).WithFields(logrus.Fields{"order_id": o.ID, "tx": tx}).Error("Could not decode transaction")
		return err
	}

	//TODO check signer with tx

	s := build.Sign{Seed: aec.R.IcoPhase.DistPostsignerSeed}
	err = s.MutateTransactionEnvelope(txe)
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("Error postsigning transaction")
		return err
	}

	txes, err := txe.Base64()
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("Error decoding txe")
		return err
	}
	return c.submitTransaction(txes)
}

//CreateAccount create the user stellar account
//it uses the configured distribution account as source and uses autosequence of distribution account
func (c *Configurator) CreateAccount(account string, order *m.UserOrder) error {
	ec, aec, err := c.DB.GetActiveExchangeCurrecnyByID(order.ExchangeCurrencyID, order.IcoPhaseID, c.log)
	if err != nil {
		return err
	}

	startingBallanceDenom, err := ec.DenomFromString(aec.R.IcoPhase.StellarStartingBalanceDenom)
	if err != nil {
		return err
	}

	tx, err := c.buildTransaction(
		aec.R.IcoPhase.DistPK,
		[]string{aec.R.IcoPhase.DistPresignerSeed, aec.R.IcoPhase.DistPostsignerSeed},
		0,
		build.CreateAccount(
			build.SourceAccount{AddressOrSeed: aec.R.IcoPhase.DistPK},
			build.Destination{AddressOrSeed: account},
			build.NativeAmount{Amount: ec.ToNativ(startingBallanceDenom)},
		),
	)

	if err != nil {
		return errors.Wrap(err, "Error building user create transaction")
	}

	err = c.submitTransaction(tx)
	if err != nil {
		return errors.Wrap(err, "Error creating user stellar account")
	}

	//transaction was submitted without error -> flag user profile
	u, err := m.UserProfiles(qm.Where("id=?", order.UserID)).One(c.DB)
	if err != nil {
		return err
	}
	u.StellarAccountCreated = true
	_, err = u.Update(c.DB, boil.Whitelist(m.UserProfileColumns.StellarAccountCreated, m.UserProfileColumns.UpdatedAt))
	if err != nil {
		return err
	}

	//now that the account was created, we can also update the order to hild this new public key
	order.StellarUserPublicKey = account
	_, err = order.Update(c.DB, boil.Whitelist(m.UserOrderColumns.StellarUserPublicKey, m.UserOrderColumns.UpdatedAt))
	if err != nil {
		return err
	}

	return nil
}

//GetAccount returns the horizon-account for the given address or false if it does not exist
func (c *Configurator) GetAccount(account string) (horizon.Account, bool, error) {
	var hAccount horizon.Account
	hAccount, err := c.Horizon.LoadAccount(account)
	if err != nil {
		if err, ok := err.(*horizon.Error); ok && err.Response.StatusCode == http.StatusNotFound {
			return hAccount, false, nil
		}
		return hAccount, false, err
	}

	return hAccount, true, nil
}
