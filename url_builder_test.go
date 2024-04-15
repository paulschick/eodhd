package eodhd

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

type testUrlProvider struct {
	BaseUrlStr string
	ApiKey     string
}

func (t *testUrlProvider) GetApiToken() string {
	return t.ApiKey
}

func (t *testUrlProvider) GetCountryCode() string {
	return "US"
}

func (t *testUrlProvider) GetBaseUrl() *url.URL {
	u, err := url.Parse(t.BaseUrlStr)
	if err != nil {
		u, _ = url.Parse("http://localhost:8080/")
	}
	return u
}

func (t *testUrlProvider) GetDefaultFormat() RequestFormat {
	return FormatCSV
}

func TestUrlBuilder_BuildUrl(t *testing.T) {
	form := FormatJson
	paramsNoTime := &UrlParams{
		Symbol: "AAPL",
		Format: &form,
		From:   nil,
		To:     nil,
	}
	urlProvider := &testUrlProvider{
		BaseUrlStr: "https://test.com/api/",
		ApiKey:     "test-key",
	}
	builder := NewUrlBuilder(urlProvider)

	expected1 := "https://test.com/api/AAPL.US?api_key=test-key&fmt=json"
	u1 := builder.BuildUrl(paramsNoTime)
	if u1 != expected1 {
		t.Errorf("expected URL: %s, got %s", expected1, u1)
	}

	// 2024-01-01
	fromTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	paramsWithTime := &UrlParams{
		Symbol: "AAPL",
		Format: nil,
		From:   &fromTime,
		To:     nil,
	}
	expected2 := "https://test.com/api/AAPL.US?api_key=test-key&fmt=csv&from=2024-01-01"

	builder.ResetUrl()
	u2 := builder.BuildUrl(paramsWithTime)
	fmt.Println(expected2)
	fmt.Println(u2)
}
