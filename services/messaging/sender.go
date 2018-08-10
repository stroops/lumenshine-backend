package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"

	"github.com/sirupsen/logrus"
)

var ticker *time.Ticker

//StartSender - starts the sender routine
func StartSender() {
	InitIosClient()

	ticker = time.NewTicker(time.Duration(cnf.IdleSeconds) * time.Second)

	taskDone := make(chan bool)
	done := make(chan bool)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-c: // got signal
			stopSender(taskDone)
		}
	}()

	go func() { defer func() { done <- true }(); sendNotifications(taskDone) }()

	<-done //main waits for task termination
}

//stopSender - stops the sender routine
func stopSender(taskDone chan bool) {
	taskDone <- true
}

//sendNotifications - sends the queued notifications
func sendNotifications(taskDone chan bool) {
	log := helpers.GetDefaultLog(ServiceName, "Sender")
	for {
		ctx := context.Background()
		response, err := dbClient.DequeueNotifications(ctx, &pb.DequeueRequest{
			Base:       &pb.BaseRequest{UpdateBy: ServiceName},
			LimitCount: int64(cnf.LimitCount)})
		if err != nil {
			log.WithError(err).Error("Could not dequeue notifications")
		} else {
			count := len(response.Notifications)
			log.Printf("Dequeued %d notifications", count)
			if count > 0 {
				sendDequeued(response.Notifications)
				continue
			}
		}

		select {
		case <-taskDone:
			return
		case <-ticker.C:
		}
	}
}

func sendDequeued(notifications []*pb.Notification) {
	log := helpers.GetDefaultLog(ServiceName, "Sender")

	var sentNotifications []*pb.NotificationArchive

	for _, notification := range notifications {
		switch notification.NotificationType {
		case pb.NotificationType_ios:
			sentNotifications = append(sentNotifications, sendIos(notification))
		case pb.NotificationType_android:
			sentNotifications = append(sentNotifications, sendAndroid(notification))
		case pb.NotificationType_mail:
			sentNotifications = append(sentNotifications, sendMail(notification))
		}
	}
	ctx := context.Background()
	_, err := dbClient.UpdateNotificationsStatus(ctx, &pb.UpdateNotificationsStatusRequest{
		Base:          &pb.BaseRequest{UpdateBy: ServiceName},
		Notifications: sentNotifications})

	if err != nil {
		log.WithError(err).Error("Could not update sent notifications")
	}
}

func sendIos(notification *pb.Notification) *pb.NotificationArchive {
	log := helpers.GetDefaultLog(ServiceName, "Sender")

	response := sendIosNotification(notification)

	archive := &pb.NotificationArchive{}
	archive.Id = notification.Id
	archive.ExternalStatus = string(response.StatusCode)
	archive.ExternalError = response.Reason

	if response.StatusCode != iosStatuscodeSuccess {
		archive.Status = pb.NotificationStatusCode_error
		log.WithFields(logrus.Fields{"token": notification.PushToken, "status": archive.ExternalStatus}).Error("Ios notification failure")
	} else {
		archive.Status = pb.NotificationStatusCode_success
	}

	if response.StatusCode == iosStatuscodeInvalidtoken {
		ctx := context.Background()
		_, err := dbClient.DeletePushToken(ctx, &pb.DeletePushTokenRequest{
			Base:      &pb.BaseRequest{UpdateBy: ServiceName},
			UserId:    notification.UserId,
			PushToken: notification.PushToken})

		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{"user_id": notification.UserId, "token": notification.PushToken}).Error("Error deleting push token")
		}
	}

	return archive
}

func sendAndroid(notification *pb.Notification) *pb.NotificationArchive {
	log := helpers.GetDefaultLog(ServiceName, "Sender")
	response := sendAndroidNotification(notification)

	archive := &pb.NotificationArchive{}
	archive.Id = notification.Id

	if response == nil {
		archive.Status = pb.NotificationStatusCode_error
		archive.InternalError = "Missing response object"
		return archive
	}

	if response.Fail == 0 && response.Success > 0 {
		archive.Status = pb.NotificationStatusCode_success
	} else {
		archive.Status = pb.NotificationStatusCode_error
	}
	archive.ExternalStatus = string(response.StatusCode)
	archive.ExternalError = ResultsString(response) + " " + response.Err

	if IsNotRegisteredToken(response) {
		ctx := context.Background()
		_, err := dbClient.DeletePushToken(ctx, &pb.DeletePushTokenRequest{
			Base:      &pb.BaseRequest{UpdateBy: ServiceName},
			UserId:    notification.UserId,
			PushToken: notification.PushToken})

		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{"user_id": notification.UserId, "token": notification.PushToken}).Error("Error deleting push token")
		}
	}

	return archive
}

func sendMail(notification *pb.Notification) *pb.NotificationArchive {
	log := helpers.GetDefaultLog(ServiceName, "Sender")
	ctx := context.Background()
	_, err := mailClient.SendMail(ctx, &pb.SendMailRequest{
		Base:    &pb.BaseRequest{UpdateBy: ServiceName},
		From:    cnf.EmailSender,
		To:      notification.UserEmail,
		Subject: notification.MailSubject,
		Body:    notification.Content,
		IsHtml:  isHTML(notification.MailType)})

	archive := &pb.NotificationArchive{}
	archive.Id = notification.Id
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"to": notification.UserEmail}).Error("Could not send email")
		archive.Status = pb.NotificationStatusCode_error
		archive.ExternalError = err.Error()
	} else {
		archive.Status = pb.NotificationStatusCode_success
	}

	return archive
}

func isHTML(emailContentType pb.EmailContentType) bool {
	if emailContentType == pb.EmailContentType_html {
		return true
	}

	return false
}
