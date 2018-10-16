package api

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/admin/client"
	"github.com/stellar/go/clients/horizon"
)

//GetHorizonAccount returns the horizon-account for the given address or false if it does not exist
func GetHorizonAccount(account string) (horizon.Account, bool, error) {
	var hAccount horizon.Account
	hAccount, err := client.HorizonClient.LoadAccount(account)
	if err != nil {
		if err, ok := err.(*horizon.Error); ok && err.Response.StatusCode == http.StatusNotFound {
			return hAccount, false, nil
		}
		return hAccount, false, err
	}
	return hAccount, true, nil
}
