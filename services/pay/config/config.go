package config

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"math/big"
	"os"
	"path/filepath"

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
	Enabled            bool
	MasterPublicKey    string
	RPCServer          string
	RPCUser            string
	RPCPass            string
	Testnet            bool
	MinimumValueBTC    *big.Int
	MinimumValueBTCStr string
	TokenPrice         float64
}

//Ethereum data for ethereum client
type Ethereum struct {
	Enabled               bool
	MasterPublicKey       string
	RPCServer             string
	NetworkID             string
	MinimumWeiValueEthStr string
	MinimumWeiValueEth    *big.Int
	TokenPrice            float64
	Testnet               bool
}

//Stellar data for stellar network
type Stellar struct {
	Enabled               bool
	IssuerPublicKey       string
	DistributionPublicKey string
	DistributionSeed      string
	TokenAssetCode        string
	Horizon               string
	NetworkPassphrase     string
	StartingBalance       string
	Testnet               bool
	TokenPrice            float64
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
	Port           int
	ApplicationDir string
	IsDevSystem    bool

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
