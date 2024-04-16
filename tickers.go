// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

type Ticker struct {
	Code     string  `json:"Code" csv:"Code"`
	Name     string  `json:"Name" csv:"Name"`
	Country  string  `json:"Country" csv:"Country"`
	Exchange string  `json:"Exchange" csv:"Exchange"`
	Currency string  `json:"Currency" csv:"Currency"`
	Type     string  `json:"Type" csv:"Type"`
	Isin     *string `json:"Isin,omitempty" csv:"Isin,omitempty"`
}

type TickerParams struct {
	ApiToken     string         `url:"api_token"`
	Format       *RequestFormat `url:"fmt"`
	ExchangeCode *string        `url:"-"`
}

func NewTickerParamsDefault(apiToken string) *TickerParams {
	return &TickerParams{
		ApiToken:     apiToken,
		Format:       GetFormatCsv(),
		ExchangeCode: GetPtrString("US"),
	}
}

func NewTickerParams(apiToken, exchangeCode string, format RequestFormat) *TickerParams {
	return &TickerParams{
		ApiToken:     apiToken,
		Format:       &format,
		ExchangeCode: &exchangeCode,
	}
}

func (t *TickerParams) GetEncoded() (string, error) {
	q, err := query.Values(t)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (t *TickerParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("exchange-symbol-list/%s", *t.ExchangeCode)
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := t.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type TickerService struct {
	c RequestClient
}

func NewTickerService(c RequestClient) *TickerService {
	return &TickerService{
		c: c,
	}
}

func (t *TickerService) GetTickers(exchangeCode string, format *RequestFormat) ([]*Ticker, *Response, error) {
	var reqForm RequestFormat
	if format == nil {
		reqForm = formatCSV
	} else {
		reqForm = *format
	}

	params := NewTickerParams(t.c.GetApiToken(), exchangeCode, reqForm)
	u, err := params.BuildPath(t.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := t.c.NewGetRequest(u, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*Ticker
	res, err := t.c.Do(req, &data)
	if err != nil {
		return nil, res, err
	}

	return data, res, nil
}
