package main

import (
	"fmt"
	"net"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/notification/cmd"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Server - our server
type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "notification-service"
)

var dbClient pb.DBServiceClient
var sseClient pb.SSEServiceClient
var sseListener *SSEListener

func main() {
	var err error

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	if err = readConfig(cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	connectServices(log)

	//listener for SSE Events
	sseListener = NewListenSSE()
	go sseListener.Run()

	//start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("Notification-Service listening")

	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}
}

func connectServices(log *logrus.Entry) {
	//connect db service
	dbURL := fmt.Sprintf("%s:%d", cnf.DBSrvHost, cnf.DBSrvPort)
	connDB, err := grpc.Dial(dbURL, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"host": cnf.DBSrvHost,
			"port": cnf.DBSrvPort,
		}).Fatalf("Dial db-srv failed")
	}
	dbClient = pb.NewDBServiceClient(connDB)

	//connect sse service
	connSSE, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.SSESrvHost, cnf.SSESrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	sseClient = pb.NewSSEServiceClient(connSSE)
}
