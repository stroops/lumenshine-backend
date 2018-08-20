package api

import (
	"net/http"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"

	"github.com/Soneso/lumenshine-backend/admin/route"
)

//init setup all the routes for the users handling
func init() {
	route.AddRoute("POST", "/add_unathorized_trustline", AddTrustline, []string{"Administrators"}, "add_unathorized_trustline", StellarAccountRoutePrefix)
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
	issuer, err := db.GetStellarAccountLight(rr.IssuingPublicKey)
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
	trustline := &models.AdminUnauthorizedTrustline{
		StellarAccountPublicKeyID: rr.TrustorPublicKey,
		IssuerPublicKeyID:         rr.IssuingPublicKey,
		AssetCode:                 rr.AssetCode,
		Status:                    rr.Status,
		Reason:                    rr.Reason}

	err = db.AddTrustline(trustline, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding trustline", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
