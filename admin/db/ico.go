package db

import (
	"database/sql"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//ExistsIcoName - true if an ico with the same name exists
func ExistsIcoName(icoName string) (bool, error) {
	ico, err := models.Icos(qm.Select(models.IcoColumns.ID),
		qm.Where(models.IcoColumns.IcoName+"=?", icoName)).One(DBC)

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
