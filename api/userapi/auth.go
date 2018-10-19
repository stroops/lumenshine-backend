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
	})
}

//LoginStep2Request is the data needed for the second step of the login
type LoginStep2Request struct {
	Key string `form:"key" json:"key" validate:"required"`
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

	match := CheckPasswordHash(uc.Log, l.Key, u.Password)
	if !match {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("key", cerr.InvalidArgument, "Can not login user, public key is invalid", "loginStep2.key.invalid"))
		return
	}

	ret := &LoginStep2Response{
		MailConfirmed:     user.MailConfirmed,
		TfaConfirmed:      user.TfaConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
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

//LoginSEP10Request is the data needed for creating the challenge
// swagger:parameters LoginSEP10Request LoginSEP10Get
type LoginSEP10Request struct {
	// Account is the stellar account wanted to log in
	// required: true
	Account string `form:"account" json:"account" validate:"required"`
}

//LoginSEP10Response to the api
type LoginSEP10Response struct {
	Transaction string `json:"transaction"`
}

// LoginSEP10Get returns the SEP10 challenge for the account
// swagger:route GET /portal/auth/login LoginSEP10Get
//     returns the SEP10 challenge for the account
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: LoginSEP10Response
func LoginSEP10Get(uc *mw.IcopContext, c *gin.Context) {
	var l LoginSEP10Request
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	now := time.Now()
	validTo := now.Add(time.Second * 300)
	var keyName = helpers.RandomString(50) + " auth"

	value := make([]byte, 64)
	_, err := rand.Read(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not create random string", cerr.GeneralError))
		return
	}
	//value := base64.StdEncoding.EncodeToString(b)

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
			build.SourceAccount{AddressOrSeed: l.Account},
		),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating transaction", cerr.GeneralError))
		return
	}

	txe, err := tx.Sign(cnf.AuthServerSigningAccountSeed)
	txeStr, err := txe.Base64()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error base64 encoding tx", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &LoginSEP10Response{Transaction: txeStr})
}

//LoginSEP10PostRequest is the data needed for doing the validation
// swagger:parameters LoginSEP10PostRequest LoginSEP10Post
type LoginSEP10PostRequest struct {
	// Transaction is the (user)signed transaction from the client
	// required: true
	Transaction string `form:"transaction" json:"transaction" validate:"required"`
}

//LoginSEP10PostResponse valid jwt token
type LoginSEP10PostResponse struct {
	Token string `json:"token"`
}

// LoginSEP10Post validates a login transaction and returns a valid full authenticated JWT token on success
// swagger:route GET /portal/auth/login LoginSEP10Get
//     validates a login transaction and returns a valid full authenticated JWT token on success
//     Consumes:
//     - application/x-www-form-urlencoded
//	   - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: LoginSEP10PostResponse
func LoginSEP10Post(uc *mw.IcopContext, c *gin.Context) {
	var l LoginSEP10PostRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	txe, err := decodeFromBase64(l.Transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error base64 decoding tx", cerr.GeneralError))
		return
	}
	var tx xdr.Transaction
	tx = txe.E.Tx
	if tx.SourceAccount.Address() != cnf.AuthServerSigningAccountPK {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "tx source invalid", cerr.GeneralError))
		return
	}

	now := xdr.Uint64(time.Now().Unix())
	if now < tx.TimeBounds.MinTime || tx.TimeBounds.MinTime == 0 {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "tx not valid yet", cerr.GeneralError))
		return
	}
	if now > tx.TimeBounds.MaxTime || tx.TimeBounds.MaxTime == 0 {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "tx not valid any more", cerr.GeneralError))
		return
	}

	if len(tx.Operations) != 1 {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "invalid operation count", cerr.GeneralError))
		return
	}

	op := tx.Operations[0]
	if op.Body.Type != xdr.OperationTypeManageData {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "invalid operation type", cerr.GeneralError))
		return
	}

	if op.SourceAccount == nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "invalid operation source", cerr.GeneralError))
		return
	}

	userPK := op.SourceAccount.Address()

	//check sgnatures
	if txe.E.Signatures == nil || len(txe.E.Signatures) != 2 {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "wrong signature amount", cerr.GeneralError))
		return
	}

	serverKeyPair, err := keypair.Parse(cnf.AuthServerSigningAccountSeed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "could not parse server key", cerr.GeneralError))
		return
	}

	userKeyPair, err := keypair.Parse(userPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "could not parse server key", cerr.GeneralError))
		return
	}

	hash32, err := network.HashTransaction(&txe.E.Tx, cnf.StellarNetworkPassphrase)
	txHash := hash32[:]

	err = serverKeyPair.Verify(txHash, txe.E.Signatures[0].Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "could not verify server signature", cerr.GeneralError))
		return
	}

	err = userKeyPair.Verify(txHash, txe.E.Signatures[1].Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "could not verify user signature", cerr.GeneralError))
		return
	}

	//generate JWT and return it

	c.JSON(http.StatusOK, &LoginSEP10PostResponse{Token: ""})
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
