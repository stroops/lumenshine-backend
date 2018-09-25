package client

import (
	"github.com/Soneso/lumenshine-backend/pb"
)

var (
	//MailClient - protobuffer mail client
	MailClient pb.MailServiceClient
)

//SetMailClient - sets the mail client
func SetMailClient(mailClient pb.MailServiceClient) {
	MailClient = mailClient
}
