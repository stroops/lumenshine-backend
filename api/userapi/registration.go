package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/Soneso/lumenshine-backend/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//RegisterUserRequest is the data needed for registration
type RegisterUserRequest struct {
	Email             string `form:"email" json:"email" validate:"required,icop_email"`
	KDFSalt           string `form:"kdf_salt" json:"kdf_salt" validate:"required,base64,len=44"`
	MnemonicMasterKey string `form:"mnemonic_master_key" json:"mnemonic_master_key" validate:"required,base64,len=44"`
	MnemonicMasterIV  string `form:"mnemonic_master_iv" json:"mnemonic_master_iv" validate:"required,base64,len=24"`
	WordlistMasterKey string `form:"wordlist_master_key" json:"wordlist_master_key" validate:"required,base64,len=44"`
	WordlistMasterIV  string `form:"wordlist_master_iv" json:"wordlist_master_iv" validate:"required,base64,len=24"`
	Mnemonic          string `form:"mnemonic" json:"mnemonic" validate:"required,base64,len=64"`
	MnemonicIV        string `form:"mnemonic_iv" json:"mnemonic_iv" validate:"required,base64,len=24"`
	Wordlist          string `form:"wordlist" json:"wordlist" validate:"required,base64"`
	WordlistIV        string `form:"wordlist_iv" json:"wordlist_iv" validate:"required,base64,len=24"`

	PublicKey0   string `form:"public_key_0" json:"public_key_0" validate:"required,base64,len=56"`
	PublicKey188 string `form:"public_key_188" json:"public_key_188" validate:"required,base64,len=56"`

	Salutation    string `form:"salutation" json:"salutation"`
	Title         string `form:"title" json:"title"`
	Forename      string `form:"forename" json:"forename"`
	Lastname      string `form:"lastname" json:"lastname"`
	Company       string `form:"company" json:"company"`
	StreetAddress string `form:"street_address" json:"street_address"`
	StreetNumber  string `form:"street_number" json:"street_number"`
	ZipCode       string `form:"zip_code" json:"zip_code"`
	City          string `form:"city" json:"city"`
	State         string `form:"state" json:"state"`
	CountryCode   string `form:"country_code" json:"country_code"`
	Nationality   string `form:"nationality" json:"nationality"`
	MobileNr      string `form:"mobile_nr" json:"mobile_nr" validate:"omitempty,icop_phone"`
	BirthDay      string `form:"birth_day" json:"birth_day"`
	BirthPlace    string `form:"birth_place" json:"birth_place"`
}

//RegisterUserResponse response for registration
type RegisterUserResponse struct {
	TFASecret  string `json:"tfa_secret,omitempty"`
	TFAQrImage string `json:"tfa_qr_image,omitempty"`
}

//RegisterUser registers and creates the user in the db
//func RegisterUser(uc *mw.IcopContext, c *gin.Context) {
func RegisterUser(uc *mw.IcopContext, c *gin.Context) {
	ur := new(RegisterUserRequest)

	if err := c.Bind(ur); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, ur); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//check email does not exist
	req := &pb.ExistsEmailRequest{
		Base:  NewBaseRequest(uc),
		Email: ur.Email,
	}
	exists, err := dbClient.ExistsEmail(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Email could not be checked", cerr.GeneralError))
		return
	}
	if exists.Exists {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.EmailExists, "Email already exists", ""))
		return
	}

	//TODO: do additional checks

	//get 2fa data for email
	req2FA := &pb.NewSecretRequest{
		Base:  NewBaseRequest(uc),
		Email: ur.Email,
	}
	twoFaResp, err := twoFAClient.NewSecret(c, req2FA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error generating 2FAcode", cerr.GeneralError))
		return
	}

	//hash the password
	pwd, err := bcrypt.GenerateFromPassword([]byte(ur.PublicKey188), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error generating password", cerr.GeneralError))
		return
	}

	//get the birthday
	birthDaty, err := time.Parse("2006-01-02", ur.BirthDay)
	if ur.BirthDay != "" && err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("birth_day", cerr.InvalidArgument, "Birthday wrong format", ""))
		return
	}

	reqC := &pb.CreateUserRequest{
		Base:                   NewBaseRequest(uc),
		Email:                  ur.Email,
		MailConfirmationKey:    helpers.RandomString(constants.DefaultMailkeyLength),
		MailConfirmationExpiry: time.Now().AddDate(0, 0, constants.DefaultMailkeyExpiryDays).Unix(),
		TfaTempSecret:          twoFaResp.Secret,
		Salutation:             ur.Salutation,
		Title:                  ur.Title,
		Forename:               ur.Forename,
		Lastname:               ur.Lastname,
		Company:                ur.Company,
		StreetAddress:          ur.StreetAddress,
		StreetNumber:           ur.StreetNumber,
		ZipCode:                ur.ZipCode,
		City:                   ur.City,
		State:                  ur.State,
		CountryCode:            ur.CountryCode,
		Nationality:            ur.Nationality,
		MobileNr:               ur.MobileNr,
		BirthDay:               birthDaty.Unix(),
		BirthPlace:             ur.BirthPlace,
		Password:               string(pwd),
		KdfSalt:                ur.KDFSalt,
		MnemonicMasterKey:      ur.MnemonicMasterKey,
		MnemonicMasterIv:       ur.MnemonicMasterIV,
		WordlistMasterKey:      ur.WordlistMasterKey,
		WordlistMasterIv:       ur.WordlistMasterIV,
		Mnemonic:               ur.Mnemonic,
		MnemonicIv:             ur.MnemonicIV,
		Wordlist:               ur.Wordlist,
		WordlistIv:             ur.WordlistIV,
		PublicKey_0:            ur.PublicKey0,
		PublicKey_188:          ur.PublicKey188,
	}
	user, err := dbClient.CreateUser(c, reqC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating user", cerr.GeneralError))
		return
	}

	msgBody := RenderTemplateToString(uc, c, "confirm_mail", gin.H{
		"Forename": ur.Forename,
		"Lastname": ur.Lastname,
		"TokeUrl":  cnf.WebLinks.ConfirmMail + reqC.MailConfirmationKey,
		"TokenValidTo": helpers.TimeToString(
			time.Unix(reqC.MailConfirmationExpiry, 0), uc.Language,
		),
	})
	msgSubject := fmt.Sprintf("%s :: Mail confirmation", cnf.Site.SiteName)

	_, err = mailClient.SendMail(c, &pb.SendMailRequest{
		Base: NewBaseRequest(uc),
		From: cnf.Site.EmailSender,
		To:   ur.Email,
		//ToMultiple: []string{ur.Email, "udo@n-wt.com"},
		Subject: msgSubject,
		Body:    msgBody,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error sending mail to user", cerr.GeneralError))
		return
	}

	authMiddlewareSimple.SetAuthHeader(c, user.Id)

	c.JSON(http.StatusOK, &RegisterUserResponse{
		TFAQrImage: base64.StdEncoding.EncodeToString(twoFaResp.Bitmap),
		TFASecret:  twoFaResp.Secret,
	})
}

//Confirm2FARequest is the data needed for the 2fa registration
type Confirm2FARequest struct {
	TfaCode string `form:"tfa_code" json:"tfa_code" validate:"required"`
}

//Confirm2FAResponse response for API
type Confirm2FAResponse struct {
	MailConfirmed     bool `json:"mail_confirmed"`
	TfaConfirmed      bool `json:"tfa_confirmed"`
	MnemonicConfirmed bool `json:"mnemonic_confirmed"`
}

//Confirm2FA checks the given 2fa code
func Confirm2FA(uc *mw.IcopContext, c *gin.Context) {
	var cd Confirm2FARequest
	if err := c.Bind(&cd); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, cd); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	// get the user from DB
	userID := mw.GetAuthUser(c).UserID

	req := &pb.GetUserByIDOrEmailRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		//should not happen, because user_id is from the jwt but just in case ...
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User from jwt could not be found in db", cerr.GeneralError))
		return
	}

	if user.TfaTempSecret == "" {
		//user already did 2fa registration
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.TfaAlreadyConfirmed, "already confirmed", "confirm_2FAregsitration.2FACode.already_done"))
		return
	}

	//do the 2fa authentication
	req2FA := &pb.AuthenticateRequest{
		Base:   NewBaseRequest(uc),
		Code:   cd.TfaCode,
		Secret: user.TfaTempSecret,
	}
	resp2FA, err := twoFAClient.Authenticate(c, req2FA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error doing 2FA authenticate", cerr.GeneralError))
		return
	}

	if !resp2FA.Result {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("tfa_code", cerr.InvalidArgument, "2FA code is invalid", "confirm_2FAregsitration.2FACode.invalid"))
		return
	}

	reqQRCode := &pb.FromSecretRequest{
		Base:   NewBaseRequest(uc),
		Email:  user.Email,
		Secret: user.TfaTempSecret,
	}
	respQRCode, err := twoFAClient.FromSecret(c, reqQRCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error doing 2FA qr code", cerr.GeneralError))
		return
	}

	_, err = dbClient.SetUserTFAConfirmed(c, &pb.SetUserTfaConfirmedRequest{
		Base:      NewBaseRequest(uc),
		UserId:    userID,
		TfaSecret: respQRCode.Secret,
		TfaQrcode: respQRCode.Bitmap,
		TfaUrl:    respQRCode.Url,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error setting user 2fa confirmed", cerr.GeneralError))
		return
	}

	authMiddlewareFull.SetAuthHeader(c, user.Id)

	c.JSON(http.StatusOK, &Confirm2FAResponse{
		MailConfirmed:     user.MailConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
		TfaConfirmed:      true,
	})
}

//ResendConfirmationMailRequest is the data needed for confirming the mail
type ResendConfirmationMailRequest struct {
	Email string `form:"email" json:"email" validate:"required,icop_email"`
}

//ResendConfirmMail resend the email confirmation mail
func ResendConfirmMail(uc *mw.IcopContext, c *gin.Context) {
	var rs ResendConfirmationMailRequest
	if err := c.Bind(&rs); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rs); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	req := &pb.GetUserByIDOrEmailRequest{
		Base:  NewBaseRequest(uc),
		Email: rs.Email,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusBadRequest,
			cerr.NewIcopError("email", cerr.InvalidArgument, "Email-address not found in database", "resend_mail.email.notFound"))
		return
	}

	if user.MailConfirmed {
		c.JSON(http.StatusBadRequest,
			cerr.NewIcopError("email", cerr.EmailAlreadyConfirmed, "Email-address already confirmed", "resend_mail.emai.already_confirmed"))

		return
	}

	ur, err := dbClient.GetUserProfile(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   user.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting userProfile", cerr.GeneralError))
		return
	}
	msgBody := RenderTemplateToString(uc, c, "confirm_mail", gin.H{
		"Forename": ur.Forename,
		"Lastname": ur.Lastname,
		"TokeUrl":  cnf.WebLinks.ConfirmMail + user.MailConfirmationKey,
		"TokenValidTo": helpers.TimeToString(
			time.Unix(user.MailConfirmationExpiry, 0), uc.Language,
		),
	})
	msgSubject := fmt.Sprintf("%s :: Mail confirmation", cnf.Site.SiteName)

	_, err = mailClient.SendMail(c, &pb.SendMailRequest{
		Base:    NewBaseRequest(uc),
		From:    cnf.Site.EmailSender,
		To:      ur.Email,
		Subject: msgSubject,
		Body:    msgBody,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error sending mail to user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ConfirmMnemonic sets the flag in the database
func ConfirmMnemonic(uc *mw.IcopContext, c *gin.Context) {
	userID := mw.GetAuthUser(c).UserID
	req := &pb.GetUserByIDOrEmailRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	user, err := dbClient.GetUserDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusBadRequest,
			cerr.NewIcopError("email", cerr.InvalidArgument, "Email-address not found in database", "resend_mail.email.notFound"))
		return
	}

	if !user.MailConfirmed {
		c.JSON(http.StatusBadRequest,
			cerr.NewIcopError("", cerr.EmailNotConfigured, "Email address not confirmed", ""))

		return
	}

	//set flag in db
	_, err = dbClient.SetUserMnemonicConfirmed(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	})

	c.JSON(http.StatusOK, &Confirm2FAResponse{
		MailConfirmed:     user.MailConfirmed,
		MnemonicConfirmed: true,
		TfaConfirmed:      user.TfaConfirmed,
	})
}
