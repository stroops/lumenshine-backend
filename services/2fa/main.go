package main

import (
	"encoding/base32"
	"fmt"
	"io/ioutil"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"

	"net"
	"net/url"

	"github.com/Soneso/lumenshine-backend/services/2fa/cmd"

	"github.com/dgryski/dgoogauth"
	qr "github.com/rsc/qr"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

const (
	//ServiceName name of this service
	ServiceName = "2fa-svc"
)

var (
	qrFilename = "qr.png"
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

	//start the service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("2FA-Service listening")

	s := grpc.NewServer()
	pb.RegisterTwoFactorAuthServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}
}

func (s *server) ConstructUser(ctx context.Context, r *pb.ConstructUserRequest) (*pb.ConstructUserResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	// now := time.Now()
	issuer := cnf.IssuerName

	secrete := []byte(helpers.RandomString(cnf.SecreteKeyLen))

	secretBase32 := base32.StdEncoding.EncodeToString(secrete)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		log.WithError(err).Error("Error parsing url")
		return nil, err
	}

	// google auth app will not recognize url encoded characters e.g. will not recognize %20 as space will just print %20
	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(r.Email)

	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		log.WithError(err).Error("Error encoding qr-image")
		return nil, err
	}

	if cnf.IsDevSystem {
		//write the image to disk on dev
		b := code.PNG()
		err = ioutil.WriteFile("qr.png", b, 0600)
		if err != nil {
			log.WithError(err).Error("Error saving qr-image")
			return nil, err
		}
	}

	return &pb.ConstructUserResponse{
		Url:    URL.String(),
		Bitmap: code.PNG(),
		Secret: secretBase32,
	}, nil
}

func (s *server) Authenticate(ctx context.Context, r *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	otpc := &dgoogauth.OTPConfig{
		Secret:      r.Secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(r.Code)
	if err != nil {
		log.WithError(err).WithField("code", r.Code).Error("Error authenticating")
		return nil, err
	}

	return &pb.AuthenticateResponse{Result: val}, nil
}
