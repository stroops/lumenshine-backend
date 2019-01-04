package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/gin-gonic/gin"
)

//SalutationResponse List of all salutations
// swagger:model
type SalutationResponse struct {
	Salutations []string `json:"salutations"`
}

//SalutationList returns a list of all salutations for the current language
// swagger:route GET /portal/user/salutation_list general SalutationList
//
// Lists all salutations for the current language
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: SalutationResponse - List of all salutations
func SalutationList(uc *mw.IcopContext, c *gin.Context) {

	//TODO: check that langcode is valid. We do this in memory
	req := &pb.IDString{
		Base: NewBaseRequest(uc),
		Id:   uc.Language,
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
// swagger:model
type Language struct {
	Code string `json:"lang_code"`
	Name string `json:"lang_name"`
}

//LanguagesResponse list for languages
// swagger:model
type LanguagesResponse struct {
	Languages []Language `json:"languages"`
}

//LanguageList returns a list of all languages
// swagger:route GET /portal/user/language_list general LanguageList
//
// Lists all languages
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: LanguagesResponse - List of all languages
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
// swagger:model
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

//CountryResponse list for countries
// swagger:model
type CountryResponse struct {
	Countries []Country `json:"countries"`
}

//CountryList returns a list of all countries for the current language
// swagger:route GET /portal/user/country_list general CountryList
//
// List all countries for the current language
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: CountryResponse - List of all countries
func CountryList(uc *mw.IcopContext, c *gin.Context) {

	//TODO: check that langcode is valid. We do this in memory
	req := &pb.Empty{
		Base: NewBaseRequest(uc),
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

//GetOccupationsRequest - filters occupations by name
//swagger:parameters GetOccupationsRequest OccupationList
type GetOccupationsRequest struct {
	//Occupation name
	//required: true
	Name string `form:"name"`
}

//Occupation represents one occupation
// swagger:model
type Occupation struct {
	Code08 string `json:"code08"`
	Code88 string `json:"code88"`
	Name   string `json:"name"`
}

//OccupationResponse list for occupations
// swagger:model
type OccupationResponse struct {
	Occupations []Occupation `json:"occupations"`
}

//OccupationList returns a list of occupations filtered by name
// swagger:route GET /portal/user/occupation_list general OccupationList
//
// List of occupations filtered by name
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: OccupationResponse - List of all occupations
func OccupationList(uc *mw.IcopContext, c *gin.Context) {
	rr := new(GetOccupationsRequest)
	if err := c.Bind(rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

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
// swagger:model
type UserDetails struct {
	MailConfirmed     bool `json:"mail_confirmed"`
	TfaConfirmed      bool `json:"tfa_confirmed"`
	MnemonicConfirmed bool `json:"mnemonic_confirmed"`
}

//GetUserRegistrationDetails returns the detail data for the user
//This function may be called either with a valid JWT, then we will use the userID to get the user data
//or with a tfa code and an email, without a valid jwt. Then we will get the user by 2fa code and email
// swagger:route GET /portal/user/get_user_details general GetUserRegistrationDetails
//
// Returns the detail data for the user
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: UserDetails - user details
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

	if user.IsClosed {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.UserIsClosed, "user closed", ""))
		return
	}

	if user.IsSuspended {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.UserIsSuspended, "user suspended", ""))
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
//swagger:parameters ConfirmTokeRequest ConfirmToken
type ConfirmTokeRequest struct {
	//Mail token to confirm
	//required: true
	Token string `form:"token" json:"token" validate:"required"`
}

//ConfirmTokenResponse response
// swagger:model
type ConfirmTokenResponse struct {
	Email                     string    `json:"email"`
	ConfirmationDate          time.Time `json:"confirmation_date"`
	TokenAlreadyConfirmed     bool      `json:"token_already_confirmed"`
	SEP10TransactionChallenge string    `json:"sep10_transaction_challenge"`
}

//ConfirmToken confirms a mail token
// swagger:route POST /portal/user/confirm_token general ConfirmToken
//
// Confirms the mail token
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: ConfirmTokenResponse
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
		c.JSON(http.StatusOK, &ConfirmTokenResponse{
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

	sep10ChallangeTX, err := getSEP10Challenge(u.PublicKey_0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "error generating challange :"+err.Error(), cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &ConfirmTokenResponse{
		Email:                     u.Email,
		SEP10TransactionChallenge: sep10ChallangeTX,
	})
}

//UserSecurityDataResponse response for API
// swagger:model
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
// swagger:route GET /portal/user/dashboard/user_auth_data general UserSecurityData
//
// Returns the security data for the authenticated user
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: UserSecurityDataResponse
//
func UserSecurityData(uc *mw.IcopContext, c *gin.Context) {
	// swagger:route GET /portal/user/auth/user_auth_data general UserSecurityData
	//
	// Returns the security data for the authenticated user
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: UserSecurityDataResponse

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
//swagger:parameters GetTfaSecretRequest GetTfaSecret
type GetTfaSecretRequest struct {
	//SEP 10 transaction
	SEP10Transaction string `form:"sep10_transaction" json:"sep10_transaction"`
}

//GetTfaSecretResponse response for the api call
// swagger:model
type GetTfaSecretResponse struct {
	TfaSecret string `json:"tfa_secret"`
}

//GetTfaSecret returns the tfa secrete for the user
// swagger:route GET /portal/user/dashboard/tfa_secret general GetTfaSecret
//
// Returns the tfa secrete for the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetTfaSecretResponse
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

	if user.IsClosed {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.UserIsClosed, "user closed", ""))
		return
	}

	if user.IsSuspended {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("email", cerr.UserIsSuspended, "user suspended", ""))
		return
	}

	valid, _, err := verifySEP10Data(l.SEP10Transaction, userID, uc, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, err.Error(), cerr.GeneralError))
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("transaction", cerr.InvalidArgument, "could not validate challange transaction", ""))
		return
	}

	c.JSON(http.StatusOK, &GetTfaSecretResponse{
		TfaSecret: user.TfaSecret,
	})
}

//GetUserMessageListRequest used for requesting the user messages
//if from_archive specified, the data will be read from the archive
//swagger:parameters GetUserMessageListRequest GetUserMessageList
type GetUserMessageListRequest struct {
	//Flag to filter the archived messages
	FromArchive bool `form:"from_archive" json:"from_archive" query:"from_archive"`
}

//MessageListItem is the list of the messages (short)
// swagger:model
type MessageListItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	DateCreated time.Time `json:"date_created"`
}

//GetUserMessageListResponse response for the api call
// swagger:model
type GetUserMessageListResponse struct {
	Messages     []MessageListItem `json:"messages"`
	CurrentCount int64             `json:"current_count"`
	ArchiveCount int64             `json:"archive_count"`
}

//GetUserMessageList returns a list of all messages for the user
// swagger:route GET /portal/user/dashboard/get_user_messages_list general GetUserMessageList
//
// Returns a list of all messages for the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetUserMessageListResponse
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
//swagger:parameters GetUserMessageRequest GetUserMessage
type GetUserMessageRequest struct {
	//Filter by message id
	//required: true
	MessageID int64 `form:"message_id" json:"message_id" query:"message_id" validate:"required"`
	//Flag to filter the archived messages
	FromArchive bool `form:"from_archive" json:"from_archive" query:"from_archive"`
}

//MessageItem is one user message
// swagger:model
type MessageItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	DateCreated time.Time `json:"date_created"`
}

//GetUserMessage returns the requested message
// swagger:route GET /portal/user/dashboard/get_user_message general GetUserMessage
//
// Returns the requested message
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: MessageItem
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

//GetStellarTomlRequest used for requesting toml file
//swagger:parameters GetStellarTomlRequest GetStellarToml
type GetStellarTomlRequest struct {
	//Filter by message id
	//required: true
	Domain string `form:"domain" json:"domain" query:"domain" validate:"required"`

	//Protocol might be http or https, uses https if not specified
	Protocol string `form:"protocol" json:"protocol" query:"protocol" validate:"omitempty,oneof=http https"`
}

//GetStellarTomlResponse is a toml file content
// swagger:model
type GetStellarTomlResponse struct {
	TomlFileContent string `json:"toml_file_content"`
	HTTPStatusCode  int    `json:"http_status_code"`
	HTTPUrl         string `json:"http_url"`
}

//GetStellarToml returns the requested toml file
// swagger:route GET /portal/user/dashboard/get_stellar_toml general GetStellarToml
//
// Returns the requested message
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetStellarTomlResponse
func GetStellarToml(uc *mw.IcopContext, c *gin.Context) {
	var l GetStellarTomlRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	protocol := "https"
	if l.Protocol != "" {
		protocol = l.Protocol
	}

	client := &http.Client{
		Timeout: (60 * time.Second), //max 60 seconds, then timeout
	}

	if protocol == "https" {
		//allow insecure configuration of selfsigned ceritficates
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	url := fmt.Sprintf("%s://%s/.well-known/stellar.toml", protocol, l.Domain)
	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading http request", cerr.GeneralError))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading http response body", cerr.GeneralError))
		return
	}
	bodyStr := string(body[:])

	c.JSON(http.StatusOK, &GetStellarTomlResponse{
		TomlFileContent: bodyStr,
		HTTPStatusCode:  resp.StatusCode,
		HTTPUrl:         url,
	})
}
