package qr

import (
	"fmt"
	"strconv"
	"strings"
)

type encoder struct{}

const alphaNumericFactor = 45

const (
	mask_DIGIT int = iota
	mask_TEN
	mask_HUNDRED
)

const (
	mask_FULL_GROUP int = iota
	mask_ONE_ONLY
)

const (
	mask_CHAR int = iota
)

const (
	byte_FIRST int = iota
	byte_SECOND
)

var splitValues = map[string]int{
	"numeric":      3,
	"alphanumeric": 2,
	"byte":         1,
}

var numericMasks = map[int]int{
	mask_DIGIT:   4,
	mask_TEN:     7,
	mask_HUNDRED: 10,
}

var alphanumericMask = map[int]int{
	mask_FULL_GROUP: 11,
	mask_ONE_ONLY:   6,
}

var byteMasks = map[int]int{
	mask_CHAR: 4,
}

var alphanumericValues = map[byte]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
	'G': 16,
	'H': 17,
	'I': 18,
	'J': 19,
	'K': 20,
	'L': 21,
	'M': 22,
	'N': 23,
	'O': 24,
	'P': 25,
	'Q': 26,
	'R': 27,
	'S': 28,
	'T': 29,
	'U': 30,
	'V': 31,
	'W': 32,
	'X': 33,
	'Y': 34,
	'Z': 35,
	' ': 36,
	'$': 37,
	'%': 38,
	'*': 39,
	'+': 40,
	'-': 41,
	'.': 42,
	'/': 43,
	':': 44,
}

var paddingBytes = map[int]string{
	byte_FIRST:  "11101100",
	byte_SECOND: "00010001",
}

func newEncoder() *encoder {
	return &encoder{}
}

func (e *encoder) encode(s string, lvl rune) (string, error) {
	v := newVersioner()

	m, err := v.getMode(s)
	if err != nil {
		return "", fmt.Errorf("error on computing the encoding string: %v", err)
	}

	mi := v.getModeIndicator(m)

	int, err := v.getVersion(s, m, lvl)
	if err != nil {
		return "", fmt.Errorf("error on computing the encoding int: %v", err)
	}

	countIndicator, err := v.getCountIndicator(s, int, m)
	if err != nil {
		return "", fmt.Errorf("error on computing the encoding count indicator: %v", err)
	}

	encodedInput := e.encodeInput(s, string(m))

	return string(mi) + countIndicator + encodedInput, nil
}

func (e *encoder) encodeNumericInput(s string) string {
	groups := splitInGroups(s, splitValues[string(mode_NUMERIC)])
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryString := strconv.FormatInt(int64(numericValue), 2)

		switch true {
		case numericValue <= 9:
			binaryString = padLeft(binaryString, "0", numericMasks[mask_DIGIT])
		case 10 <= numericValue && numericValue <= 99:
			binaryString = padLeft(binaryString, "0", numericMasks[mask_TEN])
		default:
			binaryString = padLeft(binaryString, "0", numericMasks[mask_HUNDRED])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *encoder) encodeAlphanumericInput(s string) string {
	groups := splitInGroups(s, splitValues[string(mode_ALPHANUMERIC)])
	result := make([]string, len(groups))

	for index, group := range groups {
		var binaryString string
		var firstCharValue int
		var secondCharValue int
		var groupValue int

		if len(group) == 2 {
			firstCharValue, secondCharValue = alphanumericValues[group[0]], alphanumericValues[group[1]]
			groupValue = alphaNumericFactor*firstCharValue + secondCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = padLeft(binaryString, "0", alphanumericMask[mask_FULL_GROUP])
		} else {
			firstCharValue = alphanumericValues[group[0]]
			groupValue = firstCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = padLeft(binaryString, "0", alphanumericMask[mask_ONE_ONLY])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *encoder) encodeByteInput(s string) string {
	groups := splitInGroups(s, splitValues[string(mode_BYTE)])
	result := make([]string, len(groups))

	for index, group := range groups {
		hex := strconv.FormatInt(int64(group[0]), 16)
		hex0, _ := strconv.ParseInt(string(hex[0]), 16, 64)
		hex1, _ := strconv.ParseInt(string(hex[1]), 16, 64)

		bin0 := padLeft(strconv.FormatInt(hex0, 2), "0", byteMasks[mask_CHAR])
		bin1 := padLeft(strconv.FormatInt(hex1, 2), "0", byteMasks[mask_CHAR])

		result[index] = bin0 + bin1
	}

	return strings.Join(result, "")
}

func (e *encoder) encodeInput(s string, m string) string {
	switch m {
	case mode_NUMERIC:
		return e.encodeNumericInput(s)
	case mode_ALPHANUMERIC:
		return e.encodeAlphanumericInput(s)
	case mode_BYTE:
		return e.encodeByteInput(s)
	default:
		return ""
	}
}

func (e *encoder) augmentEncodedInput(s string, v int, lvl rune) string {
	requiredBitsCount := e.getNumberOfRequiredBits(v, lvl)

	s = e.augmentWithTerminatorBits(s, requiredBitsCount)
	remainingBitsCount := requiredBitsCount - len(s)
	if remainingBitsCount == 0 {
		return s
	}

	s = e.augmentWithZeroBits(s)
	remainingBitsCount = requiredBitsCount - len(s)
	if remainingBitsCount == 0 {
		return s
	}

	return e.augmentWithPaddingBits(s, requiredBitsCount)
}

func (e *encoder) getNumberOfRequiredBits(v int, lvl rune) int {
	key := getECMappingKey(v, string(lvl))
	return codewordSize * ecInfo[key].TotalDataCodewords
}

func (e *encoder) augmentWithTerminatorBits(s string, requiredBitsCount int) string {
	remainingBitsCount := requiredBitsCount - len(s)

	if remainingBitsCount >= 4 {
		return padRight(s, "0", len(s)+4)
	}

	return padRight(s, "0", len(s)+remainingBitsCount)
}

func (e *encoder) augmentWithPaddingBits(s string, requiredBitsCount int) string {
	numberOfPadBytes := (requiredBitsCount - len(s)) / codewordSize
	paddingByteIndex := 0
	paddingSequence := ""

	for i := 0; i < numberOfPadBytes; i++ {
		if paddingByteIndex == 2 {
			paddingByteIndex = paddingByteIndex % 2
		}

		paddingSequence = paddingSequence + paddingBytes[paddingByteIndex]
		paddingByteIndex += 1
	}

	return s + paddingSequence
}

func (e *encoder) augmentWithZeroBits(s string) string {
	multiple := getClosestMultiple(len(s), codewordSize)
	return padRight(s, "0", multiple)
}
