package main

import (
	"log"
	"net/http"

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
	if respcode != http.StatusOK {
		if addsaferesp.ErrorResponse.ErrorCode == "SFWS0002" {
			safedetails, respcode, err := client.GetSafeDetails(newsafe.SafeName)
			if err != nil || respcode >= 300 {
				log.Fatalf("Error: not able to fetch details about existing safe: (%d) %s", respcode, err.Error())
			}
			addsaferesp.SafeURLID = safedetails.SafeURLID
			addsaferesp.SafeName = safedetails.SafeName
			addsaferesp.SafeNumber = safedetails.SafeNumber
			addsaferesp.Description = safedetails.Description
			addsaferesp.Location = safedetails.Location

		} else {
			log.Fatalf("Error: (%d) %s: %s", respcode, addsaferesp.ErrorCode, addsaferesp.ErrorMessage)
		}
	}
	log.Printf("Response Code: %d\nSafeURLID: %s\nSafeName: %s\nSafeNumber: %d\nDescription: %s\nLocation: %s\n",
		respcode,
		addsaferesp.SafeURLID, addsaferesp.SafeName, addsaferesp.SafeNumber, addsaferesp.Description, addsaferesp.Location)

	membername := "my-new-safe-1-safe-role"
	memberReq := pam.PostAddMemberRequest{
		MemberName: membername,
		MemberType: "Role",
		IsReadOnly: true,
		Permissions: pam.Permissions{
			ListAccounts:                           true,
			AddAccounts:                            true,
			UpdateAccountContent:                   true,
			UpdateAccountProperties:                true,
			InitiateCPMAccountManagementOperations: true,
			AccessWithoutConfirmation:              true,
			ManageSafeMembers:                      true,
		},
	}

	// Go to identityadmin-sdk-go, and use the cmd/identity-client to create a role
	addsaferoleresp, respcode, err := client.AddSafeMember(memberReq, addsaferesp.SafeURLID)
	if err != nil {
		log.Fatalf("Error: could not add safe member: %s", err.Error())
	}
	log.Printf("Response Code: %d\nSafeURLID: %s\nSafeName: %s\nSafeNumber: %d\nMemberID: %s\nMemberName: %s\nMemberType: %s\n",
		respcode,
		addsaferoleresp.SafeURLID,
		addsaferoleresp.SafeName,
		addsaferoleresp.SafeNumber,
		addsaferoleresp.MemberID,
		addsaferoleresp.MemberName,
		addsaferoleresp.MemberType)

}
