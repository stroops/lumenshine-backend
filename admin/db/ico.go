package db

import (
	"database/sql"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

//ExistsIcoName - true if an ico with the same name exists
func ExistsIcoName(icoName string, selfID *int) (bool, error) {
	q := []qm.QueryMod{
		qm.Select(
			models.IcoColumns.ID,
		),
		qm.Where(models.IcoColumns.IcoName+" ilike ?", icoName),
	}
	if selfID != nil {
		q = append(q, qm.Where("id!=?", selfID))
	}
	ico, err := models.Icos(q...).One(DBC)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if ico == nil {
		return false, nil
	}

	return true, nil
}

//GetIco - returns the ico
func GetIco(id int) (*models.Ico, error) {
	ico, err := models.Icos(qm.Where(models.IcoColumns.ID+"=?", id)).One(DBC)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ico, nil
}

//UpdateSupportedCurrencies - updates the supported currencies of an ico
func UpdateSupportedCurrencies(icoID int, supportedCurrencies models.IcoSupportedExchangeCurrencySlice) error {
	excs, err := models.IcoSupportedExchangeCurrencies(
		qm.Where(models.IcoSupportedExchangeCurrencyColumns.IcoID+"=?", icoID)).All(DBC)
	if err != nil {
		return err
	}
	tx, err := DBC.Begin()
	if err != nil {
		return err
	}
	for _, exc := range excs {
		_, err = exc.Delete(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, newCurrency := range supportedCurrencies {
		newCurrency.IcoID = icoID
		err = newCurrency.Insert(tx, boil.Whitelist(models.IcoSupportedExchangeCurrencyColumns.IcoID,
			models.IcoSupportedExchangeCurrencyColumns.ExchangeCurrencyID,
			models.IcoSupportedExchangeCurrencyColumns.UpdatedBy))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
