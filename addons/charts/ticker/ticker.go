package ticker

import (
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"
	"strconv"
	"time"
)

type exchange struct {
	source            string
	sourceIssuer      string
	destination       string
	destinationIssuer string
	rate              float64
	timestamp         int
}

// Ticker gets new exchange data every x seconds
func Ticker() {

	go coinmarketcapTicker()
	go stellarDecentralizedExchangeTicker()

}

func handleData(ex exchange) error {

	sourceCurrency, err := utils.GetCurrencyByCode(ex.source, ex.sourceIssuer, true)
	if err != nil {
		return err
	}

	destinationCurrency, err := utils.GetCurrencyByCode(ex.destination, ex.destinationIssuer, true)
	if err != nil {
		return err
	}

	var t time.Time
	if ex.timestamp != 0 {
		i, err := strconv.ParseInt(strconv.Itoa(ex.timestamp), 10, 64)
		if err != nil {
			return err
		}
		t = time.Unix(i, 0)
	} else {
		t = time.Now()
	}

	var minutelyData models.CurrentChartDataMinutely
	minutelyData.ExchangeRateTime = t
	minutelyData.ExchangeRate = ex.rate
	minutelyData.SourceCurrencyID = sourceCurrency.ID
	minutelyData.DestinationCurrencyID = destinationCurrency.ID

	// upsert set to do nothin on conflict, so we don't have duplicates in the db
	err = minutelyData.Upsert(utils.DB, false, nil, nil)
	if err != nil {
		return err
	}

	var hourlyData models.CurrentChartDataHourly
	hourlyData.ExchangeRateTime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location()) // only take date, hours and timezone. ommit rest
	hourlyData.ExchangeRate = ex.rate
	hourlyData.SourceCurrencyID = sourceCurrency.ID
	hourlyData.DestinationCurrencyID = destinationCurrency.ID

	// upsert set to update exchange rate on conflict, so we always have the latest from that hour in the db
	err = hourlyData.Upsert(utils.DB, true, []string{"exchange_rate_time", "source_currency_id", "destination_currency_id"}, []string{"exchange_rate"})
	if err != nil {
		return err
	}

	var dailyData models.HistoryChartDatum
	dailyData.ExchangeRateDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // only take date. ommit rest
	dailyData.ExchangeRate = ex.rate
	dailyData.SourceCurrencyID = sourceCurrency.ID
	dailyData.DestinationCurrencyID = destinationCurrency.ID

	// upsert set to update exchange rate on conflict, so we always have the latest from that hour in the db
	err = dailyData.Upsert(utils.DB, true, []string{"exchange_rate_date", "source_currency_id", "destination_currency_id"}, []string{"exchange_rate"})
	if err != nil {
		return err
	}

	return nil
}
