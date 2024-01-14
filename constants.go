package main

type QrNumericMask int
type QrAlphanumericMask int
type QrByteMask int
type QrPaddingByte int
type QrPolynomial []int

type QrECInformation struct {
	TotalDataCodewords          int
	ECCodewordsPerBlock         int
	NumBlocksGroup1             int
	DataCodeworkdsInGroup1Block int
	NumBlocksGroup2             int
	DataCodewordsInGroup2Block  int
}

const QR_ALPHA_NUMERIC_FACTOR = 45
const QR_CODEWORD_SIZE = 8
const QR_GALOIS_ORDER = 256
const QR_GALOIS_MOD_VALUE = 285

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

var QR_EC_INFO = map[string]QrECInformation{
	// Extend this map to support higher levels of QR codes
	"1-L": {19, 7, 1, 19, 0, 0},
	"1-M": {16, 10, 1, 16, 0, 0},
	"1-Q": {13, 13, 1, 13, 0, 0},
	"1-H": {9, 17, 1, 9, 0, 0},
	"2-L": {34, 10, 1, 34, 0, 0},
	"2-M": {28, 16, 1, 28, 0, 0},
	"2-Q": {22, 22, 1, 22, 0, 0},
	"2-H": {16, 28, 1, 16, 0, 0},
	"3-L": {55, 15, 1, 55, 0, 0},
	"3-M": {44, 26, 1, 44, 0, 0},
	"3-Q": {34, 18, 2, 17, 0, 0},
	"3-H": {26, 22, 2, 13, 0, 0},
	"4-L": {80, 20, 1, 80, 0, 0},
	"4-M": {64, 18, 2, 32, 0, 0},
	"4-Q": {48, 26, 2, 24, 0, 0},
	"4-H": {36, 16, 4, 9, 0, 0},
	"5-L": {108, 26, 1, 108, 0, 0},
	"5-M": {86, 24, 2, 43, 0, 0},
	"5-Q": {62, 18, 2, 15, 2, 16},
	"5-H": {46, 22, 2, 11, 2, 12},
}

var QR_REMAINDER_BITS = map[QrVersion]int{
	1:  0,
	2:  7,
	3:  7,
	4:  7,
	5:  7,
	6:  7,
	7:  0,
	8:  0,
	9:  0,
	10: 0,
	11: 0,
	12: 0,
	13: 0,
	14: 3,
	15: 3,
	16: 3,
	17: 3,
	18: 3,
	19: 3,
	20: 3,
	21: 4,
	22: 4,
	23: 4,
	24: 4,
	25: 4,
	26: 4,
	27: 4,
	28: 3,
	29: 3,
	30: 3,
	31: 3,
	32: 3,
	33: 3,
	34: 3,
	35: 0,
	36: 0,
	37: 0,
	38: 0,
	39: 0,
	40: 0,
}

var QR_PADDING_BYTES = map[QrPaddingByte]string{
	FIRST:  "11101100",
	SECOND: "00010001",
}

var QR_GALOIS_LOG_TABLE [QR_GALOIS_ORDER]int
var QR_GALOIS_ANTILOG_TABLE [QR_GALOIS_ORDER]int
