package main

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

//Config for the app
type Config struct {
	CustomerDB DBConfig
	HorizonDB  DBConfig

	Port            int
	SQLMigrationDir string
	ApplicationDir  string
	IsDevSystem     bool
}

var cnf *Config

func readConfig(cmd *cobra.Command) error {
	cnf = new(Config)
	helpers.ReadConfig(cmd, "db-config", "db-config-local", cnf)

	//cnf = cnftmp.(*Config)
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

	return nil
}
