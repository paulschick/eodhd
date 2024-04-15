package eodhd

import (
	"net/http"
	"net/url"
	"strings"
)

type RequestFormat string

const (
	defaultBaseUrl                   = "https://eodhd.com/api/eod/"
	defaultCountryCode               = "US"
	urlDateFormat                    = "2006-01-02"
	FormatJson         RequestFormat = "json"
	FormatCSV          RequestFormat = "csv"
)

type Client struct {
	client        *http.Client
	baseUrl       *url.URL
	apiToken      string
	countryCode   string
	defaultFormat RequestFormat

	// services
	urlBuilder *UrlBuilder
}

func NewClient(token string) (*Client, error) {
	client := &Client{
		apiToken:      token,
		client:        &http.Client{},
		countryCode:   defaultCountryCode,
		defaultFormat: FormatCSV,
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
	return c.baseUrl
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

func (c *Client) BuildUrl(params *UrlParams) string {
	c.urlBuilder.ResetUrl()
	return c.urlBuilder.BuildUrl(params)
}
