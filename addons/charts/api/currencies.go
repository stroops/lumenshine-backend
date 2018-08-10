package api

import (
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"
	"net/http"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type exchangeCurrencyPairs struct {
	SourceCurrency struct {
		AssetCode       string `json:"asset_code"`
		IssuerPublicKey string `json:"issuer_public_key"`
	} `json:"source_currency"`
	DestinationCurrencies []string `json:"destination_currencies"`
}

// ChartCurrencyPairs returns json with possible exchange rate pairs
func ChartCurrencyPairs(uc *mw.IcopContext, c *gin.Context) {

	data := []exchangeCurrencyPairs{}
	schema := utils.GetSchemaForQuery()

	currencies, _ := models.Currencies(utils.DB).All()
	for _, sourceCurrency := range currencies {

		destinationCurrencies, err := models.Currencies(utils.DB,
			qm.Select(schema+"currency.*"),
			qm.InnerJoin(schema+"history_chart_data h on h.destination_currency_id = "+schema+"currency.id"),
			qm.Where("h.source_currency_id =?", sourceCurrency.ID),
			qm.GroupBy(schema+"currency.id"),
		).All()

		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting data from db", cerr.GeneralError))
		}

		pair := exchangeCurrencyPairs{}
		pair.SourceCurrency.AssetCode = sourceCurrency.CurrencyCode
		pair.SourceCurrency.IssuerPublicKey = sourceCurrency.CurrencyIssuer

		for _, destinationCurrency := range destinationCurrencies {
			pair.DestinationCurrencies = append(pair.DestinationCurrencies, destinationCurrency.CurrencyCode)
		}

		if len(pair.DestinationCurrencies) > 0 {
			data = append(data, pair)
		}

	}

	c.JSON(http.StatusOK, data)

}
