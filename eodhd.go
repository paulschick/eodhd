// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"context"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RequestFormat string

const (
	defaultBaseUrl     = "https://eodhd.com/api/"
	defaultCountryCode = "US"
	urlDateFormat      = "2006-01-02"
	userAgent          = "go-eodhd"
)

const (
	formatJson RequestFormat = "json"
	formatCSV  RequestFormat = "csv"
)

const (
	DefaultRateLimitPercent = 0.75
	DefaultBurstPercent     = 0.25
)

func GetFormatCsv() *RequestFormat {
	f := formatCSV
	return &f
}

func GetFormatJson() *RequestFormat {
	f := formatJson
	return &f
}

type RequestClient interface {
	NewGetRequest(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error)
	Do(req *retryablehttp.Request, data interface{}) (*Response, error)
	GetApiToken() string
	GetBaseUrl() *url.URL
}

type Client struct {
	client        *retryablehttp.Client
	baseUrl       *url.URL
	apiToken      string
	countryCode   string
	defaultFormat RequestFormat
	UserAgent     string

	limiter           *rate.Limiter
	maxPercentOfLimit float64
	limiterBurst      float64
	configureOnce     sync.Once

	// services
	OhlcvService     *OhlcvService
	ExchangesService *ExchangesService
	TickerService    *TickerService
	BulkEodService   *BulkEodService
}

func NewClient(token string, options ...ClientOption) (*Client, error) {
	client := &Client{
		apiToken:          token,
		countryCode:       defaultCountryCode,
		defaultFormat:     formatCSV,
		UserAgent:         userAgent,
		maxPercentOfLimit: DefaultRateLimitPercent,
		limiterBurst:      DefaultBurstPercent,
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

	client.OhlcvService = NewOhlcvService(client)
	client.ExchangesService = NewExchangesService(client)
	client.TickerService = NewTickerService(client)
	client.BulkEodService = NewBulkEodService(client)

	err = client.applyOptions(options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) applyOptions(options ...ClientOption) error {
	for _, fn := range options {
		if fn != nil {
			err := fn(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func (c *Client) NewGetRequest(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error) {
	reqHeaders := make(http.Header)
	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	if headers != nil {
		for k, v := range *headers {
			reqHeaders.Set(k, v)
		}
	}

	req, err := retryablehttp.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

func (c *Client) configureRateLimiter(ctx context.Context, limitAmt int) {
	// Rate Limit is provided as per minute by EODHD
	rl := float64(limitAmt) / 60.0

	limit := rate.Limit(rl * c.maxPercentOfLimit)
	burst := 1

	if int(rl*c.limiterBurst) > 1 {
		burst = int(rl * c.limiterBurst)
	}

	c.limiter = rate.NewLimiter(limit, burst)

	// wait since we get the limit from the http headers of a response
	_ = c.limiter.Wait(ctx)
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
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	response := newResponse(resp)

	// Configure the limiter on the first request made
	// Only do this once.
	c.configureOnce.Do(func() { c.configureRateLimiter(req.Context(), response.RateLimit) })

	reqFormat := req.URL.Query().Get("fmt")
	format := formatCSV
	if reqFormat != "" {
		format = RequestFormat(reqFormat)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if format == formatCSV {
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
