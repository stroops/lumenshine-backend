package api

import (
	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	//KnownInflationDestinationsRoutePrefix for the known inflation destinations routes
	KnownInflationDestinationsRoutePrefix = "known_inflation_destinations"
)

//init setup all the routes for the known inflation destinations handling
func init() {
	route.AddRoute("GET", "/get/:id", GetKnownInflationDestination, []string{}, "known_inflation_destinations_get", KnownInflationDestinationsRoutePrefix)
	route.AddRoute("GET", "/all", AllKnownInflationDestinations, []string{}, "known_inflation_destinations_all", KnownInflationDestinationsRoutePrefix)
	route.AddRoute("POST", "/add", AddKnownInflationDestination, []string{}, "known_inflation_destinations_add", KnownInflationDestinationsRoutePrefix)
	route.AddRoute("POST", "/edit", EditKnownInflationDestination, []string{}, "known_inflation_destinations_edit", KnownInflationDestinationsRoutePrefix)
	route.AddRoute("POST", "/delete", DeleteKnownInflationDestination, []string{}, "known_inflation_destinations_delete", KnownInflationDestinationsRoutePrefix)
	route.AddRoute("POST", "/changeOrder", ChangeOrderKnownInflationDestination, []string{}, "known_inflation_destinations_change_order", KnownInflationDestinationsRoutePrefix)
}

//AddKnownInflationDestinationsRoutes adds all the routes for the known inflation destinations handling
func AddKnownInflationDestinationsRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(KnownInflationDestinationsRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//KnownInflationDestinationIDRequest request used in get one and delete
type KnownInflationDestinationIDRequest struct {
	ID int `form:"id" json:"id"  validate:"required"`
}

//GetKnownInflationDestinationsResponse response
type GetKnownInflationDestinationsResponse struct {
	ID               int    `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"`
	ShortDescription string `form:"short_description" json:"short_description"`
	LongDescription  string `form:"long_description" json:"long_description"`
	OrderIndex       int    `form:"order_index" json:"order_index"`
}

//GetKnownInflationDestination returns inflation destination by id
func GetKnownInflationDestination(uc *mw.AdminContext, c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error inflation destinations id", cerr.GeneralError))
		return
	}
	r := KnownInflationDestinationIDRequest{ID: id}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	inflationDestination, err := db.GetKnownInflationDestinationByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destinations", cerr.GeneralError))
		return
	}

	response := GetKnownInflationDestinationsResponse{
		ID:               inflationDestination.ID,
		Name:             inflationDestination.Name,
		IssuerPublicKey:  inflationDestination.IssuerPublicKey,
		ShortDescription: inflationDestination.ShortDescription,
		LongDescription:  inflationDestination.LongDescription,
		OrderIndex:       inflationDestination.OrderIndex,
	}

	c.JSON(http.StatusOK, response)

}

//AllKnownInflationDestinations returns all inflation destinations
func AllKnownInflationDestinations(uc *mw.AdminContext, c *gin.Context) {

	inflationDestinations, err := db.GetKnownInflationDestinations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destinations", cerr.GeneralError))
		return
	}

	var response []GetKnownInflationDestinationsResponse

	for _, c := range inflationDestinations {

		response = append(response, GetKnownInflationDestinationsResponse{
			ID:               c.ID,
			Name:             c.Name,
			IssuerPublicKey:  c.IssuerPublicKey,
			ShortDescription: c.ShortDescription,
			LongDescription:  c.LongDescription,
			OrderIndex:       c.OrderIndex,
		})

	}

	c.JSON(http.StatusOK, response)

}

//AddKnownInflationDestinationRequest request
type AddKnownInflationDestinationRequest struct {
	Name             string `form:"name" json:"name"  validate:"required,max=500"`
	IssuerPublicKey  string `form:"issuer_public_key" json:"issuer_public_key"  validate:"required,max=500"`
	ShortDescription string `form:"short_description" json:"short_description"  validate:"required,max=500"`
	LongDescription  string `form:"long_description" json:"long_description"  validate:"required,max=500"`
}

//AddKnownInflationDestination adds a new inflation destination
func AddKnownInflationDestination(uc *mw.AdminContext, c *gin.Context) {
	var r AddKnownInflationDestinationRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	existsInflationDestination, err := db.ExistsKnownInflationDestination(r.IssuerPublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destination", cerr.GeneralError))
		return
	}
	if existsInflationDestination {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key and asset_code", cerr.InvalidArgument, "Inflation destination with same public key already exists", ""))
		return
	}
	orderIndex, err := db.KnownInflationDestinationNewOrderIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error geting new order index", cerr.GeneralError))
		return
	}

	inflationDestination := &models.AdminKnownInflationDestination{
		Name:             r.Name,
		IssuerPublicKey:  r.IssuerPublicKey,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		OrderIndex:       orderIndex,
		UpdatedBy:        getUpdatedBy(c)}

	err = db.AddKnownInflationDestination(inflationDestination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding known inflation destination", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//EditKnownInflationDestinationRequest request
type EditKnownInflationDestinationRequest struct {
	ID               int     `form:"id" json:"id"  validate:"required"`
	Name             *string `form:"name" json:"name"  validate:"max=500"`
	IssuerPublicKey  *string `form:"issuer_public_key" json:"issuer_public_key"  validate:"max=500"`
	ShortDescription *string `form:"short_description" json:"short_description"  validate:"max=500"`
	LongDescription  *string `form:"long_description" json:"long_description"  validate:"max=500"`
}

//EditKnownInflationDestination edits known inflation destination details
func EditKnownInflationDestination(uc *mw.AdminContext, c *gin.Context) {
	var r EditKnownInflationDestinationRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	inflationDestination, err := db.GetKnownInflationDestinationByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destination", cerr.GeneralError))
		return
	}
	if inflationDestination == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Inflation destination not found in database", ""))
		return
	}

	if inflationDestination.IssuerPublicKey != *r.IssuerPublicKey {
		existsInflationDestination, err := db.ExistsKnownInflationDestination(*r.IssuerPublicKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destination", cerr.GeneralError))
			return
		}
		if existsInflationDestination {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "public_key and asset_code", cerr.InvalidArgument, "Inflation destination with same public key already exists", ""))
			return
		}
	}

	if r.Name != nil {
		inflationDestination.Name = *r.Name
	}
	if r.IssuerPublicKey != nil {
		inflationDestination.IssuerPublicKey = *r.IssuerPublicKey
	}
	if r.ShortDescription != nil {
		inflationDestination.ShortDescription = *r.ShortDescription
	}
	if r.LongDescription != nil {
		inflationDestination.LongDescription = *r.LongDescription
	}

	err = db.UpdateKnownInflationDestination(inflationDestination, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating known inflation destination", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//DeleteKnownInflationDestination deletes known inflation destination
func DeleteKnownInflationDestination(uc *mw.AdminContext, c *gin.Context) {
	var r KnownInflationDestinationIDRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	inflationDestination, err := db.GetKnownInflationDestinationByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destination", cerr.GeneralError))
		return
	}
	if inflationDestination == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Inflation destination not found in database", ""))
		return
	}

	err = db.DeleteKnownInflationDestination(inflationDestination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting known inflation destination", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ChangeOrderKnownInflationDestinationRequest request
type ChangeOrderKnownInflationDestinationRequest struct {
	ID            int `form:"id" json:"id"  validate:"required"`
	OrderModifier int `form:"order_modifier" json:"order_modifier" validate:"required"`
}

//ChangeOrderKnownInflationDestination alters a inflation destination and changes the order index with +-1
func ChangeOrderKnownInflationDestination(uc *mw.AdminContext, c *gin.Context) {
	var r ChangeOrderKnownInflationDestinationRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	inflationDestination, err := db.GetKnownInflationDestinationByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing inflation destination", cerr.GeneralError))
		return
	}
	if inflationDestination == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Inflation destination not found in database", ""))
		return
	}
	if r.OrderModifier != -1 && r.OrderModifier != 1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "order_modifier", cerr.InvalidArgument, "Order modifier must be -1 or 1", ""))
		return
	}

	err = db.ChangeKnownInflationDestinationOrder(inflationDestination, r.OrderModifier, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating known inflation destination order", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
