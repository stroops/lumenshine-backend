package icop_error

//global Error codes
const (
	GeneralError               = 1    //Used for global error handling to the client
	BindError                  = 2    //Used if the data from the client could not be bound
	JwtError                   = 3    //returned, if the jwt token is not valid
	InvalidArgument            = 1000 //used to signal any invalid argument
	EmailExists                = 1001 //returned on registration if email exists
	MissingMandatoryField      = 1002
	MobileNotCountry           = 1003 //returned if mobile does not match country
	InvalidLength              = 1004
	MasterKySameAsIV           = 1005 //returned if the IV is the same as the masterkey
	TokenExpiered              = 1006 //returned if mail token expiered
	JWTExpired                 = 1007
	EmailAlreadyConfirmed      = 1008 //returned if email was already confirmed before
	TfaAlreadyConfirmed        = 1009 //returned if 2fa already done
	EmailNotConfigured         = 1010
	MnemonicNotConfigured      = 1011
	TfaNotYetConfirmed         = 1012 //returned if 2fa not yet confirmed
	InvalidPassword            = 1013 //returned if passed in publickey188 does not match
	UserInactive               = 1014 //returned if the user is inactive
	NoPermission               = 1015 //returned if the user has not permission to the ressource
	UserNotExists              = 1016 //returned if user is not found in the db
	WalletFederationNameExists = 1017 //returned if federation name already exists
	WalletIsLast               = 1018 //returned if federation name already exists
	OrderWrongStatus           = 1019 //returned if order has the wrong status for the desired action
	NoActivePhase              = 1020 //returned if no active phase was found
	InsufficientCoins          = 1021 //returned if there are not suffuciant coins for the order
	StellarAccountNotExists    = 1022 //returned if the stellar account does not exists
	StellarTrustlineNotExists  = 1023 //returned if the trustline does not exists

	ValidBadInputData = "Bad input data"
)
