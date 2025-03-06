package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type PostAddMemberRequest struct {
	MemberName               string      `json:"memberName,omitempty"`
	SearchIn                 string      `json:"searchIn,omitempty"`
	MembershipExpirationDate int         `json:"membershipExpirationDate,omitempty"`
	Permissions              Permissions `json:"permissions,omitempty"`
	MemberType               string      `json:"MemberType,omitempty"`
	IsReadOnly               bool        `json:"isReadOnly,omitempty"`
}

type Permissions struct {
	UseAccounts                            bool `json:"useAccounts,omitempty"`
	RetrieveAccounts                       bool `json:"retrieveAccounts,omitempty"`
	ListAccounts                           bool `json:"listAccounts,omitempty"`
	AddAccounts                            bool `json:"addAccounts,omitempty"`
	UpdateAccountContent                   bool `json:"updateAccountContent,omitempty"`
	UpdateAccountProperties                bool `json:"updateAccountProperties,omitempty"`
	InitiateCPMAccountManagementOperations bool `json:"initiateCPMAccountManagementOperations,omitempty"`
	SpecifyNextAccountContent              bool `json:"specifyNextAccountContent,omitempty"`
	RenameAccounts                         bool `json:"renameAccounts,omitempty"`
	DeleteAccounts                         bool `json:"deleteAccounts,omitempty"`
	UnlockAccounts                         bool `json:"unlockAccounts,omitempty"`
	ManageSafe                             bool `json:"manageSafe,omitempty"`
	ManageSafeMembers                      bool `json:"manageSafeMembers,omitempty"`
	BackupSafe                             bool `json:"backupSafe,omitempty"`
	ViewAuditLog                           bool `json:"viewAuditLog,omitempty"`
	ViewSafeMembers                        bool `json:"viewSafeMembers,omitempty"`
	AccessWithoutConfirmation              bool `json:"accessWithoutConfirmation,omitempty"`
	CreateFolders                          bool `json:"createFolders,omitempty"`
	DeleteFolders                          bool `json:"deleteFolders,omitempty"`
	MoveAccountsAndFolders                 bool `json:"moveAccountsAndFolders,omitempty"`
	RequestsAuthorizationLevel1            bool `json:"requestsAuthorizationLevel1,omitempty"`
	RequestsAuthorizationLevel2            bool `json:"requestsAuthorizationLevel2,omitempty"`
}

type PostAddMemberResponse struct {
	SafeURLID                 string      `json:"safeUrlId,omitempty"`
	SafeName                  string      `json:"safeName,omitempty"`
	SafeNumber                int         `json:"safeNumber,omitempty"`
	MemberID                  string      `json:"memberId,omitempty"`
	MemberName                string      `json:"memberName,omitempty"`
	MemberType                string      `json:"memberType,omitempty"`
	MembershipExpirationDate  int         `json:"membershipExpirationDate,omitempty"`
	IsExpiredMembershipEnable bool        `json:"isExpiredMembershipEnable,omitempty"`
	IsPredefinedUser          bool        `json:"isPredefinedUser,omitempty"`
	Permissions               Permissions `json:"permissions,omitempty"`
}

func (c *Client) AddSafeMember(member PostAddMemberRequest, safeurlid string) (PostAddMemberResponse, int, error) {
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/content/webservices/add+safe+member.htm

	addMemberResponse := PostAddMemberResponse{}

	// POST /PasswordVault/API/Safes/{safeUrlId}/Members/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Safes/%s/Members/", c.Config.PcloudUrl, safeurlid)

	jsonbody, err := json.Marshal(member)
	if err != nil {
		log.Fatalf("failed to create json body for safe member: %s\n", err.Error())
	}
	req, err := http.NewRequest(http.MethodPost, apiurl, strings.NewReader(string(jsonbody)))
	if err != nil {
		return addMemberResponse, http.StatusConflict, err
	}

	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return addMemberResponse, http.StatusBadGateway, fmt.Errorf("failed to send request. %s", err)
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	err = json.Unmarshal(body, &addMemberResponse)
	if err != nil {
		return addMemberResponse, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}
	if res.StatusCode >= 300 {
		return addMemberResponse, res.StatusCode, fmt.Errorf("received non-200 status code(%d): %s", res.StatusCode, string(body))
	}

	return addMemberResponse, http.StatusOK, nil

}
