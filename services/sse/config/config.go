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

//Stellar data for stellar network
type Stellar struct {
	Horizon           string
	NetworkPassphrase string
	Testnet           bool
}

//Config for the app
type Config struct {
	Port            int
	ApplicationDir  string
	IsDevSystem     bool
	SQLMigrationDir string

	StellarHorizenDB DBConfig

	Stellar Stellar
}

//ReadConfig reads the configuration from the file
func ReadConfig(cmd *cobra.Command) (*Config, error) {
	cnf := new(Config)
	err := helpers.ReadConfig(cmd, "sse-config", "sse-config-local", cnf)
	if err != nil {
		return nil, err
	}

	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	if cnf.IsDevSystem {
		cnf.SQLMigrationDir = filepath.Join(cnf.ApplicationDir, "tmp", "migrations")
	} else {
		if cnf.SQLMigrationDir == "" {
			cnf.SQLMigrationDir = filepath.Join(cnf.ApplicationDir, "db", "migrations")
		} else {
			cnf.SQLMigrationDir = filepath.Join(cnf.ApplicationDir, cnf.SQLMigrationDir)
		}

	}
	helpers.CreateDirIfNotExists(cnf.SQLMigrationDir)
	return cnf, nil
}
