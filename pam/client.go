package pam

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

const (
	Version          = "v1.0.0"
	defaultUserAgent = "privilegeaccessmanager-sdk-go" + "/" + Version
)

type Client struct {
	BaseURL  string
	AuthType string
	Session  *Session
	Config   *Config
}

type Config struct {
	IdTenantUrl   string
	PcloudUrl     string
	User          string
	Pass          string
	TlsSkipVerify bool
}

func NewConfig(idtenanturl string, pcloudurl string, u string, p string) *Config {
	config := Config{
		IdTenantUrl:   idtenanturl, // Example: "https://EXAMPLE123.id.cyberark.cloud"
		PcloudUrl:     pcloudurl,   // Example: "https://EXAMPLE123.privilegecloud.cyberark.cloud"
		User:          u,           // Note: this must be a service account user
		Pass:          p,
		TlsSkipVerify: false,
	}
	return &config
}

// NewClient - create a client with reasonable defaults
func NewClient(baseurl string, config *Config, options ...func(*Client) error) *Client {
	client := Client{
		BaseURL:  baseurl,
		AuthType: "",
		Session:  nil,
		Config:   config,
	}
	for _, option := range options {
		option(&client)
	}
	return &client
}
func DisableTlsVerify() func(*Client) error {
	return func(c *Client) error {
		c.Config.TlsSkipVerify = true
		return nil
	}
}

func (c *Client) SendRequest(req *http.Request) (*http.Response, error) {
	// if token is provided, add header Authorization
	if c.Session != nil && c.Session.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.Session.TokenType, c.Session.Token))
	}

	client := GetHTTPClient(time.Second*30, c.Config.TlsSkipVerify)
	return client.Do(req)
}

// GetDefaultHTTPClient create http client with 30s timeout and no skip verify
func GetDefaultHTTPClient() *http.Client {
	return GetHTTPClient(time.Second*30, false)
}

// GetHTTPClient create http client for HTTPS
func GetHTTPClient(timeout time.Duration, skipverify bool) *http.Client {
	client := &http.Client{
		Timeout: timeout, /*time.Second * 30 */
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipverify, /* TLS_SKIP_VERIFY */
			},
		},
	}
	return client
}
