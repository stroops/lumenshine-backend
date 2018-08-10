package cleanup

import (
	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"
	"log"
	"time"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Cleanup will delete old data from the db
func Cleanup() {

	for {

		err := cleanupMinutelyData()
		if err != nil {
			log.Printf("Error cleaning minutely data %v", err)
		}

		err = cleanupHourlyData()
		if err != nil {
			log.Printf("Error cleaning hourly data %v", err)
		}

		waitTime := time.Minute * time.Duration(config.Cnf.Cleanup.CleanupMinutesInterval)
		time.Sleep(waitTime)

	}
}

func cleanupMinutelyData() error {
	// DELETE FROM current_chart_data_minutely where exchange_rate_time older than 1 hour
	return models.CurrentChartDataMinutelies(utils.DB, qm.Where("exchange_rate_time<?", time.Now().Add(-time.Hour*time.Duration(config.Cnf.Cleanup.HoursToKeepMinutelyData)))).DeleteAll()
}

func cleanupHourlyData() error {

	// DELETE FROM current_chart_data_hourly where exchange_rate_time older than 1 day
	return models.CurrentChartDataHourlies(utils.DB, qm.Where("exchange_rate_time<?", time.Now().Add(-time.Hour*time.Duration(config.Cnf.Cleanup.HoursToKeepHourlyData)))).DeleteAll()
}
