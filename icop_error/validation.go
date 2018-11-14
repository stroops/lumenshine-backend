package icop_error

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/sirupsen/logrus"

	"regexp"

	"github.com/fatih/structtag"
	validator "gopkg.in/go-playground/validator.v9"
)

type vD struct {
	Code    int
	Message string
}

//validationErrorCodes are the error codes used for the Data-Validation
var validationErrorCodes = map[string]vD{
	"required":         vD{1002, "Field required"},
	"isdefault":        vD{2, "isDefault"},
	"len":              vD{1004, "Length does not match"},
	"min":              vD{1000, "Length to short"},
	"max":              vD{1000, "Length to high"},
	"eq":               vD{1000, "Does not equal"},
	"ne":               vD{1000, "ne"},
	"lt":               vD{1000, "lt"},
	"lte":              vD{1000, "lte"},
	"gt":               vD{1000, "gt"},
	"gte":              vD{1000, "gte"},
	"eqfield":          vD{12, "eqfield"},
	"eqcsfield":        vD{13, "eqcsfield"},
	"necsfield":        vD{14, "necsfield"},
	"gtcsfield":        vD{15, "gtcsfield"},
	"gtecsfield":       vD{16, "gtecsfield"},
	"ltcsfield":        vD{17, "ltcsfield"},
	"ltecsfield":       vD{18, "ltecsfield"},
	"nefield":          vD{19, "nefield"},
	"gtefield":         vD{20, "gtefield"},
	"gtfield":          vD{21, "gtfield"},
	"ltefield":         vD{22, "ltefield"},
	"ltfield":          vD{23, "ltfield"},
	"alpha":            vD{1000, "alpha"},
	"alphanum":         vD{1000, "alphanum"},
	"alphaunicode":     vD{1000, "alphaunicode"},
	"alphanumunicode":  vD{1000, "alphanumunicode"},
	"numeric":          vD{1000, "numeric"},
	"number":           vD{1000, "number"},
	"hexadecimal":      vD{1000, "hexadecimal"},
	"hexcolor":         vD{1000, "hexcolor"},
	"rgb":              vD{1000, "rgb"},
	"rgba":             vD{1000, "rgba"},
	"hsl":              vD{1000, "hsl"},
	"hsla":             vD{1000, "hsla"},
	"email":            vD{1000, "email"},
	"url":              vD{1000, "url"},
	"uri":              vD{1000, "uri"},
	"base64":           vD{1000, "base64"},
	"base64url":        vD{1000, "base64url"},
	"contains":         vD{41, "contains"},
	"containsany":      vD{42, "containsany"},
	"containsrune":     vD{43, "containsrune"},
	"excludes":         vD{44, "excludes"},
	"excludesall":      vD{45, "excludesall"},
	"excludesrune":     vD{46, "excludesrune"},
	"isbn":             vD{1000, "isbn"},
	"isbn10":           vD{1000, "isbn10"},
	"isbn13":           vD{1000, "isbn13"},
	"eth_addr":         vD{1000, "eth_addr"},
	"btc_addr":         vD{1000, "btc_addr"},
	"btc_addr_bech32":  vD{1000, "btc_addr_bech32"},
	"uuid":             vD{1000, "uuid"},
	"uuid3":            vD{1000, "uuid3"},
	"uuid4":            vD{1000, "uuid4"},
	"uuid5":            vD{1000, "uuid5"},
	"ascii":            vD{1000, "ascii"},
	"printascii":       vD{1000, "printascii"},
	"multibyte":        vD{1000, "multibyte"},
	"datauri":          vD{1000, "datauri"},
	"latitude":         vD{1000, "latitude"},
	"longitude":        vD{1000, "longitude"},
	"ssn":              vD{1000, "ssn"},
	"ipv4":             vD{1000, "ipv4"},
	"ipv6":             vD{1000, "ipv6"},
	"ip":               vD{1000, "ip"},
	"cidrv4":           vD{1000, "cidrv4"},
	"cidrv6":           vD{1000, "cidrv6"},
	"cidr":             vD{1000, "cidr"},
	"tcp4_addr":        vD{1000, "tcp4_addr"},
	"tcp6_addr":        vD{1000, "tcp6_addr"},
	"tcp_addr":         vD{1000, "tcp_addr"},
	"udp4_addr":        vD{1000, "udp4_addr"},
	"udp6_addr":        vD{1000, "udp6_addr"},
	"udp_addr":         vD{1000, "udp_addr"},
	"ip4_addr":         vD{1000, "ip4_addr"},
	"ip6_addr":         vD{1000, "ip6_addr"},
	"ip_addr":          vD{1000, "ip_addr"},
	"unix_addr":        vD{1000, "unix_addr"},
	"mac":              vD{1000, "mac"},
	"hostname":         vD{1000, "hostname"},         // RFC 952
	"hostname_rfc1123": vD{1000, "hostname_rfc1123"}, // RFC 1123
	"fqdn":             vD{1000, "fqdn"},
	"unique":           vD{84, "unique"},
	"oneof":            vD{85, "oneof"},

	//custom errors
	"icop_email":      vD{1000, "wrong email format"},
	"icop_phone":      vD{1000, "wrong phone number format"},
	"icop_assetcode":  vD{1000, "wrong asset code format"},
	"icop_devicetype": vD{1000, "wrong device type value"},
	"icop_nonum":      vD{1000, "contains numbers"},
	"min_trim":        vD{1000, "Length to short"},
}

var (
	validate        = getValidator()
	emailRegExp     *regexp.Regexp
	phoneRegExp     *regexp.Regexp
	assetCodeRegExp *regexp.Regexp
	noNumRegExp     *regexp.Regexp
)

func getValidator() *validator.Validate {
	v := validator.New()
	var err error
	emailRegExp, err = regexp.Compile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])")
	if err != nil {
		log.Fatalf("Error building email regexp %v", err)
	}

	phoneRegExp, err = regexp.Compile("^[+]?[0-9]{11,16}$")
	if err != nil {
		log.Fatalf("Error building phone regexp %v", err)
	}

	assetCodeRegExp, err = regexp.Compile("^[a-zA-Z0-9]{1,12}$")
	if err != nil {
		log.Fatalf("Error building asset code regexp %v", err)
	}

	noNumRegExp, err = regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatalf("Error building no num regexp %v", err)
	}

	v.RegisterValidation("icop_email", icopEmailValidator)
	v.RegisterValidation("icop_phone", icopPhoneValidator)
	v.RegisterValidation("icop_assetcode", icopAssetCodeValidator)
	v.RegisterValidation("icop_devicetype", icopDeviceTypeValidator)
	v.RegisterValidation("icop_nonum", icopNoNumbers)
	v.RegisterValidation("min_trim", icopMinTrim)
	return v
}

func icopPhoneValidator(fl validator.FieldLevel) bool {
	return phoneRegExp.MatchString(fl.Field().String())
}

func icopEmailValidator(fl validator.FieldLevel) bool {
	return emailRegExp.MatchString(fl.Field().String())
}

func icopAssetCodeValidator(fl validator.FieldLevel) bool {
	return assetCodeRegExp.MatchString(fl.Field().String())
}

func icopDeviceTypeValidator(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	_, ok := pb.DeviceType_value[strings.ToLower(input)]
	return ok
}

func icopNoNumbers(fl validator.FieldLevel) bool {
	return noNumRegExp.FindString(fl.Field().String()) == ""
}

func icopMinTrim(fl validator.FieldLevel) bool {
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(fl.Field().String())) >= param
}

//IcopError are the details for the validations
type IcopError struct {
	ErrorCode           int    `json:"error_code,omitempty"`
	ParameterName       string `json:"parameter_name,omitempty"`
	UserErrorMessageKey string `json:"user_error_message_key,omitempty"`
	ErrorMessage        string `json:"error_message,omitempty"`
}

//IcopErrors is the list of all errors returned
type IcopErrors []IcopError

//ValidateStruct validates the struct and returns an errorlist (or empty error slice)
func ValidateStruct(log *logrus.Entry, data interface{}) (bool, IcopErrors) {
	err := validate.Struct(data)

	if err != nil {
		e := validateErrors(log, data, err)
		return false, e
	}

	return true, []IcopError{}
}

//ValidateErrors takes the data object and validation errors and creates the returning validation-data struct
func validateErrors(log *logrus.Entry, data interface{}, validateError error) IcopErrors {
	var e []IcopError

	for _, err := range validateError.(validator.ValidationErrors) {
		errorCode, errorMessage := getErrorData(log, err.ActualTag())
		fieldName := getFieldName(log, data, err)

		e = append(
			e,
			IcopError{
				ErrorCode:     errorCode,
				ErrorMessage:  errorMessage,
				ParameterName: fieldName,
			},
		)
	}

	return e
}

//GetJSON returns the JSON for the errorstruct
func (e IcopErrors) GetJSON(log *logrus.Entry) []byte {
	b, err := json.Marshal(e)
	if err != nil {
		log.WithError(err).Error("Error constructing json")
		return nil
	}

	return b
}

//getFieldName get's either the json, form or query name for the field, based on the fields-tag
func getFieldName(log *logrus.Entry, data interface{}, fieldError validator.FieldError) string {
	var field reflect.StructField
	var ok bool

	t := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		field, ok = t.Elem().FieldByName(fieldError.Field())
	} else {
		field, ok = t.FieldByName(fieldError.Field())
	}

	//t := reflect.TypeOf(data)
	//field, ok := t.Elem().FieldByName(fieldError.Field())
	//field, ok := reflect.Indirect(reflect.ValueOf(data)).Elem().FieldByName(fieldError.Field())
	if ok {
		tag := field.Tag
		tags, err := structtag.Parse(string(tag))
		if err != nil {
			log.WithError(err).WithField("tag", tag).Error("Error on structtag.Parse:")
			return fieldError.Field()
		}

		//json as default
		jsonTag, err := tags.Get("json")
		if err != nil {
			//try form
			form, err := tags.Get("form")
			if err != nil {
				//try query param
				query, err := tags.Get("query")
				if err == nil {
					return query.Name
				}
			} else {
				return form.Name
			}
		} else {
			return jsonTag.Name
		}
	}

	return fieldError.Field()
}

//getErrorCode get's the code for the validation error
func getErrorData(log *logrus.Entry, tag string) (int, string) {
	code, ok := validationErrorCodes[tag]
	if ok {
		return code.Code, code.Message
	}
	return -1, "-NA-"
}

//NewIcopErrorShort returns a new List of icop errors, initialized with one error
func NewIcopErrorShort(errorCode int, errorMessage string) IcopErrors {
	var e []IcopError
	return append(e, IcopError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	})
}

//NewIcopError returns a new List of icop errors, initialized with one error
func NewIcopError(paramNane string, errorCode int, errorMessage string, userErrorMessageKey string) IcopErrors {
	var e []IcopError
	return append(e, IcopError{
		ErrorCode:           errorCode,
		ParameterName:       paramNane,
		ErrorMessage:        errorMessage,
		UserErrorMessageKey: userErrorMessageKey,
	})
}

//AddError adds a new Error to the response
func (e *IcopErrors) AddError(paramNane string, errorCode int, errorMessage string, userErrorMessageKey string) {
	if e == nil {
		e = new(IcopErrors)
	}
	*e = append(*e, IcopError{
		ErrorCode:           errorCode,
		ParameterName:       paramNane,
		ErrorMessage:        errorMessage,
		UserErrorMessageKey: userErrorMessageKey},
	)
}
