package api

import (
	"github.com/sirupsen/logrus"

	"github.com/Soneso/lumenshine-backend/admin/db"

	"github.com/gin-gonic/gin"

	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/helpers"

	"golang.org/x/crypto/bcrypt"
)

//LoginFunc is used to wrap the gin function
//this is just an example on how to implement such functions, if the AdminContext is needed inside the api
//we could also have used a default gin handler and construct e.g. the log from the gin.Request, which
//also holds the RequestID and ServiceName
func LoginFunc(f func(email string, password string, uc *mw.AdminContext, c *gin.Context) (string, bool)) func(string, string, *gin.Context) (string, bool) {
	uc := &mw.AdminContext{}
	return func(email string, password string, c *gin.Context) (string, bool) {
		uc.Language = c.GetString("language")
		uc.RequestID = c.GetString("request_id")
		uc.Log = helpers.GetDefaultLog(c.GetString("servicename"), uc.RequestID)
		return f(email, password, uc, c)
	}
}

//Login logs in the user
// swagger:route POST /portal/admin/login auth Login
//
// Logs in the user
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func Login(email string, password string, uc *mw.AdminContext, c *gin.Context) (string, bool) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		uc.Log.WithError(err).WithField("email", email).Error("Could not get user by email")
		return email, false
	}
	if user == nil {
		uc.Log.WithField("email", email).Info("User is null")
		return email, false
	}
	if !user.Active {
		uc.Log.WithField("email", email).Info("User not active")
		return email, false
	}
	match := checkPasswordHash(uc.Log, password, user.Password)
	if !match {
		return email, false
	}

	if err := db.UpdateLastLogin(*user); err != nil {
		uc.Log.WithError(err).WithField("email", email).Error("Could not update last login")
		return email, false
	}
	return email, true
}

//checkPasswordHash check a given password to the hashed value
func checkPasswordHash(log *logrus.Entry, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.WithError(err).Warn("Error checking password")
	}
	return err == nil
}
