package config

import (
	"os"
	"path/filepath"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//DBConfig is the definition for a DB connection
type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int64
	DBName     string
}

//Service definition for a service confifuration
type Service struct {
	ServicePort int64
	ServiceHost string
}

//Bitcoin data for bitcoin client
type Bitcoin struct {
	Enabled   bool
	RPCServer string
	RPCUser   string
	RPCPass   string
	Testnet   bool
}

//Ethereum data for ethereum client
type Ethereum struct {
	Enabled   bool
	RPCServer string
	NetworkID string
	Testnet   bool
}

//Stellar data for stellar network
type Stellar struct {
	Enabled           bool
	Horizon           string
	NetworkPassphrase string
	Testnet           bool
}

//Fiat is the definition for the fiat payment
type Fiat struct {
	IBAN            string
	BIC             string
	TokenPrice      float64
	DestiantionName string
	PaymentUsage    string
}

//Config for the app
type Config struct {
	Port                  int
	ApplicationDir        string
	IsDevSystem           bool
	AllowFakeTransactions bool

	DBService        Service
	CustomerDB       DBConfig
	Bitcoin          Bitcoin
	Ethereum         Ethereum
	Stellar          Stellar
	Fiat             Fiat
	StellarHorizenDB DBConfig
}

//ReadConfig reads the configuration from the file
func ReadConfig(cmd *cobra.Command) (*Config, error) {
	cnf := new(Config)
	err := helpers.ReadConfig(cmd, "pay-config", "pay-config-local", cnf)
	if err != nil {
		return nil, err
	}

	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	return cnf, nil
}
