//go:generate swagger generate spec -o ./sse_api_swagger.yml -m
//open swagger with swagger-template:
//  swagger serve sse_api_swagger.yml -F swagger --port=8088 --host=localhost --no-open
//  browser: http://petstore.swagger.io/?url=http%3A%2F%2F127.0.0.1%3A8088%2Fswagger.json
//open swagger with redoc-template:
//  swagger serve sse_api_swagger.yml --port=8088 --host=localhost --no-open
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
	ServiceName = "sseapi"
)

var (
	dbClient    pb.DBServiceClient
	jwtClient   pb.JwtServiceClient
	sseClient   pb.SSEServiceClient
	hub         *Hub
	sseListener *SSEListener
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

	r := gin.New()

	//setup gin-logging. This should always be the first middleware
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	r.Use(mw.GinLogger(logger))

	r.Use(gin.Recovery())
	r.Use(mw.RequestID())

	ic := &mw.IcopContextMiddleware{ServiceName: ServiceName}
	r.Use(ic.MiddlewareFunc())

	//websocket hub
	hub = newHub()
	go hub.run()

	//listener for SSE Events
	sseListener = NewListenSSE()
	go sseListener.Run()

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

	r.GET("/portal/sse/test", Test)
	r.GET("/portal/sse/info", Info)

	if cnf.IsDevSystem {
		//we run this endpoints only in dev mode, because we can't use websockets in insomnia
		//we have the send_message endpoint for sending messages and the /test which will return a small vue app for doing local test
		r.GET("/get_ws", SSEContext(GetWS))
		r.POST("/remove_ws", SSEContext(RemoveWS))
		r.POST("/listen_account", SSEContext(ListenAccount))
		r.POST("/remove_account", SSEContext(RemoveAccount))

		r.POST("/send_message", SSEContext(SendMessage))

		r.GET("/test", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "templates/test.html")
		})
	}

	//this group is used, with the full authenticator, which means, userID and claim is present
	//the middleware will check for full logged in (calim must be present)
	//this is the main group, that is used to read data etc.
	auth := r.Group("/portal/sse")
	auth.Use(authMiddlewareFull.MiddlewareFunc())
	{
		auth.POST("refresh", authMiddlewareFull.RefreshHandler)
		auth.GET("/get_ws", SSEContext(GetWS))
		auth.POST("/remove_ws", SSEContext(RemoveWS))
		auth.POST("/listen_account", SSEContext(ListenAccount))
		auth.POST("/remove_account", SSEContext(RemoveAccount))
	}

	authWS := r.Group("/portal/sse/ws")
	authWS.Use(authMiddlewareFullWS.MiddlewareFunc())
	{
		authWS.POST("refresh", authMiddlewareFullWS.RefreshHandler)
		authWS.GET("/get_ws", SSEContext(GetWS))
	}

	//run the api
	if err := r.Run(fmt.Sprintf(":%d", cnf.Port)); err != nil {
		log.WithError(err).Fatalf("Failed to run server")
	}
}

//SSEContext request context for sse
func SSEContext(f func(h *Hub, uc *mw.IcopContext, c *gin.Context)) gin.HandlerFunc {
	uc := &mw.IcopContext{}
	return func(c *gin.Context) {
		uc.Language = c.GetString("language")
		uc.RequestID = c.GetString("request_id")
		uc.Log = helpers.GetDefaultLog(c.GetString("servicename"), uc.RequestID)
		f(hub, uc, c)
	}
}

func connectServices(log *logrus.Entry) {
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

	//connect sse service
	connSSE, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.SSESrvHost, cnf.Services.SSESrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	sseClient = pb.NewSSEServiceClient(connSSE)
}

var (
	buildDate  string
	gitVersion string
	gitRemote  string
)

//InfoClient info for all connected clients
type InfoClient struct {
	Key       string
	Addresses []string
}

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

	HubClients      []*InfoClient
	HubClientsCount int

	HubAddresses      map[string][]string
	HubAddressesCount int
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

	d.HubClientsCount = len(hub.clients)
	for key, c := range hub.clients {
		d.HubClients = append(d.HubClients, &InfoClient{
			Key:       key,
			Addresses: c.addresses,
		})
	}
	d.HubAddresses = hub.addresses
	d.HubAddressesCount = len(hub.addresses)

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
	rice.MustFindBox("templates")
}
