package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidString .
var ErrInvalidString = errors.New("invalid string")

// Unpack распаковывает строку с повторяющимися рунами.
func Unpack(raw string) (string, error) {
	var (
		err   error
		cache string // предыдущий символ в строке
	)

	builder := strings.Builder{}

	for _, symbol := range strings.Split(raw, "") {
		var count int // кол-во повторений предыдущего символа

		if count, err = strconv.Atoi(symbol); err == nil { // текущий символ - цифра
			if _, err = strconv.Atoi(cache); err == nil || cache == "" { // 2 цифры подряд, или строка начинается с цифры
				return "", ErrInvalidString
			}
		} else if _, err = strconv.Atoi(cache); err != nil { // предыдущий и текущий символы не цифры
			count = 1
		}

		builder.WriteString(strings.Repeat(cache, count))
		cache = symbol
	}

	if _, err = strconv.Atoi(cache); err != nil { // отдельно дописываем последний символ, если это не цифра
		builder.WriteString(cache)
	}

	return builder.String(), nil
}
