package api

import (
	"net/http"
	"strconv"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"
)

const (
	//KnownCurrenciesRoutePrefix for the known currencies routes
	KnownCurrenciesRoutePrefix = "known_currencies"
)

//init setup all the routes for the known currencies handling
func init() {
	route.AddRoute("GET", "/get/:id", GetKnownCurrency, []string{}, "known_currencies_get", KnownCurrenciesRoutePrefix)
	route.AddRoute("GET", "/all", AllKnownCurrencies, []string{}, "known_currencies_all", KnownCurrenciesRoutePrefix)
	route.AddRoute("POST", "/add", AddKnownCurrency, []string{}, "known_currencies_add", KnownCurrenciesRoutePrefix)
	route.AddRoute("POST", "/edit", EditKnownCurrency, []string{}, "known_currencies_edit", KnownCurrenciesRoutePrefix)
	route.AddRoute("POST", "/delete", DeleteKnownCurrency, []string{}, "known_currencies_delete", KnownCurrenciesRoutePrefix)
	route.AddRoute("POST", "/changeOrder", ChangeOrderKnownCurrency, []string{}, "known_currencies_change_order", KnownCurrenciesRoutePrefix)
}

//AddKnownCurrenciesRoutes adds all the routes for the known currencies handling
func AddKnownCurrenciesRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(KnownCurrenciesRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//KnownCurrencyIDRequest request used in get one and delete
//swagger:parameters KnownCurrencyIDRequest GetKnownCurrency
type KnownCurrencyIDRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
}

//GetKnownCurrenciesResponse response
// swagger:model
type GetKnownCurrenciesResponse struct {
	ID               int    `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"`
	AssetCode        string `form:"asset_code" json:"asset_code"`
	ShortDescription string `form:"short_description" json:"short_description"`
	LongDescription  string `form:"long_description" json:"long_description"`
	OrderIndex       int    `form:"order_index" json:"order_index"`
}

//GetKnownCurrency returns currency by id
// swagger:route GET /portal/admin/dash/known_currencies/get/:id KnownCurrencies GetKnownCurrency
//
// Returns currency by id
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:GetKnownCurrenciesResponse
func GetKnownCurrency(uc *mw.AdminContext, c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading currency id", cerr.GeneralError))
		return
	}
	r := KnownCurrencyIDRequest{ID: id}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	currency, err := db.GetKnownCurrencyByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currencies", cerr.GeneralError))
		return
	}

	response := GetKnownCurrenciesResponse{
		ID:               currency.ID,
		Name:             currency.Name,
		IssuerPublicKey:  currency.IssuerPublicKey,
		AssetCode:        currency.AssetCode,
		ShortDescription: currency.ShortDescription,
		LongDescription:  currency.LongDescription,
		OrderIndex:       currency.OrderIndex,
	}

	c.JSON(http.StatusOK, response)

}

//AllKnownCurrencies returns all currencies
// swagger:route GET /portal/admin/dash/known_currencies/all KnownCurrencies AllKnownCurrencies
//
// Returns all currencies
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:GetKnownCurrenciesResponse
func AllKnownCurrencies(uc *mw.AdminContext, c *gin.Context) {

	currencies, err := db.GetKnownCurrencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currencies", cerr.GeneralError))
		return
	}

	var response []GetKnownCurrenciesResponse

	for _, c := range currencies {

		response = append(response, GetKnownCurrenciesResponse{
			ID:               c.ID,
			Name:             c.Name,
			IssuerPublicKey:  c.IssuerPublicKey,
			AssetCode:        c.AssetCode,
			ShortDescription: c.ShortDescription,
			LongDescription:  c.LongDescription,
			OrderIndex:       c.OrderIndex,
		})

	}

	c.JSON(http.StatusOK, response)

}

//AddKnownCurrencyRequest request
//swagger:parameters AddKnownCurrencyRequest AddKnownCurrency
type AddKnownCurrencyRequest struct {
	//required : true
	Name string `form:"name" json:"name"  validate:"required,max=500"`
	//required : true
	IssuerPublicKey string `form:"issuer_public_key" json:"issuer_public_key"  validate:"required,max=500"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code"  validate:"required,max=500"`
	//required : true
	ShortDescription string `form:"short_description" json:"short_description"  validate:"required,max=500"`
	//required : true
	LongDescription string `form:"long_description" json:"long_description"  validate:"required,max=500"`
}

//AddKnownCurrency adds a new currency
// swagger:route POST /portal/admin/dash/known_currencies/add KnownCurrencies AddKnownCurrency
//
// Adds a new currency
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddKnownCurrency(uc *mw.AdminContext, c *gin.Context) {
	var r AddKnownCurrencyRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	existsCurrency, err := db.ExistsKnownCurrency(r.IssuerPublicKey, r.AssetCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
		return
	}
	if existsCurrency {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key and asset_code", cerr.InvalidArgument, "Currency with same public key and asset code already exists", ""))
		return
	}
	orderIndex, err := db.KnownCurrencyNewOrderIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error geting new order index", cerr.GeneralError))
		return
	}

	currency := &models.AdminKnownCurrency{
		Name:             r.Name,
		IssuerPublicKey:  r.IssuerPublicKey,
		AssetCode:        r.AssetCode,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		OrderIndex:       orderIndex,
		UpdatedBy:        getUpdatedBy(c)}

	err = db.AddKnownCurrency(currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding known currency", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//EditKnownCurrencyRequest request
//swagger:parameters EditKnownCurrencyRequest EditKnownCurrency
type EditKnownCurrencyRequest struct {
	//required : true
	ID               int     `form:"id" json:"id"  validate:"required"`
	Name             *string `form:"name" json:"name"  validate:"max=500"`
	IssuerPublicKey  *string `form:"issuer_public_key" json:"issuer_public_key"  validate:"max=500"`
	AssetCode        *string `form:"asset_code" json:"asset_code"  validate:"max=500"`
	ShortDescription *string `form:"short_description" json:"short_description"  validate:"max=500"`
	LongDescription  *string `form:"long_description" json:"long_description"  validate:"max=500"`
}

//EditKnownCurrency edits known currency details
// swagger:route POST /portal/admin/dash/known_currencies/edit KnownCurrencies EditKnownCurrency
//
// Edits known currency details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func EditKnownCurrency(uc *mw.AdminContext, c *gin.Context) {
	var r EditKnownCurrencyRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	currency, err := db.GetKnownCurrencyByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
		return
	}
	if currency == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Currency not found in database", ""))
		return
	}

	if currency.IssuerPublicKey != *r.IssuerPublicKey || currency.AssetCode != *r.AssetCode {
		existsCurrency, err := db.ExistsKnownCurrency(*r.IssuerPublicKey, *r.AssetCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
			return
		}
		if existsCurrency {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key and asset_code", cerr.InvalidArgument, "Currency with same public key and asset code already exists", ""))
			return
		}
	}

	if r.Name != nil {
		currency.Name = *r.Name
	}
	if r.IssuerPublicKey != nil {
		currency.IssuerPublicKey = *r.IssuerPublicKey
	}
	if r.AssetCode != nil {
		currency.AssetCode = *r.AssetCode
	}
	if r.ShortDescription != nil {
		currency.ShortDescription = *r.ShortDescription
	}
	if r.LongDescription != nil {
		currency.LongDescription = *r.LongDescription
	}

	err = db.UpdateKnownCurrency(currency, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating known currency", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//KnownCurrencyDeleteRequest request used in get one and delete
//swagger:parameters KnownCurrencyDeleteRequest DeleteKnownCurrency
type KnownCurrencyDeleteRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
}

//DeleteKnownCurrency deletes known currency
// swagger:route POST /portal/admin/dash/known_currencies/delete KnownCurrencies DeleteKnownCurrency
//
// Deletes known currency
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func DeleteKnownCurrency(uc *mw.AdminContext, c *gin.Context) {
	var r KnownCurrencyDeleteRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	currency, err := db.GetKnownCurrencyByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
		return
	}
	if currency == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Currency not found in database", ""))
		return
	}

	err = db.DeleteKnownCurrency(currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting known currency", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ChangeOrderKnownCurrencyRequest request
//swagger:parameters ChangeOrderKnownCurrencyRequest ChangeOrderKnownCurrency
type ChangeOrderKnownCurrencyRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
	//required : true
	OrderModifier int `form:"order_modifier" json:"order_modifier" validate:"required"`
}

//ChangeOrderKnownCurrency alters a currency and changes the order index with +-1
// swagger:route POST /portal/admin/dash/known_currencies/changeOrder KnownCurrencies ChangeOrderKnownCurrency
//
// Alters a currency and changes the order index with +-1
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func ChangeOrderKnownCurrency(uc *mw.AdminContext, c *gin.Context) {
	var r ChangeOrderKnownCurrencyRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	currency, err := db.GetKnownCurrencyByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currency", cerr.GeneralError))
		return
	}
	if currency == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Currency not found in database", ""))
		return
	}
	if r.OrderModifier != -1 && r.OrderModifier != 1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "order_modifier", cerr.InvalidArgument, "Order modifier must be -1 or 1", ""))
		return
	}

	err = db.ChangeKnownCurrencyOrder(currency, r.OrderModifier, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating known currency order", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
