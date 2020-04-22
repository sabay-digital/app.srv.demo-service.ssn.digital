package demoservice

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
)

// We should really use the stream service to call the callbackHandler - that way we can be sure of payment success
func callbackHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the URL encoded values from the request body
	in, err := ioutil.ReadAll(r.Body)
	ssn.Log(err, "Callback handler: Read request body")
	req, err := url.ParseQuery(string(in))
	ssn.Log(err, "Callback handler: Parse request body")

	resp := ssn.RedirectPayload{
		RedirectURL: "/v1/success",
		Payload: []ssn.PayloadItem{
			ssn.PayloadItem{
				Key:   "hash",
				Value: req.Get("hash"),
			},
		},
	}
	redirectTemplate := template.Must(template.ParseFiles("templates/redirect.html"))
	redirectTemplate.Execute(w, resp)
}
