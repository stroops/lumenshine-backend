//go:generate swagger generate spec -o ./pay_api_swagger.yml -m
//open swagger with swagger-template:
//  swagger serve pay_api_swagger.yml -F swagger --port=8088 --host=localhost --no-open
//  browser: http://petstore.swagger.io/?url=http%3A%2F%2F127.0.0.1%3A8088%2Fswagger.json
//open swagger with redoc-template:
//  swagger serve pay_api_swagger.yml --port=8088 --host=localhost --no-open
//  browser: http://127.0.0.1:8088/docs

package main

import (
	"fmt"
	"runtime"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"

	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/api/payapi/cmd"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	//ServiceName name of this service
	ServiceName = "payapi"
)

var (
	dbClient   pb.DBServiceClient
	jwtClient  pb.JwtServiceClient
	mailClient pb.MailServiceClient
	payClient  pb.PayServiceClient
)

func main() {
	var err error
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	if err = readConfig(log, cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	logLevel, err := logrus.ParseLevel(cnf.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	connectServices(log)
	initBoxes()
	loadTemplates(log)

	r := gin.New()

	//setup gin-logging. This should always be the first middleware
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	r.Use(mw.GinLogger(logger))

	r.Use(gin.Recovery())
	r.Use(mw.RequestID())
	r.Use(mw.Language())

	ic := &mw.IcopContextMiddleware{ServiceName: ServiceName}
	r.Use(ic.MiddlewareFunc())

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: cnf.CORSHosts,
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Credentials",
			"Cache-Control", "Accept-Language", "Accept-User-Language", "X-Request-Id"},
		ExposeHeaders:    []string{"Authorization", "X-Request-Id", "X-MessageCount"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/portal/pay/test", Test)
	r.GET("/portal/pay/info", Info)

	//this group is used, with the full authenticator, which means, userID and claim is present
	//the middleware will check for full logged in (calim must be present)
	//this is the main group, that is used to read data etc.
	auth := r.Group("/portal/pay/")
	auth.Use(authMiddlewareFull.MiddlewareFunc())
	auth.Use(mw.MessageCount()) //Messagecount only for fully logged in users
	{
		auth.POST("refresh", authMiddlewareFull.RefreshHandler)
		auth.GET("ico_phase_price_for_amount", mw.UseIcopContext(PriceForCoin))
		auth.GET("ico_phase_details", mw.UseIcopContext(IcoPhaseDetails))

		auth.POST("create_order", mw.UseIcopContext(CreateOrder))
		auth.GET("order_list", mw.UseIcopContext(OrderList))
		auth.GET("order_details", mw.UseIcopContext(OrderDetails))

		auth.GET("get_issuer_data", mw.UseIcopContext(OrderDetails))
		auth.POST("create_stellar_account", mw.UseIcopContext(ExecuteTransaction))
		auth.GET("get_payment_transaction", mw.UseIcopContext(OrderDetails))
		auth.POST("execute_payment_transaction", mw.UseIcopContext(ExecuteTransaction))

		auth.GET("get_order_trust_status", mw.UseIcopContext(OrderGetTrustStatus))
		auth.GET("get_order_payment_transaction", mw.UseIcopContext(OrderGetTransaction))
		auth.POST("execute_transaction", mw.UseIcopContext(ExecuteTransaction))
	}

	//run the api
	if err := r.Run(fmt.Sprintf(":%d", cnf.Port)); err != nil {
		log.WithError(err).Fatalf("Failed to run server")
	}
}

func connectServices(log *logrus.Entry) {
	//connect mail service
	connMail, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.MailSrvHost, cnf.Services.MailSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	mailClient = pb.NewMailServiceClient(connMail)

	//connect jwt service
	connJwt, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.JwtSrvHost, cnf.Services.JwtSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	jwtClient = pb.NewJwtServiceClient(connJwt)

	//connect db service
	connDB, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.DBSrvHost, cnf.Services.DBSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	dbClient = pb.NewDBServiceClient(connDB)

	//connect pay service
	connPay, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.PaySrvHost, cnf.Services.PaySrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	payClient = pb.NewPayServiceClient(connPay)
}

var (
	buildDate  string
	gitVersion string
	gitRemote  string
)

// InfoStruct represents the information for the application
// swagger:model InfoStruct
type InfoStruct struct {
	Version             string
	NumGoRutines        int
	MemMbUsedAlloc      uint64
	MemMbUsedTotalAlloc uint64
	BuildDate           string
	GitVersion          string
	GitRemote           string
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Info Prints some information on the binary
// swagger:route GET /info InfoPage
//
// Prints some information on the binary and runtime
//
// Responses:
//   200: InfoStruct
func Info(c *gin.Context) {
	d := new(InfoStruct)
	d.Version = "1"
	d.NumGoRutines = runtime.NumGoroutine()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	d.MemMbUsedAlloc = bToMb(m.Alloc)
	d.MemMbUsedTotalAlloc = bToMb(m.TotalAlloc)
	d.BuildDate = buildDate
	d.GitVersion = gitVersion
	d.GitRemote = gitRemote
	c.JSON(http.StatusOK, d)
}

//Test just for test purpose
func Test(c *gin.Context) {
	l, _ := c.Get("language")
	c.JSON(http.StatusOK, gin.H{
		"language": l,
	})
}

//we need this, because rice will not look for subfunctions/packages yet ...
//please add all boxes in here
func initBoxes() {
	rice.MustFindBox("templates/mail")
}
