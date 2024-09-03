package converter

import (
	"errors"
	"strconv"
)

// ErrUnsupportedType ошибка означает, что переданный тип не поддерживается конвертером.
var ErrUnsupportedType = errors.New("type is unsupported")

// Str возвращает строковое отображение для переданного аргумента v.
func Str(v any) (string, error) {
	var result string
	switch v := v.(type) {
	default:
		return "", ErrUnsupportedType
	case int:
		result = strconv.Itoa(v)
	case int64:
		result = strconv.FormatInt(v, 10)
	case float32:
		result = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		result = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		result = strconv.FormatBool(v)
	}
	return result, nil
}
