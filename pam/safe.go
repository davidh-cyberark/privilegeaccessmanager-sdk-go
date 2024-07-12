package pam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

func (c *Client) AddSafe(safereq PostAddSafeRequest) (PostAddSafeResponse, int, error) {
	// https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/Content/WebServices/Add%20Safe.htm
	newsafe := PostAddSafeResponse{}

	// POST /PasswordVault/API/Safes/
	apiurl := fmt.Sprintf("%s/PasswordVault/API/Safes/", c.Config.PcloudUrl)

	jsonbody, err := json.Marshal(safereq)
	if err != nil {
		log.Fatalf("failed to parse json body for platform token: %s\n", err.Error())
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
	if res.StatusCode >= 300 {
		return newsafe, res.StatusCode, fmt.Errorf("received non-200 status code(%d): %s", res.StatusCode, string(body))
	}

	return newsafe, http.StatusOK, nil
}
