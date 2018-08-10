package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MailConfig specifies all the parameters needed for mailing
type MailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Sender   string
	Insecure bool
	UseSSL   bool
}

//Config for the app
type Config struct {
	Port           int
	ApplicationDir string
	IsDevSystem    bool

	DBSrvPort int64
	DBSrvHost string

	MailConfig MailConfig
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "mail-config", "mail-config-local", cnf)

	//cnf = cnftmp.(*Config)
	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
