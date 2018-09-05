package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

func init() {
	route.AddRoute("GET", "/kyc_details/:id", KycDetails, []string{}, "kyc_details", CustomerRoutePrefix)
}

//KycDocument - document item
type KycDocument struct {
	ID               int       `json:"id"`
	Type             string    `json:"type"`
	Format           string    `json:"format"`
	Side             string    `json:"side"`
	UploadDate       time.Time `json:"upload_date"`
	IDCountryCode    string    `json:"id_country_code"`
	IDNumber         string    `json:"id_number"`
	IDIssueDate      time.Time `json:"id_issue_date"`
	IDExpirationDate time.Time `json:"id_expiration_date"`
}

//KycDetailsResponse - kyc details
type KycDetailsResponse struct {
	Status    string        `json:"status"`
	Documents []KycDocument `json:"documents"`
}

//KycDetails returns details of kyc status and documents
func KycDetails(uc *mw.AdminContext, c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error parsing id", cerr.GeneralError))
		return
	}

	u, err := m.UserProfiles(
		qm.Where(m.UserProfileColumns.ID+"=?", id),
		qm.Select(
			m.UserProfileColumns.ID,
			m.UserProfileColumns.KycStatus,
		),
	).One(db.DBC)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting user from db", cerr.GeneralError))
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("id", cerr.UserNotExists, "User does not exist in db", ""))
		return
	}
	docs, err := m.UserKycDocuments(
		qm.Where(m.UserKycDocumentColumns.UserID+"=?", id),
		qm.OrderBy(
			m.UserKycDocumentColumns.UploadDate,
		),
	).All(db.DBC)

	resp := KycDetailsResponse{Status: u.KycStatus}

	resp.Documents = make([]KycDocument, len(docs))
	for i, c := range docs {
		resp.Documents[i] = KycDocument{
			ID:               c.ID,
			Type:             c.Type,
			Format:           c.Format,
			Side:             c.Side,
			UploadDate:       c.UploadDate,
			IDCountryCode:    c.IDCountryCode,
			IDNumber:         c.IDNumber,
			IDIssueDate:      c.IDIssueDate,
			IDExpirationDate: c.IDExpirationDate,
		}
	}

	c.JSON(http.StatusOK, resp)
}
