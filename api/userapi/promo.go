package main

import (
	"net/http"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

//GetPromoResponse - promo response item
// swagger:model
type GetPromoResponse struct {
	ID      int           `form:"id" json:"id"`
	Name    string        `form:"name" json:"name"`
	Title   string        `form:"title" json:"title"`
	Text    string        `form:"text" json:"text"`
	Image   PromoImage    `form:"image" json:"image"`
	Active  bool          `form:"active" json:"active"`
	Type    string        `form:"type" json:"type"`
	Buttons []PromoButton `form:"buttons" json:"buttons"`
}

//PromoImage - image content and details
// swagger:model
type PromoImage struct {
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

//PromoButton - one button
// swagger:model
type PromoButton struct {
	Name string `form:"name" json:"name"`
	Link string `form:"link" json:"link"`
}

//GetPromos returns all active promo cards
// swagger:route GET /portal/user/dashboard/get_promo_cards promo GetPromos
//
// Returns all active promo cards
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []GetPromoResponse
func GetPromos(uc *mw.IcopContext, c *gin.Context) {

	res, err := adminAPIClient.GetPromos(c, &pb.Empty{Base: NewBaseRequest(uc)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting promos", cerr.GeneralError))
		return
	}

	promos := make([]GetPromoResponse, len(res.Promos))
	for i, promo := range res.Promos {
		buttons := make([]PromoButton, len(promo.Buttons))
		for j, button := range promo.Buttons {
			buttons[j] = PromoButton{Name: button.Name, Link: button.Link}
		}
		promos[i] = GetPromoResponse{
			ID:      int(promo.Id),
			Name:    promo.Name,
			Title:   promo.Title,
			Text:    promo.Text,
			Type:    promo.Type,
			Buttons: buttons,
			Image:   PromoImage{Content: promo.Image.Content, MimeType: promo.Image.MimeType},
		}
	}

	c.JSON(http.StatusOK, &promos)
}
