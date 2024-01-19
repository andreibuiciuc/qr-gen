package encoder

import (
	"fmt"
	"qr/qr-gen/util"
	"qr/qr-gen/versioner"
	"strconv"
	"strings"
)

type Encoder interface {
	EncodeNumericInput(s string) string
	EncodeAlphanumericInput(s string) string
	EncodeByteInput(s string) string
	Encode(s string, lvl versioner.QrEcLevel) (string, error)
	AugmentEncodedInput(s string, version versioner.QrVersion, lvl versioner.QrEcLevel) string
}

type QrEncoder struct{}
type QrNumericMask int
type QrAlphanumericMask int
type QrByteMask int
type QrPaddingByte int

func New() Encoder {
	return &QrEncoder{}
}

func (e *QrEncoder) Encode(s string, lvl versioner.QrEcLevel) (string, error) {
	v := versioner.New()

	mode, err := v.GetMode(s)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding mode: %v", err)
	}

	modeIndicator := v.GetModeIndicator(mode)

	version, err := v.GetVersion(s, mode, lvl)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding version: %v", err)
	}

	countIndicator, err := v.GetCountIndicator(s, version, mode)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding count indicator: %v", err)
	}

	encodedInput := e.EncodeInput(s, versioner.QrMode(mode))

	return string(modeIndicator) + countIndicator + encodedInput, nil
}

func (e *QrEncoder) EncodeNumericInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[versioner.QrMode(versioner.QrNumericMode)])
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryString := strconv.FormatInt(int64(numericValue), 2)

		switch true {
		case numericValue <= 9:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[DIGIT])
		case 10 <= numericValue && numericValue <= 99:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[TEN])
		default:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[HUNDRED])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *QrEncoder) EncodeAlphanumericInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[versioner.QrMode(versioner.QrAlphanumericMode)])
	result := make([]string, len(groups))

	for index, group := range groups {
		var binaryString string
		var firstCharValue int
		var secondCharValue int
		var groupValue int

		if len(group) == 2 {
			firstCharValue, secondCharValue = ALPHA_NUMERIC_VALUES[group[0]], ALPHA_NUMERIC_VALUES[group[1]]
			groupValue = QR_ALPHA_NUMERIC_FACTOR*firstCharValue + secondCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = util.PadLeft(binaryString, "0", QR_ALPHA_NUMERIC_MASKS[FULL_GROUP])
		} else {
			firstCharValue = ALPHA_NUMERIC_VALUES[group[0]]
			groupValue = firstCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = util.PadLeft(binaryString, "0", QR_ALPHA_NUMERIC_MASKS[ONE_ONLY])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *QrEncoder) EncodeByteInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[versioner.QrMode(versioner.QrByteMode)])
	result := make([]string, len(groups))

	for index, group := range groups {
		hex := strconv.FormatInt(int64(group[0]), 16)
		hex0, _ := strconv.ParseInt(string(hex[0]), 16, 64)
		hex1, _ := strconv.ParseInt(string(hex[1]), 16, 64)

		bin0 := util.PadLeft(strconv.FormatInt(hex0, 2), "0", QR_BYTE_MASKS[CHAR])
		bin1 := util.PadLeft(strconv.FormatInt(hex1, 2), "0", QR_BYTE_MASKS[CHAR])

		result[index] = bin0 + bin1
	}

	return strings.Join(result, "")
}

func (e *QrEncoder) EncodeInput(s string, mode versioner.QrMode) string {
	switch mode {
	case versioner.QrMode(versioner.QrNumericMode):
		return e.EncodeNumericInput(s)
	case versioner.QrMode(versioner.QrAlphanumericMode):
		return e.EncodeAlphanumericInput(s)
	case versioner.QrMode(versioner.QrByteMode):
		return e.EncodeByteInput(s)
	default:
		return ""
	}
}

func (e *QrEncoder) AugmentEncodedInput(s string, version versioner.QrVersion, lvl versioner.QrEcLevel) string {
	requiredBitsCount := e.getNumberOfRequiredBits(version, lvl)

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

func (e *QrEncoder) getNumberOfRequiredBits(version versioner.QrVersion, lvl versioner.QrEcLevel) int {
	key := util.GetECMappingKey(int(version), string(lvl))
	return util.QrCodewordSize * util.QrEcInfo[key].TotalDataCodewords
}

func (e *QrEncoder) augmentWithTerminatorBits(s string, requiredBitsCount int) string {
	remainingBitsCount := requiredBitsCount - len(s)

	if remainingBitsCount >= 4 {
		return util.PadRight(s, "0", len(s)+4)
	}

	return util.PadRight(s, "0", len(s)+remainingBitsCount)
}

func (e *QrEncoder) augmentWithPaddingBits(s string, requiredBitsCount int) string {
	numberOfPadBytes := (requiredBitsCount - len(s)) / util.QrCodewordSize
	paddingByteIndex := 0
	paddingSequence := ""

	for i := 0; i < numberOfPadBytes; i++ {
		if paddingByteIndex == 2 {
			paddingByteIndex = paddingByteIndex % 2
		}

		paddingSequence = paddingSequence + QR_PADDING_BYTES[QrPaddingByte(paddingByteIndex)]
		paddingByteIndex += 1
	}

	return s + paddingSequence
}

func (e *QrEncoder) augmentWithZeroBits(s string) string {
	multiple := util.GetClosestMultiple(len(s), util.QrCodewordSize)
	return util.PadRight(s, "0", multiple)
}

const QR_ALPHA_NUMERIC_FACTOR = 45

const (
	DIGIT QrNumericMask = iota
	TEN
	HUNDRED
)

const (
	FULL_GROUP QrAlphanumericMask = iota
	ONE_ONLY
)

const (
	CHAR QrByteMask = iota
)

const (
	FIRST QrPaddingByte = iota
	SECOND
)

var SPLIT_VALUES = map[versioner.QrMode]int{
	"numeric":      3,
	"alphanumeric": 2,
	"byte":         1,
}

var QR_NUMERIC_MASKS = map[QrNumericMask]int{
	DIGIT:   4,
	TEN:     7,
	HUNDRED: 10,
}

var QR_ALPHA_NUMERIC_MASKS = map[QrAlphanumericMask]int{
	FULL_GROUP: 11,
	ONE_ONLY:   6,
}

var QR_BYTE_MASKS = map[QrByteMask]int{
	CHAR: 4,
}

var ALPHA_NUMERIC_VALUES = map[byte]int{
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

var QR_PADDING_BYTES = map[QrPaddingByte]string{
	FIRST:  "11101100",
	SECOND: "00010001",
}
