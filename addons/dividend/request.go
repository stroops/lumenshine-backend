package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

//Req holds the current request instance
var Req *Request

//Request holds the dividend request information
type Request struct {
	AssetCode    string
	Issuer       string
	DividendMode DividendMode
	Blacklist    []string
}

//DividendMode holds the dividend mode information
type DividendMode struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

//ModeSum type of dividend distribution
const ModeSum = "SUM"

//ModePercent type of dividend distribution
const ModePercent = "PERCENT"

//InitRequest creates a new request based on the specified flags
func InitRequest(cmd *cobra.Command) error {
	Req = new(Request)

	assetCode, err := cmd.Flags().GetString("assetcode")
	if err != nil {
		log.Fatalf("Error reading asset code flag. %v", err)
	}
	if assetCode == "" {
		log.Fatalf("Missing asset code flag. Use -a flag.")
	}

	issuer, err := cmd.Flags().GetString("issuer")
	if err != nil {
		log.Fatalf("Error reading issuer flag. %v", err)
	}
	if issuer == "" {
		log.Fatalf("Missing issuer flag. Use -i flag.")
	}

	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		log.Fatalf("Error reading mode flag. %v", err)
	}
	if mode == "" {
		log.Fatalf("Missing mode flag. Use -m flag.")
	}
	mode = strings.Replace(mode, "'", "\"", -1)

	var dividendMode DividendMode
	err = json.Unmarshal([]byte(mode), &dividendMode)
	if err != nil {
		log.Fatalf("Error decoding dividend mode json. %v", err)
	}

	if !strings.EqualFold(dividendMode.Type, ModeSum) && !strings.EqualFold(dividendMode.Type, ModePercent) {
		log.Fatalf("Invlid dividende mode type: %v. Possible values: %v, %v", dividendMode.Type, ModeSum, ModePercent)
	}

	if strings.EqualFold(dividendMode.Type, ModePercent) && (dividendMode.Value < 0 || dividendMode.Value > 100) {
		log.Fatalf("The value for mode PERCENT must be between 0 and 100")
	}

	var blacklist []string
	b, err := cmd.Flags().GetString("blacklist")
	if err != nil {
		log.Fatalf("Error reading blacklist flag. %v", err)
	}

	if len(b) > 0 {
		b = strings.Replace(b, "'", "\"", -1)
		err = json.Unmarshal([]byte(b), &blacklist)
		if err != nil {
			log.Fatalf("Error decoding blacklist json. %v", err)
		}
	}
	Req.AssetCode = assetCode
	Req.Issuer = issuer
	Req.DividendMode = dividendMode
	Req.Blacklist = blacklist

	return nil
}
