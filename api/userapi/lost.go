package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/Soneso/lumenshine-backend/constants"
	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

//LostPasswordRequest is the data needed for the first step of the login
//swagger:parameters LostPasswordRequest LostPassword
type LostPasswordRequest struct {
	//required: true
	Email string `form:"email" json:"email" validate:"required,icop_email"`
}

//LostPassword called from the API when the user initiates the lost password function
// swagger:route POST /portal/user/lost_password lost LostPassword
//
// Called from the API when the user initiates the lost password function
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func LostPassword(uc *mw.IcopContext, c *gin.Context) {
	var l LostPasswordRequest
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
		Base:  NewBaseRequest(uc),
		Email: l.Email,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, user not found", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email not found", cerr.GeneralError))
		return
	}

	if !user.MailConfirmed {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.EmailNotConfigured, "Email address is not verified", ""))
		return
	}

	//set the confirmation keys
	reqConf := &pb.SetMailTokenRequest{
		Base:                   NewBaseRequest(uc),
		UserId:                 user.Id,
		MailConfirmationKey:    helpers.RandomString(constants.DefaultMailkeyLength),
		MailConfirmationExpiry: time.Now().AddDate(0, 0, constants.DefaultMailkeyExpiryDays).Unix(),
	}
	_, err = dbClient.SetUserMailToken(c, reqConf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not set token", cerr.GeneralError))
		return
	}

	//Send email
	ur, err := dbClient.GetUserProfile(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, profile not found", cerr.GeneralError))
		return
	}
	msgBody := RenderTemplateToString(uc, c, "lost_password_mail", gin.H{
		"Forename": ur.Forename,
		"Lastname": ur.Lastname,
		"TokeUrl":  cnf.WebLinks.LostPassword + reqConf.MailConfirmationKey,
		"TokenValidTo": helpers.TimeToString(
			time.Unix(reqConf.MailConfirmationExpiry, 0), uc.Language,
		),
	})
	msgSubject := fmt.Sprintf("%s :: Reset password", cnf.Site.SiteName)

	_, err = mailClient.SendMail(c, &pb.SendMailRequest{
		Base:    NewBaseRequest(uc),
		From:    cnf.Site.EmailSender,
		To:      ur.Email,
		Subject: msgSubject,
		Body:    msgBody,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, could not send mail", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//Need2FAResetPasswordResponse returned from the API
// swagger:model
type Need2FAResetPasswordResponse struct {
	Need2FAResetPassword bool `json:"need2fa_reset_pwd"`
}

//Need2FAResetPassword - returns the flag
// swagger:route GET /portal/user/auth/need_2fa_reset_pwd lost Need2FAResetPassword
//
// Returns the flag
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Need2FAResetPasswordResponse
func Need2FAResetPassword(uc *mw.IcopContext, c *gin.Context) {
	userID := mw.GetAuthUser(c).UserID

	req := &pb.GetUserByIDOrEmailRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, user not found", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email not found", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &Need2FAResetPasswordResponse{
		Need2FAResetPassword: !user.Reset2FaByAdmin,
	})
}

//LostPasswordTfaRequest is the data needed for the first step of the login
//The user clicked the confirm link and thus called the /ico/confirm_token function, which sets the simple-auth-token
//swagger:parameters LostPasswordTfaRequest LostPasswordTfa
type LostPasswordTfaRequest struct {
	//required: true
	TfaCode string `form:"tfa_code" json:"tfa_code" validate:"required"`
}

//LostPasswordTfaResponse returned from the API
// swagger:model
type LostPasswordTfaResponse struct {
	PublicKey0 string `json:"public_key_0"`
}

//LostPasswordTfa called from the API when the user initiates the lost password TFA function
// swagger:route POST /portal/user/auth/lost_password_tfa lost LostPasswordTfa
//
// Called from the API when the user initiates the lost password TFA function
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: LostPasswordTfaResponse
func LostPasswordTfa(uc *mw.IcopContext, c *gin.Context) {
	var l LostPasswordTfaRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	userID := mw.GetAuthUser(c).UserID

	req := &pb.GetUserByIDOrEmailRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, user not found", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email not found", cerr.GeneralError))
		return
	}

	if !user.TfaConfirmed {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.TfaNotYetConfirmed, "2FA secret not confirmed yet", ""))
		return
	}

	req2FA := &pb.AuthenticateRequest{
		Base:   NewBaseRequest(uc),
		Code:   l.TfaCode,
		Secret: user.TfaSecret,
	}
	resp2FA, err := twoFAClient.Authenticate(c, req2FA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, 2FA authentication failed", cerr.GeneralError))
		return
	}

	if !resp2FA.Result {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.InvalidArgument, "2FA code is invalid", "login.2FACode.invalid"))
		return
	}

	userSec, err := dbClient.GetUserSecurities(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, 2FA securities", cerr.GeneralError))
		return
	}

	if userSec.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, 2FA user not found", cerr.GeneralError))
		return
	}

	//we set the special token, because the user had the correct email and also did the 2fa registration
	//we need this to secure the next call, where the user updates the security data
	authMiddlewareLostPwd.SetAuthHeader(c, userID)

	c.JSON(http.StatusOK, &LostPasswordTfaResponse{
		PublicKey0: userSec.PublicKey_0,
	})
}

//LostPasswordUpdateRequest is the new secuurity data for the user
//swagger:parameters LostPasswordUpdateRequest LostPasswordUpdate
type LostPasswordUpdateRequest struct {
	//required: true
	KDFSalt string `form:"kdf_salt" json:"kdf_salt" validate:"required,base64,len=44"`
	//required: true
	MnemonicMasterKey string `form:"mnemonic_master_key" json:"mnemonic_master_key" validate:"required,base64,len=44"`
	//required: true
	MnemonicMasterIV string `form:"mnemonic_master_iv" json:"mnemonic_master_iv" validate:"required,base64,len=24"`
	//required: true
	EncryptedMnemonic string `form:"encrypted_mnemonic" json:"encrypted_mnemonic" validate:"required,base64,len=64"`
	//required: true
	EncryptedMnemonicIV string `form:"encryption_mnemonic_iv" json:"encryption_mnemonic_iv" validate:"required,base64,len=24"`
	//required: true
	WordlistMasterKey string `form:"wordlist_master_key" json:"wordlist_master_key" validate:"required,base64,len=44"`
	//required: true
	WordlistMasterIV string `form:"wordlist_master_iv" json:"wordlist_master_iv" validate:"required,base64,len=24"`
	//required: true
	EncryptedWordlist string `form:"encrypted_wordlist" json:"encrypted_wordlist" validate:"required,base64"`
	//required: true
	EncryptionWordlistIV string `form:"encryption_wordlist_iv" json:"encryption_wordlist_iv" validate:"required,base64,len=24"`
	//required: true
	PublicKey0       string `form:"public_key_0" json:"public_key_0" validate:"required,base64,len=56"`
	PublicKey188     string `form:"public_key_188" json:"public_key_188"`
	SEP10Transaction string `form:"sep10_transaction" json:"sep10_transaction"`
}

//LostPasswordUpdate updates the security data for the user
// swagger:route POST /portal/user/auth2/lost_password_update lost LostPasswordUpdate
//
// Updates the security data for the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func LostPasswordUpdate(uc *mw.IcopContext, c *gin.Context) {
	var l LostPasswordUpdateRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	user := mw.GetAuthUser(c)

	if l.SEP10Transaction == "" {
		//check that public key 188 is correct
		match := CheckPasswordHash(uc.Log, l.PublicKey188, user.Password)
		if !match {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_188", cerr.InvalidPassword, "Password does not match", ""))
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

	req := &pb.UserSecurityRequest{
		Base:              NewBaseRequest(uc),
		UserId:            user.UserID,
		KdfSalt:           l.KDFSalt,
		MnemonicMasterKey: l.MnemonicMasterKey,
		MnemonicMasterIv:  l.MnemonicMasterIV,
		WordlistMasterKey: l.WordlistMasterKey,
		WordlistMasterIv:  l.WordlistMasterIV,
		Mnemonic:          l.EncryptedMnemonic,
		MnemonicIv:        l.EncryptedMnemonicIV,
		Wordlist:          l.EncryptedWordlist,
		WordlistIv:        l.EncryptionWordlistIV,
		PublicKey_0:       l.PublicKey0,
		PublicKey_188:     l.PublicKey188,
	}
	_, err := dbClient.SetUserSecurities(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, update security data", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ChangePasswordUpdateRequest is the new security data for the user
//swagger:parameters ChangePasswordUpdateRequest ChangePasswordUpdate
type ChangePasswordUpdateRequest struct {
	//required: true
	KDFSalt string `form:"kdf_salt" json:"kdf_salt" validate:"required"`
	//required: true
	MnemonicMasterKey string `form:"mnemonic_master_key" json:"mnemonic_master_key" validate:"required,base64,len=44"`
	//required: true
	MnemonicMasterIV string `form:"mnemonic_master_iv" json:"mnemonic_master_iv" validate:"required,base64,len=24"`
	//required: true
	WordlistMasterKey string `form:"wordlist_master_key" json:"wordlist_master_key" validate:"required,base64,len=44"`
	//required: true
	WordlistMasterIV string `form:"wordlist_master_iv" json:"wordlist_master_iv" validate:"required,base64,len=24"`

	PublicKey188     string `form:"public_key_188" json:"public_key_188"`
	SEP10Transaction string `form:"sep10_transaction" json:"sep10_transaction"`
}

//ChangePasswordUpdate updates the security data for the user
// swagger:route POST /portal/user/dashboard/change_password lost ChangePasswordUpdate
//
// Updates the security data for the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func ChangePasswordUpdate(uc *mw.IcopContext, c *gin.Context) {
	var l ChangePasswordUpdateRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	user := mw.GetAuthUser(c)

	if l.SEP10Transaction == "" {
		//check that the password_188 was correct
		match := CheckPasswordHash(uc.Log, l.PublicKey188, user.Password)
		if !match {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_188", cerr.InvalidPassword, "password does not match", ""))
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

	req := &pb.UserSecurityRequest{
		Base:              NewBaseRequest(uc),
		UserId:            user.UserID,
		KdfSalt:           l.KDFSalt,
		MnemonicMasterKey: l.MnemonicMasterKey,
		MnemonicMasterIv:  l.MnemonicMasterIV,
		WordlistMasterKey: l.WordlistMasterKey,
		WordlistMasterIv:  l.WordlistMasterIV,
	}
	_, err := dbClient.UpdateUserSecurityPassword(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, update security data", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//LostTfaRequest initiates the reset 2fa function
//swagger:parameters LostTfaRequest LostTfa
type LostTfaRequest struct {
	//required: true
	Email string `form:"email" json:"email" validate:"required,icop_email"`
}

//LostTfa called from the API when the user initiates the reset 2fa function
// swagger:route POST /portal/user/dashboard/change_password lost LostTfa
//
// Called from the API when the user initiates the reset 2fa function
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func LostTfa(uc *mw.IcopContext, c *gin.Context) {
	var l LostTfaRequest
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
		Base:  NewBaseRequest(uc),
		Email: l.Email,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, user not found", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email not found", cerr.GeneralError))
		return
	}

	if !user.MailConfirmed {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.EmailNotConfigured, "Email address is not verified", ""))
		return
	}

	//set the confirmation keys
	reqConf := &pb.SetMailTokenRequest{
		Base:                   NewBaseRequest(uc),
		UserId:                 user.Id,
		MailConfirmationKey:    helpers.RandomString(constants.DefaultMailkeyLength),
		MailConfirmationExpiry: time.Now().AddDate(0, 0, constants.DefaultMailkeyExpiryDays).Unix(),
	}
	_, err = dbClient.SetUserMailToken(c, reqConf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not set token", cerr.GeneralError))
	}

	ur, err := dbClient.GetUserProfile(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, profile not found", cerr.GeneralError))
		return
	}
	msgBody := RenderTemplateToString(uc, c, "lost_tfa_mail", gin.H{
		"Forename": ur.Forename,
		"Lastname": ur.Lastname,
		"TokeUrl":  cnf.WebLinks.LostTFA + reqConf.MailConfirmationKey,
		"TokenValidTo": helpers.TimeToString(
			time.Unix(reqConf.MailConfirmationExpiry, 0), uc.Language,
		),
	})
	msgSubject := fmt.Sprintf("%s :: Reset 2FA Secret", cnf.Site.SiteName)

	_, err = mailClient.SendMail(c, &pb.SendMailRequest{
		Base:    NewBaseRequest(uc),
		From:    cnf.Site.EmailSender,
		To:      ur.Email,
		Subject: msgSubject,
		Body:    msgBody,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not send mail", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//NewTfaRequest for proving the password
//swagger:parameters NewTfaRequest NewTfaUpdate
type NewTfaRequest struct {
	PublicKey188     string `form:"public_key_188" json:"public_key_188"`
	SEP10Transaction string `form:"sep10_transaction" json:"sep10_transaction"`
}

//NewTfaResponse response for update
// swagger:model
type NewTfaResponse struct {
	TFASecret  string `json:"tfa_secret,omitempty"`
	TFAQrImage string `json:"tfa_qr_image,omitempty"`
}

//NewTfaUpdate called for updating the tfa data
// swagger:route POST /portal/user/dashboard/new_2fa_secret lost NewTfaUpdate
//
// Called for updating the tfa data
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: NewTfaResponse
func NewTfaUpdate(uc *mw.IcopContext, c *gin.Context) {
	var l NewTfaRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	user := mw.GetAuthUser(c)

	if l.PublicKey188 == "" {
		req := &pb.GetUserByIDOrEmailRequest{
			Base: NewBaseRequest(uc),
			Id:   user.UserID,
		}
		user, err := dbClient.GetUserDetails(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, user not found", cerr.GeneralError))
			return
		}
		if !user.Reset2FaByAdmin {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_188", cerr.MissingMandatoryField, "Missing parameter", ""))
			return
		}
	} else {
		if l.SEP10Transaction == "" {
			//check that public key 188 is correct
			match := CheckPasswordHash(uc.Log, l.PublicKey188, user.Password)
			if !match {
				c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_188", cerr.InvalidPassword, "Password does not match", ""))
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
	}

	//get 2fa data for email
	req2FA := &pb.NewSecretRequest{
		Base:  NewBaseRequest(uc),
		Email: user.Email,
	}
	tfaResp, err := twoFAClient.NewSecret(c, req2FA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not generate 2FA code", cerr.GeneralError))
		return
	}

	req := &pb.SetTempTfaSecretRequest{
		Base:      NewBaseRequest(uc),
		UserId:    user.UserID,
		TfaSecret: tfaResp.Secret,
	}
	_, err = dbClient.SetTempTfaSecret(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not set 2FA secret", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &NewTfaResponse{
		TFAQrImage: base64.StdEncoding.EncodeToString(tfaResp.Bitmap),
		TFASecret:  tfaResp.Secret,
	})
}

//UpdateSecurityDataRequest is the data needed for updateing the data
//swagger:parameters UpdateSecurityDataRequest UpdateSecurityData
type UpdateSecurityDataRequest struct {
	//required: true
	KDFSalt string `form:"kdf_salt" json:"kdf_salt" validate:"required,base64,len=44"`
	//required: true
	MnemonicMasterKey string `form:"mnemonic_master_key" json:"mnemonic_master_key" validate:"required,base64,len=44"`
	//required: true
	MnemonicMasterIV string `form:"mnemonic_master_iv" json:"mnemonic_master_iv" validate:"required,base64,len=24"`
	//required: true
	WordlistMasterKey string `form:"wordlist_master_key" json:"wordlist_master_key" validate:"required,base64,len=44"`
	//required: true
	WordlistMasterIV string `form:"wordlist_master_iv" json:"wordlist_master_iv" validate:"required,base64,len=24"`
	//required: true
	Mnemonic string `form:"mnemonic" json:"mnemonic" validate:"required,base64,len=64"`
	//required: true
	MnemonicIV string `form:"mnemonic_iv" json:"mnemonic_iv" validate:"required,base64,len=24"`
	//required: true
	Wordlist string `form:"wordlist" json:"wordlist" validate:"required,base64"`
	//required: true
	WordlistIV string `form:"wordlist_iv" json:"wordlist_iv" validate:"required,base64,len=24"`
	//required: true
	PublicKey0   string `form:"public_key_0" json:"public_key_0" validate:"required,base64,len=56"`
	PublicKey188 string `form:"public_key_188" json:"public_key_188"`

	TfaCode string `form:"tfa_code" json:"tfa_code"`
}

//UpdateSecurityData will update the security data in the DB
// swagger:route POST /portal/user/auth/update_security_data lost UpdateSecurityData
//
// Will update the security data in the DB
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateSecurityData(uc *mw.IcopContext, c *gin.Context) {
	user := mw.GetAuthUser(c)

	var l UpdateSecurityDataRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if user.MnemonicConfirmed {
		c.JSON(http.StatusForbidden, cerr.NewIcopError("menmeonic", cerr.GeneralError, "Mnemonic is already confirmed", ""))
		return
	}

	if user.TfaConfirmed {

		if l.TfaCode == "" {
			c.JSON(http.StatusForbidden, cerr.NewIcopError("tfa_code", cerr.MissingMandatoryField, "Missing parameter: 2FA code", ""))
			return
		}

		//check that TFA is correct
		req2FA := &pb.AuthenticateRequest{
			Base:   &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName},
			Code:   l.TfaCode,
			Secret: user.TfaSecret,
		}
		resp2FA, err := twoFAClient.Authenticate(c, req2FA)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not authenticate 2FA", cerr.GeneralError))
			return
		}

		if !resp2FA.Result {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.InvalidArgument, "2FA code is invalid", ""))
			return
		}
	}

	req := &pb.UserSecurityRequest{
		Base:              NewBaseRequest(uc),
		UserId:            user.UserID,
		KdfSalt:           l.KDFSalt,
		MnemonicMasterKey: l.MnemonicMasterKey,
		MnemonicMasterIv:  l.MnemonicMasterIV,
		WordlistMasterKey: l.WordlistMasterKey,
		WordlistMasterIv:  l.WordlistMasterIV,
		Mnemonic:          l.Mnemonic,
		MnemonicIv:        l.MnemonicIV,
		Wordlist:          l.Wordlist,
		WordlistIv:        l.WordlistIV,
		PublicKey_0:       l.PublicKey0,
		PublicKey_188:     l.PublicKey188,
	}
	_, err := dbClient.SetUserSecurities(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Server error, can not update data", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
