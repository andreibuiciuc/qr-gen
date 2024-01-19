package util

import (
	"math"
	"strconv"
	"strings"
)

type QrErrorCorrectionInfo struct {
	TotalDataCodewords          int
	ECCodewordsPerBlock         int
	NumBlocksGroup1             int
	DataCodeworkdsInGroup1Block int
	NumBlocksGroup2             int
	DataCodewordsInGroup2Block  int
}

const QrCodewordSize = 8

var logTable = make([]int, 256)
var antilogTable = make([]int, 256)

var QrEcInfo = map[string]QrErrorCorrectionInfo{
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

func ComputeAlphaToPower(alpha, power int) int {
	result := 1

	for i := 0; i < power; i++ {
		prod := result * alpha

		if prod >= 256 {
			prod = prod ^ 285
		}

		result = prod
	}

	return result
}

func ComputeLogAntilogTables() {
	for i := 0; i < 256; i++ {
		logTable[i] = ComputeAlphaToPower(2, i)
		antilogTable[logTable[i]] = i
	}
	antilogTable[1] = 0
}

func ConvertValueToExponent(value int) int {
	if value < 0 || value > len(antilogTable)-1 {
		return 0
	}

	return antilogTable[value]
}

func ConvertExponentToValue(exp int) int {
	if exp < 0 || exp > len(logTable)-1 {
		return 0
	}

	return logTable[exp]
}

func ConvertIntListToBin(list []int) []string {
	result := make([]string, len(list))

	for i, elem := range list {
		bin := strconv.FormatInt(int64(elem), 2)
		result[i] = PadLeft(bin, "0", 8)
	}

	return result
}

func ConvertIntListToCodewords(list []int) string {
	codewords := ConvertIntListToBin(list)
	return strings.Join(codewords, "")
}

func GetClosestMultiple(n int, multipleOf int) int {
	multiple := int(math.Round(float64(n) / float64(multipleOf)))
	return multiple * multipleOf
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func PadLeft(s string, c string, n int) string {
	return strings.Repeat(c, n-len(s)) + s
}

func PadRight(s string, c string, n int) string {
	return s + strings.Repeat(c, n-len(s))
}

func SplitInGroups(s string, n int) []string {
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

// TODO: Fix cycle import
func GetECMappingKey(version int, lvl string) string {
	return strconv.Itoa(int(version)) + "-" + string(lvl)
}
