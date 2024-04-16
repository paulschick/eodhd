// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

type ExchangesService struct {
	c RequestClient
}

type ExchangeParams struct {
	ApiToken string        `url:"api_token"`
	Format   RequestFormat `url:"fmt"`
}

func NewExchangeParams(apiToken string) *ExchangeParams {
	return &ExchangeParams{
		ApiToken: apiToken,
		Format:   formatJson,
	}
}

func (e *ExchangeParams) GetEncoded() (string, error) {
	q, err := query.Values(e)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (e *ExchangeParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("exchanges-list/")
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := e.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type Exchange struct {
	Name         string `json:"Name"`
	Code         string `json:"Code"`
	OperatingMIC string `json:"OperatingMIC"`
	Country      string `json:"Country"`
	Currency     string `json:"Currency"`
	CountryISO2  string `json:"CountryISO2"`
	CountryISO3  string `json:"CountryISO3"`
}

func NewExchangesService(c RequestClient) *ExchangesService {
	return &ExchangesService{
		c: c,
	}
}

func (e *ExchangesService) GetExchanges() ([]*Exchange, *Response, error) {
	exchangeParams := NewExchangeParams(e.c.GetApiToken())
	u, err := exchangeParams.BuildPath(e.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}
	req, err := e.c.NewGetRequest(u, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*Exchange
	res, err := e.c.Do(req, &data)
	if err != nil {
		return nil, res, err
	}

	return data, res, nil
}
