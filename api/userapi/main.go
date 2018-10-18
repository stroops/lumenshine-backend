package main

import (
	"fmt"
	"runtime"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"

	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/api/userapi/cmd"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	//ServiceName name of this service
	ServiceName = "usersvc"
)

var (
	twoFAClient    pb.TwoFactorAuthServiceClient
	dbClient       pb.DBServiceClient
	jwtClient      pb.JwtServiceClient
	mailClient     pb.MailServiceClient
	adminAPIClient pb.AdminApiServiceClient
)

func main() {
	var err error
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := helpers.GetDefaultLog(ServiceName, "Startup")

	cmd := cmd.RootCommand()
	if err = cmd.Execute(); err != nil {
		log.WithError(err).Fatalf("Error reading root command")
	}

	if err = readConfig(log, cmd); err != nil {
		log.WithError(err).Fatalf("Error reading config")
	}

	logLevel, err := logrus.ParseLevel(cnf.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	connectServices(log)
	initBoxes()
	loadTemplates(log)

	r := gin.New()

	//setup gin-logging. This should always be the first middleware
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	r.Use(mw.GinLogger(logger))

	r.Use(gin.Recovery())
	r.Use(mw.RequestID())
	r.Use(mw.Language())

	ic := &mw.IcopContextMiddleware{ServiceName: ServiceName}
	r.Use(ic.MiddlewareFunc())

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: cnf.CORSHosts,
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Credentials",
			"Cache-Control", "Accept-Language", "Accept-User-Language", "X-Request-Id"},
		ExposeHeaders:    []string{"Authorization", "X-Request-Id", "X-MessageCount"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/portal/test", Test)
	r.POST("/portal/user/register_user", mw.UseIcopContext(RegisterUser))
	r.GET("/portal/user/salutation_list", mw.UseIcopContext(SalutationList))
	r.GET("/portal/user/language_list", mw.UseIcopContext(LanguageList))
	r.GET("/portal/user/country_list", mw.UseIcopContext(CountryList))
	r.GET("/portal/user/occupation_list", mw.UseIcopContext(OccupationList))
	r.POST("/portal/user/confirm_token", mw.UseIcopContext(ConfirmToken))
	r.POST("/portal/user/resend_confirmation_mail", mw.UseIcopContext(ResendConfirmMail))
	r.GET("/portal/user/get_user_details", mw.UseIcopContext(GetUserRegistrationDetails)) //called without jwt, so email and 2fa is present
	r.POST("/portal/user/login_step1", mw.UseIcopContext(LoginStep1))
	r.POST("/portal/user/lost_password", mw.UseIcopContext(LostPassword))
	r.POST("/portal/user/lost_tfa", mw.UseIcopContext(LostTfa))
	r.GET("/portal/info", Info)

	//this group is used, with the simple authenticator, which means, only the userID is present
	//the middleware will not check for full logged in
	auth := r.Group("/portal/user/auth")
	auth.Use(authMiddlewareSimple.MiddlewareFunc())
	{
		auth.POST("/refresh", authMiddlewareSimple.RefreshHandler)
		auth.POST("/confirm_tfa_registration", mw.UseIcopContext(Confirm2FA))
		auth.POST("/login_step2", mw.UseIcopContext(LoginStep2))
		auth.POST("/lost_password_tfa", mw.UseIcopContext(LostPasswordTfa))
		auth.POST("/new_2fa_secret", mw.UseIcopContext(NewTfaUpdate))
		auth.GET("/user_auth_data", mw.UseIcopContext(UserSecurityData))
		auth.POST("/update_security_data", mw.UseIcopContext(UpdateSecurityData))
		auth.GET("/get_user_registration_status", mw.UseIcopContext(GetUserRegistrationDetails))
		auth.GET("/need_2fa_reset_pwd", mw.UseIcopContext(Need2FAResetPassword))
	}

	//this group is used, with the full authenticator, which means, userID and claim is present
	//the middleware will check for full logged in (calim must be present)
	//this is the main group, that is used to read data etc.
	authDash := r.Group("/portal/user/dashboard")
	authDash.Use(authMiddlewareFull.MiddlewareFunc())
	authDash.Use(mw.MessageCount()) //Messagecount only for fully logged in users
	{
		authDash.GET("/get_user_data", mw.UseIcopContext(GetUserData))
		authDash.POST("/update_user_data", mw.UseIcopContext(UpdateUserData))

		authDash.POST("/refresh", authMiddlewareFull.RefreshHandler)
		authDash.POST("/confirm_tfa_registration", mw.UseIcopContext(Confirm2FA))
		authDash.POST("/change_password", mw.UseIcopContext(ChangePasswordUpdate))
		authDash.GET("/user_auth_data", mw.UseIcopContext(UserSecurityData))
		authDash.GET("/get_user_registration_status", mw.UseIcopContext(GetUserRegistrationDetails)) //called with jwt, so userId is present
		authDash.POST("/confirm_mnemonic", mw.UseIcopContext(ConfirmMnemonic))
		authDash.POST("/new_2fa_secret", mw.UseIcopContext(NewTfaUpdate))
		authDash.POST("/confirm_new_2fa_secret", mw.UseIcopContext(Confirm2FA))
		authDash.POST("/tfa_secret", mw.UseIcopContext(GetTfaSecret))
		authDash.GET("/get_user_messages_list", mw.UseIcopContext(GetUserMessageList))
		authDash.GET("/get_user_message", mw.UseIcopContext(GetUserMessage))

		authDash.GET("/get_user_wallets", mw.UseIcopContext(GetUserWallets))
		authDash.POST("/add_wallet", mw.UseIcopContext(AddWallet))
		authDash.POST("/remove_wallet", mw.UseIcopContext(RemoveWallet))
		authDash.POST("/change_wallet_data", mw.UseIcopContext(WalletChangeData))
		authDash.POST("/remove_wallet_federation_address", mw.UseIcopContext(RemoveWalletFederationAddress))
		authDash.POST("/wallet_set_homescreen", mw.UseIcopContext(WalletSetHomescreen))

		authDash.POST("/subscribe_push_token", mw.UseIcopContext(SubscribeForPushNotifications))
		authDash.POST("/unsubscribe_push_token", mw.UseIcopContext(UnsubscribeFromPushNotifications))
		authDash.POST("/unsubscribe_previous_user_push_token", mw.UseIcopContext(UnsubscribePreviousUserFromPushNotifications))
		authDash.POST("/update_push_token", mw.UseIcopContext(UpdatePushToken))

		authDash.POST("/get_known_currency", mw.UseIcopContext(GetKnownCurrency))
		authDash.POST("/get_known_currencies", mw.UseIcopContext(GetKnownCurrencies))

		authDash.POST("/get_known_inflation_destination", mw.UseIcopContext(GetKnownInflationDestination))
		authDash.POST("/get_known_inflation_destinations", mw.UseIcopContext(GetKnownInflationDestinations))

		authDash.POST("/upload_kyc_document", mw.UseIcopContext(UploadKycDocument))

		authDash.GET("/contact_list", mw.UseIcopContext(ContactList))
		authDash.POST("/add_contact", mw.UseIcopContext(AddContact))
		authDash.POST("/edit_contact", mw.UseIcopContext(EditContact))
		authDash.POST("/remove_contact", mw.UseIcopContext(RemoveContact))
	}

	//this group is used only for the change password functionality. It is a special key, which is received from
	//auth/lost_password_tfa
	auth2 := r.Group("/portal/user/auth2")
	auth2.Use(authMiddlewareLostPwd.MiddlewareFunc())
	{
		auth2.POST("/refresh", authMiddlewareLostPwd.RefreshHandler)
		auth2.POST("/lost_password_update", mw.UseIcopContext(LostPasswordUpdate))
	}

	//run the api
	if err := r.Run(fmt.Sprintf(":%d", cnf.Port)); err != nil {
		log.WithError(err).Fatalf("Failed to run server")
	}
}

func connectServices(log *logrus.Entry) {
	//connect mail service
	connMail, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.MailSrvHost, cnf.Services.MailSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	mailClient = pb.NewMailServiceClient(connMail)

	//connect jwt service
	connJwt, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.JwtSrvHost, cnf.Services.JwtSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	jwtClient = pb.NewJwtServiceClient(connJwt)

	//connect 2fa service
	conn2FA, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.TfaSrvHost, cnf.Services.TfaSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	twoFAClient = pb.NewTwoFactorAuthServiceClient(conn2FA)

	//connect db service
	connDB, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.DBSrvHost, cnf.Services.DBSrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	dbClient = pb.NewDBServiceClient(connDB)

	//connect db service
	connAdminAPI, err := grpc.Dial(fmt.Sprintf("%s:%d", cnf.Services.AdminAPISrvHost, cnf.Services.AdminAPISrvPort), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("Dial failed: %v", err)
	}
	adminAPIClient = pb.NewAdminApiServiceClient(connAdminAPI)
}

var (
	buildDate  string
	gitVersion string
	gitRemote  string
)

//Info show some info
type infoStruct struct {
	Version             string
	NumGoRutines        int
	MemMbUsedAlloc      uint64
	MemMbUsedTotalAlloc uint64
	BuildDate           string
	GitVersion          string
	GitRemote           string
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//Info Prints some information on the binary
func Info(c *gin.Context) {
	d := new(infoStruct)
	d.Version = "1"
	d.NumGoRutines = runtime.NumGoroutine()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	d.MemMbUsedAlloc = bToMb(m.Alloc)
	d.MemMbUsedTotalAlloc = bToMb(m.TotalAlloc)
	d.BuildDate = buildDate
	d.GitVersion = gitVersion
	d.GitRemote = gitRemote
	c.JSON(http.StatusOK, d)
}

//Test just for test purpose
func Test(c *gin.Context) {
	l, _ := c.Get("language")
	c.JSON(http.StatusOK, gin.H{
		"language": l,
	})
}

//we need this, because rice will not look for subfunctions/packages yet ...
//please add all boxes in here
func initBoxes() {
	rice.MustFindBox("templates/mail")
}
