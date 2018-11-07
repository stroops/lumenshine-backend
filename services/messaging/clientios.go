package main

import (
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

var client *apns2.Client

const (
	iosStatuscodeSuccess      = 200
	iosStatuscodeInvalidtoken = 410
)

//InitIosClient - initializes the ios notification client
func InitIosClient() {
	log := helpers.GetDefaultLog(ServiceName, "Sender")

	if cnf.IsDevSystem {
		cert, err := certificate.FromP12File(cnf.IOSConfig.DevelopmentCertificate, cnf.IOSConfig.DevelopmentCertificatePassword)
		if err != nil {
			log.WithError(err).Fatalf("Development certificate error: %v", err)
		}
		client = apns2.NewClient(cert).Development()
	} else {
		cert, err := certificate.FromP12File(cnf.IOSConfig.ProductionCertificate, cnf.IOSConfig.ProductionCertificatePassword)
		if err != nil {
			log.WithError(err).Fatalf("Production certificate error: %v", err)
		}
		client = apns2.NewClient(cert).Production()
	}
}

func sendIosNotification(request *pb.Notification) *apns2.Response {
	log := helpers.GetDefaultLog(ServiceName, "Sender")

	notification := &apns2.Notification{}
	notification.DeviceToken = request.PushToken
	notification.Topic = cnf.IOSConfig.BundleID
	notification.Payload = []byte(request.Content)

	//fmt.Printf("Payload: %v", request.Content)

	response, err := client.Push(notification)

	if err != nil {
		log.WithError(err).Errorf("Ios push notification error: %v", err)
	}

	// if response != nil {
	// 	fmt.Printf("Code: %v Id: %v Reason: %v\n", strconv.Itoa(response.StatusCode), response.ApnsID, response.Reason)
	// }

	return response
}
