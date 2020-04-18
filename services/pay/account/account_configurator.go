package account

import (
	"fmt"
	"math/big"
	"net/http"

	"strconv"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/services/pay/channel"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/constants"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	h "github.com/Soneso/lumenshine-backend/services/pay/horizon"

	"github.com/pkg/errors"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Configurator is responsible for configuring new Stellar accounts that participate in ICO.
type Configurator struct {
	DB         *db.DB
	log        *logrus.Entry
	cnf        *config.Config
	ChannelMgr *channel.Manager
}

//NewAccountConfigurator initialises the account configurator
func NewAccountConfigurator(DB *db.DB, cnf *config.Config, cm *channel.Manager) *Configurator {
	l := new(Configurator)
	l.DB = DB
	l.cnf = cnf
	l.log = helpers.GetDefaultLog("Account-Configurator", "")
	l.ChannelMgr = cm

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
		build.AutoSequence{SequenceProvider: h.Horizon},
	}

	nc := db.NewNativeCalculator(constants.StellarDecimalPlaces) // stellar uses 8 decimals

	if !hasTrustline {
		//check if user has enough lumen to do the trustline

		//need to calculate the minimum balance from the base-reserve and entry count
		//we use +1 entries, because the new trustline also counts in here
		minB, err := nc.DenomFromString(c.cnf.StellarBaseReserveDenom)
		if err != nil {
			return "", 0, errors.Wrap(err, "Could not read stellar base-reserve")
		}
		minB.Mul(minB, big.NewInt(int64(2+acc.SubentryCount+1)))

		//get the current balance of the account, in order to check, if we need to send some money to make the trustline
		bS, err := acc.GetNativeBalance()
		if err != nil {
			return "", 0, errors.Wrap(err, "Could not read stellar balance")
		}
		curB, err := nc.DenomFromNativ(bS)
		if err != nil {
			return "", 0, errors.Wrap(err, "Could not convert balance balance")
		}

		addStroops := curB.Sub(curB, minB)

		if addStroops.Cmp(big.NewInt(0)) < 0 {
			//not enough --> transfer difference, in order to be able to do the trustline
			muts = append(muts,
				build.Payment(
					build.SourceAccount{AddressOrSeed: phase.DistPK},
					build.Destination{AddressOrSeed: o.StellarUserPublicKey},
					build.NativeAmount{Amount: nc.ToNativ(addStroops.Abs(addStroops))},
				),
			)
		}

		//add the trustline
		muts = append(muts,
			build.Trust(
				ico.AssetCode,
				ico.IssuerPK,
				build.SourceAccount{AddressOrSeed: o.StellarUserPublicKey}),
		)
	}

	//add the coin transfer to the transaction
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
func (c *Configurator) ExecuteTransaction(o *m.UserOrder, tx string) (string, error) {

	_, aec, err := c.DB.GetActiveExchangeCurrecnyByID(o.ExchangeCurrencyID, o.IcoPhaseID, c.log)
	if err != nil {
		c.log.WithError(err).WithFields(logrus.Fields{"order_id": o.ID}).Error("Could not get active ExchangeCurrency")
		return "", err
	}

	txe, err := h.DecodeFromBase64(tx)
	if err != nil {
		c.log.WithError(err).WithFields(logrus.Fields{"order_id": o.ID, "tx": tx}).Error("Could not decode transaction")
		return "", err
	}

	//we verify that the presigner was the first one signed, and that the transaction matches this signature
	if txe.E.Signatures == nil || len(txe.E.Signatures) < 2 {
		c.log.WithError(err).WithField("order_id", o.ID).Error("signatures count not correct")
		return "", errors.New("signatures count not correct")
	}

	hash32, err := network.HashTransaction(&txe.E.Tx, c.cnf.Stellar.NetworkPassphrase)
	txHash := hash32[:]
	preKeyPair, err := keypair.Parse(aec.R.IcoPhase.DistPresignerSeed)
	if err != nil {
		return "", err
	}

	err = preKeyPair.Verify(txHash, txe.E.Signatures[0].Signature)
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("could not check transation pre signature")
		return "", err
	}

	//check that autosequence from stellar user account is still correct
	acc, exists, err := c.GetAccount(o.StellarUserPublicKey)
	if err != nil {
		return "", errors.Wrap(err, "Could not read stellar account")
	}

	if !exists {
		return "", fmt.Errorf("account does not exist")
	}

	accSeq, err := strconv.ParseInt(acc.Sequence, 10, 64)
	if err != nil {
		return "", err
	}

	if fmt.Sprintf("%d", accSeq+1) != fmt.Sprintf("%d", txe.E.Tx.SeqNum) {
		return "", fmt.Errorf("sequencenumber does not match")
	}

	//check that sum of weight of signers (excluding presigner) is enough to do the transaction
	var sw int32
	for i := 1; i < len(txe.E.Signatures); i++ {
		sigH := txe.E.Signatures[i].Hint

		//get correstponding signer from account
		for _, signer := range acc.Signers {
			kp, _ := keypair.Parse(signer.Key)
			if kp.Hint() == sigH {
				sw += signer.Weight
			}
		}
	}
	if sw < int32(acc.Thresholds.MedThreshold) {
		return "", fmt.Errorf("weight does not match")
	}

	//all checks passed -> transfert fee to useraccount if not done yet
	if !o.FeePayed {
		phase := aec.R.IcoPhase
		nc := db.NewNativeCalculator(constants.StellarDecimalPlaces) // stellar uses 8 decimals
		opDenom, err := nc.DenomFromString(c.cnf.StellarOperationFeeDenom)
		if err != nil {
			return "", errors.Wrap(err, "Could not read stellar base-base")
		}
		opDenom.Mul(opDenom, big.NewInt(int64(len(txe.E.Tx.Operations))))

		ch, err := c.ChannelMgr.GetChannel(phase.DistPK, phase.DistPresignerSeed, phase.DistPostsignerSeed)
		if err != nil {
			return "", err
		}

		txFee, err := build.Transaction(
			build.SourceAccount{AddressOrSeed: ch.PK},
			build.Network{Passphrase: c.cnf.Stellar.NetworkPassphrase},
			build.AutoSequence{SequenceProvider: h.Horizon},
			build.Payment(
				build.SourceAccount{AddressOrSeed: phase.DistPK},
				build.Destination{AddressOrSeed: o.StellarUserPublicKey},
				build.NativeAmount{Amount: nc.ToNativ(opDenom)},
			),
		)
		if err != nil {
			c.ChannelMgr.ReleaseChannel(ch.PK)
			return "", errors.Wrap(err, "could not create fee transaction")
		}

		txeFee, err := txFee.Sign(ch.Seed, phase.DistPresignerSeed, phase.DistPostsignerSeed)
		txeFeeStr, err := txeFee.Base64()
		if err != nil {
			c.ChannelMgr.ReleaseChannel(ch.PK)
			return "", errors.Wrap(err, "error signing fee transaction")
		}

		_, _, err = h.SubmitTransaction(o, txeFeeStr, c.log, true)
		if err != nil {
			c.ChannelMgr.ReleaseChannel(ch.PK)
			return "", errors.Wrap(err, "error submitting fee transaction")
		}
		c.ChannelMgr.ReleaseChannel(ch.PK)

		//mark the fee payed flag on the order
		o.FeePayed = true
		_, err = o.Update(c.DB, boil.Whitelist(m.UserOrderColumns.FeePayed, m.UserOrderColumns.UpdatedAt))
		if err != nil {
			return "", errors.Wrap(err, "error updating fee payed in order")
		}
	}

	//fee transfered-> move coins
	s := build.Sign{Seed: aec.R.IcoPhase.DistPostsignerSeed}
	err = s.MutateTransactionEnvelope(txe)
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("Error postsigning transaction")
		return "", err
	}

	txes, err := txe.Base64()
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("Error decoding txe")
		return "", err
	}

	hash, _, err := h.SubmitTransaction(o, txes, c.log, true)
	if err != nil {
		c.log.WithError(err).WithField("order_id", o.ID).Error("error decoding txe")
		return "", err
	}

	//update order to reflect payment
	o.OrderStatus = m.OrderStatusFinished
	o.StellarTransactionID = hash
	_, err = o.Update(c.DB, boil.Whitelist(m.UserOrderColumns.OrderStatus, m.UserOrderColumns.StellarTransactionID, m.UserOrderColumns.UpdatedAt))
	if err != nil {
		return "", err
	}

	return hash, err
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

	ch, err := c.ChannelMgr.GetChannel(aec.R.IcoPhase.DistPK, aec.R.IcoPhase.DistPresignerSeed, aec.R.IcoPhase.DistPostsignerSeed)
	if err != nil {
		return err
	}

	phase := aec.R.IcoPhase
	nc := db.NewNativeCalculator(constants.StellarDecimalPlaces)
	startBalance := nc.ToNativ(startingBallanceDenom)

	tx, err := h.BuildTransaction(
		ch.PK,
		[]string{ch.Seed, phase.DistPresignerSeed, phase.DistPostsignerSeed},
		0,
		build.CreateAccount(
			build.SourceAccount{AddressOrSeed: aec.R.IcoPhase.DistPK},
			build.Destination{AddressOrSeed: account},
			build.NativeAmount{Amount: startBalance},
		),
	)

	if err != nil {
		c.ChannelMgr.ReleaseChannel(ch.PK)
		return errors.Wrap(err, "Error building user create transaction")
	}

	_, _, err = h.SubmitTransaction(order, tx, c.log, true)
	if err != nil {
		c.ChannelMgr.ReleaseChannel(ch.PK)
		return errors.Wrap(err, "Error creating user stellar account")
	}

	c.ChannelMgr.ReleaseChannel(ch.PK)

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
	hAccount, err := h.Horizon.LoadAccount(account)
	if err != nil {
		if err, ok := err.(*horizon.Error); ok && err.Response.StatusCode == http.StatusNotFound {
			return hAccount, false, nil
		}
		return hAccount, false, err
	}

	return hAccount, true, nil
}
