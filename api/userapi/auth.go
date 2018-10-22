package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"

	"github.com/sirupsen/logrus"

	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/build"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	assertAvailablePRNG()
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

//LoginStep1Request is the data needed for the first step of the login
type LoginStep1Request struct {
	Email   string `form:"email" json:"email" validate:"required,icop_email"`
	TfaCode string `form:"tfa_code" json:"tfa_code"`
}

//LoginStep1Response response for API
type LoginStep1Response struct {
	KdfPasswordSalt               string `json:"kdf_password_salt"`
	EncryptedMnemonicMasterKey    string `json:"encrypted_mnemonic_master_key"`
	MnemonicMasterKeyEncryptionIV string `json:"mnemonic_master_key_encryption_iv"`
	EncryptedMnemonic             string `json:"encrypted_mnemonic"`
	MnemonicEncryptionIV          string `json:"mnemonic_encryption_iv"`
	EncryptedWordlistMasterKey    string `json:"encrypted_wordlist_master_key"`
	WordlistMasterKeyEncryptionIV string `json:"wordlist_master_key_encryption_iv"`
	EncryptedWordlist             string `json:"encrypted_wordlist"`
	WordlistEncryptionIV          string `json:"wordlist_encryption_iv"`
	PublicKeyIndex0               string `json:"public_key_index0"`
	TfaConfirmed                  bool   `json:"tfa_confirmed"`
	SEP10TransactionChallenge     string `json:"sep10_transaction_challenge"`
}

//LoginStep1 is the first step of the login
func LoginStep1(uc *mw.IcopContext, c *gin.Context) {
	var l LoginStep1Request
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//get user details
	req := &pb.GetUserByIDOrEmailRequest{
		Base:  &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
		Email: l.Email,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email could not be found in db", cerr.GeneralError))
		return
	}

	if user.TfaConfirmed {
		//user already did 2fa registration, so the tfa code is mandatory
		if l.TfaCode == "" {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.TfaAlreadyConfirmed, "already registered", "confirm_2FAregsitration.2FACode.already_done"))
			return
		}

		//user did tfa registration and passed in code, so we need to check that code
		req2FA := &pb.AuthenticateRequest{
			Base:   &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
			Code:   l.TfaCode,
			Secret: user.TfaSecret,
		}
		resp2FA, err := twoFAClient.Authenticate(c, req2FA)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error doing 2FA", cerr.GeneralError))
			return
		}

		if !resp2FA.Result {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.InvalidArgument, "2FA code is invalid", "login.2FACode.invalid"))
			return
		}
	}

	//if we reach this, the user either is not TFSEnabled, or he entered a correct tfa code
	//if he did not pass a tfa code he MUST have TfaConfirmed==false. Therefore we will pass out the security data from the user
	//although there is a smal gap in here: if the user did not finish the 2fa registration but created the account (in register)
	//the data will be readable by only passing in the email

	idRequest := &pb.IDRequest{
		Base: &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
		Id:   user.Id,
	}
	s, err := dbClient.GetUserSecurities(c, idRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user data", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User from could not be found in db", cerr.GeneralError))
		return
	}

	authMiddlewareSimple.SetAuthHeader(c, user.Id)

	sep10ChallangeTX, err := getSEP10Challenge(user.PublicKey_0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "error generating challange :"+err.Error(), cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &LoginStep1Response{
		KdfPasswordSalt:               s.KdfSalt,
		EncryptedMnemonicMasterKey:    s.MnemonicMasterKey,
		MnemonicMasterKeyEncryptionIV: s.MnemonicMasterIv,
		EncryptedMnemonic:             s.Mnemonic,
		MnemonicEncryptionIV:          s.MnemonicIv,
		EncryptedWordlistMasterKey:    s.WordlistMasterKey,
		WordlistMasterKeyEncryptionIV: s.WordlistMasterIv,
		EncryptedWordlist:             s.Wordlist,
		WordlistEncryptionIV:          s.WordlistIv,
		PublicKeyIndex0:               s.PublicKey_0,
		SEP10TransactionChallenge:     sep10ChallangeTX,
	})
}

//LoginStep2Request is the data needed for the second step of the login
type LoginStep2Request struct {
	/*Key string `form:"key" json:"key" validate:"required"`
	PublicKey0 string `form:"public_key_0" json:"public_key_0" validate:"required"`*/
	Key              string `form:"key" json:"key"`
	SEP10Transaction string `form:"sep10_transaction" json:"sep10_transaction"`
}

//LoginStep2Response to the api
type LoginStep2Response struct {
	PaymentState      string `json:"payment_state"`
	TfaSecret         string `json:"tfa_secret"`
	TfaQRImage        []byte `json:"tfa_qr_image"`
	MailConfirmed     bool   `json:"mail_confirmed"`
	TfaConfirmed      bool   `json:"tfa_confirmed"`
	MnemonicConfirmed bool   `json:"mnemonic_confirmed"`
}

//LoginStep2 is the first step of the login
func LoginStep2(uc *mw.IcopContext, c *gin.Context) {
	var l LoginStep2Request
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//check that the passed in key matches the saved password in the userprofile
	user := mw.GetAuthUser(c)
	req := &pb.GetUserByIDOrEmailRequest{
		Base: &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
		Id:   user.UserID,
	}
	u, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if u.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User could not be found in db", cerr.GeneralError))
		return
	}

	if l.SEP10Transaction == "" {
		match := CheckPasswordHash(uc.Log, l.Key, u.Password)
		if !match {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("key", cerr.InvalidArgument, "Can not login user, public key is invalid", "loginStep2.key.invalid"))
			return
		}
	} else {
		valid, _, err := verifySEP10Data(l.SEP10Transaction, user.UserID, uc, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, err.Error(), cerr.GeneralError))
			return
		}
		if !valid {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("transaction", cerr.InvalidArgument, "could not validate challange transaction", ""))
			return
		}
	}

	ret := &LoginStep2Response{
		MailConfirmed:     user.MailConfirmed,
		TfaConfirmed:      user.TfaConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
		PaymentState:      u.PaymentState,
	}

	if user.TfaConfirmed {
		//if confirmed, we don't create the image, because it's not needed any longer
		//authentication of the tfa code was allready done in LoginStep1
		authMiddlewareFull.SetAuthHeader(c, u.Id)
	} else {
		//if not confirmed, recreate image and send this back to the client, which will continue the setup process
		qr2FA, err := twoFAClient.FromSecret(c, &pb.FromSecretRequest{
			Base:   &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
			Email:  user.Email,
			Secret: u.TfaTempSecret, // need the temp-secrete here
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting 2FA data", cerr.GeneralError))
			return
		}
		ret.TfaQRImage = qr2FA.Bitmap
		ret.TfaSecret = u.TfaTempSecret

		authMiddlewareSimple.SetAuthHeader(c, u.Id)
	}

	c.JSON(http.StatusOK, ret)
}

//CheckPasswordHash check a given password to the hashed value
func CheckPasswordHash(log *logrus.Entry, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.WithError(err).Error("Error checking password")
	}
	return err == nil
}

func getSEP10Challenge(account string) (string, error) {
	now := time.Now()
	validTo := now.Add(time.Second * 300)
	var keyName = helpers.RandomString(50) + " auth"

	value := make([]byte, 64)
	_, err := rand.Read(value)
	if err != nil {
		return "", fmt.Errorf("Could not create random string: %s", err.Error())
	}

	//create challange
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: cnf.AuthServerSigningAccountPK},
		build.Network{Passphrase: cnf.StellarNetworkPassphrase},
		build.Sequence{Sequence: 0},
		build.Timebounds{
			MinTime: uint64(now.Unix()), MaxTime: uint64(validTo.Unix()),
		},
		build.SetData(
			keyName, value,
			build.SourceAccount{AddressOrSeed: account},
		),
	)
	if err != nil {
		return "", fmt.Errorf("Error creating transaction: %s", err.Error())
	}

	txe, err := tx.Sign(cnf.AuthServerSigningAccountSeed)
	txeStr, err := txe.Base64()
	if err != nil {
		return "", fmt.Errorf("Error base64 encoding tx: %s", err.Error())
	}

	return txeStr, nil
}

//verifySEP10Data verifies the SEP10 transaction and returns the user publickey and validity
//also checks, that the given UserID matches to the data in the DB (based on the user account/pubKey0)
func verifySEP10Data(txStr string, userID int64, uc *mw.IcopContext, c *gin.Context) (bool, string, error) {
	var userPK string
	if cnf.SEP10FakeLoginEnabled && txStr == "I am the king of the world" {
		//read key from simple jwt
		userSec, err := dbClient.GetUserSecurities(c, &pb.IDRequest{
			Base: NewBaseRequest(uc),
			Id:   userID,
		})
		if err != nil {
			return false, "", fmt.Errorf("error reading user securities: %s", err.Error())
		}

		if userSec.UserNotFound {
			return false, "", fmt.Errorf("User-Sec could not be found in db: %s", err.Error())
		}
		userPK = userSec.PublicKey_0
	} else {
		txe, err := decodeFromBase64(txStr)
		if err != nil {
			return false, "", fmt.Errorf("Error base64 decoding tx: %s", err.Error())
		}
		var tx xdr.Transaction
		tx = txe.E.Tx
		if tx.SourceAccount.Address() != cnf.AuthServerSigningAccountPK {
			return false, "", fmt.Errorf("tx source invalid")
		}

		now := xdr.Uint64(time.Now().Unix())
		if now < tx.TimeBounds.MinTime || tx.TimeBounds.MinTime == 0 {
			return false, "", fmt.Errorf("tx not valid yet")
		}
		if now > tx.TimeBounds.MaxTime || tx.TimeBounds.MaxTime == 0 {
			return false, "", fmt.Errorf("tx not valid any more")
		}

		if len(tx.Operations) != 1 {
			return false, "", fmt.Errorf("invalid operation count")
		}

		op := tx.Operations[0]
		if op.Body.Type != xdr.OperationTypeManageData {
			return false, "", fmt.Errorf("invalid operation type")
		}

		if op.SourceAccount == nil {
			return false, "", fmt.Errorf("no source account")
		}

		userPK = op.SourceAccount.Address()

		//check sgnatures
		if txe.E.Signatures == nil || len(txe.E.Signatures) != 2 {
			return false, "", fmt.Errorf("wrong signature amount")
		}

		serverKeyPair, err := keypair.Parse(cnf.AuthServerSigningAccountSeed)
		if err != nil {
			return false, "", fmt.Errorf("could not parse server key: %s", err.Error())
		}

		userKeyPair, err := keypair.Parse(userPK)
		if err != nil {
			return false, "", fmt.Errorf("could not parse user key: %s", err.Error())
		}

		hash32, err := network.HashTransaction(&txe.E.Tx, cnf.StellarNetworkPassphrase)
		txHash := hash32[:]

		err = serverKeyPair.Verify(txHash, txe.E.Signatures[0].Signature)
		if err != nil {
			return false, "", fmt.Errorf("could not verify server signature: %s", err.Error())
		}

		err = userKeyPair.Verify(txHash, txe.E.Signatures[1].Signature)
		if err != nil {
			return false, "", fmt.Errorf("could not verify user signature: %s", err.Error())
		}
	}

	//check that userPK matches transaction data
	userSec, err := dbClient.GetUserSecurities(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	})
	if err != nil {
		return false, "", fmt.Errorf("error reading user securities: %s", err.Error())
	}

	if userSec.UserNotFound {
		return false, "", fmt.Errorf("User-Sec could not be found in db: %s", err.Error())
	}

	if userSec.PublicKey_0 != userPK {
		return false, "", fmt.Errorf("account does not match user-data")
	}

	return true, userPK, nil
}

// DecodeFromBase64 decodes the transaction from a base64 string into a TransactionEnvelopeBuilder
func decodeFromBase64(encodedXdr string) (*build.TransactionEnvelopeBuilder, error) {
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
	n := build.Network{Passphrase: cnf.StellarNetworkPassphrase}
	err = txEnvelopeBuilder.MutateTX(n)
	if err != nil {
		return nil, err
	}

	return &txEnvelopeBuilder, nil
}
