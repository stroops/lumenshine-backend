package api

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Soneso/lumenshine-backend/admin/config"
	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/route"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

func init() {
	route.AddRoute("GET", "/kyc_details/:id", KycDetails, []string{}, "kyc_details", CustomerRoutePrefix)
	route.AddRoute("GET", "/kyc_document/:id", KycDocumentDownload, []string{}, "kyc_document", CustomerRoutePrefix)
}

//KycDocument - document item
// swagger:model
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
// swagger:model
type KycDetailsResponse struct {
	Status    string        `json:"status"`
	Documents []KycDocument `json:"documents"`
}

//KycDetails returns details of kyc status and documents
// swagger:route GET /portal/admin/dash/customer/kyc_details/:id kyc KycDetails
//
// Returns details of kyc status and documents
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:KycDetailsResponse
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
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting documents from db", cerr.GeneralError))
		return
	}

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

	c.JSON(http.StatusOK, &resp)
}

//KycDocumentDownloadResponse - document content and details
// swagger:model
type KycDocumentDownloadResponse struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

//KycDocumentDownload returns the document
// swagger:route GET /portal/admin/dash/customer/kyc_document/:id kyc KycDocumentDownload
//
// Returns the document
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:KycDocumentDownloadResponse
func KycDocumentDownload(uc *mw.AdminContext, c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Error parsing id", cerr.GeneralError))
		return
	}

	doc, err := m.UserKycDocuments(qm.Where(m.UserKycDocumentColumns.ID+"=?", id)).One(db.DBC)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting document from db", cerr.GeneralError))
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("id", cerr.InvalidArgument, "Document does not exist in db", ""))
		return
	}
	fileName := fmt.Sprintf("%d.%v", id, doc.Format)
	filePath := filepath.Join(config.Cnf.Kyc.DocumentsPath, fileName)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading document file", cerr.GeneralError))
		return
	}

	str := base64.StdEncoding.EncodeToString(content)
	mimeType := MimeTypes[doc.Format]
	c.JSON(http.StatusOK, &KycDocumentDownloadResponse{
		FileName: fileName,
		Content:  str,
		MimeType: mimeType})
}
