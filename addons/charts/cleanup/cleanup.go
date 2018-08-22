package cleanup

import (
	"log"
	"time"

	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/models"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

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
	_, err := models.CurrentChartDataMinutelies(qm.Where("exchange_rate_time<?", time.Now().Add(-time.Hour*time.Duration(config.Cnf.Cleanup.HoursToKeepMinutelyData)))).DeleteAll(utils.DB)
	return err
}

func cleanupHourlyData() error {

	// DELETE FROM current_chart_data_hourly where exchange_rate_time older than 1 day
	_, err := models.CurrentChartDataHourlies(qm.Where("exchange_rate_time<?", time.Now().Add(-time.Hour*time.Duration(config.Cnf.Cleanup.HoursToKeepHourlyData)))).DeleteAll(utils.DB)
	return err
}
