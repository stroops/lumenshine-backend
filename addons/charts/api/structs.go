package api

type errorCodesData struct {
	code        int
	parameter   string
	description string
}

type errorData struct {
	ErrorCode     int    `json:"error_code"`
	ParameterName string `json:"parameter_name"`
	ErrorMessage  string `json:"error_message"`
}
