package eodhd

import (
	"context"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	client        *retryablehttp.Client
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
		countryCode:   defaultCountryCode,
		defaultFormat: FormatCSV,
		UserAgent:     userAgent,
	}
	err := client.setBaseUrl(defaultBaseUrl)
	if err != nil {
		return nil, err
	}

	client.client = &retryablehttp.Client{
		CheckRetry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			if err != nil {
				return false, err
			}
			if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				return true, nil
			}
			return false, nil
		},
		Backoff: func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			min = 1 * time.Second
			max = 2 * time.Second
			return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
		},
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   cleanhttp.DefaultClient(),
		RetryWaitMin: 250 * time.Millisecond,
		RetryWaitMax: 1 * time.Second,
		RetryMax:     5,
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

func (c *Client) BuildUrl(params UrlParamProvider) (string, error) {
	return c.urlBuilder.BuildUrl(params)
}

func (c *Client) NewEodRequest(params UrlParamProvider, headers *map[string]string) (*retryablehttp.Request, error) {
	urlStr, err := c.BuildUrl(params)
	if err != nil {
		return nil, err
	}

	reqHeaders := make(http.Header)
	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	if headers != nil {
		for k, v := range *headers {
			reqHeaders.Set(k, v)
		}
	}

	req, err := retryablehttp.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

type Response struct {
	*http.Response

	RateLimit          int
	RateLimitRemaining int
}

const (
	RateLimitHeader          = "X-RateLimit-Limit"
	RateLimitRemainingHeader = "X-RateLimit-Remaining"
)

func newResponse(response *http.Response) *Response {
	r := &Response{Response: response}
	r.SetHeaderValues()
	return r
}

func (r *Response) SetHeaderValues() {
	if limit := r.Header.Get(RateLimitHeader); limit != "" {
		r.RateLimit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(RateLimitRemainingHeader); remaining != "" {
		r.RateLimitRemaining, _ = strconv.Atoi(remaining)
	}
}

func (c *Client) Do(req *retryablehttp.Request, data interface{}) (*Response, error) {
	// TODO add limiter

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	response := newResponse(resp)

	reqFormat := req.URL.Query().Get("fmt")
	format := FormatCSV
	if reqFormat != "" {
		format = RequestFormat(reqFormat)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if format == FormatCSV {
		err = gocsv.UnmarshalBytes(bodyBytes, data)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal(bodyBytes, data)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}
