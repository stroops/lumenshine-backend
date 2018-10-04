package db

import (
	"database/sql"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//ExistsCurrency - true if a currency with the specified asset code and issuer is found
func ExistsCurrency(id int) (bool, error) {
	currency, err := models.ExchangeCurrencies(qm.Select(models.ExchangeCurrencyColumns.ID),
		qm.Where(models.ExchangeCurrencyColumns.ID+"=?", id)).One(DBC)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if currency == nil {
		return false, nil
	}

	return true, nil
}
