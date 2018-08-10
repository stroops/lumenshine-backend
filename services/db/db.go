package main

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"fmt"

	_ "github.com/lib/pq" //needed for SQL access

	"flag"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/volatiletech/sqlboiler/boil"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	db    *sql.DB
)

//createNewDB create a new DB connection
func createNewDB(log *logrus.Entry, cnf *Config) error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.DBHost, cnf.DBPort, cnf.DBUser, cnf.DBPassword, cnf.DBName)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.WithError(err).Fatalf("Failed to connect to db")
	}

	//try to ping the db
	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatalf("Failed to ping db")
	}

	if err = helpers.MigrateDB(db, cnf.SQLMigrationDir); err != nil {
		return err
	}

	boil.SetDB(db)

	return nil
}
