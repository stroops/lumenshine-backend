package helpers

import (
	"net/http"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//ValidateRequestData validates the request data and writes to the response, if not ok
func ValidateRequestData(destStruct interface{}, log *logrus.Entry, c *gin.Context) bool {
	if err := c.Bind(destStruct); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(log, err, cerr.ValidBadInputData, cerr.BindError))
		return false
	}

	if valid, validErrors := cerr.ValidateStruct(log, destStruct); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return false
	}

	return true
}
