package api

import (
	"fmt"
	"net/http"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

var errorCodesCurrent = map[string]errorCodesData{
	"source_missing": {
		code:        1500,
		parameter:   "source",
		description: "source currency parameter missing",
	},
	"source_inexistent": {
		code:        1501,
		parameter:   "source",
		description: "source currency %s not available",
	},
	"destination_missing": {
		code:        1502,
		parameter:   "destination",
		description: "destination currency parameter missing",
	},
	"destination_inexistent": {
		code:        1503,
		parameter:   "destination",
		description: "destination currency not available",
	},
	"issuer_public_key_missing": {
		code:        1504,
		parameter:   "issuer_public_key",
		description: "issuer public key is missing for currency %s",
	},
	"internal_error": {
		code: 500,
	},
}

type rateDataCurrent struct {
	SourceCurrency struct {
		AssetCode       string `json:"asset_code"`
		IssuerPublicKey string `json:"issuer_public_key"`
	} `json:"source_currency"`
	Rate        float64 `json:"rate"`
	LastUpdated string  `json:"last_updated"`
}

type exchangeDataCurrent struct {
	DestinationCurrency string            `json:"destination_currency"`
	Rates               []rateDataCurrent `json:"rates"`
}

// Binding from JSON
type requestDataCurrent struct {
	SourceCurrencies []struct {
		AssetCode       string `form:"asset_code" json:"asset_code" binding:"required"`
		IssuerPublicKey string `form:"issuer_public_key" json:"issuer_public_key" binding:"required"`
	} `form:"source_currencies" json:"source_currencies" binding:"required"`
	DestinationCurrency string `form:"destination_currency" json:"destination_currency" binding:"required"`
}

// ChartCurrentRates returns json data containing current rates between multiple source currencies and a FIAT destination currency
func ChartCurrentRates(uc *mw.IcopContext, c *gin.Context) {

	var requestData requestDataCurrent
	err := c.BindJSON(&requestData)

	//bad request
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = 400
		errData.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, errData)
		return
	}

	// check if params exist
	if requestData.SourceCurrencies == nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesCurrent["source_missing"].code
		errData.ParameterName = errorCodesCurrent["source_missing"].parameter
		errData.ErrorMessage = errorCodesCurrent["source_missing"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	if requestData.DestinationCurrency == "" {
		errData := errorData{}
		errData.ErrorCode = errorCodesCurrent["destination_missing"].code
		errData.ParameterName = errorCodesCurrent["destination_missing"].parameter
		errData.ErrorMessage = errorCodesCurrent["destination_missing"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}

	// check if destination param is correct
	destinationCurrency, err := utils.GetCurrencyByCode(requestData.DestinationCurrency, config.ExternalCurrencyIssuer, false)
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesCurrent["destination_inexistent"].code
		errData.ParameterName = errorCodesCurrent["destination_inexistent"].parameter
		errData.ErrorMessage = errorCodesCurrent["destination_inexistent"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}

	data := exchangeDataCurrent{}
	data.DestinationCurrency = destinationCurrency.CurrencyCode

	for _, source := range requestData.SourceCurrencies {
		// check if source param is correct
		if source.AssetCode != "XLM" && source.IssuerPublicKey == "" {
			errData := errorData{}
			errData.ErrorCode = errorCodesCurrent["issuer_public_key_missing"].code
			errData.ParameterName = errorCodesCurrent["issuer_public_key_missing"].parameter
			errData.ErrorMessage = fmt.Sprintf(errorCodesCurrent["issuer_public_key_missing"].description, source.AssetCode)
			c.JSON(http.StatusBadRequest, errData)
			return
		}

		sourceCurrency, err := utils.GetCurrencyByCode(source.AssetCode, source.IssuerPublicKey, false)
		if err != nil {
			errData := errorData{}
			errData.ErrorCode = errorCodesCurrent["source_inexistent"].code
			errData.ParameterName = errorCodesCurrent["source_inexistent"].parameter
			errData.ErrorMessage = fmt.Sprintf(errorCodesCurrent["source_inexistent"].description, source.AssetCode)
			c.JSON(http.StatusBadRequest, errData)
			return
		}

		// add data to response struct
		sourceData := rateDataCurrent{}
		sourceData.SourceCurrency.AssetCode = source.AssetCode
		sourceData.SourceCurrency.IssuerPublicKey = source.IssuerPublicKey
		lastTransaction, err := utils.GetCurrentRate(sourceCurrency, destinationCurrency)
		if err != nil {
			errData := errorData{}
			errData.ErrorCode = errorCodesCurrent["internal_error"].code
			errData.ErrorMessage = err.Error()
			c.JSON(http.StatusBadRequest, errData)
			return
		}
		sourceData.Rate = lastTransaction.ExchangeRate
		sourceData.LastUpdated = lastTransaction.ExchangeRateTime.Format(config.TimeFormat)

		data.Rates = append(data.Rates, sourceData)

	}

	c.JSON(http.StatusOK, data)

}
