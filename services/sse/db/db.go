package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Soneso/lumenshine-backend/services/sse/config"
	_ "github.com/lib/pq" //needed for SQL access

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Soneso/lumenshine-backend/helpers"
)

const (
	stellarLastLedgerKey = "xlm_last_ledger_id"
)

//DB general DB struct
type DB struct {
	*sql.DB
}

//CreateNewHorizonDB creates a new horizon DB connection
func CreateNewHorizonDB(cnf *config.Config) (*DB, error) {
	var err error

	//connect the customer database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.StellarHorizenDB.DBHost, cnf.StellarHorizenDB.DBPort, cnf.StellarHorizenDB.DBUser, cnf.StellarHorizenDB.DBPassword, cnf.StellarHorizenDB.DBName)

	DBH, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to horizon-db: %v", err)
	}

	err = DBH.Ping()
	if err != nil {
		log.Fatalf("Failed to ping horizon-database: %v", err)
	}

	if err = helpers.MigrateDB(DBH, cnf.SQLMigrationDir); err != nil {
		return nil, err
	}
	boil.SetDB(DBH)
	return &DB{DBH}, nil
}
