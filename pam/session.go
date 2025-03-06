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

type Session struct {
	Token      string
	TokenType  string
	Expiration time.Time
}

type IDTenantResponse struct {
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type ErrorResponse struct {
	ErrorCode    string `json:"ErrorCode,omitempty"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
}

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("Error: %s: %s", er.ErrorCode, er.ErrorMessage)
}

func NewSession(options ...func(*Session) error) *Session {
	session := Session{
		Token:      "",
		TokenType:  "",
		Expiration: time.Now(),
	}
	for _, option := range options {
		option(&session)
	}
	return &session
}
func WithTokenInfo(tok string, toktype string, exp time.Time) func(*Session) error {
	return func(s *Session) error {
		s.Token = tok
		s.TokenType = toktype
		s.Expiration = exp
		return nil
	}
}

func (c *Client) GetSession() (*Session, int, error) {
	identurl := fmt.Sprintf("%s/oauth2/platformtoken", c.Config.IdTenantUrl) // Use PCloud OAuth

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.Config.User)
	data.Set("client_secret", c.Config.Pass)
	encodedData := data.Encode()

	req, err := http.NewRequest(http.MethodPost, identurl, strings.NewReader(encodedData))
	if err != nil {
		log.Fatalf("error in request to get session token: %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	client := GetHTTPClient(time.Second*30, c.Config.TlsSkipVerify)
	response, err := client.Do(req)

	body, e := io.ReadAll(response.Body)
	if e != nil {
		log.Fatalf("error reading platform token response: %s", err.Error())
	}
	defer response.Body.Close()

	var idresp IDTenantResponse
	err = json.Unmarshal(body, &idresp)
	if err != nil {
		log.Fatalf("failed to parse json body for platform token: %s\n", err.Error())
	}

	if idresp.Error != "" {
		return nil, response.StatusCode, fmt.Errorf("error getting token: (%s) %s", idresp.Error, idresp.ErrorDescription)
	}

	sess := Session{
		Token:      idresp.AccessToken,
		TokenType:  idresp.TokenType,
		Expiration: time.Now().Add(time.Second * time.Duration(idresp.ExpiresIn)),
	}
	return &sess, response.StatusCode, nil
}

func (c *Client) RefreshSession() error {
	session, status, err := c.GetSession()
	if err == nil && status >= 300 {
		err = fmt.Errorf("failed to get session token: %d", status)
	}
	c.Session = session
	return err
}
