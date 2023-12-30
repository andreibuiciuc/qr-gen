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

type QrEncoder interface {
	GetMode(s string) (QrMode, error)
	GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error)
	GetModeIndicator(mode QrMode) QrModeIndicator
	GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error)
	EncodeNumericInput(s string) string
	EncodeAlphanumericInput(s string) string
	EncodeByteInput(s string) string
	Encode(s string) string
	splitInGroups(s string, n int) []string
}

type Encoder struct{}

func NewEncoder() QrEncoder {
	return &Encoder{}
}

func (e *Encoder) GetMode(s string) (QrMode, error) {
	if matched, _ := regexp.MatchString(PATTERNS[NUMERIC], s); matched {
		return NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(PATTERNS[ALPHA_NUMERIC], s); matched {
		return ALPHA_NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(PATTERNS[BYTE], s); matched {
		return BYTE, nil
	}

	return QrMode(EMPTY_STRING), fmt.Errorf("Invalid input pattern")
}

func (e *Encoder) GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error) {
	version := 1

	for version <= len(CAPACITIES) {
		if len(s) <= CAPACITIES[QrVersion(version)][lvl][MODE_INDICES[mode]] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return QrVersion(INVALID_IDX), fmt.Errorf("Cannot compute QR version")
}

func (e *Encoder) GetModeIndicator(mode QrMode) QrModeIndicator {
	return MODE_INDICATORS[mode]
}

func (e *Encoder) GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error) {
	cntIndicatorLen, err := e.getCountIndicatorLen(version, mode)

	if err != nil {
		return EMPTY_STRING, err
	}

	sLenBinary := strconv.FormatInt(int64(len(s)), BINARY_RADIX)
	return e.padStart(sLenBinary, DEFAULT_PAD_CHAR, cntIndicatorLen), nil
}

func (e *Encoder) EncodeNumericInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryString := strconv.FormatInt(int64(numericValue), BINARY_RADIX)

		switch true {
		case numericValue <= 9:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[DIGIT])
		case 10 <= numericValue && numericValue <= 99:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[TEN])
		default:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[HUNDRED])
		}

		result[index] = binaryString
	}

	return strings.Join(result, EMPTY_STRING)
}

func (e *Encoder) EncodeAlphanumericInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[ALPHA_NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		var binaryString string
		var firstCharValue int
		var secondCharValue int
		var groupValue int

		if len(group) == 2 {
			firstCharValue, secondCharValue = ALPHA_NUMERIC_VALUES[group[0]], ALPHA_NUMERIC_VALUES[group[1]]
			groupValue = 45*firstCharValue + secondCharValue
			binaryString = strconv.FormatInt(int64(groupValue), BINARY_RADIX)
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, ALPHA_NUMERIC_MASKS[FULL_GROUP])
		} else {
			firstCharValue = ALPHA_NUMERIC_VALUES[group[0]]
			groupValue = firstCharValue
			binaryString = strconv.FormatInt(int64(groupValue), BINARY_RADIX)
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, ALPHA_NUMERIC_MASKS[ONE_ONLY])
		}

		result[index] = binaryString
	}

	return strings.Join(result, EMPTY_STRING)
}

func (e *Encoder) EncodeByteInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[BYTE])
	result := make([]string, len(groups))

	for index, group := range groups {
		hex := strconv.FormatInt(int64(group[0]), HEXADECIMAL_RADIX)
		hex0, _ := strconv.ParseInt(string(hex[0]), HEXADECIMAL_RADIX, INTEGER_RADIX)
		hex1, _ := strconv.ParseInt(string(hex[1]), HEXADECIMAL_RADIX, INTEGER_RADIX)

		bin0 := e.padStart(strconv.FormatInt(hex0, BINARY_RADIX), DEFAULT_PAD_CHAR, BYTE_MASKS[CHAR])
		bin1 := e.padStart(strconv.FormatInt(hex1, BINARY_RADIX), DEFAULT_PAD_CHAR, BYTE_MASKS[CHAR])

		result[index] = bin0 + bin1
	}

	return strings.Join(result, EMPTY_STRING)
}

func (e *Encoder) Encode(s string) string {
	return EMPTY_STRING
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

func (e *Encoder) padStart(s string, padChar string, n int) string {
	return strings.Repeat(padChar, n-len(s)) + s
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
