package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/stellar/go/clients/horizon"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/keypair"
)

const (
	//ICORoutePrefix for the ico manangement
	ICORoutePrefix = "ico"
)

//init setup all the routes for the known currencies handling
func init() {
	route.AddRoute("GET", "/list", ICOList, []string{}, "ico_list", ICORoutePrefix)
	route.AddRoute("GET", "/exchange_currencies", ExchangeCurrencyList, []string{}, "exchange_currencies", ICORoutePrefix)
	route.AddRoute("POST", "/add", AddIco, []string{}, "add_ico", ICORoutePrefix)
	route.AddRoute("POST", "/add_phase", AddIcoPhase, []string{}, "add_ico_phase", ICORoutePrefix)
	route.AddRoute("POST", "/update_name", UpdateIcoName, []string{}, "update_ico_name", ICORoutePrefix)
	route.AddRoute("POST", "/update_kyc", UpdateIcoKyc, []string{}, "update_ico_kyc", ICORoutePrefix)
	route.AddRoute("POST", "/update_issuer_data", UpdateIcoIssuer, []string{}, "update_ico_issuer_data", ICORoutePrefix)
	route.AddRoute("POST", "/update_supported_currencies", UpdateIcoCurrencies, []string{}, "update_ico_currencies", ICORoutePrefix)
	route.AddRoute("POST", "/remove", RemoveIco, []string{}, "remove_ico", ICORoutePrefix)
}

//AddICORoutes adds all the routes for the ico management
func AddICORoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(ICORoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//ICOListResponse response
// swagger:model
type ICOListResponse struct {
	ID                 int    `json:"ico_id"`
	Name               string `json:"name"`
	Status             string `json:"status"`
	KYC                bool   `json:"kyc"`
	SaleModel          string `json:"sale_model"`
	TokensReleased     int64  `json:"tokens_released"`
	IssuingAccountName string `json:"issuing_account_name"`
	IssuingAccountPK   string `json:"issuing_account_pk"`
	AssetCode          string `json:"asset_code"`
}

//ICOList returns the list of ICOs
// swagger:route GET /portal/admin/dash/ico/list ico ICOList
//
// Returns the list of ICOs
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:ICOListResponse
func ICOList(uc *mw.AdminContext, c *gin.Context) {
	icos, err := m.Icos(
		qm.Select(m.IcoColumns.ID,
			m.IcoColumns.IcoName,
			m.IcoColumns.IcoStatus,
			m.IcoColumns.Kyc,
			m.IcoColumns.SalesModel,
			m.IcoColumns.IssuerPK,
			m.IcoColumns.AssetCode,
		),
		qm.Load(m.IcoRels.IcoPhases)).All(db.DBC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading icos", cerr.GeneralError))
		return
	}

	dbAccounts, err := db.IssuerStellarAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading issuer accounts", cerr.GeneralError))
		return
	}

	response := make([]ICOListResponse, len(icos))
	var tokencount int64
	var issuerName string
	for i, ico := range icos {
		tokencount = 0
		issuerName = ""
		for _, phase := range ico.R.IcoPhases {
			tokencount += phase.TokensReleased
		}
		for _, issuer := range dbAccounts {
			if issuer.PublicKey == ico.IssuerPK {
				issuerName = issuer.Name
				break
			}
		}

		response[i] = ICOListResponse{
			ID:                 ico.ID,
			Name:               ico.IcoName,
			Status:             ico.IcoStatus,
			KYC:                ico.Kyc,
			SaleModel:          ico.SalesModel,
			TokensReleased:     tokencount,
			IssuingAccountName: issuerName,
			IssuingAccountPK:   ico.IssuerPK,
			AssetCode:          ico.AssetCode,
		}
	}
	c.JSON(http.StatusOK, response)
}

//ExchangeCurrencyListRequest request
//swagger:parameters ExchangeCurrencyListRequest ExchangeCurrencyList
type ExchangeCurrencyListRequest struct {
	//required : true
	IcoID int `form:"ico_id" json:"ico_id"`
}

//ExchangeCurrencyListResponse response
// swagger:model
type ExchangeCurrencyListResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	AssetCode string `json:"asset_code"`
	Issuer    string `json:"issuer"`
	Enabled   string `json:"enabled"`
}

//ExchangeCurrencyList returns the list of exchange currencies
// swagger:route GET /portal/admin/dash/ico/exchange_currencies ico ExchangeCurrencyList
//
// Returns the list of exchange currencies
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:ExchangeCurrencyListResponse
func ExchangeCurrencyList(uc *mw.AdminContext, c *gin.Context) {
	var rr ExchangeCurrencyListRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	excs, err := m.ExchangeCurrencies().All(db.DBC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading exchange currencies", cerr.GeneralError))
		return
	}

	var enabledCurrencies m.IcoSupportedExchangeCurrencySlice
	if rr.IcoID != 0 {
		enabledCurrencies, err = m.IcoSupportedExchangeCurrencies(
			qm.Select(m.IcoSupportedExchangeCurrencyColumns.ExchangeCurrencyID),
			qm.Where(m.IcoSupportedExchangeCurrencyColumns.IcoID+"=?", rr.IcoID)).All(db.DBC)

		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading supported exchange currencies", cerr.GeneralError))
			return
		}
	}

	response := make([]ExchangeCurrencyListResponse, len(excs))
	for i, exc := range excs {
		response[i] = ExchangeCurrencyListResponse{
			ID:        exc.ID,
			Name:      exc.Name,
			Type:      exc.ExchangeCurrencyType,
			AssetCode: exc.AssetCode,
			Issuer:    exc.EcAssetIssuerPK,
		}
		if rr.IcoID != 0 {
			enabled := "false"
			for _, currency := range enabledCurrencies {
				if currency.ExchangeCurrencyID == exc.ID {
					enabled = "true"
					break
				}
			}
			response[i].Enabled = enabled
		}
	}
	c.JSON(http.StatusOK, response)
}

//AddIcoRequest - request
//swagger:parameters AddIcoRequest AddIco
type AddIcoRequest struct {
	//required : true
	Name string `form:"name" json:"name" validate:"required,max=256"`
	Kyc  bool   `form:"kyc" json:"kyc"`
	//required : true
	SalesModel string `form:"sales_model" json:"sales_model" validate:"required"`
	//required : true
	IssuerPublicKey string `form:"issuing_account_pk" json:"issuing_account_pk" validate:"required,base64,len=56"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
	//required : true
	SupportedCurrencies []int `form:"supported_currencies" json:"supported_currencies" validate:"required"`
}

//AddIco - adds new ico
// swagger:route POST /portal/admin/dash/ico/add ico AddIco
//
// Adds new ico
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddIco(uc *mw.AdminContext, c *gin.Context) {
	var rr AddIcoRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	existsName, err := db.ExistsIcoName(rr.Name, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if existsName {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "name", cerr.InvalidArgument, "Ico name already exists", ""))
		return
	}
	if rr.SalesModel != m.IcoSalesModelFixed {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "sales_model", cerr.InvalidArgument, "Sales model is not 'fixed'", ""))
		return
	}
	issuer, err := db.GetStellarAccount(rr.IssuerPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing issuing account", cerr.GeneralError))
		return
	}
	if issuer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_pk", cerr.InvalidArgument, "Issuing account does not exist", ""))
		return
	}
	if issuer.Type != models.StellarAccountTypeIssuing {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_pk", cerr.InvalidArgument, "Issuing public key does not belong to an issuing account", ""))
		return
	}
	existsAssetCode := false
	for _, assetCode := range issuer.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			existsAssetCode = true
			break
		}
	}
	if !existsAssetCode {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code does not exist for this issuing account", ""))
		return
	}
	if len(rr.SupportedCurrencies) == 0 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "supported_currencies", cerr.InvalidArgument, "Supported currencies is empty.", ""))
		return
	}

	supportedCurrencies := make([]m.IcoSupportedExchangeCurrency, len(rr.SupportedCurrencies))
	for i, currencyID := range rr.SupportedCurrencies {
		exists, err := db.ExistsCurrency(currencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "supported_currencies", cerr.InvalidArgument, fmt.Sprintf("A currency does not exist for the id: %d", currencyID), ""))
			return
		}
		supportedCurrencies[i] = m.IcoSupportedExchangeCurrency{ExchangeCurrencyID: currencyID, UpdatedBy: getUpdatedBy(c)}
	}

	tx, err := db.DBC.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error begining internal transaction", cerr.GeneralError))
		return
	}
	ico := m.Ico{
		IcoName:    rr.Name,
		AssetCode:  rr.AssetCode,
		IssuerPK:   rr.IssuerPublicKey,
		Kyc:        rr.Kyc,
		SalesModel: rr.SalesModel,
		IcoStatus:  m.IcoStatusPlanning,
		UpdatedBy:  getUpdatedBy(c),
	}
	err = ico.Insert(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error inserting ico", cerr.GeneralError))
		return
	}
	for _, dbCurrency := range supportedCurrencies {
		dbCurrency.IcoID = ico.ID
		err = dbCurrency.Insert(tx, boil.Whitelist(m.IcoSupportedExchangeCurrencyColumns.IcoID,
			m.IcoSupportedExchangeCurrencyColumns.ExchangeCurrencyID,
			m.IcoSupportedExchangeCurrencyColumns.UpdatedBy))
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error inserting ico supported currency", cerr.GeneralError))
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, "{}")
}

//UpdateIcoNameRequest - request
//swagger:parameters UpdateIcoNameRequest UpdateIcoName
type UpdateIcoNameRequest struct {
	//required : true
	ID int `form:"id" json:"id"`
	//required : true
	Name string `form:"name" json:"name" validate:"required,max=256"`
}

//UpdateIcoName - updates the name
// swagger:route POST /portal/admin/dash/ico/update_name ico UpdateIcoName
//
// Updates the ico name
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateIcoName(uc *mw.AdminContext, c *gin.Context) {
	var rr UpdateIcoNameRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	ico, err := db.GetIco(rr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	existsName, err := db.ExistsIcoName(rr.Name, &ico.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if existsName {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "name", cerr.InvalidArgument, "Ico name already exists", ""))
		return
	}
	ico.IcoName = rr.Name
	_, err = ico.Update(db.DBC, boil.Whitelist(m.IcoColumns.IcoName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating ico name", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}

//UpdateIcoKycRequest - request
//swagger:parameters UpdateIcoKycRequest UpdateIcoKyc
type UpdateIcoKycRequest struct {
	//required : true
	ID int `form:"id" json:"id"`
	//required : true
	Kyc bool `form:"kyc" json:"kyc"`
}

//UpdateIcoKyc - updates the kyc flag
// swagger:route POST /portal/admin/dash/ico/update_kyc ico UpdateIcoKyc
//
// Updates the KYC flag
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateIcoKyc(uc *mw.AdminContext, c *gin.Context) {
	var rr UpdateIcoKycRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	ico, err := db.GetIco(rr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	ico.Kyc = rr.Kyc
	_, err = ico.Update(db.DBC, boil.Whitelist(m.IcoColumns.Kyc))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating ico kyc flag", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}

//UpdateIcoIssuerRequest - request
//swagger:parameters UpdateIcoIssuerRequest UpdateIcoIssuer
type UpdateIcoIssuerRequest struct {
	//required : true
	ID int `form:"id" json:"id"`
	//required : true
	IssuerPublicKey string `form:"issuing_account_pk" json:"issuing_account_pk" validate:"required,base64,len=56"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
}

//UpdateIcoIssuer - updates the issuer data
// swagger:route POST /portal/admin/dash/ico/update_issuer_data ico UpdateIcoIssuer
//
// Updates the issuer data
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateIcoIssuer(uc *mw.AdminContext, c *gin.Context) {
	var rr UpdateIcoIssuerRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	issuer, err := db.GetStellarAccount(rr.IssuerPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing issuing account", cerr.GeneralError))
		return
	}
	if issuer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_pk", cerr.InvalidArgument, "Issuing account does not exist", ""))
		return
	}
	if issuer.Type != models.StellarAccountTypeIssuing {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_pk", cerr.InvalidArgument, "Issuing public key does not belong to an issuing account", ""))
		return
	}
	existsAssetCode := false
	for _, assetCode := range issuer.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			existsAssetCode = true
			break
		}
	}
	if !existsAssetCode {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code does not exist for this issuing account", ""))
		return
	}
	ico, err := db.GetIco(rr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	ico.IssuerPK = rr.IssuerPublicKey
	ico.AssetCode = rr.AssetCode
	_, err = ico.Update(db.DBC, boil.Whitelist(m.IcoColumns.IssuerPK, m.IcoColumns.AssetCode))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating ico issuer data", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}

//UpdateIcoCurrenciesRequest - request
//swagger:parameters UpdateIcoCurrenciesRequest UpdateIcoCurrencies
type UpdateIcoCurrenciesRequest struct {
	//required : true
	ID int `form:"id" json:"id"`
	//required : true
	SupportedCurrencies []int `form:"supported_currencies" json:"supported_currencies" validate:"required"`
}

//UpdateIcoCurrencies - updates the supported currencies
// swagger:route POST /portal/admin/dash/ico/update_supported_currencies ico UpdateIcoCurrencies
//
// Updates the supported currencies
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateIcoCurrencies(uc *mw.AdminContext, c *gin.Context) {
	var rr UpdateIcoCurrenciesRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	ico, err := db.GetIco(rr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	if len(rr.SupportedCurrencies) == 0 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "supported_currencies", cerr.InvalidArgument, "Supported currencies is empty.", ""))
		return
	}
	supportedCurrencies := make([]*m.IcoSupportedExchangeCurrency, len(rr.SupportedCurrencies))
	for i, currencyID := range rr.SupportedCurrencies {
		exists, err := db.ExistsCurrency(currencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "supported_currencies", cerr.InvalidArgument, fmt.Sprintf("A currency does not exist for the id: %d", currencyID), ""))
			return
		}
		supportedCurrencies[i] = &m.IcoSupportedExchangeCurrency{IcoID: ico.ID, ExchangeCurrencyID: currencyID, UpdatedBy: getUpdatedBy(c)}
	}
	err = db.UpdateSupportedCurrencies(ico.ID, supportedCurrencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating supported currencies", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//RemoveIcoRequest - request
//swagger:parameters RemoveIcoRequest RemoveIco
type RemoveIcoRequest struct {
	//required : true
	ID int `form:"id" json:"id"`
}

//RemoveIco - deletes ico from db
// swagger:route POST /portal/admin/dash/ico/remove ico RemoveIco
//
// Deletes ico from db
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveIco(uc *mw.AdminContext, c *gin.Context) {
	var rr RemoveIcoRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	ico, err := db.GetIco(rr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	if ico.IcoStatus != m.IcoStatusPlanning {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "The ico is not in 'planning' state", ""))
		return
	}
	err = db.DeleteIco(ico)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting ico", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//AddIcoPhaseRequest - request
//swagger:parameters AddIcoPhaseRequest AddIcoPhase
type AddIcoPhaseRequest struct {
	//required : true
	IcoID int `form:"ico_id" json:"ico_id"`
	//required : true
	Name string `form:"phase_name" json:"phase_name" validate:"required,max=256"`
	//required : true
	DistributionPublicKey string `form:"distribution_account_pk" json:"distribution_account_pk" validate:"required,base64,len=56"`
	//required : true
	PreSignerPublicKey string `form:"pre_signer_pk" json:"pre_signer_pk" validate:"required,base64,len=56"`
	//required : true
	PreSignerSeed string `form:"pre_signer_seed" json:"pre_signer_seed" validate:"required,base64,len=56"`
	//required : true
	PostSignerPublicKey string `form:"post_signer_pk" json:"post_signer_pk" validate:"required,base64,len=56"`
	//required : true
	PostSignerSeed string `form:"post_signer_seed" json:"post_signer_seed" validate:"required,base64,len=56"`
	//required : true
	StartDate string `form:"start" json:"start" validate:"required"`
	//required : true
	EndDate string `form:"end" json:"end" validate:"required"`
	//required : true
	TokensToDistribute int64 `form:"tokens_to_distribute" json:"tokens_to_distribute"`
	//required : true
	MinPerOrder int64 `form:"min_tokens_per_order" json:"min_tokens_per_order"`
	//required : true
	MaxPerOrder int64 `form:"max_tokens_per_order" json:"max_tokens_per_order"`
	//required : true
	ActivatedCurrencies []PhaseCurrency `form:"activated_currencies" json:"activated_currencies" validate:"required"`
}

//PhaseCurrency - request
type PhaseCurrency struct {
	CurrencyID    int   `form:"currency_id" json:"currency_id"`
	Price         int64 `form:"price" json:"price"`
	BankAccountID *int  `form:"bank_account_id" json:"bank_account_id"`
}

//AddIcoPhase - adds new ico phase
// swagger:route POST /portal/admin/dash/ico/add_phase ico AddIcoPhase
//
// Adds new ico phase
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddIcoPhase(uc *mw.AdminContext, c *gin.Context) {
	var rr AddIcoPhaseRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	ico, err := db.GetIcoEager(rr.IcoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading ico from db", cerr.GeneralError))
		return
	}
	if ico == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "ico_id", cerr.InvalidArgument, "Ico not found for the specified id", ""))
		return
	}
	startDate, err := time.Parse("2006-01-02", rr.StartDate)
	if rr.StartDate != "" && err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("start", cerr.InvalidArgument, "Start date wrong format", ""))
		return
	}
	endDate, err := time.Parse("2006-01-02", rr.EndDate)
	if rr.EndDate != "" && err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("end", cerr.InvalidArgument, "End date wrong format", ""))
		return
	}
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("end", cerr.InvalidArgument, "End date smaller than start date", ""))
		return
	}
	if ico.R.IcoPhases != nil {
		for _, phase := range ico.R.IcoPhases {
			if inTimeSpan(phase.StartTime, phase.EndTime, startDate) {
				c.JSON(http.StatusBadRequest, cerr.NewIcopError("start", cerr.InvalidArgument, "Start date overlaping with existing ico phase", ""))
				return
			}
			if inTimeSpan(phase.StartTime, phase.EndTime, endDate) {
				c.JSON(http.StatusBadRequest, cerr.NewIcopError("end", cerr.InvalidArgument, "End date overlaping with existing ico phase", ""))
				return
			}
		}
	}

	account, err := db.GetStellarAccount(rr.DistributionPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading stellar account from db", cerr.GeneralError))
		return
	}
	if account == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "distribution_account_pk", cerr.InvalidArgument, "Distribution account not found in db", ""))
		return
	}
	foundPreSigner := false
	foundPostSigner := false
	for _, signer := range account.R.StellarAccountPublicKeyAdminStellarSigners {
		if signer.SignerPublicKey == rr.PreSignerPublicKey {
			foundPreSigner = true
		}
		if signer.SignerPublicKey == rr.PostSignerPublicKey {
			foundPostSigner = true
		}
	}
	if !foundPreSigner {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_pk", cerr.InvalidArgument, "Presigner not found in db", ""))
		return
	}
	if !foundPostSigner {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "post_signer_pk", cerr.InvalidArgument, "Postsigner not found in db", ""))
		return
	}
	parsed, err := keypair.Parse(rr.PreSignerSeed)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_seed", cerr.InvalidArgument, "Error parsing presigner seed", ""))
		return
	}
	if parsed.Address() != rr.PreSignerPublicKey {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_seed", cerr.InvalidArgument, "Presigner seed does not match public key", ""))
		return
	}
	parsed, err = keypair.Parse(rr.PostSignerSeed)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "post_signer_seed", cerr.InvalidArgument, "Error parsing postsigner seed", ""))
		return
	}
	if parsed.Address() != rr.PostSignerPublicKey {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "post_signer_seed", cerr.InvalidArgument, "Postsigner seed does not match public key", ""))
		return
	}
	currencies := make([]*m.IcoPhaseActivatedExchangeCurrency, 0)
	for _, currency := range rr.ActivatedCurrencies {
		if !inSupportedCurrencies(ico, currency.CurrencyID) {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "activated_currencies", cerr.InvalidArgument, fmt.Sprintf("Currency id: %d is not supported by the ico", currency.CurrencyID), ""))
			return
		}
		dbCurrency := m.IcoPhaseActivatedExchangeCurrency{
			ExchangeCurrencyID: currency.CurrencyID,
			DenomPricePerToken: currency.Price,
			UpdatedBy:          getUpdatedBy(c),
		}
		if currency.BankAccountID != nil {
			dbCurrency.IcoPhaseBankAccountID = null.IntFrom(*currency.BankAccountID)
		}
		currencies = append(currencies, &dbCurrency)
	}
	hAccount, exists, err := GetHorizonAccount(rr.DistributionPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error calling horizon", cerr.GeneralError))
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "distribution_account_pk", cerr.InvalidArgument, "Distribution account not found in horizon", ""))
		return
	}
	balance := getAssetBalance(hAccount, ico)
	if balance == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "distribution_account_pk", cerr.InvalidArgument, "Distribution account does not trust ico's asset code", ""))
		return
	}
	totalBalance, err := strconv.ParseInt(balance.Balance, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not parse total balance", cerr.GeneralError))
		return
	}
	if rr.TokensToDistribute > totalBalance {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "tokens_to_distribute", cerr.InvalidArgument, "Total balance not sufficient for tokens to distribute", ""))
		return
	}
	preSigner := getSigner(hAccount, rr.PreSignerPublicKey)
	if preSigner == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_pk", cerr.InvalidArgument, "Presigner is not found for the specified distribution account", ""))
		return
	}
	postSigner := getSigner(hAccount, rr.PostSignerPublicKey)
	if postSigner == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "post_signer_pk", cerr.InvalidArgument, "Postsigner is not found for the specified distribution account", ""))
		return
	}
	medThreshold := int32(hAccount.Thresholds.MedThreshold)
	if preSigner.Weight >= medThreshold {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_pk", cerr.InvalidArgument, "Presigner's weight is larger than medium threshold", ""))
		return
	}
	if preSigner.Weight+postSigner.Weight < medThreshold {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "pre_signer_pk", cerr.InvalidArgument, "Weight sum of pre and post signers is smaller than medium threshold", ""))
		return
	}
	phase := m.IcoPhase{
		IcoID:               ico.ID,
		IcoPhaseName:        rr.Name,
		IcoPhaseStatus:      m.IcoPhaseStatusPlanning,
		DistPK:              rr.DistributionPublicKey,
		DistPresignerPK:     rr.PreSignerPublicKey,
		DistPresignerSeed:   rr.PreSignerSeed,
		DistPostsignerPK:    rr.PostSignerPublicKey,
		DistPostsignerSeed:  rr.PostSignerSeed,
		StartTime:           startDate,
		EndTime:             endDate,
		TokensToDistribute:  rr.TokensToDistribute,
		TokenMinOrderAmount: rr.MinPerOrder,
		TokenMaxOrderAmount: rr.MaxPerOrder,
		UpdatedBy:           getUpdatedBy(c),
	}

	err = db.AddIcoPhase(&phase, currencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding ico phase to db", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func inSupportedCurrencies(ico *m.Ico, currencyID int) bool {
	found := false
	if ico.R.IcoSupportedExchangeCurrencies != nil {
		for _, currency := range ico.R.IcoSupportedExchangeCurrencies {
			if currency.ExchangeCurrencyID == currencyID {
				found = true
				break
			}
		}
	}
	return found
}

func getAssetBalance(hAccount horizon.Account, ico *m.Ico) *horizon.Balance {
	if hAccount.Balances != nil {
		for _, balance := range hAccount.Balances {
			if balance.Asset.Code == ico.AssetCode && balance.Asset.Issuer == ico.IssuerPK {
				return &balance
			}
		}
	}
	return nil
}

func getSigner(hAccount horizon.Account, publicKey string) *horizon.Signer {
	if hAccount.Signers != nil {
		for _, signer := range hAccount.Signers {
			if signer.Key == publicKey {
				return &signer
			}
		}
	}
	return nil
}
