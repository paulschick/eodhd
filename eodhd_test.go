// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package eodhd

import (
	"math"
	"testing"
)

func TestClient_ApplyOptions(t *testing.T) {
	c, err := NewClient("test-token", SetPercentOfRateLimit(0.90))
	if err != nil {
		t.Error(err)
	}

	limitRounded := math.Round(c.maxPercentOfLimit*100) / 100
	burstRounded := math.Round(c.limiterBurst*100) / 100

	if limitRounded != 0.90 {
		t.Errorf("expected max rate limit percent to be 0.90, got %f", limitRounded)
	}
	if burstRounded != 0.10 {
		t.Errorf("expected limiter burst to be 0.10, got %f", burstRounded)
	}
}
