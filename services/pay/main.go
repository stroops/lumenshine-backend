package main

import (
	"os"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/account"
	"github.com/Soneso/lumenshine-backend/services/pay/cmd"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/Soneso/lumenshine-backend/services/pay/environment"
	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"

	"github.com/Soneso/lumenshine-backend/services/pay/config"

	"github.com/Soneso/lumenshine-backend/services/pay/bitcoin"
	"github.com/Soneso/lumenshine-backend/services/pay/ethereum"
	"github.com/Soneso/lumenshine-backend/services/pay/stellar"

	"github.com/Soneso/lumenshine-backend/services/pay/channel"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

const (
	//ServiceName name of this service
	ServiceName = "pay-svc"
)

var (
	env      *environment.Environment
	dbClient pb.DBServiceClient
)

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
	dbc, err := db.CreateNewDB(cnf)
	if err != nil {
		log.WithError(err).Fatal("Error connecting customer database")
	}

	env = new(environment.Environment)
	env.Config = cnf
	env.DBC = dbc
	env.Clients = make(map[string]paymentchannel.Channel)

	/*var bitcoinChainParams *chaincfg.Params
	bitcoinChainParams = &chaincfg.TestNet3Params
	b, _ := bitcoin.NewAddressGenerator("xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u", bitcoinChainParams)

	k, s, err := b.Generate(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(k, " ", s)*/

	/*ec := db.NewNativeCalculator(constants.StellarDecimalPlaces)

	d, _ := ec.DenomFromString("10000000")
	//d, _ := ec.DenomFromNativ("0.0000001")

	fmt.Println(ec.ToNativ(d))
	fmt.Println(d)
	return*/

	createServices()

	//The gRPC service will block the thread
	//this should be the last call in main
	StartGrpcService(env, log)
}

func createServices() {
	cfg := env.Config

	env.AccountConfigurator = account.NewAccountConfigurator(
		env.DBC,
		env.Config,
		channel.NewChanneler(env.DBC, env.Config),
	)

	//Create all clients
	env.Clients[m.PaymentNetworkEthereum] = ethereum.NewEthereumChannel(env.DBC, env.Config)
	env.Clients[m.PaymentNetworkBitcoin] = bitcoin.NewBitcoinChannel(env.DBC, env.Config)
	env.Clients[m.PaymentNetworkStellar] = stellar.NewStellarChannel(env.DBC, env.Config)

	//start the listeners if enabled
	if cfg.Bitcoin.Enabled {
		if err := env.Clients[m.PaymentNetworkBitcoin].Start(); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
		log.Info("Bitcoin not enabled")
	}

	if cfg.Ethereum.Enabled {
		if err := env.Clients[m.PaymentNetworkEthereum].Start(); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
		log.Info("Ethereum not enabled")
	}

	if cfg.Stellar.Enabled {
		if err := env.Clients[m.PaymentNetworkStellar].Start(); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
		log.Info("Stellar not enabled")
	}

}
