package icop_error

import (
	"github.com/sirupsen/logrus"
)

//LogAndReturnError logs an error to stdout and returns a struct with the code and passed message
func LogAndReturnError(log *logrus.Entry, err error, errorMessage string, errorCode int) IcopErrors {
	log.WithError(err).WithFields(logrus.Fields{"code": errorCode, "message": errorMessage}).Error()
	return NewIcopErrorShort(errorCode, errorMessage)
}

//LogAndReturnIcopError logs an error to stdout and returns a struct with the code and passed message
func LogAndReturnIcopError(log *logrus.Entry, paramName string, errorCode int, errorMessage string, userErrorMessageKey string) IcopErrors {
	log.WithFields(logrus.Fields{"code": errorCode, "message": errorMessage, "param": paramName, "messageKey": userErrorMessageKey}).Error()
	return NewIcopError(paramName, errorCode, errorMessage, userErrorMessageKey)
}
