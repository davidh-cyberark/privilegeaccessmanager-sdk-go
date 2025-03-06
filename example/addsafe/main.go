package main

import (
	"log"

	"github.com/davidh-cyberark/privilegeaccessmanager-sdk-go/pam"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

/*
Create a file, creds.toml with these parameters and fill in your values
idtenanturl = "https://YOUR-TENANT.id.cyberark.cloud"
pcloudurl = "https://YOUR-SUBDOMAIN.privilegecloud.cyberark.cloud"
user = "PAM_SERVICE_ACCOUNT_USER"
pass = "PAM_SERVICE_ACCOUNT_USER password"
*/
func main() {
	k := koanf.New(".")
	err := k.Load(file.Provider("creds.toml"), toml.Parser())
	if err != nil {
		log.Fatalf("failed to load creds.toml: %s", err.Error())
	}

	config := pam.NewConfig(k.String("idtenanturl"), k.String("pcloudurl"), k.String("user"), k.String("pass"))
	client := pam.NewClient(k.String("pcloudurl"), config)
	err = client.RefreshSession()
	if err != nil {
		log.Fatalf("Error: could not refresh session: %s", err.Error())
	}

	newsafe := pam.PostAddSafeRequest{
		SafeName:    "my-new-safe-1",    // required
		Description: "Example add safe", // not required
	}

	addsaferesp, respcode, err := client.AddSafe(newsafe)
	if err != nil {
		log.Fatalf("Error: could not add safe: %s", err.Error())
	}
	log.Printf("Response Code: %d\nSafeURLID: %s\nSafeName: %s\nSafeNumber: %d\nDescription: %s\nLocation: %s\n",
		respcode,
		addsaferesp.SafeURLID, addsaferesp.SafeName, addsaferesp.SafeNumber, addsaferesp.Description, addsaferesp.Location)
}
