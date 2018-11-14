package main

import (
	"time"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

type SSEListener struct {
	log *logrus.Entry
}

func NewListenSSE() *SSEListener {
	return &SSEListener{
		log: helpers.GetDefaultLog(ServiceName, "SSEListener"),
	}
}

//Run runs a loop and gathers the latest data to be send to the clients
//should be called as a go routine
func (s *SSEListener) Run() {
	ctx := context.Background()
	baseRequest := &pb.BaseRequest{RequestId: "0", UpdateBy: ServiceName}
	for {
		data, err := sseClient.GetData(ctx, &pb.SSEGetDataRequest{
			Base:          baseRequest,
			SourceReciver: m.SourceReceiverNotify,
			Count:         20,
		})

		if err != nil {
			s.log.WithError(err).Error("Error reading data")
			time.Sleep(5 * time.Second)
		} else {
			if data != nil && data.Data != nil {
				for _, d := range data.Data {
					err := sendPush(d.StellarAccount, "Lumenshine", "Payment Received", s.log)
					if err != nil {
						s.log.WithError(err).Error("Error sending push notification")
					}
				}
			} else {
				//max 10 messages per second
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func sendPush(publicKey string, title string, message string, log *logrus.Entry) error {
	ctx := context.Background()
	baseRequest := &pb.BaseRequest{RequestId: "0", UpdateBy: ServiceName}

	req := &pb.GetWalletByPublicKeyRequest{
		Base:      baseRequest,
		PublicKey: publicKey,
	}
	wallet, err := dbClient.GetWalletByPublicKey(ctx, req)
	if err != nil {
		return err
	}

	params := []*pb.NotificationParameter{
		&pb.NotificationParameter{Type: pb.NotificationParameterType_ios_title_localized_key, Value: "payment_received"},
		&pb.NotificationParameter{Type: pb.NotificationParameterType_ios_category, Value: "GENERAL"},
		&pb.NotificationParameter{Type: pb.NotificationParameterType_ios_wallet_key, Value: publicKey},
	}

	//TODO: implement reading of parameter from user_profile
	pushReq := &pb.PushNotificationRequest{
		Base:                       baseRequest,
		UserID:                     wallet.UserId,
		Title:                      title,
		Message:                    message,
		Parameters:                 params,
		SendAsMailIfNoTokenPresent: true,
	}
	_, err = sendPushNotification(ctx, pushReq, log)
	if err != nil {
		return err
	}

	return nil
}
