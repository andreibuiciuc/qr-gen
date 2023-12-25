package main

import (
	"fmt"
	"regexp"
)

type QrCodeMode int

const (
	INVALID QrCodeMode = iota
	NUMERIC
	ALPHA_NUMERIC
	BYTE
)

const (
	NUMERIC_REGEX = "^\\d+$"
)

var patterns = map[QrCodeMode]string{
	NUMERIC:       "^\\d+$",
	ALPHA_NUMERIC: "^[\\dA-Z $%*+\\-./:]+$",
	BYTE:          "^[\\x00-\\xff]+$",
}

func GetMode(s string) (QrCodeMode, error) {
	if matched, _ := regexp.MatchString(patterns[NUMERIC], s); matched {
		return NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(patterns[ALPHA_NUMERIC], s); matched {
		return ALPHA_NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(patterns[BYTE], s); matched {
		return BYTE, nil
	}

	return INVALID, fmt.Errorf("Invalid input pattern")
}

func main() {
}
