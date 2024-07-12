package main

import (
	"log"

	"github.com/davidh-cyberark/privilegeaccessmanager-sdk-go/v1/pam"
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
	k.Load(file.Provider("../../creds.toml"), toml.Parser())

	config := pam.NewConfig(k.String("idtenanturl"), k.String("pcloudurl"), k.String("user"), k.String("pass"))
	client := pam.NewClient(k.String("pcloudurl"), config)

	if client.Session != nil {
		log.Printf("Session-Token: %s\nSession-TokenType: %s\nSession-Expiration:%s\n", client.Session.Token, client.Session.TokenType, client.Session.Expiration.String())
	}
	log.Println("Starting PAM Refresh Session")
	client.RefreshSession()
	log.Println("Done PAM Refresh Session")
	if client.Session != nil {
		log.Printf("Session-Token: %s\nSession-TokenType: %s\nSession-Expiration:%s\n", client.Session.Token, client.Session.TokenType, client.Session.Expiration.String())
	}
}
