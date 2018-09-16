package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"math/big"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/sirupsen/logrus"
)

var (
	ecs ExchangeCurrencySlice
)

//ExchangeCurrency is one ec
type ExchangeCurrency struct {
	*m.ExchangeCurrency
	NativeCalculator
}

//ExchangeCurrencySlice list of ec's
type ExchangeCurrencySlice []*ExchangeCurrency

//NativeCalculator calculates the native and denomination amounts
type NativeCalculator struct {
	//decimal places for the given currency
	decimals int
	//calculator for the nativ value
	nativCalc *big.Rat
}

//NewNativeCalculator returns a new calculator that can be used in NativeCalculator
func NewNativeCalculator(decimals int) NativeCalculator {
	nc := NativeCalculator{
		decimals:  decimals,
		nativCalc: new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)),
	}
	return nc
}

//DenomFromString returns the denominator as *big.Int for a given denominator string, e.g. 1234567
func (nc *NativeCalculator) DenomFromString(denomValue string) (*big.Int, error) {
	i := new(big.Int)
	i, ok := i.SetString(denomValue, 10)
	if !ok {
		return nil, fmt.Errorf("Could not parse '%s' into big.Int", denomValue)
	}
	return i, nil
}

//DenomFromNativ converts a given nativ (string)value (like eth,xlm,btc,euro) to its denomination
//The value can be in format 123.45 or 123,45
//If the maximum number of accuracy is exceeded, an error is returned
func (nc *NativeCalculator) DenomFromNativ(nativValue string) (*big.Int, error) {
	valueRat := new(big.Rat)
	nativValue = strings.Replace(nativValue, ",", ".", -1)
	_, ok := valueRat.SetString(nativValue)
	if !ok {
		return nil, fmt.Errorf("Could not convert %v to *big.Rat", nativValue)
	}

	// Calculate value in denomination uom
	valueRat.Mul(valueRat, nc.nativCalc)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("Invalid precision, is value smaller than 1 denomination?")
	}

	return valueRat.Num(), nil
}

//ToNativ converts a given denomination (like wei,sat, stroop, cent) to its nativ value
func (nc *NativeCalculator) ToNativ(denom *big.Int) string {
	r := new(big.Float)
	r.SetInt(denom)

	decimals := big.NewInt(int64(nc.decimals))
	i := big.NewInt(10)
	i = i.Exp(i, decimals, big.NewInt(0))

	f := new(big.Float)
	f.SetInt(i)

	r = r.Quo(r, f)
	return r.String()
}

//readAllECs reads all ExchangeCurrencies into memory
func (db *DB) readAllECs(log *logrus.Entry) (ExchangeCurrencySlice, error) {
	if ecs != nil && len(ecs) > 0 {
		return ecs, nil
	}

	allECs, err := m.ExchangeCurrencies().All(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ExchangeCurrencySlice{}, nil
		}
		return nil, err
	}
	ecs = make(ExchangeCurrencySlice, len(allECs))
	for i, e := range allECs {
		ecs[i] = new(ExchangeCurrency)
		ecs[i].NativeCalculator = NewNativeCalculator(e.Decimals)
		ecs[i].ExchangeCurrency = e
	}
	return ecs, nil
}

//GetExchangeCurrecnyByID returns the EC for the given id
func (db *DB) GetExchangeCurrecnyByID(id int, log *logrus.Entry) (*ExchangeCurrency, error) {
	allECs, err := db.readAllECs(log)
	if err != nil {
		return nil, err
	}

	for _, e := range allECs {
		if e.ID == id {
			return e, nil
		}
	}

	return nil, fmt.Errorf("ExchangeCurrency '%d' does not exist", id)
}

//GetActiveExchangeCurrecnyByID returns the phase activated EC for the given id error if not activated
func (db *DB) GetActiveExchangeCurrecnyByID(exchangeID int, phaseID int, log *logrus.Entry) (*ExchangeCurrency, *m.IcoPhaseActivatedExchangeCurrency, error) {
	allECs, err := db.readAllECs(log)
	if err != nil {
		return nil, nil, err
	}

	phase, err := db.GetICOPhaseByID(phaseID, log)
	if err != nil {
		return nil, nil, fmt.Errorf("Phase '%d' does not exist", phaseID)
	}

	for _, e := range allECs {
		if e.ID == exchangeID {
			for _, ae := range phase.R.IcoPhaseActivatedExchangeCurrencies {
				if ae.ExchangeCurrencyID == e.ID {
					return e, ae, nil
				}
			}
		}
	}

	return nil, nil, fmt.Errorf("ExchangeCurrency or ActiveExchange '%d' does not exist", exchangeID)
}

//GetExchangeCurrecnyByCode returns the EC for the given asset code
func (db *DB) GetExchangeCurrecnyByCode(assetCode string, log *logrus.Entry) (*ExchangeCurrency, error) {
	allECs, err := db.readAllECs(log)
	if err != nil {
		return nil, err
	}

	for _, e := range allECs {
		if e.AssetCode == assetCode {
			return e, nil
		}
	}

	return nil, fmt.Errorf("ExchangeCurrency '%s' does not exist", assetCode)
}

//PriceForCoins returns the price for a given token-amount
func (db *DB) PriceForCoins(coins int64, ec *ExchangeCurrency, ph *m.IcoPhase) (denom *big.Int, nativAmount string, err error) {
	//get the phase-exchange-currency
	for _, phec := range ph.R.IcoPhaseActivatedExchangeCurrencies {
		if phec.ExchangeCurrencyID == ec.ID {
			denom = new(big.Int)
			denom = denom.Mul(big.NewInt(coins), big.NewInt(phec.DenomPricePerToken))
			nativAmount = ec.ToNativ(denom)
			return denom, nativAmount, nil
		}
	}

	return nil, "", fmt.Errorf("Could not calculate demon value for '%s/%d' and phase '%d'", ec.AssetCode, ec.ID, ph.IcoID)
}
