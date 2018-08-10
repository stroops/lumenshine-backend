package main

import (
	"github.com/Soneso/lumenshine-backend/addons/dividend/models"
)

//AddSnapshot adds a snapshot entry in the db
func AddSnapshot(assetCode string, issuer string) (*int, error) {
	snapshot := models.Snapshot{AssetCode: assetCode, Issuer: issuer}

	err := snapshot.Insert(db)
	if err != nil {
		return nil, err
	}

	return &snapshot.ID, nil
}
