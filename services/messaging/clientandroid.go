package main

import (
	"encoding/json"
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"strings"

	fcm "github.com/NaySoftware/go-fcm"
)

const (
	androidInvalidtoken = "NotRegistered"
	androidErrorKey     = "error"
)

//GooglePayload - android payload
type GooglePayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func sendAndroidNotification(request *pb.Notification) *fcm.FcmResponseStatus {
	log := helpers.GetDefaultLog(ServiceName, "Sender")

	var payload GooglePayload
	err := json.Unmarshal([]byte(request.Content), &payload)

	if err != nil {
		log.WithError(err).Errorf("Error unmarshaling android payload json: %v", err)
		return nil
	}

	apiKey := cnf.AndroidConfig.APIKey
	c := fcm.NewFcmClient(apiKey)
	c.NewFcmMsgTo(request.PushToken, payload)

	response, err := c.Send()
	if err != nil {
		log.WithError(err).Errorf("Android push notification error: %v", err)
	}

	return response
}

// IsNotRegisteredToken check whether the push token is no longer registered
func IsNotRegisteredToken(response *fcm.FcmResponseStatus) bool {
	if response.StatusCode == 200 {
		for _, val := range response.Results {
			for k, v := range val {
				if k == androidErrorKey && strings.EqualFold(v, androidInvalidtoken) {
					return true
				}
			}
		}
	}

	return false
}

//ResultsString - returns results as a concatenated string
func ResultsString(response *fcm.FcmResponseStatus) string {
	result := ""
	for _, val := range response.Results {
		for k, v := range val {
			result += fmt.Sprintf("%v : %v ", k, v)
		}
	}
	return result
}
