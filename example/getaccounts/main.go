package main

// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/content/sdk/getaccounts.htm

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
		log.Fatalf("Usage: %s <safe_name> [<account_name>]", os.Args[0])
	}
	safename := os.Args[1]
	acctname := ""
	if len(os.Args) == 3 {
		acctname = os.Args[2]
	}

	k := koanf.New(".")
	err := k.Load(file.Provider("creds.toml"), toml.Parser())
	if err != nil {
		log.Fatalf("failed to load creds.toml: %s", err.Error())
	}

	config := pam.NewConfig(k.String("idtenanturl"), k.String("pcloudurl"), k.String("user"), k.String("pass"))
	client := pam.NewClient(k.String("pcloudurl"), config)
	client.RefreshSession()

	filter := fmt.Sprintf("safeName eq %s", safename)
	var searchp *string = nil
	var searchtypep *string = nil
	if len(acctname) > 0 {
		search := fmt.Sprintf("search=%s", acctname)
		searchp = &search
		searchtype := "startswith"
		searchtypep = &searchtype
	}

	resp, respcode, err := client.GetAccounts(searchp, searchtypep, nil, &filter, nil, nil, nil)
	if err != nil {
		log.Fatalf("Error: could not get accounts: (%d) %s", respcode, err.Error())
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %s", err.Error())
	}
	fmt.Println(string(jsonData))
}
