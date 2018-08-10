package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//IOSConfig - ios configuration
type IOSConfig struct {
	DevelopmentCertificate         string
	DevelopmentCertificatePassword string
	ProductionCertificate          string
	ProductionCertificatePassword  string
	BundleID                       string
}

//AndroidConfig - android configuration
type AndroidConfig struct {
	APIKey string
}

//ServicesConfig holds the config for all services
type ServicesConfig struct {
	DBSrvPort int64
	DBSrvHost string

	MailSrvPort int64
	MailSrvHost string
}

//Config for the app
type Config struct {
	//Port           int
	ApplicationDir string
	IsDevSystem    bool

	LimitCount  int
	IdleSeconds int
	EmailSender string

	Services ServicesConfig

	IOSConfig     IOSConfig
	AndroidConfig AndroidConfig
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "messaging-config", "messaging-config-local", cnf)

	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
