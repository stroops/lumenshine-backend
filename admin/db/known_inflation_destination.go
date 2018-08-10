package db

import (
	"database/sql"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"time"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddKnownInflationDestination adds new inflation destination to known inflation destinations
func AddKnownInflationDestination(inflationDestination *models.AdminKnownInflationDestination) error {
	err := inflationDestination.InsertG()
	if err != nil {
		return err
	}

	return nil
}

//ExistsKnownInflationDestination - true if a inflation destination with the same asset code exists
func ExistsKnownInflationDestination(issuerPublicKey string) (bool, error) {
	inflationDestination, err := models.AdminKnownInflationDestinationsG(qm.Where(models.AdminKnownInflationDestinationColumns.IssuerPublicKey+"=? ", issuerPublicKey)).One()

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if inflationDestination == nil {
		return false, nil
	}

	return true, nil
}

//GetKnownInflationDestination - returns the known inflation destination
func GetKnownInflationDestination(issuerPublicKey string) (*models.AdminKnownInflationDestination, error) {
	inflationDestination, err := models.AdminKnownInflationDestinationsG(qm.Where(models.AdminKnownInflationDestinationColumns.IssuerPublicKey+"=? ", issuerPublicKey)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return inflationDestination, nil
}

//GetKnownInflationDestinationByID - returns the known inflation destination associated to the id
func GetKnownInflationDestinationByID(ID int) (*models.AdminKnownInflationDestination, error) {
	inflationDestination, err := models.AdminKnownInflationDestinationsG(qm.Where(models.AdminKnownInflationDestinationColumns.ID+"=?", ID)).One()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return inflationDestination, nil
}

//GetKnownInflationDestinations - returns all known inflation destinations
func GetKnownInflationDestinations() (models.AdminKnownInflationDestinationSlice, error) {
	inflationDestinations, err := models.AdminKnownInflationDestinationsG().All()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return inflationDestinations, nil
}

//UpdateKnownInflationDestination - updates the inflation destination in the db
func UpdateKnownInflationDestination(inflationDestination *models.AdminKnownInflationDestination, updatedBy string) error {
	inflationDestination.UpdatedBy = updatedBy
	inflationDestination.UpdatedAt = time.Now()

	err := inflationDestination.UpdateG()
	if err != nil {
		return err
	}

	return nil
}

//DeleteKnownInflationDestination - deletes the inflation destination
func DeleteKnownInflationDestination(inflationDestination *models.AdminKnownInflationDestination) error {
	orderIndex := inflationDestination.OrderIndex

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	err = inflationDestination.DeleteG()
	if err != nil {
		tx.Rollback()
		return err
	}

	inflationDestinations, err := models.AdminKnownInflationDestinationsG(qm.Where(models.AdminKnownInflationDestinationColumns.OrderIndex+">?", orderIndex)).All()
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	for _, c := range inflationDestinations {
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

//ChangeKnownInflationDestinationOrder - changes the order of the inflation destination
func ChangeKnownInflationDestinationOrder(inflationDestination *models.AdminKnownInflationDestination, OrderModifier int, updatedBy string) error {

	inflationDestination.OrderIndex += OrderModifier
	inflationDestination.UpdatedBy = updatedBy

	inflationDestination2, err := models.AdminKnownInflationDestinationsG(qm.Where(models.AdminKnownInflationDestinationColumns.OrderIndex+"=?", inflationDestination.OrderIndex)).One()
	if err != nil {
		return err
	}

	inflationDestination2.OrderIndex -= OrderModifier
	inflationDestination2.UpdatedBy = updatedBy

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	err = inflationDestination.UpdateG()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = inflationDestination2.UpdateG()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// KnownInflationDestinationNewOrderIndex returns the greatest order from the db + 1, used when inserting a new inflation destination
func KnownInflationDestinationNewOrderIndex() (int, error) {
	inflationDestination, err := models.AdminKnownInflationDestinationsG(qm.OrderBy("order_index DESC")).One()

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if inflationDestination == nil {
		return 1, nil
	}

	return inflationDestination.OrderIndex + 1, nil
}
