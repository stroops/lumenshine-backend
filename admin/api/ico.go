package api

import (
	"net/http"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"

	"github.com/gin-gonic/gin"
)

const (
	//ICORoutePrefix for the ico manangement
	ICORoutePrefix = "ico"
)

//init setup all the routes for the known currencies handling
func init() {
	route.AddRoute("GET", "/list", ICOList, []string{}, "ico_list", ICORoutePrefix)
	//route.AddRoute("GET", "/get/:id", GetKnownCurrency, []string{}, "known_currencies_get", KnownCurrenciesRoutePrefix)

}

//AddICORoutes adds all the routes for the ico management
func AddICORoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(ICORoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//ICOListResponse response
type ICOListResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//ICOList returns the list of ICOs
func ICOList(uc *mw.AdminContext, c *gin.Context) {

	icos, err := m.Icos().All(db.DBC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing currencies", cerr.GeneralError))
		return
	}

	var response []ICOListResponse

	for _, ico := range icos {

		response = append(response, ICOListResponse{
			ID:   ico.ID,
			Name: ico.IcoName,
		})

	}

	c.JSON(http.StatusOK, response)

}
