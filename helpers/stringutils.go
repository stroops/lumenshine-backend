package helpers

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

//RandomString generates a random string of len count
func RandomString(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//TimeToString converts a time to a given language time-string
func TimeToString(t time.Time, lang string) string {
	switch lang {
	case "de":
		return t.Format("02.01.2006 15:04:05")
	default: //used also for en
		return t.Format("2006/01/02 15:04:05")
	}
}
