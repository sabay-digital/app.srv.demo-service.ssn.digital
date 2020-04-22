package demoservice

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
)

// SuccessData is parsed by the success.html template
type SuccessData struct {
	Hash string
}

func successHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the URL encoded values from the request body
	in, err := ioutil.ReadAll(r.Body)
	ssn.Log(err, "Success handler: Read request body")
	req, err := url.ParseQuery(string(in))
	ssn.Log(err, "Success handler: Parse request body")

	out := SuccessData{
		Hash: req.Get("hash"),
	}

	// Load the redirect form
	successTemplate := template.Must(template.ParseFiles("templates/success.html"))
	// Serve the form
	successTemplate.Execute(w, out)
}
