package constants

import (
	"golang.org/x/text/language"
)

//defining some global constants, used in all services
const (
	DefaultMailkeyExpiryDays = 14
	DefaultMailkeyLength     = 16
)

var (
	//ServerLanguages represents all server languages
	//this is used e.g. for the templates or serverLangs in api.middleware..language.go
	ServerLanguages = []language.Tag{
		language.English, // en fallback
		language.German,  // de
	}
)
