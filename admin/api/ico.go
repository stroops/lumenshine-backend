package api

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"

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
	//route.AddRoute("GET", "/get/:id", GetKnownCurrency, []string{}, "known_currencies_get", KnownCurrenciesRoutePrefix)

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
	ID                 int     `json:"ico_id"`
	Name               string  `json:"name"`
	Status             string  `json:"status"`
	KYC                bool    `json:"kyc"`
	SaleModel          string  `json:"sale_model"`
	TokensReleased     float64 `json:"tokens_released"`
	IssuingAccountName string  `json:"issuing_account_name"`
	IssuingAccountPK   string  `json:"issuing_account_pk"`
	AssetCode          string  `json:"asset_code"`
}

//ICOList returns the list of ICOs
func ICOList(uc *mw.AdminContext, c *gin.Context) {

	icos, err := m.Icos().All(db.DBC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading icos", cerr.GeneralError))
		return
	}

	var response []ICOListResponse

	for _, ico := range icos {

		response = append(response, ICOListResponse{
			ID:                 ico.ID,
			Name:               ico.IcoName,
			Status:             ico.IcoStatus,
			KYC:                ico.Kyc,
			SaleModel:          ico.SalesModel,
			TokensReleased:     0,  //TODO
			IssuingAccountName: "", //TODO
			IssuingAccountPK:   ico.IssuerPK,
			AssetCode:          ico.AssetCode,
		})

	}

	c.JSON(http.StatusOK, response)

}

//ExchangeCurrencyListResponse response
type ExchangeCurrencyListResponse struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	AssetCode string `json:"asset_code"`
	Issuer    string `json:"issuer"`
	Enabled   bool   `json:"enabled"`
}

//ExchangeCurrencyList returns the list of ICOs
func ExchangeCurrencyList(uc *mw.AdminContext, c *gin.Context) {

	excs, err := m.ExchangeCurrencies().All(db.DBC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading exchange currencies", cerr.GeneralError))
		return
	}

	var response []ExchangeCurrencyListResponse

	for _, exc := range excs {

		response = append(response, ExchangeCurrencyListResponse{
			ID:        exc.ID,
			Type:      exc.ExchangeCurrencyType,
			AssetCode: exc.AssetCode,
			Issuer:    exc.EcAssetIssuerPK,
			Enabled:   false, // TODO
		})

	}

	c.JSON(http.StatusOK, response)

}
