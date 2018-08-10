package main

import (
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"
	"net/http"
	"strings"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

//SubscribeForPushNotificationsRequest request
type SubscribeForPushNotificationsRequest struct {
	PushToken  string `form:"push_token" json:"push_token"  validate:"required,max=500"`
	DeviceType string `form:"wallet_name" json:"wallet_name" validate:"required,max=50,icop_devicetype"`
}

//SubscribeForPushNotifications adds a new token to the user
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
type UpdatePushTokenRequest struct {
	NewPushToken string `form:"new_push_token" json:"new_push_token"  validate:"required,max=500"`
	OldPushToken string `form:"old_push_token" json:"old_push_token"  validate:"required,max=500"`
	DeviceType   string `form:"wallet_name" json:"wallet_name" validate:"required,max=50,icop_devicetype"`
}

//UpdatePushToken adds a new token to the user
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
type UnsubscribeFromPushNotificationsRequest struct {
	PushToken string `form:"push_token" json:"push_token"  validate:"required,max=500"`
}

//UnsubscribeFromPushNotifications removes the specified token from the user
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
type UnsubscribePreviousUserFromPushNotificationsRequest struct {
	Email     string `form:"user_email" json:"user_email" validate:"required,icop_email"`
	PushToken string `form:"push_token" json:"push_token"  validate:"required,max=500"`
}

//UnsubscribePreviousUserFromPushNotifications removes the specified token from the user
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
