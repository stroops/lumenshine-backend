package client

import (
	"net/http"
	"time"

	"github.com/Soneso/lumenshine-backend/admin/config"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/stellar/go/clients/horizon"
)

var (
	//MailClient - protobuffer mail client
	MailClient pb.MailServiceClient
	//HorizonClient - stellar horizon client
	HorizonClient *horizon.Client
)

//SetMailClient - sets the mail client
func SetMailClient(mailClient pb.MailServiceClient) {
	MailClient = mailClient
}

//SetHorizonClient - initializes the horizon client
func SetHorizonClient(config config.StellarNetworkConfig) {
	HorizonClient = &horizon.Client{
		URL: config.HorizonURL,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
