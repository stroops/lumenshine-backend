package main

import (
	"errors"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

//QueuePushNotification - enqueues a push notification
func (s *server) QueuePushNotification(ctx context.Context, r *pb.QueuePushNotificationRequest) (*pb.Empty, error) {
	if r.PushToken == "" {
		return nil, errors.New("need push token")
	}
	if r.DeviceType == pb.NotificationType_mail {
		return nil, errors.New("invalid device type")
	}

	u, err := models.UserProfiles(qm.Where(
		models.UserProfileColumns.ID+"=?", r.UserId,
	)).One(db)

	if err != nil {
		return nil, err
	}

	var n models.Notification
	n.UserID = u.ID
	n.PushToken = r.PushToken
	n.Type = r.DeviceType.String()
	n.Content = r.Content
	n.MailType = models.MailContentTypeText
	n.UpdatedBy = r.Base.UpdateBy

	err = n.Insert(db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

//QueueMailNotification - enqueues a mail notification
func (s *server) QueueMailNotification(ctx context.Context, r *pb.QueueMailNotificationRequest) (*pb.Empty, error) {
	if r.Content == "" && r.Subject == "" {
		return nil, errors.New("need either content or subject")
	}

	u, err := models.UserProfiles(qm.Where(
		models.UserProfileColumns.ID+"=?", r.UserId,
	)).One(db)

	if err != nil {
		return nil, err
	}

	var n models.Notification
	n.UserID = u.ID
	n.Type = pb.NotificationType_mail.String()
	n.MailSubject = r.Subject
	n.Content = r.Content
	n.MailType = r.MailType.String()
	n.UserEmail = r.UserEmail
	n.UpdatedBy = r.Base.UpdateBy

	err = n.Insert(db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

//UpdateNotificationsStatus - updates the sent status of the notifications
func (s *server) UpdateNotificationsStatus(ctx context.Context, r *pb.UpdateNotificationsStatusRequest) (*pb.Empty, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	var err error
	var n *models.NotificationArchive

	for _, rNotification := range r.Notifications {
		n, err = models.NotificationArchives(qm.Where(
			models.NotificationArchiveColumns.ID+"=?", rNotification.Id,
		)).One(db)

		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{"ID": rNotification.Id}).
				Error("Error reading archived notification")

			continue
		}

		if n == nil {
			log.WithError(err).
				WithFields(logrus.Fields{"ID": rNotification.Id}).
				Error("Could not find archived notification")

			continue
		}

		n.Status = rNotification.Status.String()
		n.InternalErrorString = rNotification.InternalError
		n.ExternalStatusCode = rNotification.ExternalStatus
		n.ExternalErrorString = rNotification.ExternalError
		n.UpdatedBy = r.Base.UpdateBy
		n.UpdatedAt = time.Now()

		_, err = n.Update(db, boil.Whitelist(
			models.NotificationArchiveColumns.Status,
			models.NotificationArchiveColumns.InternalErrorString,
			models.NotificationArchiveColumns.ExternalStatusCode,
			models.NotificationArchiveColumns.ExternalErrorString,
			models.NotificationArchiveColumns.UpdatedAt,
			models.NotificationArchiveColumns.UpdatedBy))

		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{"ID": n.ID, "userID": n.UserID, "content": n.Content, "status": n.Status}).
				Error("Error archiving notification")
		}
	}

	return &pb.Empty{}, nil
}

//DequeueNotifications - deques oldest notifications
func (s *server) DequeueNotifications(ctx context.Context, r *pb.DequeueRequest) (*pb.NotificationListResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	dbNotifications, err := models.Notifications(qm.OrderBy(models.NotificationColumns.ID), qm.Limit(int(r.LimitCount))).All(db)

	if err != nil {
		log.WithError(err).Error("Error reading notifications")
		return nil, err
	}

	var notifications []*pb.Notification
	var notification *pb.Notification
	var nArchive *models.NotificationArchive
	for _, rNotification := range dbNotifications {
		tx, err := db.Begin()
		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{"ID": rNotification.ID}).
				Error("Could not open transaction")

			continue
		}

		nArchive = &models.NotificationArchive{}
		nArchive.ID = rNotification.ID
		nArchive.UserID = rNotification.UserID
		nArchive.PushToken = rNotification.PushToken
		nArchive.Type = rNotification.Type
		nArchive.Content = rNotification.Content
		nArchive.MailSubject = rNotification.MailSubject
		nArchive.MailType = rNotification.MailType
		nArchive.UserEmail = rNotification.UserEmail
		nArchive.Status = models.NotificationStatusCodeNew
		nArchive.UpdatedBy = r.Base.UpdateBy

		err = nArchive.Insert(tx, boil.Infer())
		if err != nil {
			tx.Rollback()

			log.WithError(err).
				WithFields(logrus.Fields{"ID": rNotification.ID}).
				Error("Could not insert archived notification")

			continue
		}

		_, err = rNotification.Delete(tx)
		if err != nil {
			tx.Rollback()

			log.WithError(err).
				WithFields(logrus.Fields{"ID": rNotification.ID}).
				Error("Could not delete notification")

			continue
		}
		tx.Commit()

		notification = &pb.Notification{}
		notification.Id = int64(nArchive.ID)
		notification.UserId = int64(nArchive.UserID)
		notification.PushToken = nArchive.PushToken
		notification.NotificationType = pb.NotificationType(pb.NotificationType_value[nArchive.Type])
		notification.Content = nArchive.Content
		notification.MailSubject = nArchive.MailSubject
		notification.MailType = pb.EmailContentType(pb.EmailContentType_value[nArchive.MailType])
		notification.UserEmail = nArchive.UserEmail

		notifications = append(notifications, notification)
	}

	return &pb.NotificationListResponse{Notifications: notifications}, nil
}
