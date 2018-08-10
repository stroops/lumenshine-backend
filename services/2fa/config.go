package main

import (
	"os"
	"path/filepath"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Config for the app
type Config struct {
	Port           int
	ApplicationDir string
	IsDevSystem    bool
	IssuerName     string
	UseDemoSecrete bool
	SecreteKeyLen  int
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "2fa-config", "2fa-config-local", cnf)

	//cnf = cnftmp.(*Config)
	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
