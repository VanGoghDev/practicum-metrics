package storage

import "errors"

var (
	ErrGaugesTableNil   = errors.New("gauges table is not initialized")
	ErrCountersTableNil = errors.New("counter table is not initialized")
	ErrNotFound         = errors.New("gauge not found")
	ErrURLExists        = errors.New("url exists")
)
