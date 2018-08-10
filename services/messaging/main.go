package main

import (
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/messaging/cmd"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Server - our server
type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "messaging-service"
)

var (
	dbClient   pb.DBServiceClient
	mailClient pb.MailServiceClient
)

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
	dbURL := fmt.Sprintf("%s:%d", cnf.Services.DBSrvHost, cnf.Services.DBSrvPort)
	connDB, err := grpc.Dial(dbURL, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"host": cnf.Services.DBSrvHost,
			"port": cnf.Services.DBSrvPort,
		}).Fatalf("Dial db-srv failed")
	}
	dbClient = pb.NewDBServiceClient(connDB)

	//connect mail service
	connMail, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.MailSrvHost, cnf.Services.MailSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	mailClient = pb.NewMailServiceClient(connMail)

	log.WithFields(logrus.Fields{"idleSeconds": cnf.IdleSeconds, "limitCount": cnf.LimitCount}).Print("Messaging-Service running")

	StartSender()
}
