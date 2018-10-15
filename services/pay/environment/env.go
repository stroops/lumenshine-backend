package environment

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/services/pay/account"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"

	"github.com/Soneso/lumenshine-backend/services/pay/paymentchannel"
)

//Environment for for the service
type Environment struct {
	//Config from file
	Config *config.Config

	//connection to the customer database
	DBC *db.DB

	//this are all payment clients we can handle
	Clients map[string]paymentchannel.Channel

	AccountConfigurator *account.Configurator

	httpServer *http.Server
}
