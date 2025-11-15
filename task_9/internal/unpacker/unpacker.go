package unpacker

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const escapeSymbol = '\\'

var (
	ErrStringMustNotContainOnlyNumbers = errors.New("string must not contain only numbers")
	ErrInvalidStringFormat             = errors.New("invalid string format")
	ErrInvalidDigitConversion          = errors.New("invalid digit conversion")
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	_, err := strconv.Atoi(s)
	if err == nil {
		return "", ErrStringMustNotContainOnlyNumbers
	}

	var result strings.Builder
	currentSymbol := ""

	for _, char := range s {
		if currentSymbol == "" {
			if unicode.IsDigit(char) {
				return "", ErrInvalidStringFormat

			}

			if char == escapeSymbol {
				currentSymbol = string(char)
				continue
			}
		}

		if currentSymbol == string(escapeSymbol) {
			currentSymbol = string(char)
			continue
		}

		if char == escapeSymbol {
			result.WriteString(currentSymbol)
			currentSymbol = string(char)
			continue
		}

		if unicode.IsLetter(char) {
			if currentSymbol != "" {
				result.WriteString(currentSymbol)
			}

			currentSymbol = string(char)
		}

		if unicode.IsDigit(char) {
			count, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidDigitConversion
			}

			for range count {
				result.WriteString(currentSymbol)
			}

			currentSymbol = ""
		}
	}

	if currentSymbol != "" {
		result.WriteString(currentSymbol)
	}

	return result.String(), nil
}
