package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/Soneso/lumenshine-backend/admin/config"

	_ "github.com/lib/pq" //needed for SQL access
	"github.com/volatiletech/sqlboiler/boil"
)

//DB is the outer db connection
var DB *sql.DB
var DBC *sql.DB

//CreateNewDB creates a new DB connection
func CreateNewDB(cnf *config.Config) error {
	var err error

	//connect the admin database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.AdminDB.DBHost, cnf.AdminDB.DBPort, cnf.AdminDB.DBUser, cnf.AdminDB.DBPassword, cnf.AdminDB.DBName)

	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to admin-db: %v", err)
	}

	//try to ping the db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping admin-database: %v", err)
	}

	if err = helpers.MigrateDB(DB, cnf.SQLMigrationDir); err != nil {
		return err
	}

	boil.SetDB(DB)

	//connect the customer database
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.CustomerDB.DBHost, cnf.CustomerDB.DBPort, cnf.CustomerDB.DBUser, cnf.CustomerDB.DBPassword, cnf.CustomerDB.DBName)

	DBC, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to connect to customer-db: %v", err)
	}

	err = DBC.Ping()
	if err != nil {
		log.Fatalf("Failed to ping customer-database: %v", err)
	}

	return nil
}
