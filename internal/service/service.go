package service

import (
    "errors"
    "strings"

    "github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoDetectAndConvert(input string) (string, error) {
    normalized := strings.TrimSpace(input)
    if normalized == "" {
        return "", errors.New("empty input")
    }

    if isMorseCode(normalized) {
        return morseToText(normalized)
    }
    return textToMorse(normalized)
}

func isMorseCode(s string) bool {
    converter := morse.DefaultConverter
    result := converter.ToText(s)
    return !strings.Contains(result, "�") && 
           !strings.Contains(result, string(morse.ErrNoEncoding{}.Text))
}

func textToMorse(text string) (string, error) {
    converter := morse.NewConverter(
        morse.DefaultMorse,
        morse.WithCharSeparator(" "),
        morse.WithWordSeparator("   "),
        morse.WithLowercaseHandling(true),
        morse.WithHandler(handleError),
    )
    return converter.ToMorse(text), nil
}

func morseToText(morseStr string) (string, error) {
    converter := morse.NewConverter(
        morse.DefaultMorse,
        morse.WithCharSeparator(" "),
        morse.WithWordSeparator("   "),
        morse.WithHandler(handleError),
    )
    result := converter.ToText(morseStr)
    if strings.Contains(result, "�") {
        return "", errors.New("invalid morse code detected")
    }
    return result, nil
}

func handleError(err error) string {
    return "�"
}
