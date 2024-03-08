package qr

import (
	"fmt"
	"regexp"
	"strconv"
)

type versioner struct{}

const (
	mode_NUMERIC      string = "numeric"
	mode_ALPHANUMERIC string = "alphanumeric"
	mode_BYTE         string = "byte"
)

const (
	ec_LOW      rune = 'L'
	ec_MEDIUM   rune = 'M'
	ec_QUARTILE rune = 'Q'
	ec_HIGH     rune = 'H'
)

const (
	indicator_NUMERIC      string = "0001"
	indicator_ALPHANUMERIC string = "0010"
	indicator_BYTE         string = "0100"
)

var qrModeRegexes = map[string]string{
	mode_NUMERIC:      "^\\d+$",
	mode_ALPHANUMERIC: "^[\\dA-Z $%*+\\-./:]+$",
	mode_BYTE:         "^[\\x00-\\xff]+$",
}

var qrModeIndices = map[string]int{
	mode_NUMERIC:      0,
	mode_ALPHANUMERIC: 1,
	mode_BYTE:         2,
}

var qrModeIndicators = map[string]string{
	mode_NUMERIC:      indicator_NUMERIC,
	mode_ALPHANUMERIC: indicator_ALPHANUMERIC,
	mode_BYTE:         indicator_BYTE,
}

var qrCountIndLengths = map[string]int{
	mode_NUMERIC:      10,
	mode_ALPHANUMERIC: 9,
	mode_BYTE:         8,
}

var qrCapacities = map[int]map[rune][]int{
	1: {
		ec_LOW:      {41, 25, 17},
		ec_MEDIUM:   {34, 20, 14},
		ec_QUARTILE: {27, 16, 11},
		ec_HIGH:     {17, 10, 7},
	},
	2: {
		ec_LOW:      {77, 47, 32},
		ec_MEDIUM:   {63, 38, 26},
		ec_QUARTILE: {48, 29, 20},
		ec_HIGH:     {34, 20, 14},
	},
	3: {
		ec_LOW:      {127, 77, 53},
		ec_MEDIUM:   {101, 61, 42},
		ec_QUARTILE: {77, 47, 32},
		ec_HIGH:     {58, 35, 24},
	},
	4: {
		ec_LOW:      {187, 114, 78},
		ec_MEDIUM:   {149, 90, 62},
		ec_QUARTILE: {111, 67, 46},
		ec_HIGH:     {82, 50, 34},
	},
	5: {
		ec_LOW:      {255, 154, 106},
		ec_MEDIUM:   {202, 122, 84},
		ec_QUARTILE: {144, 87, 60},
		ec_HIGH:     {106, 64, 44},
	},
}

func newVersioner() *versioner {
	return &versioner{}
}

func (v *versioner) getMode(s string) (string, error) {
	if matched, _ := regexp.MatchString(qrModeRegexes[mode_NUMERIC], s); matched {
		return mode_NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(qrModeRegexes[mode_ALPHANUMERIC], s); matched {
		return mode_ALPHANUMERIC, nil
	}

	if matched, _ := regexp.MatchString(qrModeRegexes[mode_BYTE], s); matched {
		return mode_BYTE, nil
	}

	return string(""), fmt.Errorf("invalid input pattern")
}

func (v *versioner) getVersion(s string, string string, lvl rune) (int, error) {
	ver := 1

	for ver <= len(qrCapacities) {
		if len(s) <= qrCapacities[int(ver)][lvl][qrModeIndices[string]] {
			return int(ver), nil
		}
		ver += 1
	}

	return int(-1), fmt.Errorf("cannot compute qr int")
}

func (v *versioner) getModeIndicator(string string) string {
	return qrModeIndicators[string]
}

func (v *versioner) getCountIndicator(s string, int int, string string) (string, error) {
	sLenBin := strconv.FormatInt(int64(len(s)), 2)
	return padLeft(sLenBin, "0", qrCountIndLengths[string]), nil
}
