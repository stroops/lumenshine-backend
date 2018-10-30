package api

import (
	"net/http"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"

	"github.com/Soneso/lumenshine-backend/admin/route"
)

const (
	//StellarAccountRoutePrefix for the account routes
	StellarAccountRoutePrefix = "stellar_account"
)

//StellarAccountType holds the type of the account
type StellarAccountType string

const (
	funding StellarAccountType = "funding"
	issuing StellarAccountType = "issuing"
	worker  StellarAccountType = "worker"
)

//SignerType holds the type of a signer
type SignerType string

//init setup all the routes for the users handling
func init() {
	route.AddRoute("GET", "/details/:publickey", GetStellarAccount, []string{}, "stellar_account_details", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/accounts_list", AllStellarAccounts, []string{}, "stellar_account_list", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/asset_codes/:publickey", IssuerAssetCodes, []string{}, "stellar_asset_codes", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/signer_seed/:publickey", GetSignerSeed, []string{}, "signer_seed", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/add", AddStellarAccount, []string{}, "stellar_account_add", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/edit", EditStellarAccount, []string{}, "stellar_account_edit", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove", RemoveStellarAccount, []string{}, "remove_stellar_account", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/add_asset_code", AddIssuerAssetCode, []string{}, "issuer_add_asset_code", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove_asset_code", RemoveIssuerAssetCode, []string{}, "issuer_remove_asset_code", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/add_allowtrust_signer", AddAllowtrustSigner, []string{}, "add_allowtrust_signer", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove_allowtrust_signer", RemoveAllowtrustSigner, []string{}, "remove_allowtrust_signer", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/edit_allowtrust_signer", EditAllowtrustSigner, []string{}, "edit_allowtrust_signer", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/edit_other_signer", EditOtherSigner, []string{}, "edit_other_signer", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/add_other_signer", AddOtherSigner, []string{}, "add_other_signer", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove_other_signer", RemoveOtherSigner, []string{}, "remove_other_signer", StellarAccountRoutePrefix)
}

//AddStellarAccountRoutes adds all the routes for the account handling
func AddStellarAccountRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(StellarAccountRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//StellarAccountAddRequest new user information
//swagger:parameters StellarAccountAddRequest AddStellarAccount
type StellarAccountAddRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
	//required : true
	Name string `form:"name" json:"name" validate:"required,max=256"`
	//required : true
	Description string `form:"description" json:"description" validate:"required"`
	//required : true
	Type string `form:"type" json:"type" validate:"required"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"omitempty,icop_assetcode"`
}

//AddStellarAccount creates new account in the db
// swagger:route POST /portal/admin/dash/stellar_account/add StellarAccount AddStellarAccount
//
// Creates new account in the db
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddStellarAccount(uc *mw.AdminContext, c *gin.Context) {
	var rr StellarAccountAddRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	if !strings.EqualFold(rr.Type, string(issuing)) && !strings.EqualFold(rr.Type, string(funding)) && !strings.EqualFold(rr.Type, string(worker)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "type", cerr.InvalidArgument, "Type is not 'funding', 'issuing' or 'worker'", ""))
		return
	}
	if strings.EqualFold(rr.Type, string(issuing)) && rr.AssetCode == "" {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Issuing account must have a valid asset code", ""))
		return
	}
	if !strings.EqualFold(rr.Type, string(issuing)) && rr.AssetCode != "" {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Only an issuing account has an asset code", ""))
		return
	}
	existsAccount, err := db.ExistsStellarAccount(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if existsAccount {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Account with same public key already exists", ""))
		return
	}
	account := &models.AdminStellarAccount{
		PublicKey:   rr.PublicKey,
		Name:        rr.Name,
		Description: rr.Description,
		Type:        rr.Type,
		UpdatedBy:   getUpdatedBy(c)}

	err = db.AddStellarAccount(account, rr.AssetCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error registering account", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//StellarAccountItem new user information
// swagger:model
type StellarAccountItem struct {
	PublicKey   string `form:"public_key" json:"public_key"`
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	Type        string `form:"type" json:"type"`
}

//AllStellarAccounts - returns all stellar accounts
// swagger:route GET /portal/admin/dash/stellar_account/accounts_list StellarAccount AllStellarAccounts
//
// Returns all stellar accounts
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:[]StellarAccountItem
func AllStellarAccounts(uc *mw.AdminContext, c *gin.Context) {
	dbAccounts, err := db.AllStellarAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading accounts", cerr.GeneralError))
		return
	}

	accounts := make([]StellarAccountItem, len(dbAccounts))
	for index, account := range dbAccounts {
		accounts[index] = StellarAccountItem{
			PublicKey:   account.PublicKey,
			Name:        account.Name,
			Description: account.Description,
			Type:        account.Type}
	}

	c.JSON(http.StatusOK, &accounts)
}

//AssetCodesRequest new user information
//swagger:parameters AssetCodesRequest IssuerAssetCodes
type AssetCodesRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
}

//IssuerAssetCodes - returns the issuer asset codes
// swagger:route GET /portal/admin/dash/stellar_account/asset_codes/:publickey StellarAccount IssuerAssetCodes
//
// Returns, in a string array, the asset codes of the given issuer
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func IssuerAssetCodes(uc *mw.AdminContext, c *gin.Context) {
	publicKey := c.Param("publickey")
	rr := AssetCodesRequest{PublicKey: publicKey}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(account.Type, string(issuing)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key is not from an issuing account", ""))
		return
	}

	dbAssetCodes, err := db.IssuerAssetCodes(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading asset codes", cerr.GeneralError))
		return
	}

	assetcodes := make([]string, len(dbAssetCodes))
	for index, assetCode := range dbAssetCodes {
		assetcodes[index] = assetCode.AssetCode
	}

	c.JSON(http.StatusOK, &assetcodes)
}

//StellarAccountRequest new user information
//swagger:parameters StellarAccountRequest GetStellarAccount
type StellarAccountRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
}

//StellarSigner - signer details
// swagger:model
type StellarSigner struct {
	SignerPublicKey string `form:"public_key" json:"public_key"`
	Name            string `form:"name" json:"name"`
	Description     string `form:"description" json:"description"`
}

//StellarAccountResponse new user information
// swagger:model
type StellarAccountResponse struct {
	PublicKey         string          `form:"public_key" json:"public_key"`
	Name              string          `form:"name" json:"name"`
	Description       string          `form:"description" json:"description"`
	Type              string          `form:"type" json:"type"`
	AssetCodes        []string        `form:"asset_codes" json:"asset_codes"`
	AllowTrustSigners []StellarSigner `form:"allow_trust_signers" json:"allow_trust_signers"`
	OtherSigners      []StellarSigner `form:"other_signers" json:"other_signers"`
}

//GetStellarAccount - returns account details
// swagger:route GET /portal/admin/dash/stellar_account/details/:publickey StellarAccount GetStellarAccount
//
// Returns account details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: StellarAccountResponse
func GetStellarAccount(uc *mw.AdminContext, c *gin.Context) {
	publicKey := c.Param("publickey")
	rr := AssetCodesRequest{PublicKey: publicKey}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	account, err := db.GetStellarAccount(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	accountResponse := getAccountResponse(*account)
	c.JSON(http.StatusOK, &accountResponse)
}

func getAccountResponse(account models.AdminStellarAccount) StellarAccountResponse {
	assetcodes := make([]string, len(account.R.IssuerPublicKeyAdminStellarAssets))
	for index, assetCode := range account.R.IssuerPublicKeyAdminStellarAssets {
		assetcodes[index] = assetCode.AssetCode
	}
	allowTrustSigners := make([]StellarSigner, 0)
	otherSigners := make([]StellarSigner, 0)
	for _, signer := range account.R.StellarAccountPublicKeyAdminStellarSigners {
		signerItem := StellarSigner{
			SignerPublicKey: signer.SignerPublicKey,
			Name:            signer.Name,
			Description:     signer.Description}
		if strings.EqualFold(signer.Type, string(models.StellarSignerTypeAllowTrust)) {
			allowTrustSigners = append(allowTrustSigners, signerItem)
		}
		if strings.EqualFold(signer.Type, string(models.StellarSignerTypeOther)) {
			otherSigners = append(otherSigners, signerItem)
		}
	}

	return StellarAccountResponse{
		PublicKey:         account.PublicKey,
		Name:              account.Name,
		Description:       account.Description,
		Type:              account.Type,
		AssetCodes:        assetcodes,
		AllowTrustSigners: allowTrustSigners,
		OtherSigners:      otherSigners}
}

//StellarAccountEditRequest new user information
//swagger:parameters StellarAccountEditRequest EditStellarAccount
type StellarAccountEditRequest struct {
	//required : true
	PublicKey   string  `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
	Name        *string `form:"name" json:"name" validate:"max=256"`
	Description *string `form:"description" json:"description"`
}

//EditStellarAccount edits the given account details
// swagger:route POST /portal/admin/dash/stellar_account/edit StellarAccount EditStellarAccount
//
// Edits the given account details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func EditStellarAccount(uc *mw.AdminContext, c *gin.Context) {
	var rr StellarAccountEditRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if rr.Name != nil {
		account.Name = *rr.Name
	}
	if rr.Description != nil {
		account.Description = *rr.Description
	}

	err = db.UpdateStellarAccount(account, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error editing account", cerr.GeneralError))
		return
	}

	accountResponse := getAccountResponse(*account)
	c.JSON(http.StatusOK, &accountResponse)
}

//IssuerAssetCodeRequest new user information
//swagger:parameters IssuerAssetCodeRequest AddIssuerAssetCode
type IssuerAssetCodeRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
}

//AddIssuerAssetCode - adds an asset code
// swagger:route POST /portal/admin/dash/stellar_account/add_asset_code StellarAccount AddIssuerAssetCode
//
// Adds an asset code
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddIssuerAssetCode(uc *mw.AdminContext, c *gin.Context) {
	var rr IssuerAssetCodeRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(account.Type, string(issuing)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key is not from an issuing account", ""))
		return
	}
	exists := false
	for _, assetCode := range account.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			exists = true
			break
		}
	}
	if exists {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code already exists for this issuing account", ""))
		return
	}
	err = db.AddAssetCode(account, rr.AssetCode, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding asset code", cerr.GeneralError))
		return
	}
	accountResponse := getAccountResponse(*account)
	accountResponse.AssetCodes = append(accountResponse.AssetCodes, rr.AssetCode)
	c.JSON(http.StatusOK, &accountResponse)
}

//RemoveIssuerAssetCode - removes an asset code
// swagger:route POST /portal/admin/dash/stellar_account/remove_asset_code StellarAccount RemoveIssuerAssetCode
//
// Removes an asset code
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveIssuerAssetCode(uc *mw.AdminContext, c *gin.Context) {
	var rr IssuerAssetCodeRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(account.Type, string(issuing)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key is not from an issuing account", ""))
		return
	}
	assetCodeIndex := -1
	for index, assetCode := range account.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			assetCodeIndex = index
			break
		}
	}
	if assetCodeIndex == -1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code not found", ""))
		return
	}

	err = db.RemoveAssetCode(account, rr.AssetCode, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing asset code", cerr.GeneralError))
		return
	}

	accountResponse := getAccountResponse(*account)
	accountResponse.AssetCodes = append(accountResponse.AssetCodes[:assetCodeIndex], accountResponse.AssetCodes[assetCodeIndex+1:]...)
	c.JSON(http.StatusOK, &accountResponse)
}

//AddAllowtrustRequest signer info
//swagger:parameters AddAllowtrustRequest AddAllowtrustSigner
type AddAllowtrustRequest struct {
	//required : true
	IssuingPublicKey  string  `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	SignerName        string  `form:"signer_name" json:"signer_name" validate:"max=256"`
	SignerDescription *string `form:"signer_description" json:"signer_description"`
	//required : true
	SignerPublicKey string `form:"signer_public_key" json:"signer_public_key" validate:"required,base64,len=56"`
	//required : true
	SignerSecretSeed string `form:"signer_secret_seed" json:"signer_secret_seed" validate:"required,base64,len=56"`
}

//AddAllowtrustSigner - adds an allow trust signer
// swagger:route POST /portal/admin/dash/stellar_account/add_allowtrust_signer StellarAccount AddAllowtrustSigner
//
// Adds an allow trust signer
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddAllowtrustSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr AddAllowtrustRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing public key not found in database", ""))
		return
	}
	if !strings.EqualFold(account.Type, string(issuing)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing public key is not from an issuing account", ""))
		return
	}
	accountResponse := getAccountResponse(*account)
	exists := false
	for _, signer := range accountResponse.AllowTrustSigners {
		if strings.EqualFold(signer.SignerPublicKey, rr.SignerPublicKey) {
			exists = true
			break
		}
	}
	if exists {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "signer_public_key", cerr.InvalidArgument, "Signer already exists", ""))
		return
	}
	dbSigner := models.AdminStellarSigner{
		StellarAccountPublicKeyID: rr.IssuingPublicKey,
		Name:             rr.SignerName,
		SignerPublicKey:  rr.SignerPublicKey,
		SignerSecretSeed: rr.SignerSecretSeed,
		Type:             models.StellarSignerTypeAllowTrust}
	if rr.SignerDescription != nil {
		dbSigner.Description = *rr.SignerDescription
	}

	err = db.AddSigner(account, &dbSigner, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding signer", cerr.GeneralError))
		return
	}
	accountResponse.AllowTrustSigners = append(accountResponse.AllowTrustSigners, StellarSigner{
		SignerPublicKey: dbSigner.SignerPublicKey,
		Name:            dbSigner.Name,
		Description:     dbSigner.Description})
	c.JSON(http.StatusOK, &accountResponse)
}

//RemoveSignerRequest new user information
//swagger:parameters RemoveSignerRequest RemoveAllowtrustSigner
type RemoveSignerRequest struct {
	//required : true
	IssuingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	//required : true
	SignerPublicKey string `form:"signer_public_key" json:"signer_public_key" validate:"required,base64,len=56"`
}

//RemoveAllowtrustSigner - removes an allow trust signer
// swagger:route POST /portal/admin/dash/stellar_account/remove_allowtrust_signer StellarAccount RemoveAllowtrustSigner
//
// Removes an allow trust signer
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveAllowtrustSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr RemoveSignerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(account.Type, string(issuing)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Public key is not from an issuing account", ""))
		return
	}
	accountResponse := getAccountResponse(*account)
	signerIndex := -1
	for index, signer := range accountResponse.AllowTrustSigners {
		if strings.EqualFold(signer.SignerPublicKey, rr.SignerPublicKey) {
			signerIndex = index
			break
		}
	}
	if signerIndex == -1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "signer_public_key", cerr.InvalidArgument, "Signer not found", ""))
		return
	}

	err = db.RemoveSigner(account, rr.SignerPublicKey, string(models.StellarSignerTypeAllowTrust), getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing signer", cerr.GeneralError))
		return
	}

	accountResponse.AllowTrustSigners = append(accountResponse.AllowTrustSigners[:signerIndex], accountResponse.AllowTrustSigners[signerIndex+1:]...)
	c.JSON(http.StatusOK, &accountResponse)
}

//AddOtherSignerRequest signer info
//swagger:parameters AddOtherSignerRequest AddOtherSigner
type AddOtherSignerRequest struct {
	//required : true
	IssuingPublicKey  string  `form:"account_public_key" json:"account_public_key" validate:"required,base64,len=56"`
	SignerName        string  `form:"signer_name" json:"signer_name" validate:"max=256"`
	SignerDescription *string `form:"signer_description" json:"signer_description"`
	//required : true
	SignerPublicKey string `form:"signer_public_key" json:"signer_public_key" validate:"required,base64,len=56"`
	//required : true
	SignerSecretSeed *string `form:"signer_secret_seed" json:"signer_secret_seed" validate:"omitempty,base64,len=56"`
}

//AddOtherSigner - adds other signer
// swagger:route POST /portal/admin/dash/stellar_account/add_other_signer StellarAccount AddOtherSigner
//
// Adds other signer
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddOtherSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr AddOtherSignerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "account_public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	accountResponse := getAccountResponse(*account)
	exists := false
	for _, signer := range accountResponse.OtherSigners {
		if strings.EqualFold(signer.SignerPublicKey, rr.SignerPublicKey) {
			exists = true
			break
		}
	}
	if exists {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "signer_public_key", cerr.InvalidArgument, "Signer already exists", ""))
		return
	}
	dbSigner := models.AdminStellarSigner{
		StellarAccountPublicKeyID: rr.IssuingPublicKey,
		Name:            rr.SignerName,
		SignerPublicKey: rr.SignerPublicKey,
		Type:            models.StellarSignerTypeOther}
	if rr.SignerDescription != nil {
		dbSigner.Description = *rr.SignerDescription
	}
	if rr.SignerSecretSeed != nil {
		dbSigner.SignerSecretSeed = *rr.SignerSecretSeed
	}

	err = db.AddSigner(account, &dbSigner, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding signer", cerr.GeneralError))
		return
	}

	accountResponse.OtherSigners = append(accountResponse.OtherSigners, StellarSigner{
		SignerPublicKey: dbSigner.SignerPublicKey,
		Name:            dbSigner.Name,
		Description:     dbSigner.Description})
	c.JSON(http.StatusOK, &accountResponse)
}

//RemoveOtherSignerRequest new user information
//swagger:parameters RemoveOtherSignerRequest RemoveOtherSigner
type RemoveOtherSignerRequest struct {
	//required : true
	IssuingPublicKey string `form:"account_public_key" json:"account_public_key" validate:"required,base64,len=56"`
	//required : true
	SignerPublicKey string `form:"signer_public_key" json:"signer_public_key" validate:"required,base64,len=56"`
}

//RemoveOtherSigner - removes other signer
// swagger:route POST /portal/admin/dash/stellar_account/remove_other_signer StellarAccount RemoveOtherSigner
//
// Removes other signer
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveOtherSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr RemoveOtherSignerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "account_public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	accountResponse := getAccountResponse(*account)
	signerIndex := -1
	for index, signer := range accountResponse.OtherSigners {
		if strings.EqualFold(signer.SignerPublicKey, rr.SignerPublicKey) {
			signerIndex = index
			break
		}
	}
	if signerIndex == -1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "signer_public_key", cerr.InvalidArgument, "Signer not found", ""))
		return
	}

	err = db.RemoveSigner(account, rr.SignerPublicKey, string(models.StellarSignerTypeOther), getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing signer", cerr.GeneralError))
		return
	}

	accountResponse.OtherSigners = append(accountResponse.OtherSigners[:signerIndex], accountResponse.OtherSigners[signerIndex+1:]...)
	c.JSON(http.StatusOK, &accountResponse)
}

//RemoveStellarAccount - removes a stellar account
// swagger:route POST /portal/admin/dash/stellar_account/remove StellarAccount RemoveStellarAccount
//
// Removes a stellar account
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveStellarAccount(uc *mw.AdminContext, c *gin.Context) {
	var rr StellarAccountRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	account, err := db.GetStellarAccount(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}

	err = db.DeleteStellarAccount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting account", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}

//GetSignerSeedRequest new user information
//swagger:parameters GetSignerSeedRequest GetSignerSeed
type GetSignerSeedRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
}

//GetSignerSeed - returns the signer's seed
// swagger:route GET /portal/admin/dash/stellar_account/signer_seed/:publickey StellarAccount GetSignerSeed
//
// Returns the signer's seed
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func GetSignerSeed(uc *mw.AdminContext, c *gin.Context) {
	publicKey := c.Param("publickey")
	rr := GetSignerSeedRequest{PublicKey: publicKey}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	signer, err := db.GetStellarSigner(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing signer", cerr.GeneralError))
		return
	}
	if signer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}

	c.JSON(http.StatusOK, &signer.SignerSecretSeed)
}

//EditAllowtrustSignerRequest edits signer
//swagger:parameters EditAllowtrustSignerRequest EditAllowtrustSigner
type EditAllowtrustSignerRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
	//required : true
	Name        string  `form:"name" json:"name" validate:"required,max=256"`
	Description *string `form:"description" json:"description"`
}

//EditAllowtrustSigner - edits signer
// swagger:route POST /portal/admin/dash/stellar_account/edit_allowtrust_signer StellarAccount EditAllowtrustSigner
//
// Edits signer's details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func EditAllowtrustSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr EditAllowtrustSignerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	signer, err := db.GetStellarSigner(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing signer", cerr.GeneralError))
		return
	}
	if signer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(signer.Type, models.StellarSignerTypeAllowTrust) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Signer is not of type allow trust", ""))
		return
	}

	signer.Name = rr.Name
	if rr.Description != nil {
		signer.Description = *rr.Description
	}
	err = db.UpdateSigner(signer, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error editing signer", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//EditOtherSignerRequest edits signer
//swagger:parameters EditOtherSignerRequest EditOtherSigner
type EditOtherSignerRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
	//required : true
	Name        string  `form:"name" json:"name" validate:"required,max=256"`
	Description *string `form:"description" json:"description"`
}

//EditOtherSigner - edits signer
// swagger:route POST /portal/admin/dash/stellar_account/edit_other_signer StellarAccount EditOtherSigner
//
// Edits signer's details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func EditOtherSigner(uc *mw.AdminContext, c *gin.Context) {
	var rr EditOtherSignerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	signer, err := db.GetStellarSigner(rr.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing signer", cerr.GeneralError))
		return
	}
	if signer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key not found in database", ""))
		return
	}
	if !strings.EqualFold(signer.Type, models.StellarSignerTypeOther) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Signer is not of type other", ""))
		return
	}

	signer.Name = rr.Name
	if rr.Description != nil {
		signer.Description = *rr.Description
	}
	err = db.UpdateSigner(signer, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error editing signer", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
