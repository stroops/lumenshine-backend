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
type AddWalletRequest struct {
	PublicKey0        string `form:"public_key_0" json:"public_key_0"  validate:"required,base64,len=56"`
	WalletName        string `form:"wallet_name" json:"wallet_name" validate:"required,max=500"`
	FederationAddress string `form:"federation_address" json:"federation_address" validate:"max=255"`
	ShowOnHomescreen  bool   `form:"show_on_homescreen" json:"show_on_homescreen"`
}

//AddWalletResponse request
type AddWalletResponse struct {
	ID int64 `json:"id"`
}

//AddWallet adds a new wallet to the user
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
		UserId:      userID,
		WalletName:  l.WalletName,
		FriendlyId:  friendlyID,
		Domain:      domain,
		PublicKey_0: l.PublicKey0,
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

	if !walletStatus.PublicKey_0Ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key_0", cerr.InvalidArgument, "Publickey already exists for user", ""))
		return
	}

	if !walletStatus.NameOk {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("wallet_name", cerr.InvalidArgument, "Walletname already exists for user", ""))
		return
	}

	//add the wallet
	req := &pb.AddWalletRequest{
		Base:             NewBaseRequest(uc),
		UserId:           userID,
		PublicKey_0:      l.PublicKey0,
		WalletName:       l.WalletName,
		FriendlyId:       friendlyID,
		Domain:           domain,
		ShowOnHomescreen: l.ShowOnHomescreen,
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
type GetUserWalletsResponse struct {
	ID                int64  `json:"id"`
	PublicKey0        string `json:"public_key_0"`
	WalletName        string `json:"wallet_name"`
	FederationAddress string `json:"federation_address"`
	ShowOnHomescreen  bool   `json:"show_on_homescreen"`
}

//GetUserWallets returns all wallets for one user
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
			PublicKey0:        w.PublicKey_0,
			WalletName:        w.WalletName,
			FederationAddress: federationAddress,
			ShowOnHomescreen:  w.ShowOnHomescreen,
		}
	}

	c.JSON(http.StatusOK, ws)
}

//RemoveWalletRequest request
type RemoveWalletRequest struct {
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//RemoveWallet removes a wallet to the user
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
type WalletChangeOrderRequest struct {
	PublicKey0 string `form:"public_key_0" json:"public_key_0"  validate:"required,base64,len=56"`
	OrderNr    int    `form:"order_nr" json:"order_nr"`
}

//WalletChangeOrder changes the wallet order
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
		Base:        NewBaseRequest(uc),
		UserId:      userID,
		PublicKey_0: r.PublicKey0,
		OrderNr:     int64(r.OrderNr),
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
			PublicKey0:        w.PublicKey_0,
			WalletName:        w.WalletName,
			FederationAddress: federationAddress,
			ShowOnHomescreen:  w.ShowOnHomescreen,
		}
	}

	c.JSON(http.StatusOK, ws)
}

//WalletChangeDataRequest request
type WalletChangeDataRequest struct {
	ID                int64  `form:"id" json:"id"  validate:"required"`
	WalletName        string `form:"wallet_name" json:"wallet_name"  validate:"max=500"`
	FederationAddress string `form:"federation_address" json:"federation_address"  validate:"max=255"`
}

//WalletChangeData changes the wallet data
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
type RemoveWalletFederationAddressRequest struct {
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//RemoveWalletFederationAddress removes federation name from the wallet
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
type WalletSetHomescreenRequest struct {
	ID      int64 `form:"id" json:"id"  validate:"required"`
	Visible bool  `form:"visible" json:"visible"`
}

//WalletSetHomescreen sets the wallet visible flag
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
type GetKnownCurrencyRequest struct {
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//GetKnownCurrenciesResponse response
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
type GetKnownInflationDestinationRequest struct {
	ID int64 `form:"id" json:"id"  validate:"required"`
}

//GetKnownInflationDestinationsResponse response
type GetKnownInflationDestinationsResponse struct {
	ID               int64  `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"`
	ShortDescription string `form:"short_description" json:"short_description"`
	LongDescription  string `form:"long_description" json:"long_description"`
	OrderIndex       int64  `form:"order_index" json:"order_index"`
}

//GetKnownInflationDestination returns destination with specified id
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
type GetPaymentTemplatesRequest struct {
	WalletID int64 `form:"wallet_id" json:"wallet_id"`
}

//GetPaymentTemplateResponse result
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
type AddPaymentTemplateRequest struct {
	WalletID                int    `form:"wallet_id" json:"wallet_id"`
	RecipientStellarAddress string `form:"recipient_stellar_address" json:"recipient_stellar_address" validate:"max=256"`
	RecipientPK             string `form:"recipient_pk" json:"recipient_pk" validate:"omitempty,base64,len=56"`
	AssetCode               string `form:"asset_code" json:"asset_code" validate:"icop_assetcode"`
	IssuerPK                string `form:"issuer_pk" json:"issuer_pk" validate:"omitempty,base64,len=56"`
	Amount                  int64  `form:"amount" json:"amount"`
	MemoType                string `form:"memo_type" json:"memo_type" validate:"required,max=8"`
	Memo                    string `form:"memo" json:"memo"`
}

//AddTemplateResponse response
type AddTemplateResponse struct {
	ID int64 `json:"id"`
}

//AddPaymentTemplate adds a new template to the wallet
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
type RemovePaymentTemplateRequest struct {
	ID int64 `form:"id" json:"id"`
}

//RemovePaymentTemplate removes a template from the wallet
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
