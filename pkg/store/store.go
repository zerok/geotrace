package store

import (
	"context"
	"time"
)

type Store interface {
	Add(ts time.Time, coordinates []float64, deviceID string) error
	Open(context.Context) error
	Close(context.Context) error
	Count(context.Context) (int64, error)
}
