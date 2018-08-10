package main

import (
	"fmt"
	"github.com/Soneso/lumenshine-backend/addons/dividend/models"
	"github.com/Soneso/lumenshine-backend/addons/dividend/modelscore"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/volatiletech/sqlboiler/queries"
)

var insertChunk = 20000

//InsertDividends inserts in chunks the dividend entries
func InsertDividends(snapshotID int, dividends modelscore.TrustlineSlice) error {

	if len(dividends) == 0 {
		return nil
	}

	var whitelist = models.DividendColumns.SnapshotID + ","
	whitelist += models.DividendColumns.AccountID + ","
	whitelist += models.DividendColumns.Balance + ","
	whitelist += models.DividendColumns.BalanceLimit + ","
	whitelist += models.DividendColumns.DividendAmount

	insertHeader := "INSERT INTO dividend (" + whitelist + ") VALUES "
	valueTemplate := "(%d,'%v',%d,%d,%v)"

	var totalBalance int64
	for _, dividend := range dividends {
		totalBalance += dividend.Balance
	}

	var totalDividends = Req.DividendMode.Value
	if strings.EqualFold(Req.DividendMode.Type, ModePercent) {
		totalDividends = totalBalance * Req.DividendMode.Value / 100
	}

	var strBuilder strings.Builder
	strBuilder.WriteString(insertHeader)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	count := 0
	var dividendAmount int64
	for _, dividend := range dividends {
		if totalDividends == 0 || totalBalance == 0 {
			dividendAmount = 0
			log.Printf("Total dividend or balance is 0.")
		} else {
			balancePercent := float64(dividend.Balance) * float64(100) / float64(totalBalance)
			rawDividend := float64(totalDividends) * balancePercent / 100
			if rawDividend < 1 {
				rawDividend = 0
			} else {
				rawDividend = math.Round(rawDividend)
			}
			dividendAmount = int64(rawDividend)
		}
		dividendToInsert := strconv.FormatInt(dividendAmount, 10)
		if dividendAmount+dividend.Balance >= dividend.Tlimit {
			dividendToInsert = "null"
		}

		if count > 0 {
			strBuilder.WriteString(",")
		}
		strBuilder.WriteString(fmt.Sprintf(valueTemplate, snapshotID, dividend.Accountid, dividend.Balance, dividend.Tlimit, dividendToInsert))
		count++
		if count > insertChunk {
			query := queries.Raw(db, strBuilder.String())
			query.Exec()
			count = 0
			strBuilder.Reset()
			strBuilder.WriteString(insertHeader)
		}
	}

	if count > 0 && strBuilder.String() != "" {
		query := queries.Raw(db, strBuilder.String())
		query.Exec()
	}

	tx.Commit()
	tx.Rollback()

	return nil
}
