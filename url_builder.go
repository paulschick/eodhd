// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"time"
)

type UrlBuilder struct {
	currentURL *url.URL
	c          UrlClient
}

func NewUrlBuilder(c UrlClient) *UrlBuilder {
	return &UrlBuilder{
		c:          c,
		currentURL: c.GetBaseUrl(),
	}
}

type UrlParams struct {
	Symbol   string         `url:"-"`
	Format   *RequestFormat `url:"fmt"`
	FromTime *time.Time     `url:"-"`
	ToTime   *time.Time     `url:"-"`
	From     *string        `url:"from,omitempty"`
	To       *string        `url:"from,omitempty"`
	ApiToken string         `url:"api_token"`
}

type UrlParamProvider interface {
	GetSymbol() string
	Encode() (string, error)
}

func (u *UrlParams) GetSymbol() string {
	return u.Symbol
}

func (u *UrlParams) Encode() (string, error) {
	if u.FromTime != nil {
		from := u.FromTime.Format(urlDateFormat)
		u.From = &from
	}
	if u.ToTime != nil {
		to := u.ToTime.Format(urlDateFormat)
		u.To = &to
	}
	q, err := query.Values(u)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (u *UrlBuilder) SetOptions(options UrlParamProvider) error {
	if options == nil {
		return nil
	}
	q, err := options.Encode()
	if err != nil {
		return err
	}
	u.currentURL.RawQuery = q
	return nil
}

func (u *UrlBuilder) BuildUrl(params UrlParamProvider) (string, error) {
	u.currentURL = u.c.GetBaseUrl()
	u.currentURL = u.currentURL.JoinPath(fmt.Sprintf("%s.%s", params.GetSymbol(), u.c.GetCountryCode()))
	err := u.SetOptions(params)
	if err != nil {
		return "", err
	}
	return u.currentURL.String(), nil
}
