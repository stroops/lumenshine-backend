package main

//go:generate sqlboiler --wipe -b goose_db_version --no-tests --tinyint-as-bool=true --config $HOME/.config/sqlboiler/sqlboiler_db.toml postgres

import (
	"fmt"
	"net"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/db/cmd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

//our server
type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "db-svc"
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

	if err = createNewDB(log, cnf); err != nil {
		log.WithError(err).Fatalf("Error creating db connection")
	}
	defer db.Close()

	//start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("DB-Service listening")

	s := grpc.NewServer()
	pb.RegisterDBServiceServer(s, &server{})
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
