package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Creator struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PostAddSafeRequest struct {
	SafeName                  string `json:"safeName"` // Required
	Description               string `json:"description,omitempty"`
	Location                  string `json:"location,omitempty"`              // Default "\\"
	NumberOfDaysRetention     int    `json:"numberOfDaysRetention,omitempty"` // Default 7
	NumberOfVersionsRetention int    `json:"numberOfVersionsRetention,omitempty"`
	OlacEnabled               bool   `json:"oLACEnabled,omitempty"`
	AutoPurgeEnabled          bool   `json:"autoPurgeEnabled,omitempty"`
	ManagingCPM               string `json:"managingCPM,omitempty"`
	ErrorResponse
}

type PostAddSafeResponse struct {
	SafeURLID                 string  `json:"safeUrlId"`
	SafeName                  string  `json:"safeName"`
	SafeNumber                int     `json:"safeNumber"`
	Description               string  `json:"description,omitempty"`
	Location                  string  `json:"location"`
	Creator                   Creator `json:"creator,omitempty"`
	OlacEnabled               bool    `json:"olacEnabled,omitempty"`
	ManagingCPM               string  `json:"managingCPM,omitempty"`
	NumberOfVersionsRetention any     `json:"numberOfVersionsRetention,omitempty"`
	NumberOfDaysRetention     int     `json:"numberOfDaysRetention,omitempty"`
	AutoPurgeEnabled          bool    `json:"autoPurgeEnabled,omitempty"`
	CreationTime              int64   `json:"creationTime,omitempty"`
	LastModificationTime      int64   `json:"lastModificationTime,omitempty"`
	ErrorResponse
}

type GetSafeDetails struct {
	SafeURLID                 string  `json:"safeUrlId,omitempty"`
	SafeName                  string  `json:"safeName,omitempty"`
	SafeNumber                int     `json:"safeNumber,omitempty"`
	Description               string  `json:"description,omitempty"`
	Location                  string  `json:"location,omitempty"`
	Creator                   Creator `json:"creator,omitempty"`
	OlacEnabled               bool    `json:"olacEnabled,omitempty"`
	ManagingCPM               string  `json:"managingCPM,omitempty"`
	NumberOfVersionsRetention any     `json:"numberOfVersionsRetention,omitempty"`
	NumberOfDaysRetention     int     `json:"numberOfDaysRetention,omitempty"`
	AutoPurgeEnabled          bool    `json:"autoPurgeEnabled,omitempty"`
	CreationTime              int     `json:"creationTime,omitempty"`
	LastModificationTime      int64   `json:"lastModificationTime,omitempty"`
	Accounts                  []any   `json:"accounts,omitempty"`
	IsExpiredMember           bool    `json:"isExpiredMember,omitempty"`
	ErrorResponse
}

func (c *Client) AddSafe(safereq PostAddSafeRequest) (PostAddSafeResponse, int, error) {
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/Content/WebServices/Add%20Safe.htm
	newsafe := PostAddSafeResponse{}

	// POST /PasswordVault/API/Safes/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Safes/", c.Config.PcloudUrl)

	jsonbody, err := json.Marshal(safereq)
	if err != nil {
		log.Fatalf("failed to create json body for add safe: %s\n", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, apiurl, strings.NewReader(string(jsonbody)))
	if err != nil {
		return newsafe, http.StatusConflict, err
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return newsafe, http.StatusBadGateway, fmt.Errorf("failed to send request. %s", err)
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	err = json.Unmarshal(body, &newsafe)
	if err != nil {
		return newsafe, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}

	return newsafe, res.StatusCode, nil
}

func (c *Client) GetSafeDetails(safename string) (GetSafeDetails, int, error) {
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/content/sdk/safes+web+services+-+get+safes+details.htm

	safedetails := GetSafeDetails{}

	// GET /PasswordVault/API/Safes/{SafeUrlId}/
	safeurlid := url.QueryEscape(safename)
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Safes/%s", c.Config.PcloudUrl, safeurlid)

	req, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return safedetails, http.StatusConflict, err
	}
	// attach the header
	req.Header = make(http.Header)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.SendRequest(req)
	if err != nil {
		return safedetails, http.StatusBadGateway, fmt.Errorf("failed to send request. %s", err)
	}

	// read response body
	body, error := io.ReadAll(res.Body)
	if error != nil {
		log.Println(error)
	}
	// close response body
	defer res.Body.Close()

	err = json.Unmarshal(body, &safedetails)
	if err != nil {
		return safedetails, res.StatusCode, fmt.Errorf("response format failed to parse: %s: %s", err.Error(), string(body))
	}

	return safedetails, res.StatusCode, nil
}
