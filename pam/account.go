package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	SafeName                  string                    `json:"safeName"`       // Required
	PlatformID                string                    `json:"platformId"`     // Required
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
	CategoryModificationTime  int                       `json:"categoryModificationTime,omitempty"`
}

type GetAccountResponse struct {
	ID                        string                    `json:"id,omitempty"`
	Name                      string                    `json:"name,omitempty"`
	Address                   string                    `json:"address,omitempty"`
	UserName                  string                    `json:"userName,omitempty"`
	PlatformID                string                    `json:"platformId,omitempty"`
	SafeName                  string                    `json:"safeName,omitempty"`
	SecretType                string                    `json:"secretType,omitempty"`
	PlatformAccountProperties PlatformAccountProperties `json:"platformAccountProperties,omitempty"`
	SecretManagement          SecretManagement          `json:"secretManagement,omitempty"`
	RemoteMachinesAccess      RemoteMachinesAccess      `json:"remoteMachinesAccess,omitempty"`
	CreatedTime               int                       `json:"createdTime,omitempty"`
	CategoryModificationTime  int                       `json:"CategoryModificationTime,omitempty"`
}

type GetAccountsResponse struct {
	Value []GetAccountResponse `json:"value,omitempty"`
	Count int                  `json:"count,omitempty"`
}

func (c *Client) GetAccounts(search, searchtype, sort, filter, savedfilter, offset, limit *string) (*GetAccountsResponse, int, error) {
	// https://<subdomain>.privilegecloud.cyberark.cloud/PasswordVault/API/Accounts?search={search}&searchType={searchType}&sort={sort}&offset={offset}&limit={limit}&filter={filter}/

	accountresp := GetAccountsResponse{}

	qpathparts := map[string]string{}
	if search != nil {
		qpathparts["search"] = *search
	}
	if searchtype != nil {
		// Valid values:  contains (default) or startswith
		if *searchtype != "contains" && *searchtype != "startswith" {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid searchType: %s, must be 'contains' or 'startswith'", *searchtype)
		}
		qpathparts["searchType"] = *searchtype
	}
	if sort != nil {
		qpathparts["sort"] = *sort
	}
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/content/sdk/getaccounts.htm#Filterparameters
	if filter != nil {
		// Ex: safeName eq mysafe1
		validFilters := []string{"safeName", "modificationTime", "secretModificationTime"}
		isValid := false
		for _, validFilter := range validFilters {
			if strings.Contains(*filter, validFilter) {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid filter: %s, must contain one of: %s", *filter, strings.Join(validFilters, ", "))
		}
		qpathparts["filter"] = *filter
	}
	if savedfilter != nil {
		validSavedFilters := []string{
			"Regular", "Recently", "New", "Link", "Deleted", "PolicyFailures",
			"AccessedByUsers", "ModifiedByUsers", "ModifiedByCPM", "DisabledPasswordByUser",
			"DisabledPasswordByCPM", "ScheduledForChange", "ScheduledForVerify",
			"ScheduledForReconcile", "SuccessfullyReconciled", "FailedChange",
			"FailedVerify", "FailedReconcile", "LockedOrNew", "Locked",
			"Favorites", "DeleteInsightStatus",
		}

		isValidSavedFilter := false
		for _, validSavedFilter := range validSavedFilters {
			if *savedfilter == validSavedFilter {
				isValidSavedFilter = true
				break
			}
		}
		if !isValidSavedFilter {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid savedfilter: %s, must be one of: %s", *savedfilter, strings.Join(validSavedFilters, ", "))
		}
		qpathparts["savedfilter"] = *savedfilter
	}

	if offset != nil {
		o, e := strconv.Atoi(*limit)
		if e != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("offset is not a number, got %s", *offset)
		}
		if o < 0 {
			return nil, http.StatusBadRequest, fmt.Errorf("offset must be greater than 0, got %s", *offset)
		}

		qpathparts["offset"] = *offset
	}
	if limit != nil {
		l, e := strconv.Atoi(*limit)
		if e != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("limit is not a number, got %s", *limit)
		}
		if l < 0 || l > 1000 {
			return nil, http.StatusBadRequest, fmt.Errorf("limit valid range is 0 - 1000, got %s", *limit)
		}
		qpathparts["limit"] = *limit
	}

	qpath := ""
	var qparts []string
	if len(qpathparts) > 0 {
		for key, value := range qpathparts {
			valueenc := url.QueryEscape(value) // Escape the value parts
			qparts = append(qparts, fmt.Sprintf("%s=%s", key, valueenc))
		}
		qpath = fmt.Sprintf("?%s", strings.Join(qparts, "&"))
	}

	// https://<PVWA_Server_address>/PasswordVault/API/Accounts/{id}/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Accounts%s", c.Config.PcloudUrl, qpath)

	req, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return nil, http.StatusConflict, fmt.Errorf("failed to create new request for get account: %s", err.Error())
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to send get acount request. %s", err.Error())
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return &accountresp, res.StatusCode, fmt.Errorf("received non-200 status code(%d): %s", res.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &accountresp)
	if err != nil {
		return &accountresp, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}

	return &accountresp, http.StatusOK, nil
}

func (c *Client) GetAccount(acctid string) (GetAccountResponse, int, error) {
	accountresp := GetAccountResponse{}

	// https://<PVWA_Server_address>/PasswordVault/API/Accounts/{id}/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Accounts/%s", c.Config.PcloudUrl, acctid)

	req, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return accountresp, http.StatusConflict, fmt.Errorf("failed to create new request for get account: %s", err.Error())
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return accountresp, http.StatusBadGateway, fmt.Errorf("failed to send get acount request. %s", err.Error())
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
