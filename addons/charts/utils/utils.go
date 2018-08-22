package utils

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Global vars
var (
	Flags = flag.NewFlagSet("goose", flag.ExitOnError)
	DB    *sql.DB
)

// CreateNewDB creates a new db conn
func CreateNewDB() error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Cnf.DB.DBHost, config.Cnf.DB.DBPort, config.Cnf.DB.DBUser, config.Cnf.DB.DBPassword, config.Cnf.DB.DBName)

	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	//try to ping the db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err = helpers.MigrateDB(DB, config.Cnf.SQLMigrationDir); err != nil {
		return err
	}

	boil.SetDB(DB)

	return nil
}

//HandlePanic general panic handler
func HandlePanic(function string, message string) {
	if r := recover(); r != nil {
		log.Printf("Recovered in <%s>. Message <%s>, Rec:<%v>", function, message, r)
	}
}

// GetSchemaForQuery returns string to use as schema inside queries
func GetSchemaForQuery() string {
	if config.Cnf.DB.DBSchema != "" {
		return config.Cnf.DB.DBSchema + "."
	}
	return ""
}

// TruncateHistoryTable deletes the history data from the db
func TruncateHistoryTable() error {

	schema := GetSchemaForQuery()
	table := "history_chart_data"

	_, err := DB.Exec(`TRUNCATE TABLE ` + schema + table + ` RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}

	return nil
}

// InsertCurrencyToDB adds new currency to db
func InsertCurrencyToDB(currencyCode string, currencyIssuer string) error {

	var currency models.Currency
	currency.CurrencyCode = currencyCode
	if currencyIssuer == config.ExternalCurrencyIssuer {
		currency.CurrencyName = config.CurrencyCodeToName[currencyCode]
	}
	currency.CurrencyIssuer = currencyIssuer
	err := currency.Insert(DB, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

// GetCurrencyByCode returns currency model if exists, if not and the flag is set it adds the currency and returns it
func GetCurrencyByCode(code string, issuer string, insertIfNotFound bool) (currency *models.Currency, err error) {

	currencyCode := strings.ToUpper(code)
	currencyIssuer := strings.ToUpper(issuer)

	exists, err := models.Currencies(qm.Where("currency_code=? AND currency_issuer=?", currencyCode, currencyIssuer)).Exists(DB)
	if !exists {
		if !insertIfNotFound {
			err = errors.New("Currency not found")
			return
		}
		err = InsertCurrencyToDB(currencyCode, currencyIssuer)
		if err != nil {
			return
		}

	}
	currency, _ = models.Currencies(qm.Where("currency_code=? AND currency_issuer=?", currencyCode, currencyIssuer)).One(DB)

	return
}

// GetCurrentRate returns latest exchange available between 2 currencies
func GetCurrentRate(sourceCurrency *models.Currency, destinationCurrency *models.Currency) (lastTransaction *models.CurrentChartDataMinutely, err error) {

	lastTransaction, err = sourceCurrency.SourceCurrencyCurrentChartDataMinutelies(
		qm.Where("destination_currency_id =?", destinationCurrency.ID),
		qm.OrderBy("exchange_rate_time DESC"),
	).One(DB)

	// if err != nil {
	// 	log.Panicf("Error getting current rate %v", err)
	// }

	return

}
