package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/config"
	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"
)

const (
	//PromoRoutePrefix for the promo routes
	PromoRoutePrefix = "promo"
)

var imageExtensions = map[string]int32{
	".png":  0,
	".jpg":  1,
	".jpeg": 2,
}

//init setup all the routes for the promo handling
func init() {
	route.AddRoute("GET", "/get/:id", GetPromo, []string{}, "promo_get", PromoRoutePrefix)
	route.AddRoute("GET", "/all", AllPromos, []string{}, "promo_all", PromoRoutePrefix)
	route.AddRoute("POST", "/add", AddPromo, []string{}, "promo_add", PromoRoutePrefix)
	route.AddRoute("POST", "/edit", EditPromo, []string{}, "promo_edit", PromoRoutePrefix)
	route.AddRoute("POST", "/delete", DeletePromo, []string{}, "promo_delete", PromoRoutePrefix)
	route.AddRoute("POST", "/changeOrder", ChangeOrderPromo, []string{}, "promo_change_order", PromoRoutePrefix)
	route.AddRoute("POST", "/activate", ActivatePromo, []string{}, "promo_activate", PromoRoutePrefix)
}

//AddPromoRoutes adds all the routes for the promo handling
func AddPromoRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(PromoRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//GetPromoRequest request used in get one and delete
//swagger:parameters GetPromoRequest GetPromo
type GetPromoRequest struct {
	//required : true
	ID int `form:"id" json:"id" query:"id" validate:"required"`
}

//GetPromoResponse response
// swagger:model
type GetPromoResponse struct {
	ID      int                `form:"id" json:"id"`
	Name    string             `form:"name" json:"name"`
	Title   string             `form:"title" json:"title"`
	Text    string             `form:"text" json:"text"`
	Image   PromoImageResponse `form:"image" json:"image"`
	Active  bool               `form:"active" json:"active"`
	Type    string             `form:"type" json:"type"`
	Buttons []Button           `form:"buttons" json:"buttons"`
}

//PromoImageResponse - image content and details
// swagger:model
type PromoImageResponse struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

//GetPromo returns promo by id
// swagger:route GET /portal/admin/dash/promo/get/:id Promo GetPromo
//
// Returns promo by id
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:GetPromoResponse
func GetPromo(uc *mw.AdminContext, c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading promo id", cerr.GeneralError))
		return
	}

	r := GetPromoRequest{ID: id}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	promo, err := db.GetPromoByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promo", cerr.GeneralError))
		return
	}

	var buttons []Button
	err = json.Unmarshal([]byte(promo.Buttons), &buttons)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deserializing buttons", cerr.GeneralError))
		return
	}

	response := GetPromoResponse{
		ID:      promo.ID,
		Name:    promo.Name,
		Title:   promo.Title,
		Text:    promo.PromoText,
		Type:    promo.PromoType,
		Active:  promo.Active,
		Buttons: buttons,
	}

	image, err := getImageResponse(promo.ID, promo.ImageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not read image", cerr.GeneralError))
		return
	}
	response.Image = *image

	c.JSON(http.StatusOK, &response)
}

func getImageResponse(id int, imageType string) (*PromoImageResponse, error) {
	fileName := fmt.Sprintf("%d.%v", id, imageType)
	filePath := filepath.Join(config.Cnf.Promo.ImagesPath, fileName)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	str := base64.StdEncoding.EncodeToString(content)
	mimeType := db.MimeTypes[imageType]

	return &PromoImageResponse{
		FileName: fileName,
		Content:  str,
		MimeType: mimeType}, nil
}

//AllPromos returns all promos
// swagger:route GET /portal/admin/dash/promo/all Promo AllPromos
//
// Returns all promos
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:[]GetPromoResponse
func AllPromos(uc *mw.AdminContext, c *gin.Context) {
	promos, err := db.GetPromos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promos", cerr.GeneralError))
		return
	}

	var response []GetPromoResponse
	for _, promo := range promos {
		var buttons []Button
		err = json.Unmarshal([]byte(promo.Buttons), &buttons)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deserializing buttons", cerr.GeneralError))
			return
		}

		item := GetPromoResponse{
			ID:      promo.ID,
			Name:    promo.Name,
			Title:   promo.Title,
			Text:    promo.PromoText,
			Type:    promo.PromoType,
			Active:  promo.Active,
			Buttons: buttons,
		}
		image, err := getImageResponse(promo.ID, promo.ImageType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not read image", cerr.GeneralError))
			return
		}
		item.Image = *image

		response = append(response, item)
	}
	c.JSON(http.StatusOK, &response)
}

//AddPromoRequest request
//swagger:parameters AddPromoRequest AddPromo
type AddPromoRequest struct {
	//required : true
	Name  string `form:"name" json:"name"  validate:"required,max=256"`
	Title string `form:"title" json:"title" validate:"max=512"`
	Text  string `form:"text" json:"text"`
	//required : true
	Type string `form:"type" json:"type" validate:"required,max=32"`
	//required : true
	Buttons []Button `form:"buttons" json:"buttons" validate:"required"`
}

//Button - one button
// swagger:model
type Button struct {
	Name string `form:"name" json:"name"`
	Link string `form:"link" json:"link"`
}

//AddPromo adds a new promo
// swagger:route POST /portal/admin/dash/promo/add Promo AddPromo
//
// Adds a new promo
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func AddPromo(uc *mw.AdminContext, c *gin.Context) {
	var r AddPromoRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if r.Type != models.PromoTypeSmall && r.Type != models.PromoTypeBig {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "type", cerr.InvalidArgument, "Invalid promo type", ""))
		return
	}

	if r.Type == models.PromoTypeBig {
		if r.Title == "" {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "title", cerr.InvalidArgument, "Missing title", ""))
			return
		}
		if r.Text == "" {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "text", cerr.InvalidArgument, "Missing text", ""))
			return
		}
	}

	if r.Buttons == nil || len(r.Buttons) == 0 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "buttons", cerr.InvalidArgument, "Missing buttons", ""))
		return
	}

	buttons, err := json.Marshal(r.Buttons)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "buttons", cerr.InvalidArgument, "Error serializing buttons", ""))
		return
	}

	orderIndex, err := db.PromoNewOrderIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error geting new order index", cerr.GeneralError))
		return
	}

	file, err := c.FormFile("upload_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Error reading upload image", cerr.InvalidArgument))
		return
	}
	ext := filepath.Ext(strings.TrimSpace(file.Filename))
	if _, ok := imageExtensions[ext]; !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("upload_image", cerr.InvalidArgument, "Valid image extensions are: png, jpg, jpeg.", ""))
		return
	}
	if _, err = os.Stat(config.Cnf.Promo.ImagesPath); os.IsNotExist(err) {
		if err = os.MkdirAll(config.Cnf.Promo.ImagesPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating file path", cerr.GeneralError))
			return
		}
	}

	promo := &models.AdminPromo{
		Name:       r.Name,
		Title:      r.Title,
		PromoText:  r.Text,
		PromoType:  r.Type,
		Active:     false,
		ImageType:  ext[1:],
		Buttons:    string(buttons),
		OrderIndex: orderIndex,
		UpdatedBy:  getUpdatedBy(c)}

	err = db.AddPromo(promo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding promo", cerr.GeneralError))
		return
	}

	fileName := fmt.Sprintf("%d%v", promo.ID, ext)
	filePath := filepath.Join(config.Cnf.Promo.ImagesPath, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error saving uploaded image", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//EditPromoRequest request
//swagger:parameters EditPromoRequest EditPromo
type EditPromoRequest struct {
	//required : true
	ID    int     `form:"id" json:"id"  validate:"required"`
	Name  *string `form:"name" json:"name"  validate:"omitempty,max=256"`
	Title *string `form:"title" json:"title" validate:"omitempty,max=512"`
	Text  *string `form:"text" json:"text"`
	//required : true
	Type *string `form:"type" json:"type" validate:"omitempty,max=32"`
	//required : true
	Buttons []Button `form:"buttons" json:"buttons"`
}

//EditPromo edits promo details
// swagger:route POST /portal/admin/dash/promo/edit Promo EditPromo
//
// Edits promo details
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func EditPromo(uc *mw.AdminContext, c *gin.Context) {
	var r EditPromoRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	promo, err := db.GetPromoByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promo", cerr.GeneralError))
		return
	}
	if promo == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Promo not found in database", ""))
		return
	}

	if r.Type != nil {
		if *r.Type != models.PromoTypeSmall && *r.Type != models.PromoTypeBig {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "type", cerr.InvalidArgument, "Invalid promo type", ""))
			return
		}
		if *r.Type == models.PromoTypeBig {
			if promo.Title == "" && r.Title == nil {
				c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "title", cerr.InvalidArgument, "Missing title", ""))
				return
			}
			if promo.PromoText == "" && r.Text == nil {
				c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "text", cerr.InvalidArgument, "Missing text", ""))
				return
			}
		}
	}

	if r.Name != nil {
		promo.Name = *r.Name
	}
	if r.Title != nil {
		promo.Title = *r.Title
	}
	if r.Text != nil {
		promo.PromoText = *r.Text
	}
	if r.Type != nil {
		promo.PromoType = *r.Type
	}
	if r.Buttons != nil && len(r.Buttons) > 0 {
		buttons, err := json.Marshal(r.Buttons)
		if err != nil {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "buttons", cerr.InvalidArgument, "Error serializing buttons", ""))
			return
		}
		promo.Buttons = string(buttons)
	}

	err = db.UpdatePromo(promo, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating promo", cerr.GeneralError))
		return
	}

	file, err := c.FormFile("upload_image")
	if err == nil && file.Size > 0 {
		ext := filepath.Ext(strings.TrimSpace(file.Filename))
		if _, ok := imageExtensions[ext]; !ok {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("upload_image", cerr.InvalidArgument, "Valid image extensions are: png, jpg, jpeg.", ""))
			return
		}
		fileName := fmt.Sprintf("%d%v", promo.ID, ext)
		filePath := filepath.Join(config.Cnf.Promo.ImagesPath, fileName)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error saving uploaded image", cerr.GeneralError))
			return
		}
	}

	c.JSON(http.StatusOK, "{}")
}

//PromoDeleteRequest request used in get one and delete
//swagger:parameters PromoDeleteRequest DeletePromo
type PromoDeleteRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
}

//DeletePromo deletes the given promo
// swagger:route POST /portal/admin/dash/promo/delete Promo DeletePromo
//
// Deletes the given promo
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func DeletePromo(uc *mw.AdminContext, c *gin.Context) {
	var r PromoDeleteRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	promo, err := db.GetPromoByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promo", cerr.GeneralError))
		return
	}
	if promo == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Promo not found in database", ""))
		return
	}

	err = db.DeletePromo(promo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting promo", cerr.GeneralError))
		return
	}

	fileName := fmt.Sprintf("%d.%v", promo.ID, promo.ImageType)
	filePath := filepath.Join(config.Cnf.Promo.ImagesPath, fileName)
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting image", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ChangeOrderPromoRequest request
//swagger:parameters ChangeOrderPromoRequest ChangeOrderPromo
type ChangeOrderPromoRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
	//required : true
	OrderModifier int `form:"order_modifier" json:"order_modifier" validate:"required"`
}

//ChangeOrderPromo alters a promo and changes the order index with +-1
// swagger:route POST /portal/admin/dash/promo/changeOrder Promo ChangeOrderPromo
//
// Alters a promo and changes the order index with +-1
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func ChangeOrderPromo(uc *mw.AdminContext, c *gin.Context) {
	var r ChangeOrderPromoRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	promo, err := db.GetPromoByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promo", cerr.GeneralError))
		return
	}
	if promo == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Promo not found in database", ""))
		return
	}
	if r.OrderModifier != -1 && r.OrderModifier != 1 {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "order_modifier", cerr.InvalidArgument, "Order modifier must be -1 or 1", ""))
		return
	}

	err = db.ChangePromoOrder(promo, r.OrderModifier, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating promo order", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//ActivatePromoRequest request
//swagger:parameters ActivatePromoRequest ActivatePromo
type ActivatePromoRequest struct {
	//required : true
	ID int `form:"id" json:"id"  validate:"required"`
	//required : true
	Active bool `form:"active" json:"active"`
}

//ActivatePromo activates/deactivates the promo
// swagger:route POST /portal/admin/dash/promo/changeOrder Promo ActivatePromo
//
// Activates/deactivates the promo
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func ActivatePromo(uc *mw.AdminContext, c *gin.Context) {
	var r ActivatePromoRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	promo, err := db.GetPromoByID(r.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promo", cerr.GeneralError))
		return
	}
	if promo == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "id", cerr.InvalidArgument, "Promo not found in database", ""))
		return
	}

	promo.Active = r.Active

	err = db.UpdatePromo(promo, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error activating/deactivating promo", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
