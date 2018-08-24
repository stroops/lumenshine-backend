package db

import (
	"database/sql"

	"github.com/Soneso/lumenshine-backend/admin/models"
	coremodels "github.com/Soneso/lumenshine-backend/db/stellarcore/models"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddTrustline adds a new trustline
func AddTrustline(trustline *models.AdminUnauthorizedTrustline, updatedBy string) error {
	trustline.UpdatedBy = updatedBy

	err := trustline.InsertG(boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

//ExistsUnauthorizedTrustline - true if trustline already exists
func ExistsUnauthorizedTrustline(trustorPublicKey string, issuingPublicKey string, assetCode string) (bool, error) {
	trustline, err := models.AdminUnauthorizedTrustlines(
		qm.Where(models.AdminUnauthorizedTrustlineColumns.TrustorPublicKey+"=?", trustorPublicKey),
		qm.Where(models.AdminUnauthorizedTrustlineColumns.IssuerPublicKeyID+"=?", issuingPublicKey),
		qm.Where(models.AdminUnauthorizedTrustlineColumns.AssetCode+"=?", assetCode)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if trustline == nil {
		return false, nil
	}

	return true, nil
}

//DeleteUnauthorizedTrustline - removes the trustline
func DeleteUnauthorizedTrustline(trustorPublicKey string, issuingPublicKey string, assetCode string) error {
	trustline, err := models.AdminUnauthorizedTrustlines(
		qm.Where(models.AdminUnauthorizedTrustlineColumns.TrustorPublicKey+"=?", trustorPublicKey),
		qm.Where(models.AdminUnauthorizedTrustlineColumns.IssuerPublicKeyID+"=?", issuingPublicKey),
		qm.Where(models.AdminUnauthorizedTrustlineColumns.AssetCode+"=?", assetCode)).OneG()

	if err == sql.ErrNoRows {
		return nil
	}

	if trustline == nil {
		return nil
	}

	_, err = trustline.DeleteG()
	if err != nil {
		return err
	}

	return nil
}

//GetCoreTrustlines - gets the trustlines
func GetCoreTrustlines(accountPublicKey string) (coremodels.TrustlineSlice, error) {
	q := []qm.QueryMod{
		qm.Select(
			coremodels.TrustlineColumns.Accountid,
			coremodels.TrustlineColumns.Issuer,
			coremodels.TrustlineColumns.Assetcode,
			coremodels.TrustlineColumns.Flags,
		),
	}

	q = append(q, qm.Where(coremodels.TrustlineColumns.Accountid+"=?", accountPublicKey))

	trustlines, err := coremodels.Trustlines(q...).All(DBSC)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return trustlines, nil
}

//GetUnauthorizedTrustlines - gets the trustlines
func GetUnauthorizedTrustlines(trustorPublicKey string) (models.AdminUnauthorizedTrustlineSlice, error) {
	trustlines, err := models.AdminUnauthorizedTrustlines(
		qm.Where(models.AdminUnauthorizedTrustlineColumns.TrustorPublicKey+"=?", trustorPublicKey)).AllG()

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return trustlines, nil
}
