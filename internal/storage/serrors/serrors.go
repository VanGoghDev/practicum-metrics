package serrors

import "errors"

var (
	// ErrGaugesTableNil таблица метрик типа gauges не инициализирована.
	ErrGaugesTableNil = errors.New("gauges table is not initialized")

	// ErrCountersTableNil таблица метрик типа counter не инициализирована.
	ErrCountersTableNil = errors.New("counter table is not initialized")

	// ErrNotFound ошибка о том, что не удалось найти запрашиваем ресурс.
	ErrNotFound = errors.New("gauge not found")

	// ErrURLExists  ошибка о том, что url уже существует.
	ErrURLExists = errors.New("url exists")
)
