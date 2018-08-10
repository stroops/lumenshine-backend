package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//MessageCount sets the users message count in the header
func MessageCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := GetAuthUser(c)

		if u.MessageCount > 0 {
			c.Writer.Header().Set("X-MessageCount", fmt.Sprintf("%d", u.MessageCount))
		}
		c.Next()
	}
}
