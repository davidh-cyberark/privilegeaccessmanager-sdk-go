package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type PlatformAccountProperties map[string]string

type RemoteMachinesAccess struct {
	RemoteMachines                   string `json:"remoteMachines,omitempty"`
	AccessRestrictedToRemoteMachines bool   `json:"accessRestrictedToRemoteMachines,omitempty"`
}

type SecretManagement struct {
	AutomaticManagementEnabled bool      `json:"automaticManagementEnabled"`
	ManualManagementReason     string    `json:"manualManagementReason,omitempty"`
	Status                     string    `json:"status"`
	LastModifiedDateTime       time.Time `json:"lastModifiedDateTime,omitempty"`
	LastReconciledDateTime     time.Time `json:"lastReconciledDateTime,omitempty"`
	LastVerifiedDateTime       time.Time `json:"lastVerifiedDateTime,omitempty"`
}

// PostAddAccountRequest is used to create an account
type PostAddAccountRequest struct {
	SafeName   string `json:"safeName"`   // Required
	PlatformID string `json:"platformId"` // Required

	Name                      string                    `json:"name,omitempty"` // Account Name
	Address                   string                    `json:"address,omitempty"`
	UserName                  string                    `json:"userName,omitempty"`
	SecretType                string                    `json:"secretType,omitempty"`
	Secret                    string                    `json:"secret,omitempty"`
	SecretManagement          SecretManagement          `json:"secretManagement,omitempty"`
	PlatformAccountProperties PlatformAccountProperties `json:"platformAccountProperties,omitempty"`
	RemoteMachinesAccess      RemoteMachinesAccess      `json:"remoteMachinesAccess,omitempty"`
}

// REF: <https://docs.cyberark.com/pam-self-hosted/latest/en/Content/WebServices/Add%20Account%20v10.htm>
// PostAddAccountResponse response from getting specific account details
type PostAddAccountResponse struct {
	CategoryModificationTime  int                       `json:"categoryModificationTime,omitempty"`
	ID                        string                    `json:"id,omitempty"`
	Name                      string                    `json:"name,omitempty"`
	Address                   string                    `json:"address,omitempty"`
	UserName                  string                    `json:"userName,omitempty"`
	PlatformID                string                    `json:"platformId,omitempty"`
	SafeName                  string                    `json:"safeName,omitempty"`
	SecretType                string                    `json:"secretType,omitempty"`
	PlatformAccountProperties PlatformAccountProperties `json:"platformAccountProperties,omitempty"`
	SecretManagement          SecretManagement          `json:"secretManagement,omitempty"`
	CreatedTime               int                       `json:"createdTime,omitempty"`
}

func (c *Client) AddAccount(accountreq PostAddAccountRequest) (PostAddAccountResponse, int, error) {
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/Content/WebServices/Add%20Safe.htm
	accountresp := PostAddAccountResponse{}

	// POST /PasswordVault/API/Accounts/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Accounts/", c.Config.PcloudUrl)

	jsonbody, err := json.Marshal(accountreq)
	if err != nil {
		return accountresp,
			http.StatusConflict,
			fmt.Errorf("failed to parse json body for add account request: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, apiurl, strings.NewReader(string(jsonbody)))
	if err != nil {
		return accountresp,
			http.StatusConflict,
			fmt.Errorf("failed to create new request for add account: %s", err.Error())
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return accountresp,
			http.StatusBadGateway,
			fmt.Errorf("failed to send add acount request. %s", err.Error())
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	err = json.Unmarshal(body, &accountresp)
	if err != nil {
		return accountresp, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}
	if res.StatusCode >= 300 {
		return accountresp, res.StatusCode, fmt.Errorf("received non-200 status code(%d): %s", res.StatusCode, string(body))
	}

	return accountresp, http.StatusOK, nil
}
