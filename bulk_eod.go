// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

type BulkEod struct {
	Code          string  `csv:"Code" json:"code"`
	Exchange      string  `csv:"Ex" json:"exchange_short_name"`
	Date          string  `csv:"Date" json:"date"`
	Open          float64 `csv:"Open" json:"open"`
	High          float64 `csv:"High" json:"high"`
	Low           float64 `csv:"Low" json:"low"`
	Close         float64 `csv:"Close" json:"close"`
	AdjustedClose float64 `csv:"Adjusted_close" json:"adjusted_close"`
	Volume        int     `csv:"Volume" json:"volume"`
}

type BulkEodParams struct {
	ApiToken string        `url:"api_token"`
	Format   RequestFormat `url:"fmt"`
	Exchange string        `url:"-"`
}

func NewBulkEodParams(apiToken string, exchange *string, format *RequestFormat) *BulkEodParams {
	var (
		exch       string
		formatType RequestFormat
	)

	if exchange == nil {
		exch = "US"
	} else {
		exch = *exchange
	}
	if format == nil {
		formatType = formatCSV
	} else {
		formatType = *format
	}

	return &BulkEodParams{
		ApiToken: apiToken,
		Format:   formatType,
		Exchange: exch,
	}
}

func (b *BulkEodParams) GetEncoded() (string, error) {
	q, err := query.Values(b)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (b *BulkEodParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("eod-bulk-last-day/%s", b.Exchange)
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := b.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type BulkEodService struct {
	c RequestClient
}

func NewBulkEodService(c RequestClient) *BulkEodService {
	return &BulkEodService{
		c: c,
	}
}

func (b *BulkEodService) GetBulkEod(exchange *string, format *RequestFormat) ([]*BulkEod, *Response, error) {
	params := NewBulkEodParams(b.c.GetApiToken(), exchange, format)
	u, err := params.BuildPath(b.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := b.c.NewGetRequest(u, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*BulkEod
	res, err := b.c.Do(req, &data)
	if err != nil {
		return nil, res, err
	}

	return data, res, nil
}
