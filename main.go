package main

import (
	"net/http"
	"os"

	"git.sabay.com/payment-network/services/app.srv.demo-service.ssn.digital/demoservice"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
)

func init() {
	demoservice.Initialise(os.Getenv("SSN_API"),
		os.Getenv("SERVICE_SK"),
		os.Getenv("PREAUTH_SEED"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SLACK_WEBHOOK_URL"),
	)
}

func main() {
	err := http.ListenAndServe(":3000", demoservice.Router())
	ssn.Log(err, "Main: Server error")
}
