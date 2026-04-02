// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import "errors"

type ClientOption func(*Client) error

func SetPercentOfRateLimit(maxPercent float64) ClientOption {
	return func(c *Client) error {
		if maxPercent < 0 || maxPercent > 1{
			return errors.New("max percent must be between 0 (exclusive) and 1 (inclusive)")
		}
		burst := 1 - maxPercent
		c.maxPercentOfLimit = maxPercent
		c.limiterBurst = burst
		return nil
	}
}
