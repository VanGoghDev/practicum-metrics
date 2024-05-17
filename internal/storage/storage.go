package storage

import "errors"

var (
	ErrGaugesTableNil = errors.New("gauges table is not initialized")
	ErrURLExists      = errors.New("url exists")
)
