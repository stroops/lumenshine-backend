package main

import (
	"encoding/json"
	"errors"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	context "golang.org/x/net/context"

	"github.com/sirupsen/logrus"
)

//ApplePayload - ios payload
type ApplePayload struct {
	Aps       ApplePayloadAlert `json:"aps"`
	WalletKey string            `json:"wallet_key,omitempty"`
}

//ApplePayloadAlert - alert struct
type ApplePayloadAlert struct {
	Alert    ApplePayloadAlertContent `json:"alert"`
	Category string                   `json:"category,omitempty"`
}

//ApplePayloadAlertContent - ios payload alert
type ApplePayloadAlertContent struct {
	TitleLocKey string `json:"title-loc-key,omitempty"`
}

//GooglePayload - android payload
type GooglePayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *server) SendPushNotification(c context.Context, request *pb.PushNotificationRequest) (*pb.Empty, error) {
	log := helpers.GetDefaultLog(ServiceName, request.Base.RequestId)

	response, err := sendPushNotification(c, request, log)
	if err != nil {
		log.WithError(err).Error("Error on notification sending")
		return nil, err
	}

	return response, nil
}

//SendPushNotification - stores the notification in the db queue to be picked up and sent by the routine
func sendPushNotification(c context.Context, request *pb.PushNotificationRequest, log *logrus.Entry) (*pb.Empty, error) {
	var err error
	idRequest := &pb.IDRequest{
		Base: &pb.BaseRequest{RequestId: request.Base.RequestId, UpdateBy: request.Base.UpdateBy},
		Id:   request.UserID}

	response, err := dbClient.GetPushTokens(c, idRequest)
	if err != nil {
		log.WithFields(logrus.Fields{"userID": request.UserID}).WithError(err).Error("Error reading push tokens")
		return nil, err
	}

	var content []byte
	var notificationType pb.NotificationType
	for _, token := range response.PushTokens {
		if token.DeviceType == pb.DeviceType_apple {
			payload := ApplePayload{Aps: ApplePayloadAlert{Alert: ApplePayloadAlertContent{}}}
			for _, p := range request.Parameters {
				if p.Type == pb.NotificationParameterType_ios_title_localized_key {
					payload.Aps.Alert.TitleLocKey = p.Value
				}
				if p.Type == pb.NotificationParameterType_ios_category {
					payload.Aps.Category = p.Value
				}
				if p.Type == pb.NotificationParameterType_ios_wallet_key {
					payload.WalletKey = p.Value
				}
			}
			notificationType = pb.NotificationType_ios
			content, err = json.Marshal(payload)
			if err != nil {
				log.WithFields(logrus.Fields{"userID": request.UserID, "device": token.DeviceType.String()}).
					WithError(err).Error("Error json encoding apple payload")
				return nil, err
			}
		}
		if token.DeviceType == pb.DeviceType_google {
			gPayload := GooglePayload{Title: request.Title, Body: request.Message}
			notificationType = pb.NotificationType_android
			content, err = json.Marshal(gPayload)
			if err != nil {
				log.WithFields(logrus.Fields{"userID": request.UserID, "device": token.DeviceType.String()}).
					WithError(err).Error("Error json encoding apple payload")
				return nil, err
			}
		}
		_, err = dbClient.QueuePushNotification(c, &pb.QueuePushNotificationRequest{
			Base:       &pb.BaseRequest{RequestId: request.Base.RequestId, UpdateBy: request.Base.UpdateBy},
			Content:    string(content),
			PushToken:  token.PushToken,
			UserId:     request.UserID,
			DeviceType: notificationType})

		if err != nil {
			log.WithFields(logrus.Fields{"userID": request.UserID, "device": token.DeviceType.String()}).
				WithError(err).Error("Error json encoding apple payload")
			return nil, err
		}
	}

	if len(response.PushTokens) == 0 && request.SendAsMailIfNoTokenPresent {
		idRequest := &pb.IDRequest{
			Base: &pb.BaseRequest{RequestId: request.Base.RequestId, UpdateBy: ServiceName},
			Id:   request.UserID}

		user, err := dbClient.GetUserProfile(c, idRequest)
		if err != nil {
			log.WithFields(logrus.Fields{"userID": request.UserID}).WithError(err).Error("Error reading push tokens")
			return nil, err
		}

		if !user.MailNotifications {
			//delete sse-config for account
			req := &pb.GetWalletsRequest{
				Base:   &pb.BaseRequest{UpdateBy: ServiceName},
				UserId: user.Id,
			}
			wallets, err := dbClient.GetUserWallets(c, req)
			if err != nil {
				return nil, err
			}
			for _, w := range wallets.Wallets {
				_, err := sseClient.RemoveListening(c, &pb.SSERemoveListeningRequest{
					Base:           &pb.BaseRequest{UpdateBy: ServiceName},
					SourceReciver:  m.SourceReceiverNotify,
					StellarAccount: w.PublicKey,
				})
				if err != nil {
					return nil, err
				}
			}

			//return so we don't send mail
			return &pb.Empty{}, nil
		}

		_, err = dbClient.QueueMailNotification(c, &pb.QueueMailNotificationRequest{
			Base:      &pb.BaseRequest{RequestId: request.Base.RequestId, UpdateBy: request.Base.UpdateBy},
			UserId:    request.UserID,
			Content:   request.Message,
			Subject:   request.Title,
			MailType:  pb.EmailContentType_text,
			UserEmail: user.Email})

		if err != nil {
			log.WithFields(logrus.Fields{"userID": request.UserID, "subject": request.Title}).
				WithError(err).Error("Error queueing mail notification")
			return nil, err
		}
	}

	return &pb.Empty{}, nil
}

//SendMailNotification - stores the notification in the db queue to be picked up and sent by the routine
func (s *server) SendMailNotification(c context.Context, r *pb.MailNotificationRequest) (*pb.Empty, error) {
	var err error
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	idRequest := &pb.IDRequest{
		Base: &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: r.Base.UpdateBy},
		Id:   r.UserID}

	response, err := dbClient.GetUserProfile(c, idRequest)
	if err != nil {
		log.WithFields(logrus.Fields{"userID": r.UserID}).WithError(err).Error("Error reading push tokens")
		return nil, err
	}

	if response == nil {
		log.WithFields(logrus.Fields{"userID": r.UserID}).Error("User not found")
		return nil, errors.New("User not found")
	}

	_, err = dbClient.QueueMailNotification(c, &pb.QueueMailNotificationRequest{
		Base:      &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: r.Base.UpdateBy},
		UserId:    r.UserID,
		Content:   r.Content,
		Subject:   r.Subject,
		MailType:  r.ContentType,
		UserEmail: response.Email})

	if err != nil {
		log.WithFields(logrus.Fields{"userID": r.UserID, "subject": r.Subject}).
			WithError(err).Error("Error queueing mail notification")
		return nil, err
	}

	return &pb.Empty{}, nil
}
