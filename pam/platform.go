package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GetPlatformsResponse struct {
	Platforms []Platform `json:"Platforms,omitempty"`
	Total     int        `json:"Total,omitempty"`
}
type General struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	SystemType     string `json:"systemType,omitempty"`
	Active         bool   `json:"active,omitempty"`
	Description    string `json:"description,omitempty"`
	PlatformBaseID string `json:"platformBaseID,omitempty"`
	PlatformType   string `json:"platformType,omitempty"`
}
type Required struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
type Optional struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
type Properties struct {
	Required []Required `json:"required,omitempty"`
	Optional []Optional `json:"optional,omitempty"`
}
type LinkedAccounts struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
type CredentialsManagement struct {
	AllowedSafes                          string `json:"allowedSafes,omitempty"`
	AllowManualChange                     bool   `json:"allowManualChange,omitempty"`
	PerformPeriodicChange                 bool   `json:"performPeriodicChange,omitempty"`
	RequirePasswordChangeEveryXDays       int    `json:"requirePasswordChangeEveryXDays,omitempty"`
	AllowManualVerification               bool   `json:"allowManualVerification,omitempty"`
	PerformPeriodicVerification           bool   `json:"performPeriodicVerification,omitempty"`
	RequirePasswordVerificationEveryXDays int    `json:"requirePasswordVerificationEveryXDays,omitempty"`
	AllowManualReconciliation             bool   `json:"allowManualReconciliation,omitempty"`
	AutomaticReconcileWhenUnsynched       bool   `json:"automaticReconcileWhenUnsynched,omitempty"`
}
type SessionManagement struct {
	RequirePrivilegedSessionMonitoringAndIsolation bool   `json:"requirePrivilegedSessionMonitoringAndIsolation,omitempty"`
	RecordAndSaveSessionActivity                   bool   `json:"recordAndSaveSessionActivity,omitempty"`
	PSMServerID                                    string `json:"PSMServerID,omitempty"`
}
type PrivilegedAccessWorkflows struct {
	RequireDualControlPasswordAccessApproval bool `json:"requireDualControlPasswordAccessApproval,omitempty"`
	EnforceCheckinCheckoutExclusiveAccess    bool `json:"enforceCheckinCheckoutExclusiveAccess,omitempty"`
	EnforceOnetimePasswordAccess             bool `json:"enforceOnetimePasswordAccess,omitempty"`
}
type Platform struct {
	General                   General                   `json:"general,omitempty"`
	Properties                Properties                `json:"properties,omitempty"`
	LinkedAccounts            []LinkedAccounts          `json:"linkedAccounts,omitempty"`
	CredentialsManagement     CredentialsManagement     `json:"credentialsManagement,omitempty"`
	SessionManagement         SessionManagement         `json:"sessionManagement,omitempty"`
	PrivilegedAccessWorkflows PrivilegedAccessWorkflows `json:"privilegedAccessWorkflows,omitempty"`
}

func (c *Client) GetPlatforms() (GetPlatformsResponse, int, error) {
	resp := GetPlatformsResponse{}

	// GET /PasswordVault/API/Platforms/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Platforms/", c.Config.PcloudUrl)

	req, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return resp, http.StatusConflict, err
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return resp, http.StatusBadGateway, fmt.Errorf("failed to send request. %s", err)
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}
	if res.StatusCode >= 300 {
		return resp, res.StatusCode, fmt.Errorf("received non-200 status code(%d): %s", res.StatusCode, string(body))
	}

	return resp, http.StatusOK, nil
}
