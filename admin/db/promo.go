package db

import (
	"database/sql"
	"time"

	"github.com/Soneso/lumenshine-backend/admin/models"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

//AddPromo adds a new promo the DB
func AddPromo(promo *models.AdminPromo) error {
	err := promo.InsertG(boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

//GetPromoByID - returns the promo associated to the id
func GetPromoByID(ID int) (*models.AdminPromo, error) {
	promo, err := models.AdminPromos(qm.Where(models.AdminPromoColumns.ID+"=?", ID)).OneG()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return promo, nil
}

//GetPromos - returns all promos
func GetPromos() (models.AdminPromoSlice, error) {
	promos, err := models.AdminPromos(qm.OrderBy(models.AdminPromoColumns.OrderIndex)).AllG()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return promos, nil
}

//UpdatePromo - updates the promo in the db
func UpdatePromo(promo *models.AdminPromo, updatedBy string) error {
	promo.UpdatedBy = updatedBy
	promo.UpdatedAt = time.Now()

	_, err := promo.UpdateG(boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

//DeletePromo - deletes the promo
func DeletePromo(promo *models.AdminPromo) error {
	orderIndex := promo.OrderIndex

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = promo.Delete(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	promos, err := models.AdminPromos(qm.Where(models.AdminPromoColumns.OrderIndex+">?", orderIndex)).AllG()
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	for _, promo := range promos {
		promo.OrderIndex--
		_, err := promo.Update(tx, boil.Infer())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

//ChangePromoOrder - changes the order of the promo
func ChangePromoOrder(promo *models.AdminPromo, OrderModifier int, updatedBy string) error {

	promo.OrderIndex += OrderModifier
	promo.UpdatedBy = updatedBy

	promo2, err := models.AdminPromos(qm.Where(models.AdminPromoColumns.OrderIndex+"=?", promo.OrderIndex)).OneG()
	if err != nil {
		return err
	}

	promo2.OrderIndex -= OrderModifier
	promo2.UpdatedBy = updatedBy

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = promo.Update(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = promo2.Update(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// PromoNewOrderIndex returns the greatest order from the db + 1, used when inserting a new promo
func PromoNewOrderIndex() (int, error) {
	promo, err := models.AdminPromos(qm.OrderBy("order_index DESC")).OneG()

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if promo == nil {
		return 1, nil
	}

	return promo.OrderIndex + 1, nil
}
