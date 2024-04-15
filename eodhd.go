package eodhd

import (
	"net/http"
	"net/url"
	"strings"
)

type RequestFormat string

const (
	defaultBaseUrl     = "https://eodhd.com/api/eod/"
	defaultCountryCode = "US"
	urlDateFormat      = "2006-01-02"
	userAgent          = "go-eodhd"
)

const (
	FormatJson RequestFormat = "json"
	FormatCSV  RequestFormat = "csv"
)

type Client struct {
	client        *http.Client
	baseUrl       *url.URL
	apiToken      string
	countryCode   string
	defaultFormat RequestFormat
	UserAgent     string

	// services
	urlBuilder *UrlBuilder
}

func NewClient(token string) (*Client, error) {
	client := &Client{
		apiToken:      token,
		client:        &http.Client{},
		countryCode:   defaultCountryCode,
		defaultFormat: FormatCSV,
		UserAgent:     userAgent,
	}
	err := client.setBaseUrl(defaultBaseUrl)
	if err != nil {
		return nil, err
	}

	client.urlBuilder = NewUrlBuilder(client)

	return client, nil
}

type UrlClient interface {
	GetApiToken() string
	GetCountryCode() string
	GetBaseUrl() *url.URL
	GetDefaultFormat() RequestFormat
}

func (c *Client) GetApiToken() string {
	return c.apiToken
}

func (c *Client) GetCountryCode() string {
	return c.countryCode
}

func (c *Client) GetBaseUrl() *url.URL {
	u := *c.baseUrl
	return &u
}

func (c *Client) GetDefaultFormat() RequestFormat {
	return c.defaultFormat
}

func (c *Client) setBaseUrl(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseUrl = baseURL
	return nil
}

// BuildUrl TODO building the URL params should potentially be exposed through the client.
// Otherwise, this could be good enough for now
func (c *Client) BuildUrl(params UrlParamProvider) (string, error) {
	return c.urlBuilder.BuildUrl(params)
}
