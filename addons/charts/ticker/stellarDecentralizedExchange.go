package ticker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

type assetStruct struct {
	assetType   string
	assetCode   string
	assetIssuer string
}

// ticker for stellar decentralized exchange - gets rates of coins within stellar wrt XLM
func stellarDecentralizedExchangeTicker() {
	assets := make(map[string]assetStruct)

	// get assets information from horizon api
	for i, source := range config.Cnf.Ticker.DecentralizedExchangeIssuer {
		sourceAsset, err := getAssetInformation(source, config.Cnf.Ticker.DecentralizedExchangeCode[i])
		if err != nil {
			log.Printf("Error getting Decentralized Exchange Asset for source issuer %s %v", source, err)
			return
		}
		assets[source] = sourceAsset
	}

	destIssuer := "native"
	destCode := "XLM"

	destAsset, err := getAssetInformation(destIssuer, destCode)
	if err != nil {
		log.Printf("Error getting Decentralized Exchange Asset for native issuer %s %v", destIssuer, err)
		return
	}
	assets[destIssuer] = destAsset

	ex := exchange{}
	for {

		for _, sourceIssuer := range config.Cnf.Ticker.DecentralizedExchangeIssuer {

			ex.source = assets[sourceIssuer].assetCode
			ex.sourceIssuer = sourceIssuer
			ex.destination = assets[destIssuer].assetCode
			ex.destinationIssuer = "" // should be destIssuer which is native but we store in db with no issuer
			err := getDecentralizedExchangeData(&ex, sourceIssuer, destIssuer, assets)
			if err != nil {
				log.Printf("Error getting ticker data for Stellar Decentralized Exchange %v", err)
				continue
			}

			err = handleDecentralizedExchangeData(ex)
			if err != nil {
				log.Printf("Error on handling Stellar Decentralized Exchange Data %v", err)
			}

		}

		waitTime := time.Second * time.Duration(config.Cnf.Ticker.TickSeconds)
		time.Sleep(waitTime)

	}
}

type horizonOrderBookOffer struct {
	PriceR struct {
		N int `json:"n"`
		D int `json:"d"`
	} `json:"price_r"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type horizonOrderBookAPIResponse struct {
	Bids []horizonOrderBookOffer `json:"bids"`
	Asks []horizonOrderBookOffer `json:"asks"`
}

func handleDecentralizedExchangeData(ex exchange) error {

	// insert stellar asset to XLM rate
	err := handleData(ex)
	if err != nil {
		return err
	}

	// get latest rates from XLM to external currencies
	baseCurrency, err := utils.GetCurrencyByCode("XLM", config.ExternalCurrencyIssuer, false)
	if err != nil {
		return err
	}

	schema := utils.GetSchemaForQuery()
	externalCurrencies, err := models.Currencies(
		qm.Select(schema+"currency.*"),
		qm.InnerJoin(schema+"history_chart_data h on h.destination_currency_id = "+schema+"currency.id"),
		qm.Where("h.source_currency_id =? AND "+schema+"currency.currency_issuer =?", baseCurrency.ID, config.ExternalCurrencyIssuer),
		qm.GroupBy(schema+"currency.id"),
	).All(utils.DB)
	if err != nil {
		return err
	}

	for _, externalCurrency := range externalCurrencies {

		currentRate, err := utils.GetCurrentRate(baseCurrency, externalCurrency)
		if err != nil {
			// return err
		}

		if currentRate == nil {
			log.Printf("No current exchange rate for %s and %s", baseCurrency.CurrencyCode, externalCurrency.CurrencyCode)
			break
		}

		// compute the rate from internal to external currency wrt XLM and insert to DB
		exToExternal := exchange{}
		exToExternal.source = ex.source
		exToExternal.sourceIssuer = ex.sourceIssuer
		exToExternal.destination = externalCurrency.CurrencyCode
		exToExternal.destinationIssuer = externalCurrency.CurrencyIssuer
		exToExternal.rate = currentRate.ExchangeRate * ex.rate

		err = handleData(exToExternal)
		if err != nil {
			return err
		}

	}

	return nil
}

// performs a get request to the of the following form in order to retrieve a json with the exchange rate info
// https://horizon.stellar.org/order_book?selling_asset_type=credit_alphanum4&selling_asset_code=MOBI&selling_asset_issuer=GDCIUCGL7VEMMF6VYJOW75KQ5ZCLHAQBRM6EPFTKCRWUYVUOOYQCKC5A&buying_asset_type=native
func getDecentralizedExchangeData(ex *exchange, sourceIssuer string, destIssuer string, assets map[string]assetStruct) error {

	// we assume that we always convert from something non native to native(XLM)
	var url string
	url = config.Cnf.HorizonURL + "order_book?selling_asset_type=" + assets[sourceIssuer].assetType + "&selling_asset_code=" + assets[sourceIssuer].assetCode + "&selling_asset_issuer=" + sourceIssuer + "&buying_asset_type=native"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := horizonOrderBookAPIResponse{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return err
	}

	var lowestPrice float64 = -1
	for _, ask := range result.Asks {
		price := float64(ask.PriceR.N) / float64(ask.PriceR.D)
		if price < lowestPrice || lowestPrice == -1 {
			lowestPrice = price
		}
	}
	ex.rate = lowestPrice
	ex.timestamp = 0

	return nil
}

type horizonAssetAPIResponse struct {
	Embedded struct {
		Records []struct {
			AssetType   string `json:"asset_type"`
			AssetCode   string `json:"asset_code"`
			AssetIssuer string `json:"asset_issuer"`
			PagingToken string `json:"paging_token"`
			Amount      string `json:"amount"`
			NumAccounts int    `json:"num_accounts"`
		} `json:"records"`
	} `json:"_embedded"`
}

// performs a get request to the of the following form in order to retrieve all assets information
// https://horizon.stellar.org/assets?asset_issuer=GDCIUCGL7VEMMF6VYJOW75KQ5ZCLHAQBRM6EPFTKCRWUYVUOOYQCKC5A&asset_code=MOBI
func getAssetInformation(aIssuer string, aCode string) (assetStruct, error) {

	asset := assetStruct{
		assetIssuer: aIssuer,
	}

	if aIssuer == "native" {
		asset.assetCode = "XLM"
		asset.assetType = "native"
	} else {
		url := config.Cnf.HorizonURL + "assets?asset_issuer=" + aIssuer + "&asset_code=" + aCode

		resp, err := http.Get(url)
		if err != nil {
			return asset, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return asset, err
		}

		result := horizonAssetAPIResponse{}
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			return asset, err
		}

		asset.assetCode = result.Embedded.Records[0].AssetCode
		asset.assetType = result.Embedded.Records[0].AssetType
	}

	return asset, nil
}
