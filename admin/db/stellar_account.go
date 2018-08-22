package db

import (
	"database/sql"
	"time"

	"github.com/Soneso/lumenshine-backend/admin/models"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddStellarAccount creates a new account
func AddStellarAccount(account *models.AdminStellarAccount, assetCode string) error {
	err := account.InsertG(boil.Infer())
	if err != nil {
		return err
	}

	if assetCode != "" {
		accountAsset := models.AdminStellarAsset{
			AssetCode:         assetCode,
			IssuerPublicKeyID: account.PublicKey,
			UpdatedBy:         account.UpdatedBy}

		err := accountAsset.InsertG(boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

//ExistsStellarAccount - true if an account with the public key already exists
func ExistsStellarAccount(publicKey string) (bool, error) {
	account, err := models.AdminStellarAccounts(qm.Where(models.AdminStellarAccountColumns.PublicKey+"=?", publicKey)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if account == nil {
		return false, nil
	}

	return true, nil
}

//GetStellarAccount - returns the stellar account
func GetStellarAccount(publicKey string) (*models.AdminStellarAccount, error) {
	account, err := models.AdminStellarAccounts(
		qm.Load("IssuerPublicKeyAdminStellarAssets"),
		qm.Load("StellarAccountPublicKeyAdminStellarSigners"),
		qm.Where(models.AdminStellarAccountColumns.PublicKey+"=?", publicKey)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return account, nil
}

//AllStellarAccounts - returns all stellar accounts
func AllStellarAccounts() (models.AdminStellarAccountSlice, error) {
	accounts, err := models.AdminStellarAccounts().AllG()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

//IssuerAssetCodes - returns the asset codes
func IssuerAssetCodes(issuerPublicKey string) (models.AdminStellarAssetSlice, error) {
	assetcodes, err := models.AdminStellarAssets(qm.Where(models.AdminStellarAssetColumns.IssuerPublicKeyID+"=?", issuerPublicKey)).AllG()
	if err != nil {
		return nil, err
	}

	return assetcodes, nil
}

//UpdateStellarAccount - updates the account in the db
func UpdateStellarAccount(account *models.AdminStellarAccount, updatedBy string) error {
	account.UpdatedBy = updatedBy
	account.UpdatedAt = time.Now()

	_, err := account.UpdateG(boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

//AddAssetCode - adds new asset code to issuer account
func AddAssetCode(account *models.AdminStellarAccount, assetCode string, updatedBy string) error {
	dbAssetCode := models.AdminStellarAsset{
		IssuerPublicKeyID: account.PublicKey,
		AssetCode:         assetCode,
		UpdatedBy:         updatedBy}

	err := dbAssetCode.InsertG(boil.Infer())
	if err != nil {
		return err
	}

	err = UpdateStellarAccount(account, updatedBy)
	if err != nil {
		return err
	}

	return nil
}

//RemoveAssetCode - removes  asset code from issuer account
func RemoveAssetCode(account *models.AdminStellarAccount, assetCode string, updatedBy string) error {
	dbAssetCode, err := models.AdminStellarAssets(
		qm.Where(models.AdminStellarAssetColumns.IssuerPublicKeyID+"=?", account.PublicKey),
		qm.Where(models.AdminStellarAssetColumns.AssetCode+"=?", assetCode)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if dbAssetCode == nil {
		return nil
	}

	_, err = dbAssetCode.DeleteG()
	if err != nil {
		return err
	}

	err = UpdateStellarAccount(account, updatedBy)
	if err != nil {
		return err
	}

	return nil
}

//AddSigner - adds new signer to issuer account
func AddSigner(account *models.AdminStellarAccount, signer *models.AdminStellarSigner, updatedBy string) error {
	signer.UpdatedBy = updatedBy

	err := signer.InsertG(boil.Infer())
	if err != nil {
		return err
	}

	err = UpdateStellarAccount(account, updatedBy)
	if err != nil {
		return err
	}

	return nil
}

//RemoveSigner - removes  signer from issuer account
func RemoveSigner(account *models.AdminStellarAccount, signerPublicKey string, signerType string, updatedBy string) error {
	dbSigner, err := models.AdminStellarSigners(
		qm.Where(models.AdminStellarSignerColumns.StellarAccountPublicKeyID+"=?", account.PublicKey),
		qm.Where(models.AdminStellarSignerColumns.SignerPublicKey+"=?", signerPublicKey),
		qm.Where(models.AdminStellarSignerColumns.Type+"=?", signerType)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if dbSigner == nil {
		return nil
	}

	_, err = dbSigner.DeleteG()
	if err != nil {
		return err
	}

	err = UpdateStellarAccount(account, updatedBy)
	if err != nil {
		return err
	}

	return nil
}

//DeleteStellarAccount - deletes the account
func DeleteStellarAccount(account *models.AdminStellarAccount) error {
	if len(account.R.IssuerPublicKeyAdminStellarAssets) > 0 {
		_, err := account.R.IssuerPublicKeyAdminStellarAssets.DeleteAllG()
		if err != nil {
			return err
		}
	}
	if len(account.R.StellarAccountPublicKeyAdminStellarSigners) > 0 {
		_, err := account.R.StellarAccountPublicKeyAdminStellarSigners.DeleteAllG()
		if err != nil {
			return err
		}
	}
	_, err := account.DeleteG()
	if err != nil {
		return err
	}
	return nil
}

//GetStellarSigner - returns the stellar signer
func GetStellarSigner(publicKey string) (*models.AdminStellarSigner, error) {
	signer, err := models.AdminStellarSigners(
		qm.Where(models.AdminStellarSignerColumns.SignerPublicKey+"=?", publicKey)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return signer, nil
}

//UpdateSigner - updates the signer
func UpdateSigner(signer *models.AdminStellarSigner, updatedBy string) error {
	signer.UpdatedBy = updatedBy
	signer.UpdatedAt = time.Now()

	_, err := signer.UpdateG(boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

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
