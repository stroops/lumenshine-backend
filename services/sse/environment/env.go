package environment

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/services/sse/config"
	"github.com/Soneso/lumenshine-backend/services/sse/db"
)

//Environment for for the service
type Environment struct {
	//Config from file
	Config *config.Config

	//connection to the horizon database
	DBH *db.DB

	httpServer *http.Server
}

//Env holds the global Environment
//the struct is initialized on service startup and hold globalby
var Env = new(Environment)
