package main

import (
	"fmt"
	"regexp"
)

type QrCodeMode int
type QrErrCorrectionLvl int
type QrVersion int

const (
	INVALID_MODE QrCodeMode = iota
	NUMERIC
	ALPHA_NUMERIC
	BYTE
)

const (
	LOW QrErrCorrectionLvl = iota
	MEDIUM
	QUARTILE
	HIGH
)

const (
	INVALID_VERSION QrVersion = iota
	VERSION_1
	VERSION_2
	VERSION_3
	VERSION_4
	VERSION_5
)

var patterns = map[QrCodeMode]string{
	NUMERIC:       "^\\d+$",
	ALPHA_NUMERIC: "^[\\dA-Z $%*+\\-./:]+$",
	BYTE:          "^[\\x00-\\xff]+$",
}

var capacities = map[QrVersion]map[QrErrCorrectionLvl][]int{
	VERSION_1: {
		LOW:      {41, 25, 17},
		MEDIUM:   {34, 20, 14},
		QUARTILE: {27, 16, 11},
		HIGH:     {17, 10, 7},
	},
	VERSION_2: {
		LOW:      {77, 47, 32},
		MEDIUM:   {63, 38, 26},
		QUARTILE: {48, 29, 20},
		HIGH:     {34, 20, 14},
	},
	VERSION_3: {
		LOW:      {127, 77, 53},
		MEDIUM:   {101, 61, 42},
		QUARTILE: {77, 47, 32},
		HIGH:     {58, 35, 24},
	},
	VERSION_4: {
		LOW:      {187, 114, 78},
		MEDIUM:   {149, 90, 62},
		QUARTILE: {111, 67, 46},
		HIGH:     {82, 50, 34},
	},
	VERSION_5: {
		LOW:      {255, 154, 106},
		MEDIUM:   {202, 122, 84},
		QUARTILE: {144, 87, 60},
		HIGH:     {106, 64, 44},
	},
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

	return QrCodeMode(0), fmt.Errorf("Invalid input pattern")
}

func GetSmallestVersion(s string, level QrErrCorrectionLvl) (QrVersion, error) {
	version := 1
	mode, err := GetMode(s)

	if err != nil {
		return 0, err
	}

	for version <= len(capacities) {
		if len(s) <= capacities[QrVersion(version)][level][mode-1] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return INVALID_VERSION, fmt.Errorf("Cannot compute QR version")
}

func main() {
}
