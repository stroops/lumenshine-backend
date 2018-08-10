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
	IsDevSystem    bool

	DBSrvPort            int64
	DBSrvHost            string
	MemcachedServerURL   string
	ValidMinutes1        int
	ValidMinutes2        int
	LengthJWTKey         int
	SleepTimeLoopSeconds int
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "jwt-config", "jwt-config-local", cnf)

	//cnf = cnftmp.(*Config)
	if viper.GetString("ApplicationDirectory") != "" {
		cnf.ApplicationDir = viper.GetString("ApplicationDirectory")
	} else {
		cnf.ApplicationDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return nil
}
