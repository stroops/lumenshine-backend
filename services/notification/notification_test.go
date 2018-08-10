package main

import (
	"context"
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/messaging/cmd"
	"testing"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	dbClientTest           pb.DBServiceClient
	notificationClientTest pb.NotificationServiceClient
)

func init() {
	var err error

	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	if err = readConfig(cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	//connect db service
	dbURL := fmt.Sprintf("%s:%d", cnf.DBSrvHost, cnf.DBSrvPort)
	connDB, err := grpc.Dial(dbURL, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"host": cnf.DBSrvHost,
			"port": cnf.DBSrvPort,
		}).Fatalf("Dial db-srv failed")
	}
	dbClientTest = pb.NewDBServiceClient(connDB)

	//connect notification service
	connNotification, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", cnf.Port), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	notificationClientTest = pb.NewNotificationServiceClient(connNotification)
}

func TestPushMessage(t *testing.T) {
	ctx := context.Background()

	_, err := notificationClientTest.SendPushNotification(ctx, &pb.PushNotificationRequest{
		Base:                       &pb.BaseRequest{UpdateBy: "TestService"},
		UserID:                     2,
		Title:                      "Title push",
		Message:                    "Message push",
		SendAsMailIfNoTokenPresent: true})

	if err != nil {
		t.Error("Error sending push notification", err)
	}

	_, err = notificationClientTest.SendPushNotification(ctx, &pb.PushNotificationRequest{
		Base:                       &pb.BaseRequest{UpdateBy: "TestService"},
		UserID:                     5,
		Title:                      "Title push",
		Message:                    "Message push",
		SendAsMailIfNoTokenPresent: true})

	if err != nil {
		t.Error("Error sending push notification", err)
	}

	_, err = notificationClientTest.SendMailNotification(ctx, &pb.MailNotificationRequest{
		Base:        &pb.BaseRequest{UpdateBy: "TestService"},
		UserID:      5,
		Subject:     "Test mail subject",
		Content:     "Test mail content",
		ContentType: pb.EmailContentType_html})

	if err != nil {
		t.Error("Error sending push notification", err)
	}

	response, err := dbClientTest.DequeueNotifications(ctx, &pb.DequeueRequest{
		Base:       &pb.BaseRequest{UpdateBy: "TestService"},
		LimitCount: 10})

	if response.Notifications[0].UserId != 2 {
		t.Errorf("Expected first notifitaion user id %d, actual: %d", 2, response.Notifications[0].UserId)
	}

	if len(response.Notifications) != 4 {
		t.Errorf("Expected %d notifications, actual: %d", 4, len(response.Notifications))
	}
}
