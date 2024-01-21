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

// ComputeAlphaToPower computes the power of a to p
// in the Galois field of order 256.
func ComputeAlphaToPower(a, p int) int {
	result := 1

	for i := 0; i < p; i++ {
		prod := result * a

		if prod >= 256 {
			prod = prod ^ 285
		}

		result = prod
	}

	return result
}

// ComputeLogAntilogTables computes the Log and Antilog tables
// in the Galois field of order 256.
func ComputeLogAntilogTables() {
	for i := 0; i < 256; i++ {
		logTable[i] = ComputeAlphaToPower(2, i)
		antilogTable[logTable[i]] = i
	}
	antilogTable[1] = 0
}

// ConvertValueToExponent converts the value n into the exponent
// from the Antilog table in the Galois field of order 256.
func ConvertValueToExponent(n int) int {
	if n < 0 || n > len(antilogTable)-1 {
		return 0
	}

	return antilogTable[n]
}

// ConvertExponentToValue converts the exponent n into the base
// from the Log table in the Galois field of order 256.
func ConvertExponentToValue(n int) int {
	if n < 0 || n > len(logTable)-1 {
		return 0
	}

	return logTable[n]
}

// ConvertIntListToBin converts a list of integers into a list
// of 8 bit binary strings.
func ConvertIntListToBin(list []int) []string {
	result := make([]string, len(list))

	for i, elem := range list {
		bin := strconv.FormatInt(int64(elem), 2)
		result[i] = PadLeft(bin, "0", 8)
	}

	return result
}

// ConvertIntListToCodewords converts a list of integers into
// a binary string of codewords.
func ConvertIntListToCodewords(list []int) string {
	return strings.Join(ConvertIntListToBin(list), "")
}

// GetClosestMultiple computes the closest to n mutiple of m.
func GetClosestMultiple(n int, m int) int {
	multiple := int(math.Round(float64(n) / float64(m)))
	return multiple * m
}

// Max computes the maximum value between two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// PadLeft applies a padding to the left with the character c
// such that the padded string has length n.
func PadLeft(s string, c string, n int) string {
	return strings.Repeat(c, n-len(s)) + s
}

// PadRight applies a padding to the right with the character c
// such that the padded string has length n.
func PadRight(s string, c string, n int) string {
	return s + strings.Repeat(c, n-len(s))
}

// SplitInGroups splits a string into groups of at least n characters.
func SplitInGroups(s string, n int) []string {
	if s == "" {
		return nil
	}

	var result []string
	for {
		result = append(result, s[:n])
		s = s[n:]

		if n > len(s) {
			break
		}
	}

	if len(s) < n && s != "" {
		result = append(result, s)
	}

	return result
}

func GetECMappingKey(version int, lvl string) string {
	return strconv.Itoa(int(version)) + "-" + string(lvl)
}

func testAndrew() {

}
