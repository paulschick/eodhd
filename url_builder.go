package eodhd

import (
	"fmt"
	"net/url"
	"time"
)

type UrlParams struct {
	Symbol string
	Format *RequestFormat
	From   *time.Time
	To     *time.Time
}

type UrlBuilder struct {
	currentUrl *url.URL
	client     UrlClient
}

func NewUrlBuilder(c UrlClient) *UrlBuilder {
	return &UrlBuilder{
		currentUrl: c.GetBaseUrl(),
		client:     c,
	}
}

func (u *UrlBuilder) SetRequestSymbol(symbol string) {
	u.currentUrl = u.currentUrl.JoinPath(fmt.Sprintf("%s.%s", symbol, u.client.GetCountryCode()))
}

func (u *UrlBuilder) SetApiKeyParam() {
	values := u.currentUrl.Query()
	values.Add("api_key", u.client.GetApiToken())
	u.currentUrl.RawQuery = values.Encode()
}

func (u *UrlBuilder) SetFormat(format RequestFormat) {
	values := u.currentUrl.Query()
	values.Add("fmt", string(format))
	u.currentUrl.RawQuery = values.Encode()
}

func (u *UrlBuilder) SetFrom(from time.Time) {
	fromStr := from.Format(urlDateFormat)
	values := u.currentUrl.Query()
	values.Add("from", fromStr)
	u.currentUrl.RawQuery = values.Encode()
}

func (u *UrlBuilder) SetTo(to time.Time) {
	toStr := to.Format(urlDateFormat)
	values := u.currentUrl.Query()
	values.Add("to", toStr)
	u.currentUrl.RawQuery = values.Encode()
}

func (u *UrlBuilder) ResetUrl() {
	u.currentUrl = u.client.GetBaseUrl()
}

func (u *UrlBuilder) BuildUrl(urlParams *UrlParams) string {
	if urlParams.Symbol == "" {
		return u.currentUrl.String()
	}

	var format RequestFormat
	if urlParams.Format == nil {
		format = u.client.GetDefaultFormat()
	} else {
		format = *urlParams.Format
	}

	u.SetRequestSymbol(urlParams.Symbol)
	u.SetApiKeyParam()
	u.SetFormat(format)

	if urlParams.From != nil {
		u.SetFrom(*urlParams.From)
	}
	if urlParams.To != nil {
		u.SetTo(*urlParams.To)
	}

	return u.currentUrl.String()
}
