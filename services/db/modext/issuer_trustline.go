package modext

// IssuerTrustline is an object representing the database view.
type IssuerTrustline struct {
	Name            string `boil:"name" json:"name" toml:"name" yaml:"name"`
	PublicKey       string `boil:"public_key" json:"public_key" toml:"public_key" yaml:"public_key"`
	IssuerPublicKey string `boil:"issuer_public_key" json:"issuer_public_key" toml:"issuer_public_key" yaml:"issuer_public_key"`
	AssetCode       string `boil:"asset_code" json:"asset_code" toml:"asset_code" yaml:"asset_code"`
	Status          string `boil:"status" json:"status" toml:"status" yaml:"status"`
	Reason          string `boil:"reason" json:"reason" toml:"reason" yaml:"reason"`
}

//IssuerTrustlineColumns - column names
var IssuerTrustlineColumns = struct {
	Name            string
	PublicKey       string
	IssuerPublicKey string
	AssetCode       string
	Authorized      string
	Status          string
	Reason          string
}{
	Name:            "name",
	PublicKey:       "public_key",
	IssuerPublicKey: "issuer_public_key",
	AssetCode:       "asset_code",
	Authorized:      "authorized",
	Status:          "status",
	Reason:          "reason",
}

var AdminTrustlinesViewName = "admin_trustlines"
var CustomerTrustlinesViewName = "customer_trustlines"
