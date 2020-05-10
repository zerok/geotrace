package store

import "time"

type Store interface {
	Add(ts time.Time, coordinates []float64, deviceID string) error
}
