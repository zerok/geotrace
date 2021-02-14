package store

import (
	"context"
	"time"
)

type Trace struct {
	Time string  `json:"time"`
	Lon  float64 `json:"lon"`
	Lat  float64 `json:"lat"`
}

type Store interface {
	Add(ts time.Time, coordinates []float64, deviceID string) error
	Open(context.Context) error
	Close(context.Context) error
	Count(context.Context) (int64, error)
	GetTracesSince(context.Context, time.Time) ([]Trace, error)
}
