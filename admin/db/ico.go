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
	tx, err := DBC.Begin()
	if err != nil {
		return err
	}
	_, err = models.IcoSupportedExchangeCurrencies(
		qm.Where(models.IcoSupportedExchangeCurrencyColumns.IcoID+"=?", icoID)).DeleteAll(tx)
	if err != nil {
		tx.Rollback()
		return err
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

//DeleteIco - deletes the ico
func DeleteIco(ico *models.Ico) error {
	phases, err := models.IcoPhases(qm.Where(models.IcoPhaseColumns.IcoID+"=?", ico.ID),
		qm.Load(models.IcoPhaseRels.IcoPhaseActivatedExchangeCurrencies),
		qm.Load(models.IcoPhaseRels.IcoPhaseActivatedExchangeCurrencies+"."+models.IcoPhaseActivatedExchangeCurrencyRels.IcoPhaseBankAccount)).All(DBC)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := DBC.Begin()
	if err != nil {
		return err
	}

	if phases != nil && len(phases) > 0 {
		for _, phase := range phases {
			if phase.R != nil && phase.R.IcoPhaseActivatedExchangeCurrencies != nil {
				_, err = phase.R.IcoPhaseActivatedExchangeCurrencies.DeleteAll(tx)
				if err != nil {
					tx.Rollback()
					return err
				}
				for _, currency := range phase.R.IcoPhaseActivatedExchangeCurrencies {
					if currency.R != nil && currency.R.IcoPhaseBankAccount != nil {
						_, err = currency.R.IcoPhaseBankAccount.Delete(tx)
						if err != nil {
							tx.Rollback()
							return err
						}
					}
				}
			}
		}
		_, err = phases.DeleteAll(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = models.IcoSupportedExchangeCurrencies(
		qm.Where(models.IcoSupportedExchangeCurrencyColumns.IcoID+"=?", ico.ID)).DeleteAll(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = ico.Delete(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
