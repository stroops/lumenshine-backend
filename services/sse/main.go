package main

import (
	"fmt"
	"net"

	rice "github.com/GeertJohan/go.rice"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/sse/cmd"
	"github.com/Soneso/lumenshine-backend/services/sse/db"
	"github.com/Soneso/lumenshine-backend/services/sse/environment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Soneso/lumenshine-backend/services/sse/config"

	"github.com/sirupsen/logrus"
)

const (
	//ServiceName name of this service
	ServiceName = "sse-svc"
)

//our server
type server struct {
	Env *environment.Environment
}

func main() {
	var err error
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	var cnf *config.Config
	if cnf, err = config.ReadConfig(cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	//create DB Connection
	dbh, err := db.CreateNewHorizonDB(cnf)
	if err != nil {
		log.WithError(err).Fatal("Error connecting horizon database")
	}

	environment.Env.Config = cnf
	environment.Env.DBH = dbh

	//start stellar processor
	sp := NewStellarProcessor(log)
	go sp.StartProcessing()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("SSE-Service listening")

	s := grpc.NewServer()
	pb.RegisterSSEServiceServer(s, &server{
		Env: environment.Env,
	})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}
}

//we need this, in order for rice to find the box
//rice will not call into the subpackages (e.g. helpers) but only into this package
func initRiceBoxes() {
	rice.MustFindBox("db-files/migrations_src")
}
