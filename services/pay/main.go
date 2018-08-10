package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/pay/account"
	"github.com/Soneso/lumenshine-backend/services/pay/cmd"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/Soneso/lumenshine-backend/services/pay/environment"
	"github.com/Soneso/lumenshine-backend/services/pay/stellar"
	"os"

	"github.com/Soneso/lumenshine-backend/services/pay/config"

	"github.com/Soneso/lumenshine-backend/services/pay/bitcoin"
	"github.com/Soneso/lumenshine-backend/services/pay/ethereum"

	"github.com/btcsuite/btcd/chaincfg"
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

	createServices()

	//The gRPC service will block the thread
	//this should be the last call in main
	StartGrpcService(env, log)
}

func createServices() {
	cfg := env.Config

	env.AccountConfigurator = account.NewAccountConfigurator(env.DBC, env.Config)

	//create bitcoin address generator
	var bitcoinChainParams *chaincfg.Params
	if cfg.Bitcoin.Testnet {
		bitcoinChainParams = &chaincfg.TestNet3Params
	} else {
		bitcoinChainParams = &chaincfg.MainNetParams
	}

	bitcoinAddressGenerator, err := bitcoin.NewAddressGenerator(cfg.Bitcoin.MasterPublicKey, bitcoinChainParams)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	env.BitcoinAddressGenerator = bitcoinAddressGenerator

	//create Ethereum address generator
	ethereumAddressGenerator, err := ethereum.NewAddressGenerator(cfg.Bitcoin.MasterPublicKey)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	env.EthereumAddressGenerator = ethereumAddressGenerator

	//create Stellar address generator
	stellarAddressGenerator, err := stellar.NewAddressGenerator()
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	env.StellarAddressGenerator = stellarAddressGenerator

	//start the listener
	if cfg.Bitcoin.Enabled {
		env.BitcoinListener = bitcoin.NewListener(env.DBC, env.Config)
		if err := env.BitcoinListener.Start(); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
		log.Info("Bitcoin not enabled")
	}

	if cfg.Ethereum.Enabled {

		env.EthereumListener = ethereum.NewListener(env.DBC, env.Config)
		if err := env.EthereumListener.Start(); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
		log.Info("Ethereum not enabled")
	}

	if cfg.Stellar.Enabled {
		//start the listener?
	} else {
		log.Info("Stellar not enabled")
	}

}
