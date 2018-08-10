package ticker

import (
	"encoding/json"
	"fmt"
	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ticker for coinmarketcap
func coinmarketcapTicker() {
	coinmarketcapSourceCurrencyID, err := getCoinmarketcapListing(config.Cnf.Ticker.SourceCurrency)
	if err != nil {
		log.Printf("Error getting internalID from Coinmarketcap %v", err)
		return
	}

	ex := exchange{}
	ex.source = config.Cnf.Ticker.SourceCurrency
	ex.sourceIssuer = config.ExternalCurrencyIssuer

	for {

		for _, dest := range config.Cnf.Ticker.DestinationCurrency {

			ex.destination = dest
			ex.destinationIssuer = config.ExternalCurrencyIssuer
			err := getCoinmarketcapData(&ex, coinmarketcapSourceCurrencyID)
			if err != nil {
				log.Printf("Error getting ticker data from Coinmarketcap %v", err)
				continue
			}

			err = handleData(ex)
			if err != nil {
				log.Printf("Error inserting data to minutely db table %v", err)
			}

		}

		waitTime := time.Second * time.Duration(config.Cnf.Ticker.TickSeconds)
		time.Sleep(waitTime)

	}
}

type quotes struct {
	Price            float64 `json:"price"`
	Volume24H        float64 `json:"volume_24h"`
	MarketCap        float64 `json:"market_cap"`
	PercentChange1H  float32 `json:"percent_change_1h"`
	PercentChange24H float32 `json:"percent_change_24h"`
	PercentChange7D  float32 `json:"percent_change_7d"`
}

type data struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Symbol            string            `json:"symbol"`
	WebsiteSlug       string            `json:"website_slug"`
	Rank              int64             `json:"rank"`
	CirculatingSupply float64           `json:"circulating_supply"`
	TotalSupply       float64           `json:"total_supply"`
	MaxSupply         float64           `json:"max_supply"`
	Quotes            map[string]quotes `json:"quotes"`
	LastUpdated       int               `json:"last_updated"`
}

type metadata struct {
	Timestamp int64  `json:"timestamp"`
	Error     string `json:"error"`
}

// json response from Coinmarketcap ticker API
type coinmarketcapTickerAPIResponse struct {
	Data     data     `json:"data"`
	Metadata metadata `json:"metadata"`
}

// performs a get request to the of the following form in order to retrieve a json with the exchange rate info. params are the internal id of the coin (required) and currency to convert to(optional)
// https://api.coinmarketcap.com/v2/ticker/512/?convert=BTC
func getCoinmarketcapData(ex *exchange, internalID int) error {

	url := "https://api.coinmarketcap.com/v2/ticker/" + strconv.Itoa(internalID) + "/?convert=" + ex.destination

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := coinmarketcapTickerAPIResponse{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return err
	}

	if result.Metadata.Error != "" {
		return fmt.Errorf(result.Metadata.Error)
	}

	ex.rate = result.Data.Quotes[ex.destination].Price
	ex.timestamp = result.Data.LastUpdated

	return nil
}

// json response from Coinmarketcap ticker API
type coinmarketcapListingAPIResponse struct {
	Data     []data   `json:"data"`
	Metadata metadata `json:"metadata"`
}

// performs a get request to the of the following form in order to retrieve all currencies available and get the id for our source currency
// https://api.coinmarketcap.com/v2/listings/
func getCoinmarketcapListing(currencySymbol string) (int, error) {

	url := "https://api.coinmarketcap.com/v2/listings/"

	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	result := coinmarketcapListingAPIResponse{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return -1, err
	}

	if result.Metadata.Error != "" {
		return -1, fmt.Errorf(result.Metadata.Error)
	}

	for _, coin := range result.Data {
		if coin.Symbol == currencySymbol {
			return coin.ID, nil
		}
	}

	return -1, fmt.Errorf("Currency " + currencySymbol + " not found")
}
