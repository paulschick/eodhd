package eodhd

import (
	"net/url"
	"testing"
)

func TestExchangeParams_BuildPath(t *testing.T) {
	e := &ExchangeParams{
		ApiToken: "test-token",
		Format:   formatJson,
	}
	u, _ := url.Parse("https://eodhd.com/api")
	expected := "https://eodhd.com/api/exchanges-list/?api_token=test-token&fmt=json"
	result, _ := e.BuildPath(u)
	if expected != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
