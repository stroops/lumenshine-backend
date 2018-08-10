package db

import (
	"database/sql"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"time"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddKnownCurrency adds new currency to known currencies
func AddKnownCurrency(currency *models.AdminKnownCurrency) error {
	err := currency.InsertG()
	if err != nil {
		return err
	}

	return nil
}

//ExistsKnownCurrency - true if a currency with the same public key and asset code exists
func ExistsKnownCurrency(issuerPublicKey, assetCode string) (bool, error) {
	currency, err := models.AdminKnownCurrenciesG(qm.Where(models.AdminKnownCurrencyColumns.IssuerPublicKey+"=?  AND "+models.AdminKnownCurrencyColumns.AssetCode+"=?", issuerPublicKey, assetCode)).One()

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if currency == nil {
		return false, nil
	}

	return true, nil
}

//GetKnownCurrency - returns the known currency
func GetKnownCurrency(issuerPublicKey, assetCode string) (*models.AdminKnownCurrency, error) {
	currency, err := models.AdminKnownCurrenciesG(qm.Where(models.AdminKnownCurrencyColumns.IssuerPublicKey+"=?  AND "+models.AdminKnownCurrencyColumns.AssetCode+"=?", issuerPublicKey, assetCode)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return currency, nil
}

//GetKnownCurrencyByID - returns the known currency associated to the id
func GetKnownCurrencyByID(ID int) (*models.AdminKnownCurrency, error) {
	currency, err := models.AdminKnownCurrenciesG(qm.Where(models.AdminKnownCurrencyColumns.ID+"=?", ID)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return currency, nil
}

//GetKnownCurrencies - returns all known currencies
func GetKnownCurrencies() (models.AdminKnownCurrencySlice, error) {
	currencies, err := models.AdminKnownCurrenciesG().All()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return currencies, nil
}

//UpdateKnownCurrency - updates the currency in the db
func UpdateKnownCurrency(currency *models.AdminKnownCurrency, updatedBy string) error {
	currency.UpdatedBy = updatedBy
	currency.UpdatedAt = time.Now()

	err := currency.UpdateG()
	if err != nil {
		return err
	}

	return nil
}

//DeleteKnownCurrency - deletes the currency
func DeleteKnownCurrency(currency *models.AdminKnownCurrency) error {
	orderIndex := currency.OrderIndex

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	err = currency.DeleteG()
	if err != nil {
		tx.Rollback()
		return err
	}

	currencies, err := models.AdminKnownCurrenciesG(qm.Where(models.AdminKnownCurrencyColumns.OrderIndex+">?", orderIndex)).All()
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	for _, c := range currencies {
		c.OrderIndex--
		err := c.UpdateG()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

//ChangeKnownCurrencyOrder - changes the order of the currency
func ChangeKnownCurrencyOrder(currency *models.AdminKnownCurrency, OrderModifier int, updatedBy string) error {

	currency.OrderIndex += OrderModifier
	currency.UpdatedBy = updatedBy

	currency2, err := models.AdminKnownCurrenciesG(qm.Where(models.AdminKnownCurrencyColumns.OrderIndex+"=?", currency.OrderIndex)).One()
	if err != nil {
		return err
	}

	currency2.OrderIndex -= OrderModifier
	currency2.UpdatedBy = updatedBy

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	err = currency.UpdateG()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = currency2.UpdateG()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// KnownCurrencyNewOrderIndex returns the greatest order from the db + 1, used when inserting a new currency
func KnownCurrencyNewOrderIndex() (int, error) {
	currency, err := models.AdminKnownCurrenciesG(qm.OrderBy("order_index DESC")).One()

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if currency == nil {
		return 1, nil
	}

	return currency.OrderIndex + 1, nil
}
