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
