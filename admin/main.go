package main

//go:generate sqlboiler --wipe --no-tests --no-context --add-global-variants --config $HOME/.config/sqlboiler/sqlboiler_admin.toml psql
//go:generate sqlboiler --wipe --no-tests --no-context --add-global-variants --output=../db/stellarcore/models --pkgname=stellarcore --config $HOME/.config/sqlboiler/sqlboiler_stellar_core.toml psql

import (
	"errors"
	"fmt"
	"net"
	"time"

	context "golang.org/x/net/context"

	"github.com/Soneso/lumenshine-backend/admin/api"
	"github.com/Soneso/lumenshine-backend/admin/client"
	"github.com/Soneso/lumenshine-backend/admin/cmd"
	tt "github.com/Soneso/lumenshine-backend/admin/templates"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/gin-contrib/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Soneso/lumenshine-backend/admin/config"
	"github.com/Soneso/lumenshine-backend/admin/db"
	"github.com/Soneso/lumenshine-backend/admin/middleware"
	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	rice "github.com/GeertJohan/go.rice"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "adminsvc"
)

func main() {

	var err error

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	if err = config.ReadConfig(cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	if err = db.CreateNewDB(config.Cnf); err != nil {
		log.WithError(err).Fatalf("Error creating db connection")
	}

	connectServices(log)
	initBoxes()
	tt.LoadTemplates(log)

	go startGRPC()

	r := gin.New()
	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: config.Cnf.CORSHosts,
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Credentials",
			"Cache-Control", "Accept-Language", "Accept-User-Language", "X-Request-Id"},
		ExposeHeaders:    []string{"Authorization", "X-Request-Id", "X-MessageCount", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	r.Use(mw.GinLogger(logger))
	r.Use(gin.Recovery())
	r.Use(mw.RequestID())
	r.Use(mw.Language())

	ic := &mw.IcopContextMiddleware{ServiceName: ServiceName}
	r.Use(ic.MiddlewareFunc())

	authMiddleware := &middleware.GinJWTMiddleware{
		Realm:         "admin",
		Key:           []byte("86a111a1-9072-4c06-95fa-9d6f80c025f5"),
		Authenticator: api.LoginFunc(api.Login),
		Timeout:       time.Duration(999 * time.Hour),
		MaxRefresh:    time.Duration(999 * time.Hour),
	}

	// routes wihtout jwt
	r.POST(api.BaseAPIUrl+"login", authMiddleware.LoginHandler)

	auth := r.Group(api.BaseAPIUrlProt)
	auth.Use(authMiddleware.MiddlewareFunc())
	api.AddUserRoutes(auth)
	api.AddStellarAccountRoutes(auth)
	api.AddCustomerRoutes(auth)
	api.AddKnownCurrenciesRoutes(auth)
	api.AddKnownInflationDestinationsRoutes(auth)
	api.AddICORoutes(auth)

	//special handling for the refresh
	auth.GET("refresh", authMiddleware.RefreshHandler)

	//run the api
	if err := r.Run(fmt.Sprintf(":%d", config.Cnf.Port)); err != nil {
		log.WithError(err).Fatalf("Failed to run server")
	}

}

func startGRPC() {

	log := helpers.GetDefaultLog(ServiceName, "Startup")

	//start the grpc server for api endpoints accessible from outside admin
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Cnf.GRPCPort))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": config.Cnf.GRPCPort}).Fatalf("Failed to listen")
		panic(err)
	}
	log.WithFields(logrus.Fields{"port": config.Cnf.GRPCPort}).Print("AdminAPI-Service listening")

	s := grpc.NewServer()
	pb.RegisterAdminApiServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
		panic(err)
	}
}

func connectServices(log *logrus.Entry) {
	//connect mail service
	connMail, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Cnf.Services.MailSrvHost, config.Cnf.Services.MailSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	client.SetMailClient(pb.NewMailServiceClient(connMail))
	client.SetHorizonClient(config.Cnf.StellarNetwork)
}

//GetKnownCurrency returns the currency for the id
func (s *server) GetKnownCurrency(c context.Context, r *pb.GetKnownCurrencyRequest) (*pb.GetKnownCurrencyResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	currency, err := db.GetKnownCurrencyByID(int(r.Id))
	if err != nil {
		log.WithError(err).WithField("ID", r.Id).Error("Error getting known currency by id")
		return nil, err
	}

	if currency == nil {
		return nil, errors.New("Currency inexistent")
	}

	return &pb.GetKnownCurrencyResponse{
		Id:               int64(currency.ID),
		Name:             currency.Name,
		IssuerPublicKey:  currency.IssuerPublicKey,
		AssetCode:        currency.AssetCode,
		ShortDescription: currency.ShortDescription,
		LongDescription:  currency.LongDescription,
		OrderIndex:       int64(currency.OrderIndex),
	}, nil
}

//GetKnownCurrencies returns all currencies
func (s *server) GetKnownCurrencies(c context.Context, r *pb.Empty) (*pb.GetKnownCurrenciesResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	currencies, err := db.GetKnownCurrencies()
	if err != nil {
		log.WithError(err).Error("Error getting known currencies")
		return nil, err
	}

	if currencies == nil {
		return nil, errors.New("No currencies")
	}

	res := new(pb.GetKnownCurrenciesResponse)
	for _, cr := range currencies {
		res.Currencies = append(res.Currencies, &pb.GetKnownCurrencyResponse{
			Id:               int64(cr.ID),
			Name:             cr.Name,
			IssuerPublicKey:  cr.IssuerPublicKey,
			AssetCode:        cr.AssetCode,
			ShortDescription: cr.ShortDescription,
			LongDescription:  cr.LongDescription,
			OrderIndex:       int64(cr.OrderIndex),
		})
	}

	return res, nil
}

//GetKnownInflationDestination returns the destination for the id
func (s *server) GetKnownInflationDestination(c context.Context, r *pb.GetKnownInflationDestinationRequest) (*pb.GetKnownInflationDestinationResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	dest, err := db.GetKnownInflationDestinationByID(int(r.Id))
	if err != nil {
		log.WithError(err).WithField("ID", r.Id).Error("Error getting known inflation destination by id")
		return nil, err
	}

	if dest == nil {
		return nil, errors.New("Destination inexistent")
	}

	return &pb.GetKnownInflationDestinationResponse{
		Id:               int64(dest.ID),
		Name:             dest.Name,
		IssuerPublicKey:  dest.IssuerPublicKey,
		ShortDescription: dest.ShortDescription,
		LongDescription:  dest.LongDescription,
		OrderIndex:       int64(dest.OrderIndex),
	}, nil

}

//GetKnownInflationDestinations returns all destinations
func (s *server) GetKnownInflationDestinations(c context.Context, r *pb.Empty) (*pb.GetKnownInflationDestinationsResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	dest, err := db.GetKnownInflationDestinations()
	if err != nil {
		log.WithError(err).Error("Error getting known inflation destinations")
		return nil, err
	}

	if dest == nil {
		return nil, errors.New("No destinations")
	}

	res := new(pb.GetKnownInflationDestinationsResponse)
	for _, cr := range dest {
		res.Destinations = append(res.Destinations, &pb.GetKnownInflationDestinationResponse{
			Id:               int64(cr.ID),
			Name:             cr.Name,
			IssuerPublicKey:  cr.IssuerPublicKey,
			ShortDescription: cr.ShortDescription,
			LongDescription:  cr.LongDescription,
			OrderIndex:       int64(cr.OrderIndex),
		})
	}

	return res, nil

}

//we need this, because rice will not look for subfunctions/packages yet ...
//please add all boxes in here
func initBoxes() {
	rice.MustFindBox("db-files/migrations_src")
	rice.MustFindBox("templates/mail")
}
