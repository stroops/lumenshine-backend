package api

// Some constants used in the APIs
const (
	BaseAPIUrl     = "/portal/admin/"
	BaseAPIUrlProt = "/portal/admin/dash/"
)

//MimeTypes - dictionary of mime types used
var MimeTypes map[string]string

func init() {
	MimeTypes = map[string]string{"pdf": "application/pdf", "png": "image/png", "jpg": "image/jpeg", "jpeg": "image/jpeg"}
}
