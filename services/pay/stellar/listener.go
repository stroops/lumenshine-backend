package stellar

import (
	"github.com/Soneso/lumenshine-backend/services/pay/db"
)

type Listener struct {
	DB *db.DB
}

type Config struct {
	DB *db.DB
}

func NewListener(cnf Config) *Listener {
	l := new(Listener)
	l.DB = cnf.DB

	//l.log.Info("Stellar-Listener created")
	return l
}
