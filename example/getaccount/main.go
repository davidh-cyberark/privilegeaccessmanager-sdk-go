package main

// https://docs.cyberark.com/pam-self-hosted/latest/en/content/webservices/get+account+details.htm

// The user who runs this web service requires List Accounts permissions in
// the Safe where the account is located inside the Vault.

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <account_id>", os.Args[0])
	}
	acctid := os.Args[1]

	k := koanf.New(".")
	err := k.Load(file.Provider("creds.toml"), toml.Parser())
	if err != nil {
		log.Fatalf("failed to load creds.toml: %s", err.Error())
	}

	config := pam.NewConfig(k.String("idtenanturl"), k.String("pcloudurl"), k.String("user"), k.String("pass"))
	client := pam.NewClient(k.String("pcloudurl"), config)
	client.RefreshSession()

	resp, respcode, err := client.GetAccount(acctid)
	if err != nil {
		log.Fatalf("Error: could not add account: (%d) %s", respcode, err.Error())
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %s", err.Error())
	}
	// fmt.Println("Account data:")
	fmt.Println(string(jsonData))
}
