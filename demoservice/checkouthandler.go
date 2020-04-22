package demoservice

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
	"github.com/stellar/go/exp/crypto/derivation"
	"github.com/stellar/go/keypair"
)

type CheckoutResponse struct {
	ExistingPreauth []PP
	Preauth         []PP
	Onetime         []PP
}

type PP struct {
	Pp_name     string
	Pp_pk       string
	User_ID     string
	Request_URL string
	Hash        string
	Signature   string
	Public_key  string
	Redirect    string
}

func checkoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the URL encoded values from the request body
	in, err := ioutil.ReadAll(r.Body)
	ssn.Log(err, "addPreauthHandler: Read request body")
	req, err := url.ParseQuery(string(in))
	ssn.Log(err, "addPreauthHandler: Parse URL encoded values")

	// Get PPs
	ppReq, _ := http.NewRequest(http.MethodGet, ssnAPI+"/accounts/"+serviceSK.Address(), nil)

	// Execute the request
	res, _ := http.DefaultClient.Do(ppReq)
	defer res.Body.Close()

	// Read the request response
	body, _ := ioutil.ReadAll(res.Body)

	account := ssn.Account{}
	// Take the JSON apart
	json.Unmarshal(body, &account)

	// Get users PK
	userID, err := strconv.Atoi(req.Get("user_id"))
	ssn.Log(err, "addPreAuthHandler: Convert user ID to int")
	path := fmt.Sprintf(derivation.StellarAccountPathFormat, userID)
	key, err := derivation.DeriveForPath(path, preauthSeed)
	ssn.Log(err, "addPreAuthHandler: Derive BIP44 keypair seed")
	kp, err := keypair.FromRawSeed(key.RawSeed())
	ssn.Log(err, "addPreAuthHandler: Derive BIP44 keypair")

	existingpreauth := []PP{}
	// Get existing preauth's
	var existingPreauth []PreAuthorization
	db.Where("user_id = ?", req.Get("user_id")).Find(&existingPreauth)
	fmt.Println(len(existingPreauth))
	fmt.Println(existingPreauth)
	if len(existingPreauth) != 0 {
		for i := range existingPreauth {
			// Get PPs
			ppReq, _ := http.NewRequest(http.MethodGet, ssnAPI+"/accounts/"+existingPreauth[i].PPPubkey, nil)

			// Execute the request
			res, _ := http.DefaultClient.Do(ppReq)
			defer res.Body.Close()

			// Read the request response
			body, _ := ioutil.ReadAll(res.Body)

			account := ssn.Account{}
			// Take the JSON apart
			json.Unmarshal(body, &account)

			existingpreauthURL := "https://" + account.Home_domain + "/v1/charge/auth/" + req.Get("payment_address")

			// Hash the URI
			existingpreauthMesg := sha256.New()
			existingpreauthMesg.Write([]byte(existingpreauthURL))

			// Sign the hash
			existingpreauthSig, err := kp.Sign(existingpreauthMesg.Sum(nil))
			ssn.Log(err, "checkoutHandler: Sign message")

			existingpreauth = append(existingpreauth, PP{
				Pp_name:     account.Data.Service_name,
				Request_URL: existingpreauthURL,
				Hash:        hex.EncodeToString(existingpreauthMesg.Sum(nil)),
				Signature:   hex.EncodeToString(existingpreauthSig),
				Public_key:  kp.Address(),
				Redirect:    "https://demo-service.testing.ssn.digital/v1/callback",
			})
		}
	}

	preauth := []PP{}
	onetime := []PP{}
	for i := range account.Balances {
		if account.Balances[i].Is_authorized && account.Balances[i].Asset_type != "native" {
			// Preauth payload
			preauthURL := "https://" + account.Balances[i].Asset_issuer_home_domain + "/v1/authorize/" + kp.Address() + "/" + serviceSK.Address()

			// Hash the URI
			preauthMesg := sha256.New()
			preauthMesg.Write([]byte(preauthURL))

			// Sign the hash
			preauthSig, err := kp.Sign(preauthMesg.Sum(nil))
			ssn.Log(err, "checkoutHandler: Sign message")

			preauth = append(preauth, PP{
				Pp_name:     account.Balances[i].Asset_issuer_service_name,
				Pp_pk:       account.Balances[i].Asset_issuer,
				User_ID:     req.Get("user_id"),
				Request_URL: preauthURL,
				Hash:        hex.EncodeToString(preauthMesg.Sum(nil)),
				Signature:   hex.EncodeToString(preauthSig),
				Public_key:  kp.Address(),
				Redirect:    "https://demo-service.testing.ssn.digital/v1/callback",
			})

			// Onetime payload
			onetimeURL := "https://" + account.Balances[i].Asset_issuer_home_domain + "/v1/charge/onetime/" + req.Get("payment_address")

			// Hash the URI
			onetimeMesg := sha256.New()
			onetimeMesg.Write([]byte(onetimeURL))

			// Sign the hash
			onetimeSig, err := serviceSK.Sign(onetimeMesg.Sum(nil))
			ssn.Log(err, "checkoutHandler: Sign message")

			onetime = append(onetime, PP{
				Pp_name:     account.Balances[i].Asset_issuer_service_name,
				Request_URL: onetimeURL,
				Hash:        hex.EncodeToString(onetimeMesg.Sum(nil)),
				Signature:   hex.EncodeToString(onetimeSig),
				Public_key:  serviceSK.Address(),
				Redirect:    "https://demo-service.testing.ssn.digital/v1/callback",
			})
		}
	}

	resp := CheckoutResponse{
		ExistingPreauth: existingpreauth,
		Preauth:         preauth,
		Onetime:         onetime,
	}

	checkoutTemplate := template.Must(template.ParseFiles("templates/checkout.html"))
	checkoutTemplate.Execute(w, resp)
}
