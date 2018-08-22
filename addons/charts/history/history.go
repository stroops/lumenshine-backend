package history

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Exchange holds the exchange rates from source to destination
type Exchange struct {
	Source            string
	SourceIssuer      string
	Destination       string
	DestinationIssuer string
	History           []ExchangeHistory
}

// ExchangeHistory holds the history of exchange rates
type ExchangeHistory struct {
	Date         string
	ExchangeRate string
}

// GetHistoricalData gathers history data
func GetHistoricalData() error {
	var err error

	// empty history table
	if config.Cnf.TruncateHistoryTable {
		err = utils.TruncateHistoryTable()
		if err != nil {
			return err
		}
	}

	for _, data := range config.Cnf.ImportTransactions {

		var localPath string

		switch data.Model {
		case "coinmetrics":
			localPath, err = downloadFile(data.FileURL)
		case "sauder":
			localPath, err = postData(data)
		default:
			err = errors.New("data model not recognized")
			return err

		}

		if err != nil {
			return err
		}

		exchange := Exchange{}
		exchange.Source = data.SourceCurrency
		exchange.SourceIssuer = data.SourceCurrencyIssuer
		exchange.Destination = data.DestinationCurrency
		exchange.DestinationIssuer = data.DestinationCurrencyIssuer

		err = parseFile(localPath, &exchange.History, data)
		if err != nil {
			return err
		}

		err = sortAndFillGaps(&exchange.History)
		if err != nil {
			return err
		}

		// err = os.Remove(localPath)
		// if err != nil {
		// 	return err
		// }

		// handles the data and inserts to db
		err = handleData(exchange)
		if err != nil {
			return err
		}

	}

	return nil
}

// parseFile will parse the csv file and store it into a predefined structure
func parseFile(path string, history *[]ExchangeHistory, data config.TransactionData) error {

	csvFile, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for i := 0; ; i++ {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}

		if i < config.ModelCols[data.Model].SkipLines {
			continue
		}

		if len(line) < config.ModelCols[data.Model].DateCol || len(line) < config.ModelCols[data.Model].ExchangeRateCol || line[config.ModelCols[data.Model].DateCol] == "" || line[config.ModelCols[data.Model].ExchangeRateCol] == "" {
			log.Printf("Found empty relevant columns at line %d %+v", i, data)
			break
		}

		*history = append(*history, ExchangeHistory{
			Date:         strings.Replace(line[config.ModelCols[data.Model].DateCol], "/", "-", -1), //in sauder model the date is of form 2014/08/08 we need it 2014-08-08
			ExchangeRate: line[config.ModelCols[data.Model].ExchangeRateCol],
		})

	}

	return nil
}

// By is the type of a "less" function that defines the ordering of its ExchangeHistory arguments.
type By func(h1, h2 *ExchangeHistory) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(history []ExchangeHistory) {
	hs := &historySorter{
		history: history,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(hs)
}

// historySorter joins a By function and a slice of ExchangeHistory to be sorted.
type historySorter struct {
	history []ExchangeHistory
	by      func(h1, g2 *ExchangeHistory) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *historySorter) Len() int {
	return len(s.history)
}

// Swap is part of sort.Interface.
func (s *historySorter) Swap(i, j int) {
	s.history[i], s.history[j] = s.history[j], s.history[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *historySorter) Less(i, j int) bool {
	return s.by(&s.history[i], &s.history[j])
}

// sorts the struct by date asc. and fills gaps with latest value(this is used to fill weekends with firday values)
func sortAndFillGaps(history *[]ExchangeHistory) error {

	// Closures that order the ExchangeHistory structure.
	date := func(h1, h2 *ExchangeHistory) bool {
		return h1.Date < h2.Date
	}

	// Sort the history by date asc
	By(date).Sort(*history)

	for i := 0; i < len(*history)-1; {

		date, _ := time.Parse("2006-01-02", (*history)[i].Date)
		nextDate, _ := time.Parse("2006-01-02", (*history)[i+1].Date)

		if !(date.AddDate(0, 0, 1).Equal(nextDate)) {

			newEl := ExchangeHistory{
				Date:         date.AddDate(0, 0, 1).Format(config.DateFormat),
				ExchangeRate: (*history)[i].ExchangeRate,
			}

			aux := append([]ExchangeHistory(nil), (*history)[:i+1]...)
			aux = append(aux, newEl)
			*history = append(aux, (*history)[i+1:]...)

		}

		i++

	}

	return nil
}

// handleData gets the data from the structure and inserts it into db
func handleData(exchange Exchange) error {

	sourceCurrency, err := utils.GetCurrencyByCode(exchange.Source, exchange.SourceIssuer, true)
	if err != nil {
		return err
	}

	destinationCurrency, err := utils.GetCurrencyByCode(exchange.Destination, exchange.DestinationIssuer, true)
	if err != nil {
		return err
	}

	if exchange.Source == "XLM" {

		for _, v := range exchange.History {
			date, _ := time.Parse(config.DateFormat, v.Date)
			rate, _ := strconv.ParseFloat(v.ExchangeRate, 64)

			var history models.HistoryChartDatum
			history.ExchangeRateDate = date
			history.ExchangeRate = rate
			history.SourceCurrencyID = sourceCurrency.ID
			history.DestinationCurrencyID = destinationCurrency.ID

			err = history.Insert(utils.DB, boil.Infer())
			if err != nil {
				return err
			}
		}

	} else {

		// you have to compute a price from the new source w.r.t. the xlm to usd price from that day

		xlm, _ := models.Currencies(qm.Where("currency_code=?", "XLM")).One(utils.DB)
		usd, _ := models.Currencies(qm.Where("currency_code=?", "USD")).One(utils.DB)

		// get the earliest date you have in the db with a transaction from xlm to usd. for other currencies the history may go earlier and we would not have a reference point
		earliestXlmToUsd, _ := models.HistoryChartData(qm.Where("source_currency_id=?", xlm.ID), qm.And("destination_currency_id=?", usd.ID), qm.OrderBy("exchange_rate_date")).One(utils.DB)
		earliestXlmToUsdExchange := earliestXlmToUsd.ExchangeRateDate.Format(config.DateFormat)

		// store exchange rates history from xlm to usd
		xlmToUsd, _ := models.HistoryChartData(qm.Where("source_currency_id=?", xlm.ID), qm.And("destination_currency_id=?", usd.ID)).All(utils.DB)
		var xlmToUsdExchangeHistory = make(map[string]float64)
		for _, v := range xlmToUsd {
			xlmToUsdExchangeHistory[v.ExchangeRateDate.Format(config.DateFormat)] = v.ExchangeRate
		}

		for _, v := range exchange.History {

			if v.Date < earliestXlmToUsdExchange {
				continue
			}

			date, _ := time.Parse(config.DateFormat, v.Date)
			rate, _ := strconv.ParseFloat(v.ExchangeRate, 64)

			xlmToSourceRate := (1 / rate) * xlmToUsdExchangeHistory[date.Format(config.DateFormat)]

			var history models.HistoryChartDatum
			history.ExchangeRateDate = date
			history.ExchangeRate = xlmToSourceRate
			history.SourceCurrencyID = xlm.ID
			history.DestinationCurrencyID = sourceCurrency.ID

			err = history.Insert(utils.DB, boil.Infer())
			if err != nil {
				return err
			}
		}

	}

	return nil
}
