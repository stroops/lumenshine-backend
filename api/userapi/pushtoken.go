package main

import (
	"net/http"
	"strings"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

//SubscribeForPushNotificationsRequest request
//swagger:parameters SubscribeForPushNotificationsRequest SubscribeForPushNotifications
type SubscribeForPushNotificationsRequest struct {
	//required: true
	PushToken string `form:"push_token" json:"push_token"  validate:"required,max=500"`
	//Device type, e.g. apple, google
	//required: true
	DeviceType string `form:"device_type" json:"device_type" validate:"required,max=50,icop_devicetype"`
}

//SubscribeForPushNotifications adds a new token to the user
// swagger:route POST /portal/user/dashboard/subscribe_push_token pushtoken SubscribeForPushNotifications
//
// Adds a new token to the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func SubscribeForPushNotifications(uc *mw.IcopContext, c *gin.Context) {
	var l SubscribeForPushNotificationsRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	reqData := &pb.AddPushTokenRequest{
		Base:       NewBaseRequest(uc),
		UserId:     userID,
		PushToken:  l.PushToken,
		DeviceType: pb.DeviceType(pb.DeviceType_value[strings.ToLower(l.DeviceType)]),
	}
	_, err := dbClient.AddPushToken(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding push token", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//UpdatePushTokenRequest request
//swagger:parameters UpdatePushTokenRequest UpdatePushToken
type UpdatePushTokenRequest struct {
	//required: true
	NewPushToken string `form:"new_push_token" json:"new_push_token"  validate:"required,max=500"`
	//required: true
	OldPushToken string `form:"old_push_token" json:"old_push_token"  validate:"required,max=500"`
	//Device type, e.g. apple, google
	//required: true
	DeviceType string `form:"device_type" json:"device_type" validate:"required,max=50,icop_devicetype"`
}

//UpdatePushToken updates the push token
// swagger:route POST /portal/user/dashboard/update_push_token pushtoken UpdatePushToken
//
// Updates the push token
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdatePushToken(uc *mw.IcopContext, c *gin.Context) {
	var l UpdatePushTokenRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	reqData := &pb.UpdatePushTokenRequest{
		Base:         NewBaseRequest(uc),
		UserId:       userID,
		NewPushToken: l.NewPushToken,
		OldPushToken: l.OldPushToken,
		DeviceType:   pb.DeviceType(pb.DeviceType_value[strings.ToLower(l.DeviceType)]),
	}
	_, err := dbClient.UpdatePushToken(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error replacing push token", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//UnsubscribeFromPushNotificationsRequest request
//swagger:parameters UnsubscribeFromPushNotificationsRequest UnsubscribeFromPushNotifications
type UnsubscribeFromPushNotificationsRequest struct {
	//required: true
	PushToken string `form:"push_token" json:"push_token" validate:"required,max=500"`
}

//UnsubscribeFromPushNotifications removes the specified token from the user
// swagger:route POST /portal/user/dashboard/unsubscribe_push_token pushtoken UnsubscribeFromPushNotifications
//
// Removes the specified token from the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UnsubscribeFromPushNotifications(uc *mw.IcopContext, c *gin.Context) {
	var l UnsubscribeFromPushNotificationsRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	reqData := &pb.DeletePushTokenRequest{
		Base:      NewBaseRequest(uc),
		UserId:    userID,
		PushToken: l.PushToken,
	}
	_, err := dbClient.DeletePushToken(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting push token", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//UnsubscribePreviousUserFromPushNotificationsRequest request
//swagger:parameters UnsubscribePreviousUserFromPushNotificationsRequest UnsubscribePreviousUserFromPushNotifications
type UnsubscribePreviousUserFromPushNotificationsRequest struct {
	//required: true
	Email string `form:"user_email" json:"user_email" validate:"required,icop_email"`
	//required: true
	PushToken string `form:"push_token" json:"push_token"  validate:"required,max=500"`
}

//UnsubscribePreviousUserFromPushNotifications removes the specified token from the user
// swagger:route POST /portal/user/dashboard/unsubscribe_previous_user_push_token pushtoken UnsubscribePreviousUserFromPushNotifications
//
// Removes the specified token from the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UnsubscribePreviousUserFromPushNotifications(uc *mw.IcopContext, c *gin.Context) {
	var l UnsubscribePreviousUserFromPushNotificationsRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userRequest := &pb.GetUserByIDOrEmailRequest{Email: l.Email}
	userResponse, err := dbClient.GetUserDetails(c, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading user", cerr.GeneralError))
		return
	}
	if userResponse.UserNotFound {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("user_email", cerr.InvalidArgument, "User email not found", ""))
		return
	}

	reqData := &pb.DeletePushTokenRequest{
		Base:      NewBaseRequest(uc),
		UserId:    userResponse.Id,
		PushToken: l.PushToken,
	}
	_, err = dbClient.DeletePushToken(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting push token", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//TestPushNotifications - for testing
func TestPushNotifications(uc *mw.IcopContext, c *gin.Context) {
	publicKey := c.Param("publickey")
	req := &pb.GetWalletByPublicKeyRequest{
		Base:      NewBaseRequest(uc),
		PublicKey: publicKey,
	}
	wallet, err := dbClient.GetWalletByPublicKey(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading wallet", cerr.GeneralError))
		return
	}
	pushReq := &pb.PushNotificationRequest{
		Base:                       NewBaseRequest(uc),
		UserID:                     wallet.UserId,
		Title:                      "Test title",
		Message:                    "Test message",
		SendAsMailIfNoTokenPresent: true,
	}

	_, err = notificationClient.SendPushNotification(c, pushReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error sending push notification", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{Success : true}")
}
