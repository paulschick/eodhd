// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"net/url"
	"testing"
)

func TestTickerParams_BuildPath(t *testing.T) {
	p := NewTickerParamsDefault("test-token")
	baseURL, _ := url.Parse("https://eodhd.com/api")
	expected := "https://eodhd.com/api/exchange-symbol-list/US?api_token=test-token&fmt=csv"
	u, _ := p.BuildPath(baseURL)
	if u != expected {
		t.Errorf("expected %s, got %s", expected, u)
	}
}
