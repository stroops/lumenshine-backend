package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

var (
	validExtensions = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".pdf": true}
	maxFileSize     = 2100000
)

//UploadKycDocumentRequest is the data needed for the kyc document
type UploadKycDocumentRequest struct {
	DocumentType     string    `form:"document_type" json:"document_type" validate:"required"`
	DocumentSide     string    `form:"document_side" json:"document_side" validate:"required"`
	IDCountryCode    string    `form:"id_country_code" json:"id_country_code" validate:"max=20"`
	IDIssueDate      time.Time `form:"id_issue_date" json:"id_issue_date" time_format:"2006-01-02"`
	IDExpirationDate time.Time `form:"id_expiration_date" json:"id_expiration_date" time_format:"2006-01-02"`
	IDNumber         string    `form:"id_number" json:"id_number" validate:"max=100"`
}

//UploadKycDocumentResponse - response
type UploadKycDocumentResponse struct {
	DocumentID int64 `form:"document_id" json:"document_id"`
}

//UploadKycDocument - stores the document on the server
func UploadKycDocument(uc *mw.IcopContext, c *gin.Context) {
	var rr UploadKycDocumentRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if _, ok := pb.DocumentType_value[rr.DocumentType]; !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("document_type", cerr.InvalidArgument, "Invalid document type value.", ""))
		return
	}
	if _, ok := pb.DocumentSide_value[rr.DocumentSide]; !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("document_side", cerr.InvalidArgument, "Valid document side values are: front, back.", ""))
		return
	}
	file, err := c.FormFile("upload_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Error reading upload file", cerr.InvalidArgument))
		return
	}
	ext := trimLeftChar(filepath.Ext(strings.TrimSpace(file.Filename)))
	if _, ok := pb.DocumentFormat_value[ext]; !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("upload_file", cerr.InvalidArgument, "Valid file extensions are: png, jpg, jpeg, pdf.", ""))
		return
	}
	if file.Size > int64(maxFileSize) {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("upload_file", cerr.InvalidArgument, "Max file size is 2MB", ""))
		return
	}

	if rr.DocumentType == pb.DocumentType_passport.String() ||
		rr.DocumentType == pb.DocumentType_drivers_license.String() ||
		rr.DocumentType == pb.DocumentType_id_card.String() {

		if rr.IDNumber == "" {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("id_number", cerr.InvalidArgument, "Id number is required", ""))
			return
		}
		if rr.IDCountryCode == "" {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("id_country_code", cerr.InvalidArgument, "Id country code is required", ""))
			return
		}
		if rr.IDIssueDate.IsZero() {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("id_issue_date", cerr.InvalidArgument, "Id issue date is required", ""))
			return
		}
		if rr.IDExpirationDate.IsZero() {
			c.JSON(http.StatusBadRequest, cerr.NewIcopError("id_expiration_date", cerr.InvalidArgument, "Id expiration date is required", ""))
			return
		}
	}

	userID := mw.GetAuthUser(c).UserID
	if _, err = os.Stat(cnf.Kyc.DocumentsPath); os.IsNotExist(err) {
		if err = os.MkdirAll(cnf.Kyc.DocumentsPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating file path", cerr.GeneralError))
			return
		}
	}

	reqData := &pb.AddKycDocumentRequest{
		Base:             NewBaseRequest(uc),
		UserId:           userID,
		DocumentType:     pb.DocumentType(pb.DocumentType_value[rr.DocumentType]),
		DocumentFormat:   pb.DocumentFormat(pb.DocumentFormat_value[ext]),
		DocumentSide:     pb.DocumentSide(pb.DocumentSide_value[rr.DocumentSide]),
		IdCountryCode:    rr.IDCountryCode,
		IdNumber:         rr.IDNumber,
		IdIssueDate:      rr.IDIssueDate.Unix(),
		IdExpirationDate: rr.IDExpirationDate.Unix(),
	}
	resp, err := dbClient.AddKycDocument(c, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error adding kyc document", cerr.GeneralError))
		return
	}

	fileName := fmt.Sprintf("%d.%v", resp.DocumentId, ext)
	filePath := filepath.Join(cnf.Kyc.DocumentsPath, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error saving upload file", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &UploadKycDocumentResponse{DocumentID: resp.DocumentId})
}

func trimLeftChar(s string) string {
	if s == "" {
		return ""
	}
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
