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
	client     *Client
}

func NewUrlBuilder(c *Client) *UrlBuilder {
	return &UrlBuilder{
		currentUrl: c.baseUrl,
		client:     c,
	}
}

func (u *UrlBuilder) SetRequestSymbol(symbol string) {
	u.currentUrl.JoinPath(fmt.Sprintf("%s.%s", symbol, u.client.countryCode))
}

func (u *UrlBuilder) SetApiKeyParam() {
	values := u.currentUrl.Query()
	values.Add("api_key", u.client.apiToken)
}

func (u *UrlBuilder) SetFormat(format RequestFormat) {
	values := u.currentUrl.Query()
	values.Add("fmt", string(format))
}

func (u *UrlBuilder) SetFrom(from time.Time) {
	fromStr := from.Format(urlDateFormat)
	values := u.currentUrl.Query()
	values.Add("from", fromStr)
}

func (u *UrlBuilder) SetTo(to time.Time) {
	toStr := to.Format(urlDateFormat)
	values := u.currentUrl.Query()
	values.Add("to", toStr)
}

func (u *UrlBuilder) ResetUrl() {
	u.currentUrl = u.client.baseUrl
}

func (u *UrlBuilder) BuildUrl(urlParams *UrlParams) string {
	if urlParams.Symbol == "" {
		return u.currentUrl.String()
	}

	var format RequestFormat
	if urlParams.Format == nil {
		format = u.client.defaultFormat
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
