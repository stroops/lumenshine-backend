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
type GetUserDataResponse struct {
	Email             string     `form:"email" json:"email"`
	Forename          string     `form:"forename" json:"forename"`
	Lastname          string     `form:"lastname" json:"lastname"`
	Company           string     `form:"company" json:"company"`
	Salutation        string     `form:"salutation" json:"salutation"`
	Title             string     `form:"title" json:"title"`
	Address           string     `form:"address" json:"address"`
	ZipCode           string     `form:"zip_code" json:"zip_code"`
	City              string     `form:"city" json:"city"`
	State             string     `form:"state" json:"state"`
	CountryCode       string     `form:"country_code" json:"country_code"`
	Nationality       string     `form:"nationality" json:"nationality"`
	MobileNR          string     `form:"mobile_nr" json:"mobile_nr"`
	BirthDay          *time.Time `form:"birth_day" json:"birth_day"`
	BirthPlace        string     `form:"birth_place" json:"birth_place"`
	AdditionalName    string     `form:"additional_name" json:"additional_name"`
	BirthCountryCode  string     `form:"birth_country_code" json:"birth_country_code"`
	BankAccountNumber string     `form:"bank_account_number" json:"bank_account_number"`
	BankNumber        string     `form:"bank_number" json:"bank_number"`
	BankPhoneNumber   string     `form:"bank_phone_number" json:"bank_phone_number"`
	TaxID             string     `form:"tax_id" json:"tax_id"`
	TaxIDName         string     `form:"tax_id_name" json:"tax_id_name"`
	Occupation        string     `form:"occupation" json:"occupation"`
	EmployerName      string     `form:"employer_name" json:"employer_name"`
	EmployerAddress   string     `form:"employer_address" json:"employer_address"`
	LanguageCode      string     `form:"language_code" json:"language_code"`
	RegistrationDate  time.Time  `form:"registration_date" json:"registration_date"`
}

//GetUserData - returns the authenticated user's data
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
		Company:          u.Company,
		Salutation:       u.Salutation,
		Title:            u.Title,
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
		BankNumber:       u.BankNumber,
		BankPhoneNumber:  u.BankPhoneNumber,
		TaxID:            u.TaxId,
		TaxIDName:        u.TaxIdName,
		Occupation:       u.Occupation,
		EmployerName:     u.EmployerName,
		EmployerAddress:  u.EmployerAddress,
		LanguageCode:     u.LanguageCode,
		RegistrationDate: time.Unix(u.CreatedAt, 0),
	}
	birthDay := time.Unix(u.BirthDay, 0)
	if !birthDay.IsZero() {
		response.BirthDay = &birthDay
	}
	banlength := len(u.BankAccountNumber)
	if banlength > 4 {
		response.BankAccountNumber = u.BankAccountNumber[banlength-4:]
	}
	c.JSON(http.StatusOK, &response)
}

//UpdateUserDataRequest - edit user request
type UpdateUserDataRequest struct {
	Forename          string  `form:"forename" json:"forename" validate:"max=64"`
	Lastname          string  `form:"lastname" json:"lastname" validate:"max=64"`
	Company           string  `form:"company" json:"company" validate:"max=128"`
	Salutation        string  `form:"salutation" json:"salutation" validate:"max=64"`
	Title             string  `form:"title" json:"title" validate:"max=64"`
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
	BankAccountNumber *string `form:"bank_account_number" json:"bank_account_number" validate:"max=256"`
	BankNumber        *string `form:"bank_number" json:"bank_number" validate:"max=256"`
	BankPhoneNumber   *string `form:"bank_phone_number" json:"bank_phone_number" validate:"max=256"`
	TaxID             string  `form:"tax_id" json:"tax_id" validate:"max=256"`
	TaxIDName         string  `form:"tax_id_name" json:"tax_id_name" validate:"max=256"`
	Occupation        string  `form:"occupation" json:"occupation" validate:"max=8"`
	EmployerName      string  `form:"employer_name" json:"employer_name" validate:"max=512"`
	EmployerAddress   string  `form:"employer_address" json:"employer_address" validate:"max=512"`
	LanguageCode      string  `form:"language_code" json:"language_code" validate:"max=16"`
}

//UpdateUserData - updates the authenticated user's data
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
	u, err := dbClient.GetUserProfile(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   user.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting userProfile", cerr.GeneralError))
		return
	}

	reqC := &pb.UpdateUserProfileRequest{
		Base:              NewBaseRequest(uc),
		Id:                user.UserID,
		Forename:          rr.Forename,
		Lastname:          rr.Lastname,
		Company:           rr.Company,
		Salutation:        rr.Salutation,
		Title:             rr.Title,
		Address:           rr.Address,
		ZipCode:           rr.ZipCode,
		City:              rr.City,
		State:             rr.State,
		CountryCode:       rr.CountryCode,
		Nationality:       rr.Nationality,
		MobileNr:          rr.MobileNR,
		BirthDay:          birthDay.Unix(),
		BirthPlace:        rr.BirthPlace,
		AdditionalName:    rr.AdditionalName,
		BirthCountryCode:  rr.BirthCountryCode,
		BankAccountNumber: u.BankAccountNumber,
		BankNumber:        u.BankNumber,
		BankPhoneNumber:   u.BankPhoneNumber,
		TaxId:             rr.TaxID,
		TaxIdName:         rr.TaxIDName,
		Occupation:        rr.Occupation,
		EmployerName:      rr.EmployerName,
		EmployerAddress:   rr.EmployerAddress,
		LanguageCode:      rr.LanguageCode,
	}
	if rr.BankAccountNumber != nil {
		reqC.BankAccountNumber = *rr.BankAccountNumber
	}
	if rr.BankNumber != nil {
		reqC.BankNumber = *rr.BankNumber
	}
	if rr.BankPhoneNumber != nil {
		reqC.BankPhoneNumber = *rr.BankPhoneNumber
	}
	_, err = dbClient.UpdateUserProfile(c, reqC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
