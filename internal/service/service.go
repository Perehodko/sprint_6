package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoDetectAndConvert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}

	converter := morse.NewConverter(
		morse.DefaultMorse,
		morse.WithCharSeparator(" "),
		morse.WithWordSeparator("   "),
		morse.WithLowercaseHandling(true),
	)

	if isMorseCode(input) {
		result := converter.ToText(input)
		if strings.Contains(result, "�") {
			return "", errors.New("invalid morse code")
		}
		return result, nil
	}

	for _, r := range input {
		if _, exists := morse.DefaultMorse[unicode.ToUpper(r)]; !exists && !unicode.IsSpace(r) {
			return "", errors.New("unsupported character: " + string(r))
		}
	}

	return converter.ToMorse(input), nil
}

func isMorseCode(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch != '.' && ch != '-' && ch != ' ' && ch != '/' {
			return false
		}
	}
	return true
}
