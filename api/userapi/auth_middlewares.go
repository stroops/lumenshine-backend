package main

import (
	"github.com/Soneso/lumenshine-backend/pb"
	"strconv"
	//jwt "github.com/appleboy/gin-jwt"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

var authMiddlewareLostPwd *mw.IcopJWTMiddleware
var authMiddlewareSimple *mw.IcopJWTMiddleware
var authMiddlewareFull *mw.IcopJWTMiddleware

func init() {
	authMiddlewareLostPwd = &mw.IcopJWTMiddleware{
		ServiceName: "usersvc",
		AuthDBKey:   "lostpwd",
		JwtClient:   getJwtClient,
		DbClient:    getDbClient,
	}

	authMiddlewareSimple = &mw.IcopJWTMiddleware{
		ServiceName: "usersvc",
		AuthDBKey:   "simple",
		JwtClient:   getJwtClient,
		DbClient:    getDbClient,
	}

	authMiddlewareFull = &mw.IcopJWTMiddleware{
		ServiceName: "usersvc",
		AuthDBKey:   "full",
		JwtClient:   getJwtClient,
		DbClient:    getDbClient,
		PayloadFunc: func(userID string) map[string]interface{} {
			// We set the full_authorized flag, that will be checked in the Authorizator method
			return map[string]interface{}{"full_authorized": true}
		},
		Authorizator: func(userID string, c *gin.Context) bool {
			id, err := strconv.ParseInt(userID, 10, 64)
			if err == nil {
				claims := mw.ExtractClaims(c)
				isFullAuthorized, ok := claims["full_authorized"]
				if isFullAuthorized.(bool) && ok {
					if authMiddlewareFull.SetAuthUserData(c, id) {
						return true
					}
				}
			}
			return false
		},
	}
}

func getJwtClient() pb.JwtServiceClient {
	return jwtClient
}

func getDbClient() pb.DBServiceClient {
	return dbClient
}
