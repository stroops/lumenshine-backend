package api

import (
	"net/http"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	coremodels "github.com/Soneso/lumenshine-backend/db/stellarcore/models"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"

	"github.com/Soneso/lumenshine-backend/admin/route"
)

//init setup all the routes for the users handling
func init() {
	route.AddRoute("POST", "/add_unathorized_trustline", AddTrustline, []string{"Administrators"}, "add_unathorized_trustline", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove_unathorized_trustline", RemoveTrustline, []string{"Administrators"}, "remove_unathorized_trustline", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/worker_account_trustlines/:publickey", WorkerAccountTrustlines, []string{"Administrators"}, "worker_account_trustlines", StellarAccountRoutePrefix)
}

//AddTrustlineRequest - info
type AddTrustlineRequest struct {
	TrustorPublicKey string `form:"trusting_account_public_key" json:"trusting_account_public_key" validate:"required,base64,len=56"`
	IssuingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	AssetCode        string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
	Status           string `form:"status" json:"status" validate:"required,max=50"`
	Reason           string `form:"reason" json:"reason" validate:"required,max=1000"`
}

//AddTrustline creates new entry in the db
func AddTrustline(uc *mw.AdminContext, c *gin.Context) {
	var rr AddTrustlineRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	if !strings.EqualFold(rr.Status, string(models.StellarTrustlineStatusDenied)) && !strings.EqualFold(rr.Status, string(models.StellarTrustlineStatusRevoked)) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "status", cerr.InvalidArgument, "Status is not 'denied' or 'revoked'", ""))
		return
	}
	existsAccount, err := db.ExistsStellarAccount(rr.TrustorPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if !existsAccount {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "trusting_account_public_key", cerr.InvalidArgument, "Trustor account does not exists", ""))
		return
	}
	issuer, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if issuer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing account does not exists", ""))
		return
	}
	if issuer.Type != models.StellarAccountTypeIssuing {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing public key does not belong to an issuing account", ""))
		return
	}
	existsAssetCode := false
	for _, assetCode := range issuer.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			existsAssetCode = true
			break
		}
	}
	if !existsAssetCode {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code does not exist for this issuing account", ""))
		return
	}
	existsTrustline, err := db.ExistsUnauthorizedTrustline(rr.TrustorPublicKey, rr.IssuingPublicKey, rr.AssetCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing trustline", cerr.GeneralError))
		return
	}
	if existsTrustline {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Trustline already exists for this issuer and asset code.", ""))
		return
	}

	trustline := &models.AdminUnauthorizedTrustline{
		TrustorPublicKey:  rr.TrustorPublicKey,
		IssuerPublicKeyID: rr.IssuingPublicKey,
		AssetCode:         rr.AssetCode,
		Status:            rr.Status,
		Reason:            rr.Reason}

	err = db.AddTrustline(trustline, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding trustline", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//RemoveTrustlineRequest - info
type RemoveTrustlineRequest struct {
	TrustorPublicKey string `form:"trusting_account_public_key" json:"trusting_account_public_key" validate:"required,base64,len=56"`
	IssuingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	AssetCode        string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
}

//RemoveTrustline creates new entry in the db
func RemoveTrustline(uc *mw.AdminContext, c *gin.Context) {
	var rr RemoveTrustlineRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	existsAccount, err := db.ExistsStellarAccount(rr.TrustorPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if !existsAccount {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "trusting_account_public_key", cerr.InvalidArgument, "Trustor account does not exists", ""))
		return
	}
	issuer, err := db.GetStellarAccount(rr.IssuingPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if issuer == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing account does not exists", ""))
		return
	}
	if issuer.Type != models.StellarAccountTypeIssuing {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "issuing_account_public_key", cerr.InvalidArgument, "Issuing public key does not belong to an issuing account", ""))
		return
	}
	existsAssetCode := false
	for _, assetCode := range issuer.R.IssuerPublicKeyAdminStellarAssets {
		if strings.EqualFold(assetCode.AssetCode, rr.AssetCode) {
			existsAssetCode = true
			break
		}
	}
	if !existsAssetCode {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "asset_code", cerr.InvalidArgument, "Asset code does not exist for this issuing account", ""))
		return
	}

	err = db.DeleteUnauthorizedTrustline(rr.TrustorPublicKey, rr.IssuingPublicKey, rr.AssetCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error removing trustline", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//WorkerAccountTrustlinesRequest - info
type WorkerAccountTrustlinesRequest struct {
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
}

//WorkerTrustlineItem - response item
type WorkerTrustlineItem struct {
	AssetCode string `json:"asset_code"`
	Issuer    string `json:"asset_issuer"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

//WorkerAccountTrustlines - returns worker account trustlines
func WorkerAccountTrustlines(uc *mw.AdminContext, c *gin.Context) {
	publicKey := c.Param("publickey")
	rr := WorkerAccountTrustlinesRequest{PublicKey: publicKey}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	workerAccount, err := db.GetStellarAccount(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing account", cerr.GeneralError))
		return
	}
	if workerAccount == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Worker account does not exists", ""))
		return
	}
	if workerAccount.Type != models.StellarAccountTypeWorker {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key", cerr.InvalidArgument, "Public key does not belong to a worker account", ""))
		return
	}

	coreTrustlines, err := db.GetCoreTrustlines(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading core trustlines", cerr.GeneralError))
		return
	}

	issuers, err := db.GetStellarIssuerAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading stellar issuer accounts", cerr.GeneralError))
		return
	}

	uaTrustlines, err := db.GetUnauthorizedTrustlines(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading unauthorized trustlines", cerr.GeneralError))
		return
	}

	internalTrustlines := make([]*coremodels.Trustline, 0)
	for _, trustline := range coreTrustlines {
		for _, issuer := range issuers {
			if strings.EqualFold(trustline.Issuer, issuer.PublicKey) {
				for _, assetCode := range issuer.R.IssuerPublicKeyAdminStellarAssets {
					if strings.EqualFold(trustline.Assetcode, assetCode.AssetCode) {
						internalTrustlines = append(internalTrustlines, trustline)
						break
					}
				}
			}
		}
	}

	resultTrustlines := make([]*WorkerTrustlineItem, 0)
	for _, trustline := range internalTrustlines {
		status := "waiting_for_authorization"
		reason := ""
		if trustline.Flags == 1 {
			status = "allowed"
		}
		for _, uaTrustline := range uaTrustlines {
			if strings.EqualFold(trustline.Issuer, uaTrustline.IssuerPublicKeyID) &&
				strings.EqualFold(trustline.Assetcode, uaTrustline.AssetCode) {
				status = uaTrustline.Status
				reason = uaTrustline.Reason
				break
			}
		}
		resultItem := WorkerTrustlineItem{
			AssetCode: trustline.Assetcode,
			Issuer:    trustline.Issuer,
			Status:    status,
			Reason:    reason,
		}
		resultTrustlines = append(resultTrustlines, &resultItem)
	}

	c.JSON(http.StatusOK, &resultTrustlines)
}