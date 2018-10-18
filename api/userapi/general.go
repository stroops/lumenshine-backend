package main

import (
	"net/http"
	"strconv"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/gin-gonic/gin"
)

//SalutationResponse List of all salutations
type SalutationResponse struct {
	Salutations []string `json:"salutations"`
}

//SalutationList returns a list of all salutations for the given LanguageCode
func SalutationList(uc *mw.IcopContext, c *gin.Context) {

	//TODO: check that langcode is valid. We do this in memory
	req := &pb.LanguageCodeRequest{
		Base:         NewBaseRequest(uc),
		LanguageCode: uc.Language,
	}
	salutations, err := dbClient.GetSalutationList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading salutations", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &SalutationResponse{
		Salutations: salutations.Salutation,
	})
}

//Language represents one language
type Language struct {
	Code string `json:"lang_code"`
	Name string `json:"lang_name"`
}

//LanguagesResponse list for languages
type LanguagesResponse struct {
	Languages []Language `json:"languages"`
}

//LanguageList returns a list of all languages
func LanguageList(uc *mw.IcopContext, c *gin.Context) {
	languages, err := dbClient.GetLanguageList(c, &pb.Empty{Base: NewBaseRequest(uc)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading languages", cerr.GeneralError))
		return
	}
	var retLanguages []Language
	for _, language := range languages.Languages {
		retLanguages = append(retLanguages, Language{Code: language.Code, Name: language.Name})
	}
	c.JSON(http.StatusOK, &LanguagesResponse{
		Languages: retLanguages,
	})
}

//Country represents one country
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

//CountryResponse list for countries
type CountryResponse struct {
	Countries []Country `json:"countries"`
}

//CountryList returns a list of all countries for the given LanguageCode
func CountryList(uc *mw.IcopContext, c *gin.Context) {

	//TODO: check that langcode is valid. We do this in memory
	req := &pb.LanguageCodeRequest{
		Base:         NewBaseRequest(uc),
		LanguageCode: uc.Language,
	}
	countries, err := dbClient.GetCountryList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading countries", cerr.GeneralError))
		return
	}

	var retCountries []Country
	for _, country := range countries.Countries {
		retCountries = append(retCountries, Country{Code: country.Code, Name: country.Name})
	}
	c.JSON(http.StatusOK, &CountryResponse{
		Countries: retCountries,
	})
}

//GetOccupationsRequest - filtered by name
type GetOccupationsRequest struct {
	Name string `form:"name"`
}

//Occupation represents one occupation
type Occupation struct {
	Code08 string `json:"code08"`
	Code88 string `json:"code88"`
	Name   string `json:"name"`
}

//OccupationResponse list for occupations
type OccupationResponse struct {
	Occupations []Occupation `json:"occupations"`
}

//OccupationList returns a list of occupations filtered by name
func OccupationList(uc *mw.IcopContext, c *gin.Context) {

	rr := new(GetOccupationsRequest)
	if err := c.Bind(rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	uc.Log.Print("request name: " + rr.Name)

	req := &pb.OccupationListRequest{
		Base:       NewBaseRequest(uc),
		Name:       rr.Name,
		LimitCount: 50,
	}
	occupations, err := dbClient.GetOccupationList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading occupations", cerr.GeneralError))
		return
	}

	var retOccupations []Occupation
	for _, occupation := range occupations.Occupations {
		retOccupations = append(retOccupations, Occupation{Code08: strconv.FormatInt(occupation.Code08, 10), Code88: strconv.FormatInt(occupation.Code88, 10), Name: occupation.Name})
	}
	c.JSON(http.StatusOK, &OccupationResponse{
		Occupations: retOccupations,
	})
}

//UserDetails response for the API
type UserDetails struct {
	MailConfirmed     bool `json:"mail_confirmed"`
	TfaConfirmed      bool `json:"tfa_confirmed"`
	MnemonicConfirmed bool `json:"mnemonic_confirmed"`
}

//GetUserRegistrationDetails returns the detail data for the user
//This function may be called either with a valid JWT, then we will use the userID to get the user data
//or with a tfa code and an email, without a valid jwt. Then we will get the user by 2fa code and email
func GetUserRegistrationDetails(uc *mw.IcopContext, c *gin.Context) {
	// get the user from DB
	userID := mw.GetAuthUser(c).UserID

	user, err := dbClient.GetUserDetails(c, &pb.GetUserByIDOrEmailRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User from jwt could not be found in db", cerr.GeneralError))
		return
	}

	//reached this point --> return requested data
	c.JSON(http.StatusOK, &UserDetails{
		MailConfirmed:     user.MailConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
		TfaConfirmed:      user.TfaConfirmed,
	})
}

//ConfirmTokeRequest is the data needed for confirming the mail
type ConfirmTokeRequest struct {
	Token string `form:"token" json:"token" validate:"required"`
}

//ConfirmTokeResponse response
type ConfirmTokeResponse struct {
	Email                 string    `json:"email"`
	ConfirmationDate      time.Time `json:"confirmation_date"`
	TokenAlreadyConfirmed bool      `json:"token_already_confirmed"`
}

//ConfirmToken confirms a mail token
func ConfirmToken(uc *mw.IcopContext, c *gin.Context) {
	var cm ConfirmTokeRequest
	if err := c.Bind(&cm); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, cm); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userMailRequest := &pb.UserMailTokenRequest{
		Base:  NewBaseRequest(uc),
		Token: cm.Token,
	}
	u, err := dbClient.GetUserByMailtoken(c, userMailRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting user by token", cerr.GeneralError))
		return
	}

	if u.TokenAlreadyConfirmed {
		c.JSON(http.StatusOK, &ConfirmTokeResponse{
			TokenAlreadyConfirmed: true,
			ConfirmationDate:      time.Unix(u.ConfirmedDate, 0),
		})
		return
	}

	if u.TokenNotFound {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("token", cerr.InvalidArgument, "Token not found", ""))
		return
	}

	//check that token not expiered
	t := time.Unix(u.MailConfirmationExpiry, 0)
	if t.Before(time.Now()) {
		c.JSON(http.StatusBadRequest,
			cerr.NewIcopError("token", cerr.TokenExpiered, "Email confirmation token expired", "confirm_mail.token.expired"))
		return
	}

	if !u.MailConfirmed {
		_, err := dbClient.SetUserMailConfirmed(c, &pb.IDRequest{
			Base: NewBaseRequest(uc),
			Id:   u.UserId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error setting user mail confirmed", cerr.GeneralError))
			return
		}
	}

	//delete the current token from the database
	_, err = dbClient.SetUserMailToken(c, &pb.SetMailTokenRequest{
		Base:                NewBaseRequest(uc),
		UserId:              u.UserId,
		MailConfirmationKey: "",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error setting token null in DB", cerr.GeneralError))
	}

	authMiddlewareSimple.SetAuthHeader(c, u.UserId)

	c.JSON(http.StatusOK, &ConfirmTokeResponse{
		Email: u.Email,
	})
}

//UserSecurityDataResponse response for API
type UserSecurityDataResponse struct {
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
}

//UserSecurityData returns the security data for the user
func UserSecurityData(uc *mw.IcopContext, c *gin.Context) {

	idRequest := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   mw.GetAuthUser(c).UserID,
	}
	s, err := dbClient.GetUserSecurities(c, idRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user security data", cerr.GeneralError))
		return
	}

	if s.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User from could not be found in db", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &UserSecurityDataResponse{
		KdfPasswordSalt: s.KdfSalt,

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

//GetTfaSecretRequest for requesting the tfa secrete
type GetTfaSecretRequest struct {
	PublicKey188 string `form:"public_key_188" json:"public_key_188" validate:"required"`
}

//GetTfaSecretResponse response for the api call
type GetTfaSecretResponse struct {
	TfaSecret string `json:"tfa_secret"`
}

//GetTfaSecret returns the tfa secrete for the user
func GetTfaSecret(uc *mw.IcopContext, c *gin.Context) {
	var l GetTfaSecretRequest
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
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}

	if user.UserNotFound {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "User for email could not be found in db", cerr.GeneralError))
		return
	}

	if !CheckPasswordHash(uc.Log, l.PublicKey188, user.Password) {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_188", cerr.InvalidPassword, "Invalid public key", ""))
		return
	}

	c.JSON(http.StatusOK, &GetTfaSecretResponse{
		TfaSecret: user.TfaSecret,
	})
}

//GetUserMessageListRequest used for requesting the user messages
//if from_archive specified, the data will be read from the archive
type GetUserMessageListRequest struct {
	FromArchive bool `form:"from_archive" json:"from_archive" query:"from_archive"`
}

//MessageListItem is the list of the messages (short)
type MessageListItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	DateCreated time.Time `json:"date_created"`
}

//GetUserMessageListResponse response for the api call
type GetUserMessageListResponse struct {
	Messages     []MessageListItem `json:"messages"`
	CurrentCount int64             `json:"current_count"`
	ArchiveCount int64             `json:"archive_count"`
}

//GetUserMessageList returns a list of all messages for the user
func GetUserMessageList(uc *mw.IcopContext, c *gin.Context) {
	var l GetUserMessageListRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	messages, err := dbClient.GetUserMessages(c, &pb.UserMessageListRequest{
		Base:    NewBaseRequest(uc),
		UserId:  userID,
		Archive: l.FromArchive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user messages", cerr.GeneralError))
		return
	}

	ret := new(GetUserMessageListResponse)
	for _, m := range messages.MessageListItems {
		ret.Messages = append(ret.Messages, MessageListItem{
			ID:          m.Id,
			Title:       m.Title,
			DateCreated: time.Unix(m.DateCreated, 0),
		})
	}
	ret.ArchiveCount = messages.ArchiveCount
	ret.CurrentCount = messages.CurrentCount

	c.JSON(http.StatusOK, ret)
}

//GetUserMessageRequest used for requesting the user message
type GetUserMessageRequest struct {
	MessageID   int64 `form:"message_id" json:"message_id" query:"message_id" validate:"required"`
	FromArchive bool  `form:"from_archive" json:"from_archive" query:"from_archive"`
}

//MessageItem is one user message
type MessageItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	DateCreated time.Time `json:"date_created"`
}

//GetUserMessage returns the requested message
func GetUserMessage(uc *mw.IcopContext, c *gin.Context) {
	var l GetUserMessageRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	m, err := dbClient.GetUserMessage(c, &pb.UserMessageRequest{
		Base:      NewBaseRequest(uc),
		MessageId: l.MessageID,
		Archive:   l.FromArchive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user message", cerr.GeneralError))
		return
	}

	//check that user is owner of message
	if userID != m.UserId {
		c.JSON(http.StatusUnauthorized, cerr.NewIcopErrorShort(cerr.NoPermission, "Not owner of message"))
		return
	}

	if !l.FromArchive {
		//when we reach this, we assume, that the message will be send to the client in the next step.
		//thus we can move the message away to the archive
		dbClient.MoveMessageToArchive(c, &pb.IDRequest{
			Base: NewBaseRequest(uc),
			Id:   l.MessageID,
		})
	}

	c.JSON(http.StatusOK, &MessageItem{
		ID:          m.Id,
		Title:       m.Title,
		Message:     m.Message,
		DateCreated: time.Unix(m.DateCreated, 0),
	})
}
