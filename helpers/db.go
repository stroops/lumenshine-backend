package helpers

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/pressly/goose"
)

var flags = flag.NewFlagSet("goose", flag.ExitOnError)

//MigrateDB migrates the database to the specified version
//migrate must be the only param for the binary
//only goose params may be attached to it
func MigrateDB(db *sql.DB, migrationDir string) error {
	//check if goose called
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] != "migrate") {
		return nil
	}

	flags.Usage = usage
	flags.Parse(os.Args[2:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		os.Exit(0)
	}

	//if migrating down, and not silent mode, we ask user if he realy wants to downgrade
	if (StringInSliceI("down", args) || StringInSliceI("down-to", args)) && !StringInSliceI("silent-down", args) {
		var err error
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you realy want to downgrade? [y/n]: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		if text != "y" && text != "Y" {
			log.Println("Stoping (down)migration")
			os.Exit(1)
		}
	}

	command := args[0]
	if len(args) > 1 && command == "create" {
		if err := goose.Run("create", nil, migrationDir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
			os.Exit(1)
		}
		log.Println("Created migration")
		os.Exit(0)
	}

	if len(args) < 1 {
		flags.Usage()
		os.Exit(1)
	}

	if command == "-h" || command == "--help" {
		flags.Usage()
		os.Exit(0)
	}

	log.Println("Starting migration")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose select-dialect: %v", err)
		os.Exit(1)
	}

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	// check if migration dir exists and if yes delete all files in there. Then unpack migration files to that dir
	err := os.RemoveAll(migrationDir)
	if err != nil {
		err := os.RemoveAll(migrationDir)
		log.Fatalf("error removing dir: %s <%v>", migrationDir, err)
		os.Exit(1)
	}

	err = CreateDirIfNotExists(migrationDir)
	if err != nil {
		log.Fatalf("error creating dir: %s <%v>", migrationDir, err)
		os.Exit(1)
	}

	box := rice.MustFindBox("db-files/migrations_src")
	err = box.Walk("", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sql" {
			fmt.Println("unpack file:", path)

			//write file to path
			fileName := filepath.Join(migrationDir, path)
			data, err := box.Bytes(path)
			if err != nil {
				log.Fatalf("error unboxing file: %s <%v>", path, err)
				os.Exit(1)
			}

			err = ioutil.WriteFile(fileName, data, 0644)
			if err != nil {
				log.Fatalf("error writing file: %s <%v>", path, err)
				os.Exit(1)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walk error [%v]\n", err)
		os.Exit(1)
	}

	if err := goose.Run(command, db, migrationDir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
		os.Exit(1)
	}

	log.Println("DB-Migration done")
	os.Exit(0)

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func usage() {
	log.Print(usageStr)
}

var (
	usageStr = `Usage: migrate COMMAND
Commands:
    up                                Migrate the DB to the most recent version available
    up-to VERSION                     Migrate the DB to a specific VERSION
    down [silent-down]                Roll back the version by 1 (silent-down without confirmation)
    down-to VERSION [silent-down]     Roll back to a specific VERSION (silent-down without confirmation)
    redo                              Re-run the latest migration
    status                            Dump the migration status for the current DB
    version                           Print the current version of the database
    create NAME [sql|go]              Creates new migration file with next version
`
)
