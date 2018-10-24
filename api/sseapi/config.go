package main

import (
	"os"
	"path/filepath"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//ServicesConfig holds the config for all services
type ServicesConfig struct {
	DBSrvPort int64
	DBSrvHost string

	TfaSrvPort int64
	TfaSrvHost string

	JwtSrvPort int64
	JwtSrvHost string

	MailSrvPort int64
	MailSrvHost string

	PaySrvPort int64
	PaySrvHost string
}

//Config for the app
type Config struct {
	Port           int
	ApplicationDir string
	IsDevSystem    bool
	LogLevel       string
	CORSHosts      []string

	Services ServicesConfig
}

var cnf *Config

func readConfig(log *logrus.Entry, cmd *cobra.Command) error {
	cnf = new(Config)
	err := helpers.ReadConfig(cmd, "sseapi-config", "sseapi-config-local", cnf)
	if err != nil {
		log.WithError(err).Fatal("Error reading config")
	}

	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
