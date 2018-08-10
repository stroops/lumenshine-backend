package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Config for the app
type Config struct {
	Port           int
	ApplicationDir string

	DBSrvPort int64
	DBSrvHost string
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "notification-config", "notification-config-local", cnf)

	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
