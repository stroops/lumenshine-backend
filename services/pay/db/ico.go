package db

import (
	"database/sql"
	"fmt"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/sirupsen/logrus"

	// "github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	// "github.com/volatiletech/sqlboiler/queries"
)

var (
	//icos holds a list of all ICOs. This list will be read only once and hold in memory
	//we will use a memcached key to synchronise updateing the list
	icos m.IcoSlice
)

//GetAllICOs returns a list of all ICOs
func (db *DB) GetAllICOs(log *logrus.Entry) (m.IcoSlice, error) {
	return db.readAllICOs(log)
}

//readAllICOs reads all ICOs including all realtions into memory
func (db *DB) readAllICOs(log *logrus.Entry) (m.IcoSlice, error) {
	if icos != nil && len(icos) > 0 {
		return icos, nil
	}

	allIcos, err := m.Icos(
		//read the Phase data
		qm.Load(m.IcoRels.IcoPhases),
		qm.Load(m.IcoRels.IcoPhases+"."+m.IcoPhaseRels.IcoPhaseActivatedExchangeCurrencies),
		qm.Load(m.IcoRels.IcoPhases+"."+m.IcoPhaseRels.IcoPhaseActivatedExchangeCurrencies+"."+m.IcoPhaseActivatedExchangeCurrencyRels.ExchangeCurrency),
		qm.Load(m.IcoRels.IcoPhases+"."+m.IcoPhaseRels.IcoPhaseActivatedExchangeCurrencies+"."+m.IcoPhaseActivatedExchangeCurrencyRels.IcoPhaseBankAccount),

		//read the ExchangeCurrencies data
		qm.Load(m.IcoRels.IcoSupportedExchangeCurrencies),
		qm.Load(m.IcoRels.IcoSupportedExchangeCurrencies+"."+m.IcoSupportedExchangeCurrencyRels.ExchangeCurrency),
	).All(db)

	if err != nil {
		if err == sql.ErrNoRows {
			return m.IcoSlice{}, nil
		}
		return nil, err
	}
	icos = allIcos

	_, err = db.readAllECs(log)

	return icos, err
}

//GetICOByID returns the given ICO by ID
func (db *DB) GetICOByID(id int, log *logrus.Entry) (*m.Ico, error) {
	allICOs, err := db.readAllICOs(log)
	if err != nil {
		return nil, err
	}

	for _, i := range allICOs {
		if i.ID == id {
			return i, nil
		}
	}

	return nil, fmt.Errorf("ICO <%d> does not exist", id)
}

//GetICOPhaseByID returns the given ICOPhase by ID
func (db *DB) GetICOPhaseByID(id int, log *logrus.Entry) (*m.IcoPhase, error) {
	allICOs, err := db.readAllICOs(log)
	if err != nil {
		return nil, err
	}

	for _, i := range allICOs {
		for _, p := range i.R.IcoPhases {
			if p.ID == id {
				return p, nil
			}
		}
	}

	return nil, fmt.Errorf("ICOPhase <%d> does not exist", id)
}
