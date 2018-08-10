package main

import (
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"
	"net/http"

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

	//first check the walletdata
	reqData := &pb.CheckWalletRequest{
		UserId:            userID,
		WalletName:        l.WalletName,
		FederationAddress: l.FederationAddress,
		PublicKey_0:       l.PublicKey0,
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
		Base:              NewBaseRequest(uc),
		UserId:            userID,
		PublicKey_0:       l.PublicKey0,
		WalletName:        l.WalletName,
		FederationAddress: l.FederationAddress,
		ShowOnHomescreen:  l.ShowOnHomescreen,
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
		ws[i] = GetUserWalletsResponse{
			ID:                w.Id,
			PublicKey0:        w.PublicKey_0,
			WalletName:        w.WalletName,
			FederationAddress: w.FederationAddress,
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

	//check new Walletdata
	checkData, err := dbClient.CheckWalletData(c, &pb.CheckWalletRequest{
		Base:              NewBaseRequest(uc),
		UserId:            userID,
		WalletName:        l.WalletName,
		FederationAddress: l.FederationAddress,
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
			Base:              NewBaseRequest(uc),
			Id:                l.ID,
			UserId:            userID,
			FederationAddress: l.FederationAddress,
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
		Base:              NewBaseRequest(uc),
		Id:                l.ID,
		UserId:            userID,
		FederationAddress: "",
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
