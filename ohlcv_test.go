package eodhd

import (
	"net/url"
	"testing"
	"time"
)

func TestOhlcvParams_BuildPath(t *testing.T) {
	p := &OhlcvParams{
		ApiToken:    "test-token",
		Format:      GetFormatCsv(),
		FromTime:    GetPtrTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		Symbol:      "AAPL",
		CountryCode: "US",
	}

	u, _ := url.Parse("https://eodhd.com/api")
	expected := "https://eodhd.com/api/eod/AAPL.US?api_token=test-token&fmt=csv&from=2024-01-01"

	result, err := p.BuildPath(u)
	if err != nil {
		t.Error(err)
	}

	if expected != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
