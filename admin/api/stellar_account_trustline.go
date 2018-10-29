package api

import (
	"net/http"
	"strings"

	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	"github.com/Soneso/lumenshine-backend/admin/route"
	"github.com/Soneso/lumenshine-backend/db/pageinate"
	qq "github.com/Soneso/lumenshine-backend/db/querying"
	coremodels "github.com/Soneso/lumenshine-backend/db/stellarcore/models"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/services/db/modext"

	"github.com/gin-gonic/gin"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

//init setup all the routes for the users handling
func init() {
	route.AddRoute("POST", "/add_unathorized_trustline", AddTrustline, []string{"Administrators"}, "add_unathorized_trustline", StellarAccountRoutePrefix)
	route.AddRoute("POST", "/remove_unathorized_trustline", RemoveTrustline, []string{"Administrators"}, "remove_unathorized_trustline", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/worker_account_trustlines/:publickey", WorkerAccountTrustlines, []string{"Administrators"}, "worker_account_trustlines", StellarAccountRoutePrefix)
	route.AddRoute("GET", "/search_trusting_accounts", IssuerAccountTrustlines, []string{"Administrators"}, "search_trusting_accounts", StellarAccountRoutePrefix)
}

//AddTrustlineRequest - info
//swagger:parameters AddTrustlineRequest AddTrustline
type AddTrustlineRequest struct {
	//required : true
	TrustorPublicKey string `form:"trusting_account_public_key" json:"trusting_account_public_key" validate:"required,base64,len=56"`
	//required : true
	IssuingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
	//required : true
	Status string `form:"status" json:"status" validate:"required,max=50"`
	//required : true
	Reason string `form:"reason" json:"reason" validate:"required,max=1000"`
}

//AddTrustline creates new entry in the db
// swagger:route POST /portal/admin/dash/stellar_account/add_unathorized_trustline StellarAccount AddTrustline
//
// Creates new entry in the db
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
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
//swagger:parameters RemoveTrustlineRequest RemoveTrustline
type RemoveTrustlineRequest struct {
	//required : true
	TrustorPublicKey string `form:"trusting_account_public_key" json:"trusting_account_public_key" validate:"required,base64,len=56"`
	//required : true
	IssuingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	//required : true
	AssetCode string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
}

//RemoveTrustline deletes trustline
// swagger:route POST /portal/admin/dash/stellar_account/remove_unathorized_trustline StellarAccount RemoveTrustline
//
// Deletes trustline from the DB
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
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
//swagger:parameters WorkerAccountTrustlinesRequest WorkerAccountTrustlines
type WorkerAccountTrustlinesRequest struct {
	//required : true
	PublicKey string `form:"public_key" json:"public_key" validate:"required,base64,len=56"`
}

//WorkerTrustlineItem - response item
// swagger:model
type WorkerTrustlineItem struct {
	AssetCode string `json:"asset_code"`
	Issuer    string `json:"asset_issuer"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

//WorkerAccountTrustlines - returns worker account trustlines
// swagger:route GET /portal/admin/dash/stellar_account/worker_account_trustlines/:publickey StellarAccount WorkerAccountTrustlines
//
// Returns worker account trustlines
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:[]WorkerTrustlineItem
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
		status := models.StellarTrustlineStatusWaiting
		reason := ""
		if trustline.Flags == 1 {
			status = models.StellarTrustlineStatusOk
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

//SearchAccountsRequest for filtering the issuer's trusting accounts
//swagger:parameters SearchAccountsRequest IssuerAccountTrustlines
type SearchAccountsRequest struct {
	pageinate.PaginationRequestStruct
	//required : true
	IssueingPublicKey string `form:"issuing_account_public_key" json:"issuing_account_public_key" validate:"required,base64,len=56"`
	//required : true
	AssetCode       string `form:"asset_code" json:"asset_code" validate:"required,icop_assetcode"`
	FilterName      string `form:"filter_name" json:"filter_name"`
	FilterPublicKey string `form:"filter_public_key" json:"filter_public_key" validate:"omitempty,base64,len=56"`
	//required : true
	FilterType     string   `form:"filter_type" json:"filter_type" validate:"required"`
	FilterStatuses []string `form:"filter_statuses" json:"filter_statuses"`
}

//SearchAccountsItem is one item in the list
// swagger:model
type SearchAccountsItem struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	AssetCode string `json:"asset_code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

//SearchAccountsResponse list of accounts
// swagger:model
type SearchAccountsResponse struct {
	pageinate.PaginationResponseStruct
	Items []SearchAccountsItem `json:"items"`
}

//IssuerAccountTrustlines returns list of all accounts, filtered by given params
// swagger:route GET /portal/admin/dash/stellar_account/search_trusting_accounts StellarAccount IssuerAccountTrustlines
//
// Returns list of all accounts, filtered by given params
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:SearchAccountsResponse
func IssuerAccountTrustlines(uc *mw.AdminContext, c *gin.Context) {
	var err error
	var rr SearchAccountsRequest
	if err = c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	customerType := "customer"
	workerType := "worker"

	if !strings.EqualFold(rr.FilterType, customerType) && !strings.EqualFold(rr.FilterType, workerType) {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnIcopError(uc.Log, "filter_type", cerr.InvalidArgument, "Filter type is not customer or worker.", ""))
		return
	}

	viewName := modext.CustomerTrustlinesViewName
	accountType := customerType
	if strings.EqualFold(rr.FilterType, workerType) {
		viewName = modext.AdminTrustlinesViewName
		accountType = workerType
	}

	selectQ := []qm.QueryMod{
		qm.Select(
			modext.IssuerTrustlineColumns.Name,
			modext.IssuerTrustlineColumns.PublicKey,
			modext.IssuerTrustlineColumns.IssuerPublicKey,
			modext.IssuerTrustlineColumns.AssetCode,
			modext.IssuerTrustlineColumns.Status,
			modext.IssuerTrustlineColumns.Reason,
		),
		qm.From(viewName),
	}
	countQ := []qm.QueryMod{
		qm.Select("count(*) as total_count"),
		qm.From(viewName),
	}
	filterQ := []qm.QueryMod{qm.Where(modext.IssuerTrustlineColumns.IssuerPublicKey+" = ?", rr.IssueingPublicKey),
		qm.Where(modext.IssuerTrustlineColumns.AssetCode+" = ?", rr.AssetCode),
	}

	if rr.FilterName != "" {
		filterQ = append(filterQ, qm.Where(modext.IssuerTrustlineColumns.Name+" ilike ?", qq.Like(rr.FilterName)))
	}
	if rr.FilterPublicKey != "" {
		filterQ = append(filterQ, qm.Where(modext.IssuerTrustlineColumns.PublicKey+" ilike ?", qq.Like(rr.FilterPublicKey)))
	}
	if len(rr.FilterStatuses) > 0 {
		filter := make([]interface{}, len(rr.FilterStatuses))
		for i := range rr.FilterStatuses {
			filter[i] = rr.FilterStatuses[i]
		}
		filterQ = append(filterQ, qm.WhereIn(modext.IssuerTrustlineColumns.Status+" in ? ", filter...))
	}

	response := new(SearchAccountsResponse)

	count := pageinate.CountStruct{}
	countQ = append(countQ, filterQ...)
	err = models.NewQuery(countQ...).Bind(nil, db.DB, &count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error counting trustlines", cerr.GeneralError))
		return
	}
	response.TotalCount = count.TotalCount

	selectQ = append(selectQ, filterQ...)
	qP := pageinate.Paginate(selectQ, &rr.PaginationRequestStruct, &response.PaginationResponseStruct)
	itemList := make([]modext.IssuerTrustline, 0)
	err = models.NewQuery(qP...).Bind(nil, db.DB, &itemList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading trustlines", cerr.GeneralError))
		return
	}

	response.Items = make([]SearchAccountsItem, len(itemList))
	for i, item := range itemList {
		response.Items[i] = SearchAccountsItem{
			Name:      item.Name,
			PublicKey: item.PublicKey,
			AssetCode: item.AssetCode,
			Type:      accountType,
			Status:    item.Status,
			Reason:    item.Reason,
		}
	}

	c.JSON(http.StatusOK, response)
}
