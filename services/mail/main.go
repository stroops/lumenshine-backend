package main

import (
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/mail/cmd"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//our server
type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "mail-svc"
)

var dbClient pb.DBServiceClient

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

	//start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("Mail-Service listening")

	s := grpc.NewServer()
	pb.RegisterMailServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}
}
