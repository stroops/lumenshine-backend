package history

import (
	"github.com/Soneso/lumenshine-backend/addons/charts/config"
	"io"
	"net/http"
	"net/url"
	"os"
)

// postData posts data of the following form in order to retrieve csv files having the exchange history from source to dest
// http://fx.sauder.ubc.ca/cgi/fxdata
// b: USD - from currency
// c: EUR - to currency
// rd:
// fd: 8 - from day
// fm: 8 - from month
// fy: 2014 - from year
// ld: 31 - to day
// lm: 12 - to month
// ly: 2018 - to year
// y: daily - data frequency
// q: volume - notation
// f: csv - file format
// o:
func postData(data config.TransactionData) (localPath string, err error) {

	// Create the file
	localPath = config.Cnf.LocalDownloadDir + data.SourceCurrency + ".csv"
	out, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer out.Close()

	resp, err := http.PostForm(data.FileURL,
		url.Values{"b": {data.SourceCurrency}, "c": {data.DestinationCurrency}, "fd": {"8"}, "fm": {"8"}, "fy": {"2014"}, "ld": {"31"}, "lm": {"12"}, "ly": {"2018"}, "y": {"daily"}, "q": {"volume"}, "f": {"csv"}})

	if nil != err {
		return
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}
