package main

import (
	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/gin-gonic/gin"
)

//GetUserDataResponse - user data response
// swagger:model
type GetUserDataResponse struct {
	Email             string     `json:"email"`
	Forename          string     `json:"forename"`
	Lastname          string     `json:"lastname"`
	Salutation        string     `json:"salutation"`
	Address           string     `json:"address"`
	ZipCode           string     `json:"zip_code"`
	City              string     `json:"city"`
	State             string     `json:"state"`
	CountryCode       string     `json:"country_code"`
	Nationality       string     `json:"nationality"`
	MobileNR          string     `json:"mobile_nr"`
	BirthDay          *time.Time `json:"birth_day"`
	BirthPlace        string     `json:"birth_place"`
	AdditionalName    string     `json:"additional_name"`
	BirthCountryCode  string     `json:"birth_country_code"`
	BankAccountNumber string     `json:"bank_account_number"`
	BankNumber        string     `json:"bank_number"`
	BankPhoneNumber   string     `json:"bank_phone_number"`
	TaxID             string     `json:"tax_id"`
	TaxIDName         string     `json:"tax_id_name"`
	OccupationName    string     `json:"occupation_name"`
	OccupationCode08  string     `json:"occupation_code08"`
	OccupationCode88  string     `json:"occupation_code88"`
	EmployerName      string     `json:"employer_name"`
	EmployerAddress   string     `json:"employer_address"`
	LanguageCode      string     `json:"language_code"`
	RegistrationDate  time.Time  `json:"registration_date"`
	ShowMemos         bool       `json:"show_memos"`
}

//GetUserData - returns the authenticated user's data
// swagger:route GET /portal/user/dashboard/get_user_data userdata GetUserData
//
// Returns the authenticated user's data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: GetUserDataResponse
func GetUserData(uc *mw.IcopContext, c *gin.Context) {
	user := mw.GetAuthUser(c)
	u, err := dbClient.GetUserProfile(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   user.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting userProfile", cerr.GeneralError))
		return
	}
	response := GetUserDataResponse{
		Email:            u.Email,
		Forename:         u.Forename,
		Lastname:         u.Lastname,
		Salutation:       u.Salutation,
		Address:          u.Address,
		ZipCode:          u.ZipCode,
		City:             u.City,
		State:            u.State,
		CountryCode:      u.CountryCode,
		Nationality:      u.Nationality,
		MobileNR:         u.MobileNr,
		BirthPlace:       u.BirthPlace,
		AdditionalName:   u.AdditionalName,
		BirthCountryCode: u.BirthCountryCode,
		TaxID:            u.TaxId,
		TaxIDName:        u.TaxIdName,
		OccupationName:   u.OccupationName,
		OccupationCode08: u.OccupationCode08,
		OccupationCode88: u.OccupationCode88,
		EmployerName:     u.EmployerName,
		EmployerAddress:  u.EmployerAddress,
		LanguageCode:     u.LanguageCode,
		RegistrationDate: time.Unix(u.CreatedAt, 0),
		ShowMemos:        u.ShowMemos,
	}
	birthDay := time.Unix(u.BirthDay, 0)
	if !birthDay.IsZero() {
		response.BirthDay = &birthDay
	}
	size := len(u.BankAccountNumber)
	if size > 4 {
		response.BankAccountNumber = u.BankAccountNumber[size-4:]
	} else {
		response.BankAccountNumber = u.BankAccountNumber
	}
	size = len(u.BankNumber)
	if size > 4 {
		response.BankNumber = u.BankNumber[size-4:]
	} else {
		response.BankNumber = u.BankNumber
	}
	size = len(u.BankPhoneNumber)
	if size > 4 {
		response.BankPhoneNumber = u.BankPhoneNumber[size-4:]
	} else {
		response.BankPhoneNumber = u.BankPhoneNumber
	}
	c.JSON(http.StatusOK, &response)
}

//UpdateUserDataRequest - edit user request
//swagger:parameters UpdateUserDataRequest UpdateUserData
type UpdateUserDataRequest struct {
	Forename          string  `form:"forename" json:"forename" validate:"required,icop_nonum,min_trim=2,max=64"`
	Lastname          string  `form:"lastname" json:"lastname" validate:"required,icop_nonum,min_trim=2,max=64"`
	Salutation        string  `form:"salutation" json:"salutation" validate:"max=64"`
	Address           string  `form:"address" json:"address" validate:"max=512"`
	ZipCode           string  `form:"zip_code" json:"zip_code" validate:"max=32"`
	City              string  `form:"city" json:"city" validate:"max=128"`
	State             string  `form:"state" json:"state" validate:"max=128"`
	CountryCode       string  `form:"country_code" json:"country_code" validate:"max=2"`
	Nationality       string  `form:"nationality" json:"nationality" validate:"max=128"`
	MobileNR          string  `form:"mobile_nr" json:"mobile_nr" validate:"max=64"`
	BirthDay          string  `form:"birth_day" json:"birth_day"`
	BirthPlace        string  `form:"birth_place" json:"birth_place" validate:"max=128"`
	AdditionalName    string  `form:"additional_name" json:"additional_name" validate:"max=256"`
	BirthCountryCode  string  `form:"birth_country_code" json:"birth_country_code" validate:"max=2"`
	BankAccountNumber *string `form:"bank_account_number" json:"bank_account_number" validate:"omitempty,max=256"`
	BankNumber        *string `form:"bank_number" json:"bank_number" validate:"omitempty,max=256"`
	BankPhoneNumber   *string `form:"bank_phone_number" json:"bank_phone_number" validate:"omitempty,max=256"`
	TaxID             string  `form:"tax_id" json:"tax_id" validate:"max=256"`
	TaxIDName         string  `form:"tax_id_name" json:"tax_id_name" validate:"max=256"`
	OccupationName    string  `form:"occupation_name" json:"occupation_name" validate:"max=256"`
	OccupationCode08  string  `form:"occupation_code08" json:"occupation_code08" validate:"max=8"`
	OccupationCode88  string  `form:"occupation_code88" json:"occupation_code88" validate:"max=8"`
	EmployerName      string  `form:"employer_name" json:"employer_name" validate:"max=512"`
	EmployerAddress   string  `form:"employer_address" json:"employer_address" validate:"max=512"`
	LanguageCode      string  `form:"language_code" json:"language_code" validate:"max=16"`
}

//UpdateUserData - updates the authenticated user's data
// swagger:route POST /portal/user/dashboard/update_user_data userdata UpdateUserData
//
// Updates the authenticated user's data
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateUserData(uc *mw.IcopContext, c *gin.Context) {
	rr := new(UpdateUserDataRequest)
	if err := c.Bind(rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	user := mw.GetAuthUser(c)
	//get the birthday
	birthDay, err := time.Parse("2006-01-02", rr.BirthDay)
	if rr.BirthDay != "" && err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("birth_day", cerr.InvalidArgument, "Birthday wrong format", ""))
		return
	}

	reqC := &pb.UpdateUserProfileRequest{
		Base:             NewBaseRequest(uc),
		Id:               user.UserID,
		Forename:         rr.Forename,
		Lastname:         rr.Lastname,
		Salutation:       rr.Salutation,
		Address:          rr.Address,
		ZipCode:          rr.ZipCode,
		City:             rr.City,
		State:            rr.State,
		CountryCode:      rr.CountryCode,
		Nationality:      rr.Nationality,
		MobileNr:         rr.MobileNR,
		BirthDay:         birthDay.Unix(),
		BirthPlace:       rr.BirthPlace,
		AdditionalName:   rr.AdditionalName,
		BirthCountryCode: rr.BirthCountryCode,
		TaxId:            rr.TaxID,
		TaxIdName:        rr.TaxIDName,
		OccupationName:   rr.OccupationName,
		OccupationCode08: rr.OccupationCode08,
		OccupationCode88: rr.OccupationCode88,
		EmployerName:     rr.EmployerName,
		EmployerAddress:  rr.EmployerAddress,
		LanguageCode:     rr.LanguageCode,
	}
	if rr.BankAccountNumber != nil {
		reqC.BankAccountNumber = *rr.BankAccountNumber
	} else {
		reqC.BankAccountNumber = pb.StringNotSet
	}
	if rr.BankNumber != nil {
		reqC.BankNumber = *rr.BankNumber
	} else {
		reqC.BankNumber = pb.StringNotSet
	}
	if rr.BankPhoneNumber != nil {
		reqC.BankPhoneNumber = *rr.BankPhoneNumber
	} else {
		reqC.BankPhoneNumber = pb.StringNotSet
	}
	_, err = dbClient.UpdateUserProfile(c, reqC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating user", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}

//UpdateUserShowMemoRequest - edit user show memo field request
//swagger:parameters UpdateUserShowMemoRequest UpdateUserShowMemo
type UpdateUserShowMemoRequest struct {
	ShowMemo bool `form:"show_memo" json:"show_memo"`
}

//UpdateUserShowMemo - updates the user's show_memo field
// swagger:route POST /portal/user/dashboard/update_user_show_memo userdata UpdateUserShowMemo
//
// Updates the authenticated user's data show_memo field
//
// 	  Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200:
func UpdateUserShowMemo(uc *mw.IcopContext, c *gin.Context) {
	rr := new(UpdateUserShowMemoRequest)
	if err := c.Bind(rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}
	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}
	user := mw.GetAuthUser(c)

	reqC := &pb.UpdateProfileShowMemoRequest{
		Base:   NewBaseRequest(uc),
		UserId: user.UserID,
		Value:  rr.ShowMemo,
	}

	_, err := dbClient.UpdateUserShowMemos(c, reqC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating user show-memo", cerr.GeneralError))
		return
	}
	c.JSON(http.StatusOK, "{}")
}
