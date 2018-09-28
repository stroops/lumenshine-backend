package api

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
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
