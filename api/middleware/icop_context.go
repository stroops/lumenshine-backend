package middleware

import "github.com/gin-gonic/gin"

//IcopContextMiddleware general middleware for all endpint services
//we store global values in here and reread them in the explizit middleware in order to store them in typed structs
type IcopContextMiddleware struct {
	ServiceName string
}

//MiddlewareFunc handler func for the middleware
func (mw *IcopContextMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("servicename", mw.ServiceName)
		c.Next()
	}
}
