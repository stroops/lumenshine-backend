package main

//go:generate sqlboiler --wipe --no-tests --no-context --config $HOME/.config/sqlboiler/sqlboiler_dividend.toml psql

import (
	"github.com/Soneso/lumenshine-backend/addons/dividend/cmd"

	"log"

	_ "github.com/lib/pq"
)

func main() {
	var err error

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.Fatalf("Error reading root command %v", err)
	}

	if err = ReadConfig(cmd); err != nil {
		log.Fatalf("Error reading configuration. %v", err)
	}

	if err = InitRequest(cmd); err != nil {
		log.Fatalf("Error creating request. %v", err)
	}

	if err = CreateNewDB(Cnf); err != nil {
		log.Fatalf("Error creating db connection. %v", err)
	}

	if err = CreateNewCoreDB(Cnf); err != nil {
		log.Fatalf("Error creating core db connection. %v", err)
	}

	trustlines, err := GetTrustlines(Req.AssetCode, Req.Issuer)
	if err != nil {
		log.Fatalf("Error getting trustlines. %v", err)
	}
	trustlines = RemoveBlacklisted(trustlines, Req.Blacklist)
	log.Printf("\n Trustlines count: %d", len(trustlines))

	snapshotID, err := AddSnapshot(Req.AssetCode, Req.Issuer)
	if err != nil {
		log.Fatalf("\n Error adding snapshot entry. %v", err)
	}
	log.Printf("\n Snapshot id: %d", *snapshotID)

	err = InsertDividends(*snapshotID, trustlines)
	if err != nil {
		log.Fatalf("\n Error inserting dividends. %v", err)
	}
}
