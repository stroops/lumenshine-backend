package middleware

import (
	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//IcopContext context used in the apis to store some default values
type IcopContext struct {
	RequestID string
	Language  string
	Log       *logrus.Entry
}

//UseIcopContext is used to wrap the gin context
func UseIcopContext(f func(uc *IcopContext, c *gin.Context)) gin.HandlerFunc {
	uc := &IcopContext{}

	return func(c *gin.Context) {
		uc.Language = c.GetString("language")
		uc.RequestID = c.GetString("request_id")
		uc.Log = helpers.GetDefaultLog(c.GetString("servicename"), uc.RequestID)
		f(uc, c)
	}
}
