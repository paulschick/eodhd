package eodhd

import "time"

func GetPtrString(v string) *string {
	return &v
}

func GetPtrTime(v time.Time) *time.Time {
	return &v
}
