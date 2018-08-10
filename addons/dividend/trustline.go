package main

import (
	"database/sql"
	"github.com/Soneso/lumenshine-backend/addons/dividend/modelscore"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

// GetTrustlines reads from the db core.
func GetTrustlines(assetCode string, issuer string) (modelscore.TrustlineSlice, error) {
	var err error

	q := []qm.QueryMod{}
	q = append(q, qm.Select(modelscore.TrustlineColumns.Accountid,
		modelscore.TrustlineColumns.Issuer,
		modelscore.TrustlineColumns.Assetcode,
		modelscore.TrustlineColumns.Tlimit,
		modelscore.TrustlineColumns.Balance))
	q = append(q, qm.Where(modelscore.TrustlineColumns.Assetcode+"=?", assetCode))
	q = append(q, qm.Where(modelscore.TrustlineColumns.Issuer+"=?", issuer))

	p, err := modelscore.Trustlines(dbCore, q...).All()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return p, nil
}

//RemoveBlacklisted removes blacklisted accounts
func RemoveBlacklisted(trustlines modelscore.TrustlineSlice, blacklist []string) modelscore.TrustlineSlice {
	for _, v := range blacklist {
		trustlines = remove(trustlines, v)
	}
	return trustlines
}

func remove(trustlines modelscore.TrustlineSlice, accountid string) modelscore.TrustlineSlice {
	for i, v := range trustlines {
		if v.Accountid == accountid {
			return append(trustlines[:i], trustlines[i+1:]...)
		}
	}
	return trustlines
}
