package main

import (
	"context"
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/jwt/cmd"
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//our server
type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "jwt-svc"
)

var dbClient pb.DBServiceClient
var mc *memcache.Client

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
		log.WithError(err).WithFields(logrus.Fields{"db-url": dbURL}).Fatalf("Dial failed")
	}
	dbClient = pb.NewDBServiceClient(connDB)

	mc = memcache.New(cnf.MemcachedServerURL)

	//start processing the jwt secretes
	go manageJWTs()

	//start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("JWT-Service listening")

	s := grpc.NewServer()
	pb.RegisterJwtServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}
}

//manageJWTs reads all jtw secretes from the database and stores them in the memcached server
//if a key is expired, the function will calculate a new one and store it in the DB and memcached server
func manageJWTs() {
	log := helpers.GetDefaultLog(ServiceName, "manageJWTs")
	c := context.Background()

	for {
		//get keys from database
		keys, err := dbClient.GetAllJwtKeys(c, &pb.Empty{})
		if err != nil {
			log.WithError(err).Error("Error on calling dbClient.GetAllJwtKeys")
			continue // nothing to do
		} else {
			for _, key := range keys.KeyValues {
				_, err := mc.Get(key.KeyName + "1") //we check if the first key is in there, if not, we create both
				if err != nil {
					if err == memcache.ErrCacheMiss {
						//add first key
						mc.Add(&memcache.Item{
							Key:   key.KeyName + "1",
							Value: []byte(key.KeyValue1),
						})

						//add second key
						mc.Add(&memcache.Item{
							Key:   key.KeyName + "2",
							Value: []byte(key.KeyValue2),
						})
					} else {
						log.WithError(err).WithField("key", key.KeyName).Error("Error on getting value from memcached")
					}
				}

				t1 := time.Unix(key.Valid1To, 0)
				t2 := time.Unix(key.Valid2To, 0)
				now := time.Now()

				if t1.Before(now) {
					//expired. We have to move the second one to the first one and create a new one for the second key
					token1 := key.KeyValue2
					token2 := helpers.RandomString(cnf.LengthJWTKey)

					//valid1 := time.Now().Add(int(cnf.ValidMinutes1) * time.Minute).Unix()
					valid1 := now.Add(time.Duration(float64(cnf.ValidMinutes1) * float64(time.Minute))).Unix()
					valid2 := now.Add(time.Duration(float64(cnf.ValidMinutes2) * float64(time.Minute))).Unix()

					if t2.Before(now) {
						//second also expire-->create new token also for one
						token1 = helpers.RandomString(cnf.LengthJWTKey)
					}

					_, err := dbClient.SetJwtKey(c, &pb.JwtSetKeyRequest{
						Base:    &pb.BaseRequest{RequestId: "JWT-Task", UpdateBy: ServiceName},
						Key:     key.KeyName,
						Value1:  token1,
						Expiry1: valid1,
						Value2:  token2,
						Expiry2: valid2,
					})
					if err != nil {
						log.WithError(err).WithField("key", key.KeyName).Error("Error on setting DB value")
						continue
					}

					//set the new values in memcached
					err = mc.Replace(&memcache.Item{
						Key:   key.KeyName + "1",
						Value: []byte(token1),
					})
					if err != nil {
						log.WithError(err).WithField("key", key.KeyName).Error("Error on setting memcahced1 value")
					}

					//set the new values in memcached
					err = mc.Replace(&memcache.Item{
						Key:   key.KeyName + "2",
						Value: []byte(token2),
					})
					if err != nil {
						log.WithError(err).WithField("key", key.KeyName).Error("Error on setting memcahced2 value")
					}
				}
			}
		}

		time.Sleep(time.Duration(float64(cnf.SleepTimeLoopSeconds) * float64(time.Second)))
	}
}

//GetJwtValue reads the jwt secrete for the key from the memcached server
func (s *server) GetJwtValue(c context.Context, r *pb.KeyRequest) (*pb.JwtKeysResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	mk1, err := mc.Get(r.Key + "1")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"value": r.Key}).Error("Error on getting value1 from memcached")
	}

	mk2, err := mc.Get(r.Key + "2")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"value": r.Key}).Error("Error on getting value2 from memcached")
	}

	return &pb.JwtKeysResponse{
		Key1: string(mk1.Value),
		Key2: string(mk2.Value),
	}, nil
}
