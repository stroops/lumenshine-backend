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
	RegistrationDate  time.Time  `form:"registration_date" json:"registration_date"`
	MobileNR          string     `form:"mobile_nr" json:"mobile_nr"`
	StreetAddress     string     `form:"street_address" json:"street_address"`
	StreetNumber      string     `form:"street_number" json:"street_number"`
	ZipCode           string     `form:"zip_code" json:"zip_code"`
	City              string     `form:"city" json:"city"`
	State             string     `form:"state" json:"state"`
	CountryCode       string     `form:"country_code" json:"country_code"`
	Nationality       string     `form:"nationality" json:"nationality"`
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
		Email:             u.Email,
		Forename:          u.Forename,
		Lastname:          u.Lastname,
		RegistrationDate:  time.Unix(u.CreatedAt, 0),
		MobileNR:          u.MobileNr,
		StreetAddress:     u.StreetAddress,
		StreetNumber:      u.StreetNumber,
		ZipCode:           u.ZipCode,
		City:              u.City,
		State:             u.State,
		CountryCode:       u.CountryCode,
		Nationality:       u.Nationality,
		BirthPlace:        u.BirthPlace,
		AdditionalName:    u.AdditionalName,
		BirthCountryCode:  u.BirthCountryCode,
		BankAccountNumber: u.BankAccountNumber,
		BankNumber:        u.BankNumber,
		BankPhoneNumber:   u.BankPhoneNumber,
		TaxID:             u.TaxId,
		TaxIDName:         u.TaxIdName,
		Occupation:        u.Occupation,
		EmployerName:      u.EmployerName,
		EmployerAddress:   u.EmployerAddress,
		LanguageCode:      u.LanguageCode,
	}
	birthDay := time.Unix(u.BirthDay, 0)
	if !birthDay.IsZero() {
		response.BirthDay = &birthDay
	}
	c.JSON(http.StatusOK, &response)
}
