package middleware

// source comes from https://github.com/appleboy/gin-jwt
// with some modification for icop

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// IcopJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type IcopJWTMiddleware struct {
	ServiceName string
	JwtClient   func() pb.JwtServiceClient //the jwtClient is used for getting the current jwts
	DbClient    func() pb.DBServiceClient  //the dbClient is used for getting the current userData

	// Key name in database
	AuthDBKey string

	// used to backup the current key, in case, memcached is not working
	backupKey1, backupKey2 string

	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userID string, password string, c *gin.Context) (string, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, c *gin.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userID string) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// Name of the header key. Default value "Authorization".
	TokenLookupName string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time
}

//AuthUser is the userdata that is stored in every request
type AuthUser struct {
	UserID            int64
	MailConfirmed     bool
	MnemonicConfirmed bool
	TfaConfirmed      bool
	TfaSecret         string
	Email             string
	MessageCount      int
	PublicKey0        string
}

// Login form structure.
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//SetAuthHeader sets the token into the gin header
func (mw *IcopJWTMiddleware) SetAuthHeader(c *gin.Context, userID int64) {
	token := mw.TokenHeadName + " " + mw.TokenGenerator(c, fmt.Sprintf("%d", userID))
	c.Header(mw.TokenLookupName, token)
}

//getAuthKeys returns the current auth keys from memcached
func (mw *IcopJWTMiddleware) getAuthKeys(c *gin.Context) ([]byte, []byte) {
	if mw.JwtClient != nil {
		//cc := context.Background()
		keys, err := mw.JwtClient().GetJwtValue(c, &pb.KeyRequest{
			Base: &pb.BaseRequest{RequestId: c.GetString("request_id"), UpdateBy: c.GetString("servicename")},
			Key:  mw.AuthDBKey,
		})

		if err == nil {
			//need to save the keys for backup
			mw.backupKey1 = keys.Key1
			mw.backupKey2 = keys.Key2
		}
	}

	return []byte(mw.backupKey1), []byte(mw.backupKey2)
}

// MiddlewareInit initialize jwt configs.
func (mw *IcopJWTMiddleware) MiddlewareInit() error {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.TokenLookupName == "" {
		mw.TokenLookupName = "Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(userID string, c *gin.Context) bool {
			id, err := strconv.ParseInt(userID, 10, 64)
			if err == nil {
				if mw.SetAuthUserData(c, id) {
					return true
				}
			}
			return false
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":           code,
				"error_message":  message,
				"parameter_name": "jwt-token",
				"error_code":     cerr.JwtError,
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) string {
			return claims["id"].(string)
		}
	}

	if mw.Realm == "" {
		mw.Realm = "ICOP"
	}

	return nil
}

// MiddlewareFunc makes IcopJWTMiddleware implement the Middleware interface.
func (mw *IcopJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *IcopJWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	c.Set("JWT_PAYLOAD", claims)
	c.Set("userID", id)

	if !mw.Authorizator(id, c) {
		mw.unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return
	}

	c.Next()
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the IcopJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *IcopJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	//generate new token for user
	userID := claims["id"].(string)
	newToken := mw.TokenHeadName + " " + mw.TokenGenerator(c, userID)
	c.Header(mw.TokenLookupName, newToken)

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}

// TokenGenerator handler that clients can use to get a jwt token.
func (mw *IcopJWTMiddleware) TokenGenerator(c *gin.Context, userID string) string {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	claims["id"] = userID

	_, key2 := mw.getAuthKeys(c)
	tokenString, _ := token.SignedString(key2)

	return tokenString
}

func (mw *IcopJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (mw *IcopJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		return "", errors.New("Query token empty")
	}

	return token, nil
}

func (mw *IcopJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)

	if cookie == "" {
		return "", errors.New("Cookie token empty")
	}

	return cookie, nil
}

func (mw *IcopJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = mw.jwtFromHeader(c, parts[1])
	case "query":
		token, err = mw.jwtFromQuery(c, parts[1])
	case "cookie":
		token, err = mw.jwtFromCookie(c, parts[1])
	}

	if err != nil {
		return nil, err
	}

	key1, key2 := mw.getAuthKeys(c)
	t, e := jwt.Parse(token, mw.intParse(key1)) //try first key
	if e != nil {
		t, e = jwt.Parse(token, mw.intParse(key2)) //try second key
	}
	return t, e
}

func (mw *IcopJWTMiddleware) intParse(k []byte) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return k, nil
	}
}

func (mw *IcopJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {

	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)

	return
}

//SetAuthUserData general functions for setting the user in the middleware
func (mw *IcopJWTMiddleware) SetAuthUserData(c *gin.Context, userID int64) bool {
	l := helpers.GetDefaultLog(mw.ServiceName, c.GetString("RequestId"))

	//get the user-details from the DB
	req := &pb.GetUserByIDOrEmailRequest{
		Base: &pb.BaseRequest{RequestId: c.GetString("request_id"), UpdateBy: c.GetString("servicename")},
		Id:   userID,
	}
	user, err := mw.DbClient().GetUserDetails(c, req)
	if err != nil {
		l.WithError(err).WithField("user-id", userID).Error("Error reading user")
		c.Set("user", &AuthUser{}) //on error we set an empty user
		return false
	}

	if user.UserNotFound {
		l.WithField("user-id", userID).Error("User from jwt could not be found in db")
		c.Set("user", &AuthUser{}) //on error we set an empty user
		return false
	}

	if user.IsClosed {
		l.WithField("user-id", userID).Error("User from jwt is closed")
		c.Set("user", &AuthUser{}) //on error we set an empty user
		return false
	}

	if user.IsSuspended {
		l.WithField("user-id", userID).Error("User from jwt is suspended")
		c.Set("user", &AuthUser{}) //on error we set an empty user
		return false
	}

	c.Set("user", &AuthUser{
		UserID:            userID,
		MailConfirmed:     user.MailConfirmed,
		MnemonicConfirmed: user.MnemonicConfirmed,
		TfaConfirmed:      user.TfaConfirmed,
		TfaSecret:         user.TfaSecret,
		Email:             user.Email,
		MessageCount:      int(user.MessageCount),
		PublicKey0:        user.PublicKey_0,
	})
	return true
}

//GetAuthUser returns the stored authUser, or an empty one
func GetAuthUser(c *gin.Context) *AuthUser {
	user, exists := c.Get("user")
	if exists {
		return user.(*AuthUser)
	}
	return &AuthUser{}
}
