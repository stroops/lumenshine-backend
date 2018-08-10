package config

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DBConfig imported from config file
type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int64
	DBName     string
	DBSchema   string
}

// TransactionData holds data relevant for a transaction between two currencies
type TransactionData struct {
	Model                     string
	FileURL                   string
	SourceCurrency            string
	DestinationCurrency       string
	SourceCurrencyName        string
	DestinationCurrencyName   string
	SourceCurrencyIssuer      string
	DestinationCurrencyIssuer string
}

// TickerConfig imported from config file
type TickerConfig struct {
	TickSeconds                 int
	SourceCurrency              string
	DestinationCurrency         []string
	DecentralizedExchangeCode   []string
	DecentralizedExchangeDomain []string
	DecentralizedExchangeIssuer []string
}

// CleanupConfig imported from config file
type CleanupConfig struct {
	CleanupMinutesInterval  int
	HoursToKeepMinutelyData int
	HoursToKeepHourlyData   int
}

// Config holds the data imported from the config file and command line
type Config struct {
	LocalDownloadDir     string
	HorizonURL           string
	TruncateHistoryTable bool // this is not set from the toml file but from a command line arg

	CORSHosts []string

	DB                 DBConfig
	ImportTransactions []TransactionData
	Ticker             TickerConfig
	Cleanup            CleanupConfig

	Port            string
	SQLMigrationDir string
	ApplicationDir  string
	IsDevSystem     bool
}

// DataColumns struct
type DataColumns struct {
	DateCol         int
	ExchangeRateCol int
	SkipLines       int
}

var tx = []TransactionData{
	{
		Model:                     "coinmetrics",
		FileURL:                   "https://coinmetrics.io/data/xlm.csv",
		SourceCurrency:            "XLM",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "Lumen",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
	{
		Model:                     "coinmetrics",
		FileURL:                   "https://coinmetrics.io/data/btc.csv",
		SourceCurrency:            "BTC",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "Bitcoin",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
	{
		Model:                     "coinmetrics",
		FileURL:                   "https://coinmetrics.io/data/eth.csv",
		SourceCurrency:            "ETH",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "Ether",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
	{
		Model:                     "sauder",
		FileURL:                   "http://fx.sauder.ubc.ca/cgi/fxdata",
		SourceCurrency:            "EUR",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "Euro",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
	{
		Model:                     "sauder",
		FileURL:                   "http://fx.sauder.ubc.ca/cgi/fxdata",
		SourceCurrency:            "CNY",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "Chinese Yuan Renminbi",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
	{
		Model:                     "sauder",
		FileURL:                   "http://fx.sauder.ubc.ca/cgi/fxdata",
		SourceCurrency:            "KRW",
		DestinationCurrency:       "USD",
		SourceCurrencyName:        "South Korean Won",
		DestinationCurrencyName:   "United States Dollar",
		SourceCurrencyIssuer:      "",
		DestinationCurrencyIssuer: "",
	},
}

// DateFormat is the format of the date used in the app
const DateFormat = "2006-01-02"

// TimeFormat is the format of the time used in the app
const TimeFormat = "2006-01-02 15:04:05"

// ExternalCurrencyIssuer is the placeholder in the DB for currencies external to stellar which have to no issuer
const ExternalCurrencyIssuer = ""

// Cnf holds configs
var Cnf *Config

//CurrencyCodeToName holds relations from currency code to currency name
var CurrencyCodeToName = make(map[string]string)

// ModelCols are the columns in the csv where the respective data is found
var ModelCols = map[string]DataColumns{
	"coinmetrics": {
		DateCol:         0,
		ExchangeRateCol: 4,
		SkipLines:       1,
	},
	"sauder": {
		DateCol:         1,
		ExchangeRateCol: 3,
		SkipLines:       2,
	},
}

// ReadConfig parses and reads config file
func ReadConfig(cmd *cobra.Command) error {
	Cnf = new(Config)
	helpers.ReadConfig(cmd, "charts-config", "charts-config-local", Cnf)

	//this should be set by reading the config toml
	Cnf.ImportTransactions = tx

	// set currency code to currency name relation
	for _, v := range Cnf.ImportTransactions {
		CurrencyCodeToName[v.SourceCurrency] = v.SourceCurrencyName
		CurrencyCodeToName[v.DestinationCurrency] = v.DestinationCurrencyName
	}

	if viper.GetString("ApplicationDirectory") != "" {
		Cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		Cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	if Cnf.IsDevSystem {
		Cnf.SQLMigrationDir = filepath.Join(Cnf.ApplicationDir, "tmp", "migrations")
	} else {
		Cnf.SQLMigrationDir = filepath.Join(Cnf.ApplicationDir, "db", "migrations")
	}
	helpers.CreateDirIfNotExists(Cnf.SQLMigrationDir)

	return nil
}
