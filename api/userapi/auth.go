package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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
		Base:       &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
		Id:         user.UserID,
		WithQrData: true,
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

	if user.TfaConfirmed {
		authMiddlewareFull.SetAuthHeader(c, u.Id)
	} else {
		authMiddlewareSimple.SetAuthHeader(c, u.Id)
	}

	c.JSON(http.StatusOK, &LoginStep2Response{
		TfaSecret:         u.TfaSecret,
		TfaQRImage:        u.TfaQrcode,
		MailConfirmed:     user.MailConfirmed,
		TfaConfirmed:      user.TfaConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
	})
}

//CheckPasswordHash check a given password to the hashed value
func CheckPasswordHash(log *logrus.Entry, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.WithError(err).Error("Error checking password")
	}
	return err == nil
}
