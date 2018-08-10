package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/volatiletech/sqlboiler/boil"
)

//DB is the outer db connection
var db *sql.DB
var dbCore *sql.DB

//CreateNewDB creates a new DB connection
func CreateNewDB(cnf *Config) error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.DBHost, cnf.DBPort, cnf.DBUser, cnf.DBPassword, cnf.DBName)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	//try to ping the db
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err = helpers.MigrateDB(db, cnf.SQLMigrationDir); err != nil {
		return err
	}

	boil.SetDB(db)

	return nil
}

//CreateNewCoreDB creates a new DB connection
func CreateNewCoreDB(cnf *Config) error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.CoreDBHost, cnf.CoreDBPort, cnf.CoreDBUser, cnf.CoreDBPassword, cnf.CoreDBName)

	dbCore, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to core db: %v", err)
	}

	//try to ping the db
	err = dbCore.Ping()
	if err != nil {
		log.Fatalf("Failed to ping core database: %v", err)
	}

	return nil
}
