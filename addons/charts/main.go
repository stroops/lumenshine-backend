package main

//go:generate sqlboiler --wipe --no-tests --no-context --config $HOME/.config/sqlboiler/sqlboiler_charts.toml psql

import (
	"fmt"
	"log"
	"time"

	"github.com/Soneso/lumenshine-backend/addons/charts/api"
	"github.com/Soneso/lumenshine-backend/addons/charts/cleanup"
	"github.com/Soneso/lumenshine-backend/addons/charts/cmd"
	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"github.com/Soneso/lumenshine-backend/addons/charts/history"
	"github.com/Soneso/lumenshine-backend/addons/charts/ticker"
	"github.com/Soneso/lumenshine-backend/addons/charts/utils"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	var err error

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.Fatalf("Error reading root command %v", err)
	}

	if err = config.ReadConfig(cmd); err != nil {
		log.Fatalf("Error reading config. %v", err)
	}

	if err = utils.CreateNewDB(); err != nil {
		log.Fatalf("Error creating db connection. %v", err)
		fmt.Printf("\n Error creating db connection. %v", err)
	}

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		log.Fatalf("Error binding flags %v", err)
	}
	historyPtr, err := cmd.Flags().GetBool("history")
	if err != nil {
		log.Fatalf("Error reading flag %v", err)
	}
	truncateHistoryPtr, err := cmd.Flags().GetBool("truncate")
	if err != nil {
		log.Fatalf("Error reading flag %v", err)
	}

	config.Cnf.TruncateHistoryTable = truncateHistoryPtr

	if historyPtr {
		err = history.GetHistoricalData()
		if err != nil {
			log.Panicf("Error on getHistoricalData %v", err)
		}
	}

	ticker.Ticker()
	go cleanup.Cleanup()

	// gin instance
	r := gin.New()
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	r.Use(mw.GinLogger(logger))
	r.Use(gin.Recovery())
	r.Use(mw.RequestID())
	r.Use(mw.Language())

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: config.Cnf.CORSHosts,
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Credentials",
			"Cache-Control", "Accept-Language", "Accept-User-Language", "X-Request-Id"},
		ExposeHeaders:    []string{"Authorization", "X-Request-Id", "X-MessageCount"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.GET("/chart_currency_pairs", mw.UseIcopContext(api.ChartCurrencyPairs))
	r.POST("/chart_currency_pairs", mw.UseIcopContext(api.ChartCurrencyPairs))

	r.GET("/chart_exchange_rates", mw.UseIcopContext(api.ChartExchangeData))
	r.POST("/chart_exchange_rates", mw.UseIcopContext(api.ChartExchangeData))

	r.GET("/chart_current_rates", mw.UseIcopContext(api.ChartCurrentRates))
	r.POST("/chart_current_rates", mw.UseIcopContext(api.ChartCurrentRates))

	//run the api
	if err := r.Run(":" + config.Cnf.Port); err != nil {
		log.Panicf("Failed to run server %v", err)
	}

}

//we need this, in order for rice to find the box
//rice will not call into the subpackages (e.g. helpers) but only into this package
func initRiceBoxes() {
	rice.MustFindBox("db-files/migrations_src")
}
