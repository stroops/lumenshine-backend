package main

import (
	"encoding/json"
	"errors"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	context "golang.org/x/net/context"

	"github.com/sirupsen/logrus"
)

//ApplePayload - ios payload
type ApplePayload struct {
	Aps ApplePayloadAlert `json:"aps"`
}

//ApplePayloadAlert - ios payload alert
type ApplePayloadAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

//GooglePayload - android payload
type GooglePayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

//SendPushNotification - stores the notification in the db queue to be picked up and sent by the routine
func (s *server) SendPushNotification(c context.Context, r *pb.PushNotificationRequest) (*pb.Empty, error) {
	var err error
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	idRequest := &pb.IDRequest{
		Base: &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: r.Base.UpdateBy},
		Id:   r.UserID}

	response, err := dbClient.GetPushTokens(c, idRequest)
	if err != nil {
		log.WithFields(logrus.Fields{"userID": r.UserID}).WithError(err).Error("Error reading push tokens")
		return nil, err
	}

	var content []byte
	var notificationType pb.NotificationType
	for _, token := range response.PushTokens {
		if token.DeviceType == pb.DeviceType_apple {
			payloadAlert := ApplePayloadAlert{Title: r.Title, Body: r.Message}
			payload := ApplePayload{Aps: payloadAlert}
			notificationType = pb.NotificationType_ios
			content, err = json.Marshal(payload)
			if err != nil {
				log.WithFields(logrus.Fields{"userID": r.UserID, "device": token.DeviceType.String()}).
					WithError(err).Error("Error json encoding apple payload")
				return nil, err
			}
		}
		if token.DeviceType == pb.DeviceType_google {
			gPayload := GooglePayload{Title: r.Title, Body: r.Message}
			notificationType = pb.NotificationType_android
			content, err = json.Marshal(gPayload)
			if err != nil {
				log.WithFields(logrus.Fields{"userID": r.UserID, "device": token.DeviceType.String()}).
					WithError(err).Error("Error json encoding apple payload")
				return nil, err
			}
		}
		_, err = dbClient.QueuePushNotification(c, &pb.QueuePushNotificationRequest{
			Base:       &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: r.Base.UpdateBy},
			Content:    string(content),
			PushToken:  token.PushToken,
			UserId:     r.UserID,
			DeviceType: notificationType})

		if err != nil {
			log.WithFields(logrus.Fields{"userID": r.UserID, "device": token.DeviceType.String()}).
				WithError(err).Error("Error json encoding apple payload")
			return nil, err
		}
	}

	if len(response.PushTokens) == 0 && r.SendAsMailIfNoTokenPresent {
		idRequest := &pb.IDRequest{
			Base: &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: ServiceName},
			Id:   r.UserID}

		user, err := dbClient.GetUserProfile(c, idRequest)
		if err != nil {
			log.WithFields(logrus.Fields{"userID": r.UserID}).WithError(err).Error("Error reading push tokens")
			return nil, err
		}

		_, err = dbClient.QueueMailNotification(c, &pb.QueueMailNotificationRequest{
			Base:      &pb.BaseRequest{RequestId: r.Base.RequestId, UpdateBy: r.Base.UpdateBy},
			UserId:    r.UserID,
			Content:   r.Message,
			Subject:   r.Title,
			MailType:  pb.EmailContentType_text,
			UserEmail: user.Email})

		if err != nil {
			log.WithFields(logrus.Fields{"userID": r.UserID, "subject": r.Title}).
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
