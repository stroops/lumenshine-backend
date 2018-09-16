package environment

import (
	"math/big"
	"net/http"

	"github.com/Soneso/lumenshine-backend/services/pay/account"
	"github.com/Soneso/lumenshine-backend/services/pay/bitcoin"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/Soneso/lumenshine-backend/services/pay/ethereum"
	"github.com/Soneso/lumenshine-backend/services/pay/stellar"
)

//Environment for for the service
type Environment struct {
	//Config from file
	Config *config.Config

	//connection to the customer database
	DBC *db.DB

	BitcoinListener *bitcoin.Listener

	EthereumListener *ethereum.Listener

	StellarListener *stellar.Listener

	AccountConfigurator *account.Configurator

	MinimumValueBtc string
	MinimumValueEth string
	SignerPublicKey string

	minimumValueSat int64
	minimumValueWei *big.Int
	httpServer      *http.Server
}
