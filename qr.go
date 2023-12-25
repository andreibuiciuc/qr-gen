package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type QrMode string
type QrErrCorrectionLvl string
type QrVersion int
type QrModeIndicator byte

const INVALID_IDX = -1
const EMPTY_STRING = ""

const (
	NUMERIC       QrMode = "numeric"
	ALPHA_NUMERIC QrMode = "alpha"
	BYTE          QrMode = "byte"
)

const (
	LOW      QrErrCorrectionLvl = "low"
	MEDIUM   QrErrCorrectionLvl = "medium"
	QUARTILE QrErrCorrectionLvl = "quartile"
	HIGH     QrErrCorrectionLvl = "high"
)

const (
	_ QrVersion = iota
	VERSION_1
	VERSION_2
	VERSION_3
	VERSION_4
	VERSION_5
)

const (
	NUMERIC_INDICATOR       QrModeIndicator = 0b0001
	ALPHA_NUMERIC_INDICATOR QrModeIndicator = 0b0010
	BYTE_INDICATOR          QrModeIndicator = 0b0100
)

var patterns = map[QrMode]string{
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

var modeIndicators = map[QrMode]QrModeIndicator{
	NUMERIC:       NUMERIC_INDICATOR,
	ALPHA_NUMERIC: ALPHA_NUMERIC_INDICATOR,
	BYTE:          BYTE_INDICATOR,
}

var modeToIndex = map[QrMode]int{
	NUMERIC:       0,
	ALPHA_NUMERIC: 1,
	BYTE:          2,
}

type QrEncoder interface {
	GetMode(s string) (QrMode, error)
	GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error)
	GetModeIndicator(mode QrMode) QrModeIndicator
	GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error)
	Encode(s string) string
}

type Encoder struct {
	splitCount int
}

func NewEncoder() QrEncoder {
	return &Encoder{
		splitCount: 3,
	}
}

func (e *Encoder) GetMode(s string) (QrMode, error) {
	if matched, _ := regexp.MatchString(patterns[NUMERIC], s); matched {
		return NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(patterns[ALPHA_NUMERIC], s); matched {
		return ALPHA_NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(patterns[BYTE], s); matched {
		return BYTE, nil
	}

	return QrMode(EMPTY_STRING), fmt.Errorf("Invalid input pattern")
}

func (e *Encoder) GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error) {
	version := 1

	for version <= len(capacities) {
		if len(s) <= capacities[QrVersion(version)][lvl][modeToIndex[mode]] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return QrVersion(INVALID_IDX), fmt.Errorf("Cannot compute QR version")
}

func (e *Encoder) GetModeIndicator(mode QrMode) QrModeIndicator {
	return modeIndicators[mode]
}

func (e *Encoder) GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error) {
	cntIndicatorLen, err := e.getCountIndicatorLen(version, mode)

	if err != nil {
		return EMPTY_STRING, err
	}

	sLenBinary := strconv.FormatInt(int64(len(s)), 2)
	return strings.Repeat("0", cntIndicatorLen-len(sLenBinary)), nil
}

func (e *Encoder) Encode(s string) string {
	groups := e.splitInGroups(s, e.splitCount)
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryValue := strconv.FormatInt(int64(numericValue), 2)
		result[index] = binaryValue
	}

	return strings.Join(result, "")
}

func (e *Encoder) getCountIndicatorLen(version QrVersion, mode QrMode) (int, error) {
	// Extend this functionality for further versions support
	if VERSION_1 <= version && version <= VERSION_5 {
		switch mode {
		case NUMERIC:
			return 10, nil
		case ALPHA_NUMERIC:
			return 9, nil
		case BYTE:
			return 8, nil
		}
	}

	return 0, fmt.Errorf("Cannot compute QR Count Indicator length")
}

func (e *Encoder) splitInGroups(s string, n int) []string {
	var result []string
	count := 0
	start := 0

	for i := 0; i < len(s); i++ {
		count += 1

		if count == n {
			result = append(result, s[start:i+1])
			start = i + 1
			count = 0
		}
	}

	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}

func main() {
}
