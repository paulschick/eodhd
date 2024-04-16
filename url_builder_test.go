// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
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

func TestUrlBuilder_BuildUrl_NoTime(t *testing.T) {
	urlProvider := &testUrlProvider{
		BaseUrlStr: "https://test.com/api/",
		ApiKey:     "test-key",
	}
	expected1 := "https://test.com/api/AAPL.US?api_key=test-key&fmt=json"
	builder := NewUrlBuilder(urlProvider)
	format1 := FormatJson
	params := &UrlParams{
		Symbol: "AAPL",
		ApiKey: "test-key",
		Format: &format1,
	}

	urlStr1, err := builder.BuildUrl(params)
	if err != nil {
		t.Error("error getting url: ", err)
	}
	if urlStr1 != expected1 {
		t.Errorf("expected %s, got %s", expected1, urlStr1)
	}
}

func TestUrlBuilder_BuildUrl_FromTime(t *testing.T) {
	urlProvider := &testUrlProvider{
		BaseUrlStr: "https://test.com/api/",
		ApiKey:     "test-key",
	}
	fromTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := "https://test.com/api/AAPL.US?api_key=test-key&fmt=csv&from=2024-01-01"
	builder := NewUrlBuilder(urlProvider)

	format := FormatCSV
	params := &UrlParams{
		Symbol:   "AAPL",
		ApiKey:   "test-key",
		Format:   &format,
		FromTime: &fromTime,
	}
	result, err := builder.BuildUrl(params)
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
