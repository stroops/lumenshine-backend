package middleware

import (
	"github.com/Soneso/lumenshine-backend/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

var serverLangs = language.NewMatcher(constants.ServerLanguages)

//Language sets the language inside the gin context, based on the passed in Accept-[User-]Language header
//if nothing specified, english will be used
func Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check if header is present, use english if not
		userLang := c.Request.Header.Get("Accept-User-Language")
		accept := c.Request.Header.Get("Accept-Language")
		if accept == "" {
			accept = language.English.String()
		}

		tag, _ := language.MatchStrings(serverLangs, userLang, accept)
		lang := tag.String()
		if lang == "" {
			lang = language.English.String()
		}
		c.Set("language", tag.String())

		c.Next()
	}
}
