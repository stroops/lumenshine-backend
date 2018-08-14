package api

import (
	"database/sql"
	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/gin-gonic/gin"

	"github.com/Soneso/lumenshine-backend/admin/db"
	"github.com/Soneso/lumenshine-backend/admin/route"
	"github.com/Soneso/lumenshine-backend/db/pageinate"
	qq "github.com/Soneso/lumenshine-backend/db/querying"

	"strconv"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
)

const (
	//CustomerRoutePrefix is the prefix for the customer group. We need this in order to get all the routes for this base url
	CustomerRoutePrefix = "customer"
)

//init setup all the routes for the users handling
func init() {
	route.AddRoute("GET", "/list", CustomerList, []string{}, "customer_list", CustomerRoutePrefix)
	route.AddRoute("GET", "/details/:id", CustomerDetails, []string{}, "customer_details", CustomerRoutePrefix)
	route.AddRoute("POST", "/update_personal_data", CustomerEdit, []string{}, "customer_update_personal_data", CustomerRoutePrefix)
	route.AddRoute("GET", "/orders/:id", CustomerOrders, []string{}, "customer_orders", CustomerRoutePrefix)
	route.AddRoute("GET", "/wallets/:id", CustomerWallets, []string{}, "customer_wallets", CustomerRoutePrefix)
}

//AddCustomerRoutes adds all the routes for the user handling
func AddCustomerRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(CustomerRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

//CustomerListRequest for filtering the customers
type CustomerListRequest struct {
	pageinate.PaginationRequestStruct

	FilterCustomerID int      `form:"filter_customer_id"`
	FilterForeName   string   `form:"filter_forename"`
	FilterLastName   string   `form:"filter_lastname"`
	FilterEmail      string   `form:"filter_email"`
	FilterKycStatus  []string `form:"filter_kyc_status"`

	SortCustomerID       string `form:"sort_customer_id"`
	SortForeName         string `form:"sort_forename"`
	SortLastName         string `form:"sort_lastname"`
	SortEmail            string `form:"sort_email"`
	SortRegistrationDate string `form:"sort_registration_date"`
	SortLastLogin        string `form:"sort_last_login"`
}

//CustomerListItem is one item in the list
type CustomerListItem struct {
	ID               int       `json:"id"`
	Forename         string    `json:"forename"`
	Lastname         string    `json:"last_name"`
	Email            string    `json:"email"`
	KycStatus        string    `json:"kyc_status"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLogin        time.Time `json:"last_login"`
}

//CustomerListResponse list of customers
type CustomerListResponse struct {
	pageinate.PaginationResponseStruct
	Items []CustomerListItem `json:"items"`
}

//CustomerList returns list of all customers, filtered by given params
func CustomerList(uc *mw.AdminContext, c *gin.Context) {
	var err error
	var rr CustomerListRequest
	if err = c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	//this is the initial queryMod
	//we will append queries and sorting to it
	q := []qm.QueryMod{
		qm.Select(
			m.UserProfileColumns.ID,
			m.UserProfileColumns.Forename,
			m.UserProfileColumns.Lastname,
			m.UserProfileColumns.Email,
			m.UserProfileColumns.KycStatus,
			m.UserProfileColumns.CreatedAt,
		),
	}

	if rr.FilterCustomerID != 0 {
		q = append(q, qm.Where("id=?", rr.FilterCustomerID))
	}

	if rr.FilterForeName != "" {
		q = append(q, qm.Where("forename ilike ?", qq.Like(rr.FilterForeName)))
	}

	if rr.FilterLastName != "" {
		q = append(q, qm.Where("lastname ilike ?", qq.Like(rr.FilterLastName)))
	}

	if rr.FilterEmail != "" {
		q = append(q, qm.Where("email ilike ? ", qq.Like(rr.FilterEmail)))
	}

	if len(rr.FilterKycStatus) > 0 {
		filter := make([]interface{}, len(rr.FilterKycStatus))
		for i := range rr.FilterKycStatus {
			filter[i] = rr.FilterKycStatus[i]
		}
		q = append(q, qm.WhereIn("kyc_status in ? ", filter...))
	}

	r := new(CustomerListResponse)

	//we need to get the total count before sorting and applying the pagination
	r.TotalCount, err = m.UserProfiles(db.DBC, q...).Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting total count", cerr.GeneralError))
		return
	}

	//apply sorting
	q = qq.AddSorting(rr.SortCustomerID, m.UserProfileColumns.ID, q)
	q = qq.AddSorting(rr.SortForeName, m.UserProfileColumns.Forename, q)
	q = qq.AddSorting(rr.SortLastName, m.UserProfileColumns.Lastname, q)
	q = qq.AddSorting(rr.SortEmail, m.UserProfileColumns.Email, q)
	q = qq.AddSorting(rr.SortRegistrationDate, m.UserProfileColumns.CreatedAt, q)

	qP := pageinate.Paginate(q, &rr.PaginationRequestStruct, &r.PaginationResponseStruct)
	customers, err := m.UserProfiles(db.DBC, qP...).All()
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting users", cerr.GeneralError))
		return
	}

	r.Items = make([]CustomerListItem, len(customers))
	for i, c := range customers {
		r.Items[i] = CustomerListItem{
			ID:               c.ID,
			Forename:         c.Forename,
			Lastname:         c.Lastname,
			Email:            c.Email,
			KycStatus:        c.KycStatus,
			RegistrationDate: c.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, r)
}

type CustomerDetailsResponse struct {
	ID               int       `json:"id"`
	Forename         string    `json:"forename"`
	Lastname         string    `json:"lastname"`
	Email            string    `json:"email"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLogin        time.Time `json:"last_login"`
	MobileNR         string    `json:"mobile_nr"`
	StreetAddress    string    `json:"street_address"`
	StreetNumber     string    `json:"street_number"`
	ZipCode          string    `json:"zip_code"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	CountryCode      string    `json:"country_code"`
	Nationality      string    `json:"nationality"`
	BirthDay         time.Time `json:"birth_day"`
	BirthPlace       string    `json:"birth_place"`
}

//CustomerDetails returns details of stefiied customer
func CustomerDetails(uc *mw.AdminContext, c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error parsing id", cerr.GeneralError))
		return
	}

	u, err := m.UserProfiles(db.DBC,
		qm.Where("id=?", id),
		qm.Select(
			m.UserProfileColumns.ID,
			m.UserProfileColumns.Forename,
			m.UserProfileColumns.Lastname,
			m.UserProfileColumns.Email,
			m.UserProfileColumns.CreatedAt,
			m.UserProfileColumns.MobileNR,
			m.UserProfileColumns.StreetAddress,
			m.UserProfileColumns.StreetNumber,
			m.UserProfileColumns.ZipCode,
			m.UserProfileColumns.City,
			m.UserProfileColumns.State,
			m.UserProfileColumns.CountryCode,
			m.UserProfileColumns.Nationality,
			m.UserProfileColumns.BirthDay,
			m.UserProfileColumns.BirthPlace,
		),
	).One()
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting user from db", cerr.GeneralError))
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("id", cerr.UserNotExists, "User does not exist in db", ""))
		return
	}

	c.JSON(http.StatusOK, &CustomerDetailsResponse{
		ID:               u.ID,
		Forename:         u.Forename,
		Lastname:         u.Lastname,
		Email:            u.Email,
		RegistrationDate: u.CreatedAt,
		MobileNR:         u.MobileNR,
		StreetAddress:    u.StreetAddress,
		StreetNumber:     u.StreetNumber,
		ZipCode:          u.ZipCode,
		City:             u.City,
		State:            u.State,
		CountryCode:      u.CountryCode,
		Nationality:      u.Nationality,
		BirthDay:         u.BirthDay,
		BirthPlace:       u.BirthPlace,
	})
}

type CustomerEditRequest struct {
	ID            int    `form:"id" json:"id"`
	Forename      string `form:"forename" json:"forename" validate:"required,max=64"`
	Lastname      string `form:"lastname" json:"lastname" validate:"required,max=64"`
	MobileNR      string `form:"mobile_nr" json:"mobile_nr" validate:"required,max=64"`
	StreetAddress string `form:"street_address" json:"street_address" validate:"required,max=128"`
	StreetNumber  string `form:"street_number" json:"street_number" validate:"required,max=128"`
	ZipCode       string `form:"zip_code" json:"zip_code" validate:"required,max=32"`
	City          string `form:"city" json:"city" validate:"required,max=128"`
	State         string `form:"state" json:"state" validate:"required,max=128"`
	CountryCode   string `form:"country_code" json:"country_code" validate:"required,max=128"`
	Nationality   string `form:"nationality" json:"nationality" validate:"required,max=2"`
	BirthPlace    string `form:"birth_place" json:"birth_place" validate:"required,max=128"`
}

//CustomerEdit updates customer details and returns customer
func CustomerEdit(uc *mw.AdminContext, c *gin.Context) {
	var err error
	var rr CustomerEditRequest
	if err = c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	u, err := m.UserProfiles(db.DBC,
		qm.Where("id=?", rr.ID),
		qm.Select(
			m.UserProfileColumns.ID,
			m.UserProfileColumns.Forename,
			m.UserProfileColumns.Lastname,
			m.UserProfileColumns.Email,
			m.UserProfileColumns.CreatedAt,
			m.UserProfileColumns.MobileNR,
			m.UserProfileColumns.StreetAddress,
			m.UserProfileColumns.StreetNumber,
			m.UserProfileColumns.ZipCode,
			m.UserProfileColumns.City,
			m.UserProfileColumns.State,
			m.UserProfileColumns.CountryCode,
			m.UserProfileColumns.Nationality,
			m.UserProfileColumns.BirthDay,
			m.UserProfileColumns.BirthPlace,
		),
	).One()
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting user from db", cerr.GeneralError))
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("id", cerr.UserNotExists, "User does not exist in db", ""))
		return
	}

	u.Forename = rr.Forename
	u.Lastname = rr.Lastname
	u.MobileNR = rr.MobileNR
	u.StreetAddress = rr.StreetAddress
	u.StreetNumber = rr.StreetNumber
	u.ZipCode = rr.ZipCode
	u.City = rr.City
	u.State = rr.State
	u.CountryCode = rr.CountryCode
	u.Nationality = rr.Nationality
	u.BirthPlace = rr.BirthPlace

	err = u.Update(db.DBC, m.UserProfileColumns.ID,
		m.UserProfileColumns.Forename,
		m.UserProfileColumns.Lastname,
		m.UserProfileColumns.MobileNR,
		m.UserProfileColumns.StreetAddress,
		m.UserProfileColumns.StreetNumber,
		m.UserProfileColumns.ZipCode,
		m.UserProfileColumns.City,
		m.UserProfileColumns.State,
		m.UserProfileColumns.CountryCode,
		m.UserProfileColumns.Nationality,
		m.UserProfileColumns.BirthPlace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &CustomerDetailsResponse{
		ID:               u.ID,
		Forename:         u.Forename,
		Lastname:         u.Lastname,
		Email:            u.Email,
		RegistrationDate: u.CreatedAt,
		MobileNR:         u.MobileNR,
		StreetAddress:    u.StreetAddress,
		StreetNumber:     u.StreetNumber,
		ZipCode:          u.ZipCode,
		City:             u.City,
		State:            u.State,
		CountryCode:      u.CountryCode,
		Nationality:      u.Nationality,
		BirthDay:         u.BirthDay,
		BirthPlace:       u.BirthPlace,
	})
}

//CustomerOrdersRequest to get the orders
type CustomerOrdersRequest struct {
	pageinate.PaginationRequestStruct
}

//OrderListItem is one item in the list
type OrderListItem struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	Amount int64     `json:"amount"`
	Price  float64   `json:"price"`
	Chain  string    `json:"chain"`
	Status string    `json:"status"`
}

//OrderListResponse list of orders
type OrderListResponse struct {
	pageinate.PaginationResponseStruct
	Items []OrderListItem `json:"items"`
}

//CustomerOrders returns list of all customer's orders
func CustomerOrders(uc *mw.AdminContext, c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error parsing id", cerr.GeneralError))
		return
	}

	var rr CustomerOrdersRequest
	if err = c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	q := []qm.QueryMod{
		qm.Select(
			m.UserOrderColumns.ID,
			m.UserOrderColumns.CoinAmount,
			m.UserOrderColumns.ChainAmount,
			m.UserOrderColumns.Chain,
			m.UserOrderColumns.OrderStatus,
			m.UserOrderColumns.CreatedAt,
		),
	}
	q = append(q, qm.Where(m.UserOrderColumns.UserID+"=?", id))

	r := new(OrderListResponse)

	//we need to get the total count before sorting and applying the pagination
	r.TotalCount, err = m.UserOrders(db.DBC, q...).Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting total count", cerr.GeneralError))
		return
	}

	q = append(q, qm.OrderBy(m.UserOrderColumns.CreatedAt))

	qP := pageinate.Paginate(q, &rr.PaginationRequestStruct, &r.PaginationResponseStruct)
	orders, err := m.UserOrders(db.DBC, qP...).All()
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting orders", cerr.GeneralError))
		return
	}

	r.Items = make([]OrderListItem, len(orders))
	for i, o := range orders {
		r.Items[i] = OrderListItem{
			ID:     o.ID,
			Date:   o.CreatedAt,
			Amount: o.CoinAmount,
			Price:  o.ChainAmount,
			Chain:  o.Chain,
			Status: o.OrderStatus,
		}
	}
	c.JSON(http.StatusOK, r)
}

//CustomerWalletsRequest to get the wallets
type CustomerWalletsRequest struct {
	pageinate.PaginationRequestStruct
}

//WalletListItem is one item in the list
type WalletListItem struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

//WalletListResponse list of wallets
type WalletListResponse struct {
	pageinate.PaginationResponseStruct
	Items []WalletListItem `json:"items"`
}

//CustomerWallets returns list of all customer's orders
func CustomerWallets(uc *mw.AdminContext, c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error parsing id", cerr.GeneralError))
		return
	}

	var rr CustomerWalletsRequest
	if err = c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	q := []qm.QueryMod{
		qm.Select(
			m.UserWalletColumns.WalletName,
			m.UserWalletColumns.PublicKey0,
		),
	}
	q = append(q, qm.Where(m.UserWalletColumns.UserID+"=?", id))

	r := new(WalletListResponse)

	//we need to get the total count before sorting and applying the pagination
	r.TotalCount, err = m.UserWallets(db.DBC, q...).Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting total count", cerr.GeneralError))
		return
	}

	q = append(q, qm.OrderBy(m.UserWalletColumns.CreatedAt))

	qP := pageinate.Paginate(q, &rr.PaginationRequestStruct, &r.PaginationResponseStruct)
	wallets, err := m.UserWallets(db.DBC, qP...).All()
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting wallets", cerr.GeneralError))
		return
	}

	r.Items = make([]WalletListItem, len(wallets))
	for i, w := range wallets {
		r.Items[i] = WalletListItem{
			Name:      w.WalletName,
			PublicKey: w.PublicKey0,
		}
	}
	c.JSON(http.StatusOK, r)
}
