package helpers

import "github.com/sirupsen/logrus"

//GetDefaultLog returns a default Log Entry to be used
func GetDefaultLog(servicename string, requestID string) *logrus.Entry {
	l := logrus.WithFields(logrus.Fields{"request_id": requestID, "srv": servicename})
	return l
}
