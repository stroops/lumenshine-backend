package db

import (
	"database/sql"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"time"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddStellarAccount creates a new account
func AddStellarAccount(account *models.AdminStellarAccount, assetCode string) error {
	err := account.InsertG()
	if err != nil {
		return err
	}

	if assetCode != "" {
		accountAsset := models.AdminStellarAsset{
			AssetCode:         assetCode,
			IssuerPublicKeyID: account.PublicKey,
			UpdatedBy:         account.UpdatedBy}

		err := accountAsset.InsertG()
		if err != nil {
			return err
		}
	}

	return nil
}

//ExistsStellarAccount - true if an account with the public key already exists
func ExistsStellarAccount(publicKey string) (bool, error) {
	account, err := models.AdminStellarAccountsG(qm.Where(models.AdminStellarAccountColumns.PublicKey+"=?", publicKey)).One()

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
	account, err := models.AdminStellarAccountsG(
		qm.Load("IssuerPublicKeyAdminStellarAssets"),
		qm.Load("StellarAccountPublicKeyAdminStellarSigners"),
		qm.Where(models.AdminStellarAccountColumns.PublicKey+"=?", publicKey)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return account, nil
}

//AllStellarAccounts - returns all stellar accounts
func AllStellarAccounts() (models.AdminStellarAccountSlice, error) {
	accounts, err := models.AdminStellarAccountsG().All()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

//IssuerAssetCodes - returns the asset codes
func IssuerAssetCodes(issuerPublicKey string) (models.AdminStellarAssetSlice, error) {
	assetcodes, err := models.AdminStellarAssetsG(qm.Where(models.AdminStellarAssetColumns.IssuerPublicKeyID+"=?", issuerPublicKey)).All()
	if err != nil {
		return nil, err
	}

	return assetcodes, nil
}

//UpdateStellarAccount - updates the account in the db
func UpdateStellarAccount(account *models.AdminStellarAccount, updatedBy string) error {
	account.UpdatedBy = updatedBy
	account.UpdatedAt = time.Now()

	err := account.UpdateG()
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

	err := dbAssetCode.InsertG()
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
	dbAssetCode, err := models.AdminStellarAssetsG(
		qm.Where(models.AdminStellarAssetColumns.IssuerPublicKeyID+"=?", account.PublicKey),
		qm.Where(models.AdminStellarAssetColumns.AssetCode+"=?", assetCode)).One()

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if dbAssetCode == nil {
		return nil
	}

	err = dbAssetCode.DeleteG()
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

	err := signer.InsertG()
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
	dbSigner, err := models.AdminStellarSignersG(
		qm.Where(models.AdminStellarSignerColumns.StellarAccountPublicKeyID+"=?", account.PublicKey),
		qm.Where(models.AdminStellarSignerColumns.SignerPublicKey+"=?", signerPublicKey),
		qm.Where(models.AdminStellarSignerColumns.Type+"=?", signerType)).One()

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if dbSigner == nil {
		return nil
	}

	err = dbSigner.DeleteG()
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
		err := account.R.IssuerPublicKeyAdminStellarAssets.DeleteAllG()
		if err != nil {
			return err
		}
	}
	if len(account.R.StellarAccountPublicKeyAdminStellarSigners) > 0 {
		err := account.R.StellarAccountPublicKeyAdminStellarSigners.DeleteAllG()
		if err != nil {
			return err
		}
	}
	err := account.DeleteG()
	if err != nil {
		return err
	}
	return nil
}

//GetStellarSigner - returns the stellar signer
func GetStellarSigner(publicKey string) (*models.AdminStellarSigner, error) {
	signer, err := models.AdminStellarSignersG(
		qm.Where(models.AdminStellarSignerColumns.SignerPublicKey+"=?", publicKey)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return signer, nil
}

//UpdateSigner - updates the signer
func UpdateSigner(signer *models.AdminStellarSigner, updatedBy string) error {
	signer.UpdatedBy = updatedBy
	signer.UpdatedAt = time.Now()

	err := signer.UpdateG()
	if err != nil {
		return err
	}

	return nil
}