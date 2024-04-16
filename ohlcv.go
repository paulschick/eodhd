// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import "time"

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

func (o *OhlcvService) GetCandles(params UrlParamProvider) ([]*Ohlcv, *Response, error) {
	req, err := o.c.NewEodRequest(params, nil)
	if err != nil {
		return nil, nil, err
	}

	var ohlcv []*Ohlcv
	resp, err := o.c.Do(req, &ohlcv)
	if err != nil {
		return nil, resp, err
	}

	return ohlcv, resp, nil
}
