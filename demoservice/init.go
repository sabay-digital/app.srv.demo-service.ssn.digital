package demoservice

import (
	"encoding/hex"
	"fmt"

	"github.com/jinzhu/gorm"
	// DB Driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sabay-digital/sdk.golang.ssn.digital/ssn"
	"github.com/stellar/go/keypair"
)

var (
	ssnAPI       string
	preauthSeed  []byte
	serviceSK    keypair.KP
	slackWebhook string
	db           *gorm.DB
)

// PreAuthorization table model
type PreAuthorization struct {
	gorm.Model
	UserID     string `gorm:"column:user_id"`
	UserPubkey string `gorm:"column:user_pk"`
	PPPubkey   string `gorm:"column:pp_pk"`
	Currencies string `gorm:"column:currencies"`
}

// Initialise ensures all package wide variables are correctly set at startup
func Initialise(apiURL, service, pas, dbHost, dbName, dbUser, dbPassword, slackWH string) {
	ssnAPI = apiURL
	serviceSK = keypair.MustParse(service)
	slackWebhook = slackWH
	var err error

	// Decode the BIP39 seed
	preauthSeed, err = hex.DecodeString(pas)
	ssn.Log(err, "Init: Decode seed string to hex")

	db, _ = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName))

	db.AutoMigrate(&PreAuthorization{})
}
