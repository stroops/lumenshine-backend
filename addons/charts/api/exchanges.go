package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var errorCodesExchange = map[string]errorCodesData{
	"source_missing": {
		code:        1500,
		parameter:   "source_currency",
		description: "source currency parameter missing",
	},
	"source_inexistent": {
		code:        1501,
		parameter:   "source_currency",
		description: "source currency not available",
	},
	"destination_missing": {
		code:        1502,
		parameter:   "destination_currency",
		description: "destination currency parameter missing",
	},
	"destination_inexistent": {
		code:        1503,
		parameter:   "destination_currency",
		description: "destination currency not available",
	},
	"range_hours_missing": {
		code:        1504,
		parameter:   "range_hours",
		description: "range of hours is missing",
	},
	"range_hours_wrong_format": {
		code:        1505,
		parameter:   "range_hours",
		description: "range of hours has wrong format (use whole positive numbers)",
	},
	"internal_error": {
		code: 500,
	},
}

type rateDataExchange struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

type exchangeDataExchange struct {
	SourceCurrency struct {
		AssetCode       string `json:"asset_code"`
		IssuerPublicKey string `json:"issuer_public_key"`
	} `json:"source_currency"`
	DestinationCurrency string             `json:"destination_currency"`
	CurrentRate         float64            `json:"current_rate"`
	LastUpdated         string             `json:"last_updated"`
	Rates               []rateDataExchange `json:"rates"`
}

// Binding from JSON
type requestDataExchange struct {
	SourceCurrency struct {
		AssetCode       string `form:"asset_code" json:"asset_code" binding:"required"`
		IssuerPublicKey string `form:"issuer_public_key" json:"issuer_public_key" binding:"required"`
	} `form:"source_currency" json:"source_currency" binding:"required"`
	DestinationCurrency string `form:"destination_currency" json:"destination_currency" binding:"required"`
	RangeHours          int    `form:"range_hours" json:"range_hours" binding:"required"`
}

// ChartExchangeData returns json data containg exchanges between two currencies within a given period of time
func ChartExchangeData(uc *mw.IcopContext, c *gin.Context) {

	var requestData requestDataExchange
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
	if requestData.SourceCurrency.AssetCode == "" {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["source_missing"].code
		errData.ParameterName = errorCodesExchange["source_missing"].parameter
		errData.ErrorMessage = errorCodesExchange["source_missing"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	// check if source param is correct
	if requestData.SourceCurrency.AssetCode != "XLM" && requestData.SourceCurrency.IssuerPublicKey == "" {
		errData := errorData{}
		errData.ErrorCode = errorCodesCurrent["issuer_public_key_missing"].code
		errData.ParameterName = errorCodesCurrent["issuer_public_key_missing"].parameter
		errData.ErrorMessage = fmt.Sprintf(errorCodesCurrent["issuer_public_key_missing"].description, requestData.SourceCurrency.AssetCode)
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	if requestData.DestinationCurrency == "" {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["destination_missing"].code
		errData.ParameterName = errorCodesExchange["destination_missing"].parameter
		errData.ErrorMessage = errorCodesExchange["destination_missing"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	// if range hours is missing it will be taken as 0
	// if requestData.RangeHours == nil {
	// 	errData := errorData{}
	// 	errData.ErrorCode = errorCodesExchange["range_hours_missing"].code
	// 	errData.ParameterName = errorCodesExchange["range_hours_missing"].parameter
	// 	errData.ErrorMessage = errorCodesExchange["range_hours_missing"].description
	// 	c.JSON(http.StatusBadRequest, errData)
	// 	return
	// }

	// check if params are correct
	sourceCurrency, err := utils.GetCurrencyByCode(requestData.SourceCurrency.AssetCode, requestData.SourceCurrency.IssuerPublicKey, false)
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["source_inexistent"].code
		errData.ParameterName = errorCodesExchange["source_inexistent"].parameter
		errData.ErrorMessage = errorCodesExchange["source_inexistent"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	destinationCurrency, err := utils.GetCurrencyByCode(requestData.DestinationCurrency, config.ExternalCurrencyIssuer, false)
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["destination_inexistent"].code
		errData.ParameterName = errorCodesExchange["destination_inexistent"].parameter
		errData.ErrorMessage = errorCodesExchange["destination_inexistent"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	if requestData.RangeHours < 0 {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["range_hours_wrong_format"].code
		errData.ParameterName = errorCodesExchange["range_hours_wrong_format"].parameter
		errData.ErrorMessage = errorCodesExchange["range_hours_wrong_format"].description
		c.JSON(http.StatusBadRequest, errData)
		return
	}

	data := exchangeDataExchange{}
	data.SourceCurrency.AssetCode = sourceCurrency.CurrencyCode
	data.SourceCurrency.IssuerPublicKey = requestData.SourceCurrency.IssuerPublicKey
	data.DestinationCurrency = destinationCurrency.CurrencyCode
	lastTransaction, err := utils.GetCurrentRate(sourceCurrency, destinationCurrency)
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["internal_error"].code
		errData.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, errData)
		return
	}
	data.CurrentRate = lastTransaction.ExchangeRate
	data.LastUpdated = lastTransaction.ExchangeRateTime.Format(config.TimeFormat)

	data.Rates, err = getRates(sourceCurrency, destinationCurrency, requestData.RangeHours)
	if err != nil {
		errData := errorData{}
		errData.ErrorCode = errorCodesExchange["internal_error"].code
		errData.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, errData)
		return
	}

	c.JSON(http.StatusOK, data)

}

type chartData struct {
	ExchangeRate     float64
	ExchangeRateDate time.Time
	ExchangeRateTime time.Time
}

// 1. if number of hours is larger than available data, return available data
// 2. if number of hours is 0 than only return the current rate
// 3. if number of hours is equal 1 hours return data per minute
// 4. if number of hours is more than 1 hour but smaller or equal 3 hours return data for every 5 minutes
// 5. if number of hours is more than 3 hours but smaller or equal 6 hours return data for every 10 minutes
// 6. if number of hours is more than 6 hours but smaller or equal 12 hours return data for every 15 minutes
// 7. if number of hours is more than 12 hours but smaller or equal 24 hours return data for every 30 minutes
// 8. if number of hours is more than 24 hours but smaller or equal 48 hours return data for every 1 hour
// 9. if number of hours is more that 48 hours but smaller or equal 96 hours return data for every 2 hours
// 10. if number of hours is more that 96 hours but smaller or equal 192 hours return data for every 4 hours
// 11. if number of hours is more that 192 hours but smaller or equal 384 hours return data for every 8 hours
// 12. if number of hours is more that 384 hours but smaller or equal 720 hours return data for every 16 hours
// 13. if number of hours is more than 720 return data for every day from history table

func getRates(sourceCurrency *models.Currency, destinationCurrency *models.Currency, rangeHours int) ([]rateDataExchange, error) {

	data := []rateDataExchange{}
	schema := utils.GetSchemaForQuery()
	var chartD []chartData
	var table, dateCol string
	step := 1

	if rangeHours <= 0 {
		return data, nil
	} else if rangeHours <= 24 {
		//take from minutely table
		table = "current_chart_data_minutely"
		dateCol = "exchange_rate_time"

		if rangeHours == 1 {
			step = 1
		} else if rangeHours <= 3 {
			// step = 5
			step = 1
		} else if rangeHours <= 6 {
			// step = 10
			step = 2
		} else if rangeHours <= 12 {
			// step = 15
			step = 3
		} else {
			// step = 30
			step = 6
		}

	} else if rangeHours <= 720 {
		//take from hourly table
		table = "current_chart_data_hourly"
		dateCol = "exchange_rate_time"

		if rangeHours <= 48 {
			step = 1
		} else if rangeHours <= 96 {
			step = 2
		} else if rangeHours <= 192 {
			step = 4
		} else if rangeHours <= 384 {
			step = 8
		} else {
			step = 16
		}
	} else {
		//take from history table
		table = "history_chart_data"
		dateCol = "exchange_rate_date"
	}

	fromTime := time.Now().Add(-time.Hour * time.Duration(rangeHours))

	// Use query building
	//TODO: Theo check
	err := models.NewQuery(
		qm.Select(schema+table+".exchange_rate, "+schema+table+"."+dateCol),
		qm.From(schema+table),
		qm.Where("source_currency_id=? AND destination_currency_id=? AND "+dateCol+">?", sourceCurrency.ID, destinationCurrency.ID, fromTime),
		qm.OrderBy(dateCol+" DESC"),
	).Bind(nil, utils.DB, &chartD)

	if err != nil {
		return data, err
	}

	for i, x := range chartD {

		if step != 1 && i%step != 0 {
			continue
		}

		t := x.ExchangeRateTime.Format(config.TimeFormat)
		if table == "history_chart_data" {
			t = x.ExchangeRateDate.Format(config.DateFormat)
		}

		data = append(data, rateDataExchange{
			Date: t,
			Rate: x.ExchangeRate,
		})
	}

	return data, nil
}
