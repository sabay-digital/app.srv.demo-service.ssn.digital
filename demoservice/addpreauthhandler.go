package demoservice

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
)

func addPreauthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the URL encoded values from the request body
	in, err := ioutil.ReadAll(r.Body)
	ssn.Log(err, "Callback handler: Read request body")
	req, err := url.ParseQuery(string(in))
	ssn.Log(err, "Callback handler: Parse request body")

	var existingPreauth []PreAuthorization
	db.Where("user_pk = ? AND pp_pk = ?", req.Get("public_key"), req.Get("pp_pk")).Find(&existingPreauth)
	fmt.Println(len(existingPreauth))
	fmt.Println(existingPreauth)

	if len(existingPreauth) == 0 {
		newPreAuth := PreAuthorization{
			UserID:     req.Get("user_id"),
			UserPubkey: req.Get("public_key"),
			PPPubkey:   req.Get("pp_pk"),
		}
		db.Create(&newPreAuth)
	}

	resp := ssn.RedirectPayload{
		RedirectURL: req.Get("request_url"),
		Payload: []ssn.PayloadItem{
			ssn.PayloadItem{
				Key:   "hash",
				Value: req.Get("hash"),
			},
			ssn.PayloadItem{
				Key:   "signature",
				Value: req.Get("signature"),
			},
			ssn.PayloadItem{
				Key:   "public_key",
				Value: req.Get("public_key"),
			},
			ssn.PayloadItem{
				Key:   "redirect",
				Value: req.Get("redirect"),
			},
		},
	}
	redirectTemplate := template.Must(template.ParseFiles("templates/redirect.html"))
	redirectTemplate.Execute(w, resp)
}
