package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/gin-gonic/gin"
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
	route.AddRoute("POST", "/update_name", UpdateIcoName, []string{}, "update_ico_name", ICORoutePrefix)
	route.AddRoute("POST", "/update_kyc", UpdateIcoKyc, []string{}, "update_ico_kyc", ICORoutePrefix)
	route.AddRoute("POST", "/update_issuer_data", UpdateIcoIssuer, []string{}, "update_ico_issuer_data", ICORoutePrefix)
	route.AddRoute("POST", "/update_supported_currencies", UpdateIcoCurrencies, []string{}, "update_ico_currencies", ICORoutePrefix)
}

//AddICORoutes adds all the routes for the ico management
func AddICORoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(ICORoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//ICOListResponse response
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
type ExchangeCurrencyListRequest struct {
	IcoID int `form:"ico_id" json:"ico_id"`
}

//ExchangeCurrencyListResponse response
type ExchangeCurrencyListResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	AssetCode string `json:"asset_code"`
	Issuer    string `json:"issuer"`
	Enabled   string `json:"enabled"`
}

//ExchangeCurrencyList returns the list of ICOs
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
type AddIcoRequest struct {
	Name                string `form:"name" json:"name" validate:"required,max=256"`
	Kyc                 bool   `form:"kyc" json:"kyc"`
	SalesModel          string `form:"sales_model" json:"sales_model" validate:"required"`
	IssuerPublicKey     string `form:"issuing_account_pk" json:"issuing_account_pk" validate:"required,base64,len=56"`
	AssetCode           string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
	SupportedCurrencies []int  `form:"supported_currencies" json:"supported_currencies" validate:"required"`
}

//AddIco - adds new ico
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
type UpdateIcoNameRequest struct {
	ID   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name" validate:"required,max=256"`
}

//UpdateIcoName - update the name
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
type UpdateIcoKycRequest struct {
	ID  int  `form:"id" json:"id"`
	Kyc bool `form:"kyc" json:"kyc"`
}

//UpdateIcoKyc - update the kyc flag
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
type UpdateIcoIssuerRequest struct {
	ID              int    `form:"id" json:"id"`
	IssuerPublicKey string `form:"issuing_account_pk" json:"issuing_account_pk" validate:"required,base64,len=56"`
	AssetCode       string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
}

//UpdateIcoIssuer - updates the issuer data
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
type UpdateIcoCurrenciesRequest struct {
	ID                  int   `form:"id" json:"id"`
	SupportedCurrencies []int `form:"supported_currencies" json:"supported_currencies" validate:"required"`
}

//UpdateIcoCurrencies - updates the supported currencies
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
