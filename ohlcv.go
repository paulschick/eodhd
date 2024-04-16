// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"time"
)

type OhlcvParams struct {
	Symbol      string         `url:"-"`
	CountryCode string         `url:"-"`
	Format      *RequestFormat `url:"fmt"`
	FromTime    *time.Time     `url:"-"`
	ToTime      *time.Time     `url:"-"`
	From        *string        `url:"from,omitempty"`
	To          *string        `url:"from,omitempty"`
	ApiToken    string         `url:"api_token"`
}

func (o *OhlcvParams) GetEncoded() (string, error) {
	if o.FromTime != nil {
		from := o.FromTime.Format(urlDateFormat)
		o.From = &from
	}
	if o.ToTime != nil {
		to := o.ToTime.Format(urlDateFormat)
		o.To = &to
	}
	q, err := query.Values(o)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (o *OhlcvParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("eod/%s.%s", o.Symbol, o.CountryCode)
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := o.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

// Ohlcv represents the OHLCV response for both JSON and CSV formats for historical EOD data.
type Ohlcv struct {
	Date     string  `csv:"Date" json:"date"`
	Open     float64 `csv:"Open" json:"open"`
	High     float64 `csv:"High" json:"high"`
	Low      float64 `csv:"Low" json:"low"`
	Close    float64 `csv:"Close" json:"close"`
	AdjClose float64 `csv:"Adjusted_close" json:"adjusted_close"`
	Volume   float64 `csv:"Volume" json:"volume"`

	// Set via ParseDate from standard urlDateFormat
	DateParsed *time.Time `csv:"-" json:"-"`
}

func (o *Ohlcv) ParseDate() error {
	// responses for EOD data use the same format as the url date, which is YYYY-MM-DD
	t, err := time.Parse(urlDateFormat, o.Date)
	if err != nil {
		return err
	}
	o.DateParsed = &t
	return nil
}

type OhlcvService struct {
	c RequestClient
}

func NewOhlcvService(c RequestClient) *OhlcvService {
	return &OhlcvService{
		c: c,
	}
}

func (o *OhlcvService) GetOhlcv(symbol string, countryCode *string, format *RequestFormat, from, to *time.Time) ([]*Ohlcv, *Response, error) {
	var f *RequestFormat
	if format == nil {
		f = GetFormatCsv()
	} else {
		f = format
	}
	var country *string
	if countryCode == nil {
		country = GetPtrString("US")
	} else {
		country = countryCode
	}

	params := &OhlcvParams{
		ApiToken:    o.c.GetApiToken(),
		FromTime:    from,
		ToTime:      to,
		Format:      f,
		Symbol:      symbol,
		CountryCode: *country,
	}

	reqUrl, err := params.BuildPath(o.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := o.c.NewGetRequest(reqUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*Ohlcv
	resp, err := o.c.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, nil
}
