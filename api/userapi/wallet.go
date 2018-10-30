package main

import (
	"net/http"
	"strings"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

//AddWalletRequest request
//swagger:parameters AddWalletRequest AddWallet
type AddWalletRequest struct {
	// required: true
	PublicKey string `form:"public_key" json:"public_key"  validate:"required,base64,len=56"`
	// required: true
	WalletName        string `form:"wallet_name" json:"wallet_name" validate:"required,max=500"`
	FederationAddress string `form:"federation_address" json:"federation_address" validate:"max=255"`
	ShowOnHomescreen  bool   `form:"show_on_homescreen" json:"show_on_homescreen"`
	WalletType        string `form:"wallet_type" json:"wallet_type"`
}

//AddWalletResponse request
// swagger:model
type AddWalletResponse struct {
	ID int64 `json:"id"`
}

//AddWallet adds a new wallet to the user
// swagger:route POST /portal/user/dashboard/add_wallet wallet AddWallet
//
// AddWallet adds a new wallet to the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:AddWalletResponse
func AddWallet(uc *mw.IcopContext, c *gin.Context) {
	var l AddWalletRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	friendlyID := ""
	domain := ""
	if l.FederationAddress != "" {
		fedS := strings.Split(l.FederationAddress, "*")
		if len(fedS) != 2 {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("federation_address", cerr.InvalidArgument, "Federation address incorrect format", ""))
			return
		}
		friendlyID = fedS[0]
		domain = fedS[1]
	}

	//first check the walletdata
	reqData := &pb.CheckWalletRequest{
		UserId:     userID,
		WalletName: l.WalletName,
		FriendlyId: friendlyID,
		Domain:     domain,
		PublicKey:  l.PublicKey,
	}
	walletStatus, err := dbClient.CheckWalletData(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error checking wallet data", cerr.GeneralError))
		return
	}
	if l.FederationAddress != "" && !walletStatus.FederationAddressOk {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("federation_address", cerr.WalletFederationNameExists, "Federation name already exists", ""))
		return
	}

	if !walletStatus.PublicKeyOk {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key", cerr.InvalidArgument, "Publickey already exists for user", ""))
		return
	}

	if !walletStatus.NameOk {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("wallet_name", cerr.InvalidArgument, "Walletname already exists for user", ""))
		return
	}
	walletType := int32(pb.WalletType_internal)
	if l.WalletType != "" {
		ok := false
		walletType, ok = pb.WalletType_value[l.WalletType]
		if !ok {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("wallet_type", cerr.InvalidArgument, "Invalid wallet type.", ""))
			return
		}
	}

	//add the wallet
	req := &pb.AddWalletRequest{
		Base:             NewBaseRequest(uc),
		UserId:           userID,
		PublicKey:        l.PublicKey,
		WalletName:       l.WalletName,
		FriendlyId:       friendlyID,
		Domain:           domain,
		ShowOnHomescreen: l.ShowOnHomescreen,
		WalletType:       pb.WalletType(walletType),
	}
	id, err := dbClient.AddWallet(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding wallet", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &AddWalletResponse{
		ID: id.Id,
	})
}

//GetUserWalletsResponse result
// swagger:model
type GetUserWalletsResponse struct {
	ID                int64  `json:"id"`
	PublicKey         string `json:"public_key"`
	WalletName        string `json:"wallet_name"`
	FederationAddress string `json:"federation_address"`
	ShowOnHomescreen  bool   `json:"show_on_homescreen"`
	WalletType        string `json:"wallet_type"`
}

//GetUserWallets returns all wallets for one user
// swagger:route GET /portal/user/dashboard/get_user_wallets wallet GetUserWallets
//
// Returns all wallets for one user
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:[]GetUserWalletsResponse
func GetUserWallets(uc *mw.IcopContext, c *gin.Context) {
	userID := mw.GetAuthUser(c).UserID
	req := &pb.GetWalletsRequest{
		Base:   NewBaseRequest(uc),
		UserId: userID,
	}
	wallets, err := dbClient.GetUserWallets(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading wallets", cerr.GeneralError))
		return
	}

	ws := make([]GetUserWalletsResponse, len(wallets.Wallets))
	for i, w := range wallets.Wallets {
		federationAddress := ""
		if w.FriendlyId != "" && w.Domain != "" {
			federationAddress = w.FriendlyId + "*" + w.Domain
		}
		ws[i] = GetUserWalletsResponse{
			ID:                w.Id,
			PublicKey:         w.PublicKey,
			WalletName:        w.WalletName,
			FederationAddress: federationAddress,
			ShowOnHomescreen:  w.ShowOnHomescreen,
			WalletType:        w.WalletType.String(),
		}
	}

	c.JSON(http.StatusOK, ws)
}

//RemoveWalletRequest request
//swagger:parameters RemoveWalletRequest RemoveWallet
type RemoveWalletRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//RemoveWallet removes a wallet from the user
// swagger:route POST /portal/user/dashboard/remove_wallet wallet RemoveWallet
//
// Removes a wallet from the user
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveWallet(uc *mw.IcopContext, c *gin.Context) {
	var l RemoveWalletRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	//check if
	isLast, err := dbClient.WalletIsLast(c, &pb.WalletIsLastRequest{
		Base:   NewBaseRequest(uc),
		Id:     l.ID,
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error checking isLast wallet", cerr.GeneralError))
		return
	}

	if isLast.Value {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("id", cerr.WalletIsLast, "Can't remove last wallet", ""))
		return
	}

	//remove the wallet
	req := &pb.RemoveWalletRequest{
		Base:   NewBaseRequest(uc),
		Id:     l.ID,
		UserId: userID,
	}
	_, err = dbClient.RemoveWallet(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing wallet", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//WalletChangeOrderRequest - request
//swagger:parameters WalletChangeOrderRequest WalletChangeOrder
type WalletChangeOrderRequest struct {
	// required: true
	PublicKey string `form:"public_key" json:"public_key"  validate:"required,base64,len=56"`
	// required: true
	OrderNr int `form:"order_nr" json:"order_nr"`
}

//WalletChangeOrder changes the wallet orders
// swagger:route POST /portal/user/dashboard/change_wallet_order wallet WalletChangeOrder
//
// Changes the wallet orders
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func WalletChangeOrder(uc *mw.IcopContext, c *gin.Context) {
	var r WalletChangeOrderRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	_, err := dbClient.WalletChangeOrder(c, &pb.WalletChangeOrderRequest{
		Base:      NewBaseRequest(uc),
		UserId:    userID,
		PublicKey: r.PublicKey,
		OrderNr:   int64(r.OrderNr),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating wallet order", cerr.GeneralError))
		return
	}

	req := &pb.GetWalletsRequest{
		Base:   NewBaseRequest(uc),
		UserId: userID,
	}
	wallets, err := dbClient.GetUserWallets(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading wallets", cerr.GeneralError))
		return
	}

	ws := make([]GetUserWalletsResponse, len(wallets.Wallets))
	for i, w := range wallets.Wallets {
		federationAddress := ""
		if w.FriendlyId != "" && w.Domain != "" {
			federationAddress = w.FriendlyId + "*" + w.Domain
		}
		ws[i] = GetUserWalletsResponse{
			ID:                w.Id,
			PublicKey:         w.PublicKey,
			WalletName:        w.WalletName,
			FederationAddress: federationAddress,
			ShowOnHomescreen:  w.ShowOnHomescreen,
		}
	}

	c.JSON(http.StatusOK, ws)
}

//WalletChangeDataRequest request
//swagger:parameters WalletChangeDataRequest WalletChangeData
type WalletChangeDataRequest struct {
	// required: true
	ID                int64  `form:"id" json:"id"  validate:"required"`
	WalletName        string `form:"wallet_name" json:"wallet_name"  validate:"max=500"`
	FederationAddress string `form:"federation_address" json:"federation_address"  validate:"max=255"`
}

//WalletChangeData changes the wallet data
// swagger:route POST /portal/user/dashboard/change_wallet_data wallet WalletChangeData
//
// Changes the wallet data
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func WalletChangeData(uc *mw.IcopContext, c *gin.Context) {
	var l WalletChangeDataRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	friendlyID := ""
	domain := ""
	if l.FederationAddress != "" {
		fedS := strings.Split(l.FederationAddress, "*")
		if len(fedS) != 2 {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("federation_address", cerr.InvalidArgument, "Federation address incorrect format", ""))
			return
		}
		friendlyID = fedS[0]
		domain = fedS[1]
	}

	//check new Walletdata
	checkData, err := dbClient.CheckWalletData(c, &pb.CheckWalletRequest{
		Base:       NewBaseRequest(uc),
		UserId:     userID,
		WalletName: l.WalletName,
		FriendlyId: friendlyID,
		Domain:     domain,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error checking wallet data", cerr.GeneralError))
		return
	}
	if (l.WalletName != "" && !checkData.NameOk) || (l.FederationAddress != "" && !checkData.FederationAddressOk) {
		errors := new(cerr.IcopErrors)

		if l.WalletName != "" && !checkData.NameOk {
			errors.AddError("wallet_name", cerr.InvalidArgument, "Name already exists", "")
		}

		if l.FederationAddress != "" && !checkData.FederationAddressOk {
			errors.AddError("federation_address", cerr.WalletFederationNameExists, "Federation-address already exists", "")
		}

		c.JSON(http.StatusBadRequest, errors)
		return
	}

	//change the wallet name
	if l.WalletName != "" {
		req := &pb.WalletChangeNameRequest{
			Base:   NewBaseRequest(uc),
			Id:     l.ID,
			UserId: userID,
			Name:   l.WalletName,
		}
		_, err = dbClient.WalletChangeName(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error changing wallet name", cerr.GeneralError))
			return
		}
	}

	//change the fed name
	if l.FederationAddress != "" {
		req := &pb.WalletChangeFederationAddressRequest{
			Base:       NewBaseRequest(uc),
			Id:         l.ID,
			UserId:     userID,
			FriendlyId: friendlyID,
			Domain:     domain,
		}
		_, err = dbClient.WalletChangeFederationAddress(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error changing federation address", cerr.GeneralError))
			return
		}
	}

	c.JSON(http.StatusOK, "{}")
}

//RemoveWalletFederationAddressRequest request
//swagger:parameters RemoveWalletFederationAddressRequest RemoveWalletFederationAddress
type RemoveWalletFederationAddressRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//RemoveWalletFederationAddress removes federation name from the wallet
// swagger:route POST /portal/user/dashboard/remove_wallet_federation_address wallet RemoveWalletFederationAddress
//
// Removes federation name from the wallet
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemoveWalletFederationAddress(uc *mw.IcopContext, c *gin.Context) {
	var l RemoveWalletFederationAddressRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	//remove the wallet
	req := &pb.WalletChangeFederationAddressRequest{
		Base:       NewBaseRequest(uc),
		Id:         l.ID,
		UserId:     userID,
		FriendlyId: "",
		Domain:     "",
	}
	_, err := dbClient.WalletChangeFederationAddress(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing wallet federation address", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//WalletSetHomescreenRequest request
//swagger:parameters WalletSetHomescreenRequest WalletSetHomescreen
type WalletSetHomescreenRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"  validate:"required"`
	// required: true
	Visible bool `form:"visible" json:"visible"`
}

//WalletSetHomescreen sets the wallet visible flag
// swagger:route POST /portal/user/dashboard/wallet_set_homescreen wallet WalletSetHomescreen
//
// Sets the wallet visible flag
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func WalletSetHomescreen(uc *mw.IcopContext, c *gin.Context) {
	var l WalletSetHomescreenRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	//remove the wallet
	req := &pb.WalletSetHomescreenRequest{
		Base:    NewBaseRequest(uc),
		Id:      l.ID,
		UserId:  userID,
		Visible: l.Visible,
	}
	_, err := dbClient.WalletSetHomescreen(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error setting wallet homescreen flag", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//GetKnownCurrencyRequest request
//swagger:parameters GetKnownCurrencyRequest GetKnownCurrency
type GetKnownCurrencyRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//GetKnownCurrenciesResponse response
// swagger:model
type GetKnownCurrenciesResponse struct {
	ID               int64  `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"`
	AssetCode        string `form:"asset_code" json:"asset_code"`
	ShortDescription string `form:"short_description" json:"short_description"`
	LongDescription  string `form:"long_description" json:"long_description"`
	OrderIndex       int64  `form:"order_index" json:"order_index"`
}

//GetKnownCurrency returns currency with specified id
// swagger:route GET /portal/user/dashboard/get_known_currency wallet GetKnownCurrency
//
// Returns currency with specified id
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetKnownCurrenciesResponse
func GetKnownCurrency(uc *mw.IcopContext, c *gin.Context) {
	var l GetKnownCurrencyRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//get the currency
	req := &pb.GetKnownCurrencyRequest{
		Base: NewBaseRequest(uc),
		Id:   l.ID,
	}
	res, err := adminAPIClient.GetKnownCurrency(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting known currency", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &GetKnownCurrenciesResponse{
		ID:               res.Id,
		Name:             res.Name,
		IssuerPublicKey:  res.IssuerPublicKey,
		AssetCode:        res.AssetCode,
		ShortDescription: res.ShortDescription,
		LongDescription:  res.LongDescription,
		OrderIndex:       res.OrderIndex,
	})
}

//GetKnownCurrencies returns all known currencies
// swagger:route GET /portal/user/dashboard/get_known_currencies wallet GetKnownCurrencies
//
// Returns all known currencies
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetKnownCurrenciesResponse
func GetKnownCurrencies(uc *mw.IcopContext, c *gin.Context) {

	//get the currency
	res, err := adminAPIClient.GetKnownCurrencies(c, &pb.Empty{Base: NewBaseRequest(uc)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting known currencies", cerr.GeneralError))
		return
	}

	crs := make([]GetKnownCurrenciesResponse, len(res.Currencies))
	for i, cr := range res.Currencies {
		crs[i] = GetKnownCurrenciesResponse{
			ID:               cr.Id,
			Name:             cr.Name,
			IssuerPublicKey:  cr.IssuerPublicKey,
			AssetCode:        cr.AssetCode,
			ShortDescription: cr.ShortDescription,
			LongDescription:  cr.LongDescription,
			OrderIndex:       cr.OrderIndex,
		}
	}

	c.JSON(http.StatusOK, crs)

}

//GetKnownInflationDestinationRequest request
//swagger:parameters GetKnownInflationDestinationRequest GetKnownInflationDestination
type GetKnownInflationDestinationRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//GetKnownInflationDestinationsResponse response
// swagger:model
type GetKnownInflationDestinationsResponse struct {
	ID               int64  `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"`
	ShortDescription string `form:"short_description" json:"short_description"`
	LongDescription  string `form:"long_description" json:"long_description"`
	OrderIndex       int64  `form:"order_index" json:"order_index"`
}

//GetKnownInflationDestination returns destination with specified id
// swagger:route GET /portal/user/dashboard/get_known_inflation_destination wallet GetKnownInflationDestination
//
// Returns destination with specified id
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetKnownInflationDestinationsResponse
func GetKnownInflationDestination(uc *mw.IcopContext, c *gin.Context) {
	var l GetKnownInflationDestinationRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//get the destination
	req := &pb.GetKnownInflationDestinationRequest{
		Base: NewBaseRequest(uc),
		Id:   l.ID,
	}
	res, err := adminAPIClient.GetKnownInflationDestination(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting known inflation destination", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &GetKnownInflationDestinationsResponse{
		ID:               res.Id,
		Name:             res.Name,
		IssuerPublicKey:  res.IssuerPublicKey,
		ShortDescription: res.ShortDescription,
		LongDescription:  res.LongDescription,
		OrderIndex:       res.OrderIndex,
	})
}

//GetKnownInflationDestinations returns all known destinations
// swagger:route GET /portal/user/dashboard/get_known_inflation_destinations wallet GetKnownInflationDestinations
//
// Returns all known destinations
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetKnownInflationDestinationsResponse
func GetKnownInflationDestinations(uc *mw.IcopContext, c *gin.Context) {

	//get the destinations
	res, err := adminAPIClient.GetKnownInflationDestinations(c, &pb.Empty{Base: NewBaseRequest(uc)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting known inflation destinations", cerr.GeneralError))
		return
	}

	crs := make([]GetKnownInflationDestinationsResponse, len(res.Destinations))
	for i, cr := range res.Destinations {
		crs[i] = GetKnownInflationDestinationsResponse{
			ID:               cr.Id,
			Name:             cr.Name,
			IssuerPublicKey:  cr.IssuerPublicKey,
			ShortDescription: cr.ShortDescription,
			LongDescription:  cr.LongDescription,
			OrderIndex:       cr.OrderIndex,
		}
	}
	c.JSON(http.StatusOK, crs)
}

//GetPaymentTemplatesRequest request
//swagger:parameters GetPaymentTemplatesRequest GetPaymentTemplates
type GetPaymentTemplatesRequest struct {
	// required: true
	WalletID int64 `form:"wallet_id" json:"wallet_id"`
}

//GetPaymentTemplateResponse result
// swagger:model
type GetPaymentTemplateResponse struct {
	ID                      int    `json:"id"`
	RecipientStellarAddress string `json:"recipient_stellar_address"`
	RecipientPK             string `json:"recipien_pk"`
	AssetCode               string `json:"asset_code"`
	IssuerPK                string `json:"issuer_pk"`
	Amount                  int64  `json:"amount"`
	MemoType                string `json:"memo_type"`
	Memo                    string `json:"memo"`
}

//GetPaymentTemplates returns all templates for a wallet
// swagger:route GET /portal/user/dashboard/get_payment_templates wallet GetPaymentTemplates
//
// Returns all templates for a wallet
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []GetPaymentTemplateResponse
func GetPaymentTemplates(uc *mw.IcopContext, c *gin.Context) {
	var r GetPaymentTemplatesRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	userID := mw.GetAuthUser(c).UserID
	req := &pb.GetTemplatesRequest{
		Base:     NewBaseRequest(uc),
		WalletId: r.WalletID,
		UserId:   userID,
	}
	templates, err := dbClient.GetPaymentTemplates(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading templates", cerr.GeneralError))
		return
	}

	ws := make([]GetPaymentTemplateResponse, len(templates.Templates))
	for i, t := range templates.Templates {
		ws[i] = GetPaymentTemplateResponse{
			ID: int(t.Id),
			RecipientStellarAddress: t.RecipientStellarAddress,
			RecipientPK:             t.RecipientPublickey,
			AssetCode:               t.AssetCode,
			IssuerPK:                t.IssuerPublickey,
			Amount:                  t.Amount,
			MemoType:                t.MemoType.String(),
			Memo:                    t.Memo,
		}
	}

	c.JSON(http.StatusOK, ws)
}

//AddPaymentTemplateRequest request
//swagger:parameters AddPaymentTemplateRequest AddPaymentTemplate
type AddPaymentTemplateRequest struct {
	// required: true
	WalletID int `form:"wallet_id" json:"wallet_id"`
	// required: if recipient public key is not specified
	RecipientStellarAddress string `form:"recipient_stellar_address" json:"recipient_stellar_address" validate:"max=256"`
	// required: if recipient stellar address is not specified
	RecipientPK string `form:"recipient_pk" json:"recipient_pk" validate:"omitempty,base64,len=56"`
	// required: true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"icop_assetcode"`
	// required: asset code is not XLM
	IssuerPK string `form:"issuer_pk" json:"issuer_pk" validate:"omitempty,base64,len=56"`
	// required: true
	Amount int64 `form:"amount" json:"amount"`
	// required: true
	MemoType string `form:"memo_type" json:"memo_type" validate:"required,max=8"`
	Memo     string `form:"memo" json:"memo"`
}

//AddTemplateResponse response
// swagger:model
type AddTemplateResponse struct {
	//newly added template id
	ID int64 `json:"id"`
}

//AddPaymentTemplate adds a new template to the wallet
// swagger:route POST /portal/user/dashboard/add_payment_template wallet AddPaymentTemplate
//
// Adds a new template to the wallet
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: AddTemplateResponse
func AddPaymentTemplate(uc *mw.IcopContext, c *gin.Context) {
	var r AddPaymentTemplateRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	if r.RecipientPK == "" && r.RecipientStellarAddress == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("recipient_pk", cerr.InvalidArgument, "At least one recipient field must be specified.", ""))
		return
	}
	if r.AssetCode != "XLM" && r.IssuerPK == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("issuer_pk", cerr.InvalidArgument, "Issuer must be specified.", ""))
		return
	}
	if r.AssetCode == "XLM" && r.IssuerPK != "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("issuer_pk", cerr.InvalidArgument, "Issuer must not be set for native XLM.", ""))
		return
	}
	if _, ok := pb.MemoType_value[r.MemoType]; !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("memo_type", cerr.InvalidArgument, "Invalid memo type.", ""))
		return
	}

	userID := mw.GetAuthUser(c).UserID
	id, err := dbClient.AddPaymentTemplate(c, &pb.AddPaymentTemplateRequest{
		Base:                    NewBaseRequest(uc),
		UserId:                  userID,
		WalletId:                int64(r.WalletID),
		RecipientStellarAddress: r.RecipientStellarAddress,
		RecipientPublickey:      r.RecipientPK,
		AssetCode:               r.AssetCode,
		IssuerPublickey:         r.IssuerPK,
		Amount:                  r.Amount,
		MemoType:                pb.MemoType(pb.MemoType_value[r.MemoType]),
		Memo:                    r.Memo,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding payment template", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, &AddTemplateResponse{
		ID: id.Id,
	})
}

//RemovePaymentTemplateRequest request
//swagger:parameters RemovePaymentTemplateRequest RemovePaymentTemplate
type RemovePaymentTemplateRequest struct {
	// required: true
	ID int64 `form:"id" json:"id"`
}

//RemovePaymentTemplate removes a template from the wallet
// swagger:route POST /portal/user/dashboard/remove_payment_template wallet RemovePaymentTemplate
//
// Removes a template from the wallet
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func RemovePaymentTemplate(uc *mw.IcopContext, c *gin.Context) {
	var r RemovePaymentTemplateRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	userID := mw.GetAuthUser(c).UserID
	req := &pb.RemovePaymentTemplateRequest{
		Base:   NewBaseRequest(uc),
		Id:     r.ID,
		UserId: userID,
	}
	_, err := dbClient.RemovePaymentTemplate(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing template", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
