package config

import (
	"os"
	"path/filepath"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//DBConfig - database configuration
type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int64
	DBName     string
}

//ServicesConfig holds the config for all services
type ServicesConfig struct {
	MailSrvPort int64
	MailSrvHost string
}

//StellarNetworkConfig - stellar network config
type StellarNetworkConfig struct {
	HorizonURL string
}

//PromoConfig - promo configuration
type PromoConfig struct {
	ImagesPath string
}

//SiteConfig globals that define the site behaviour
type SiteConfig struct {
	SiteName    string
	EmailSender string
}

//WebLinksConfig links that are used in the clients (e.g. token confirm)
type WebLinksConfig struct {
	LostTFA   string
	ImagesUrl string
}

//KycConfig - kyc configuration
type KycConfig struct {
	DocumentsPath string
}

//Config for the app
type Config struct {
	AppName   string
	Port      int64
	GRPCPort  int64
	CORSHosts []string

	AdminDB       DBConfig
	CustomerDB    DBConfig
	StellarCoreDB DBConfig

	Services       ServicesConfig
	StellarNetwork StellarNetworkConfig
	Site           SiteConfig
	WebLinks       WebLinksConfig

	SQLMigrationDir string
	ApplicationDir  string
	IsDevSystem     bool
	ConnectToCoreDB bool

	Kyc   KycConfig
	Promo PromoConfig
}

//Cnf holds the application configuration
var Cnf *Config

//ReadConfig reads configuration from command
func ReadConfig(cmd *cobra.Command) error {
	Cnf = new(Config)

	if err := helpers.ReadConfig(cmd, "admin-config", "admin-config-local", Cnf); err != nil {
		return err
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
