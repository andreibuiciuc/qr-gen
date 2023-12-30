package main

const INVALID_IDX = -1
const EMPTY_STRING = ""
const ALPHA_NUMERIC_FACTOR = 45
const BINARY_RADIX = 2
const HEXADECIMAL_RADIX = 16
const INTEGER_RADIX = 64
const DEFAULT_PAD_CHAR = "0"
const CODEWORD_BITS = 8

const (
	NUMERIC       QrMode = "numeric"
	ALPHA_NUMERIC QrMode = "alpha"
	BYTE          QrMode = "byte"
)

const (
	LOW      QrErrCorrectionLvl = "L"
	MEDIUM   QrErrCorrectionLvl = "M"
	QUARTILE QrErrCorrectionLvl = "Q"
	HIGH     QrErrCorrectionLvl = "H"
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
	NUMERIC_INDICATOR       QrModeIndicator = "0001"
	ALPHA_NUMERIC_INDICATOR QrModeIndicator = "0010"
	BYTE_INDICATOR          QrModeIndicator = "0100"
)

type QrNumericMask int
type QrAlphanumericMask int
type QrByteMask int
type QrPaddingByte int

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

var PATTERNS = map[QrMode]string{
	NUMERIC:       "^\\d+$",
	ALPHA_NUMERIC: "^[\\dA-Z $%*+\\-./:]+$",
	BYTE:          "^[\\x00-\\xff]+$",
}

var CAPACITIES = map[QrVersion]map[QrErrCorrectionLvl][]int{
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

var MODE_INDICATORS = map[QrMode]QrModeIndicator{
	NUMERIC:       NUMERIC_INDICATOR,
	ALPHA_NUMERIC: ALPHA_NUMERIC_INDICATOR,
	BYTE:          BYTE_INDICATOR,
}

var MODE_INDICES = map[QrMode]int{
	NUMERIC:       0,
	ALPHA_NUMERIC: 1,
	BYTE:          2,
}

var SPLIT_VALUES = map[QrMode]int{
	NUMERIC:       3,
	ALPHA_NUMERIC: 2,
	BYTE:          1,
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

var NUMERIC_MASKS = map[QrNumericMask]int{
	DIGIT:   4,
	TEN:     7,
	HUNDRED: 10,
}

var ALPHA_NUMERIC_MASKS = map[QrAlphanumericMask]int{
	FULL_GROUP: 11,
	ONE_ONLY:   6,
}

var BYTE_MASKS = map[QrByteMask]int{
	CHAR: 4,
}

var ERR_CORR_TOTAL_DATA = map[string]int{
	"1-L": 19,
	"1-M": 16,
	"1-Q": 13,
	"1-H": 9,
	"2-L": 34,
	"2-M": 28,
	"2-Q": 22,
	"2-H": 16,
	"3-L": 55,
	"3-M": 44,
	"3-Q": 34,
	"3-H": 26,
	"4-L": 80,
	"4-M": 64,
	"4-Q": 48,
	"4-H": 36,
	"5-L": 108,
	"5-M": 86,
	"5-Q": 62,
	"5-H": 46,
}

var PADDING_BYTES = map[QrPaddingByte]string{
	FIRST:  "11101100",
	SECOND: "00010001",
}
