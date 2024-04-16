// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"net/url"
	"testing"
)

func TestBulkEodParams_BuildPath(t *testing.T) {
	uStr := "https://eodhd.com/api"
	baseURL, _ := url.Parse(uStr)
	p := NewBulkEodParams("test-token", GetPtrString("US"), GetFormatJson())
	u, err := p.BuildPath(baseURL)
	if err != nil {
		t.Error(err)
	}
	expected := "https://eodhd.com/api/eod-bulk-last-day/US?api_token=test-token&fmt=json"
	if expected != u {
		t.Errorf("expected %s, got %s", expected, u)
	}
}
