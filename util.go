package qr

import (
	"math"
	"strconv"
	"strings"
)

type module int
type qrEcInfo struct {
	TotalDataCodewords          int
	ECCodewordsPerBlock         int
	NumBlocksGroup1             int
	DataCodeworkdsInGroup1Block int
	NumBlocksGroup2             int
	DataCodewordsInGroup2Block  int
}

const codewordSize = 8
const (
	module_LIGHTEN module = iota
	module_DARKEN
	module_FINDER_LIGHTEN
	module_FINDER_DARKEN
	module_SEPARATOR
	module_ALIGNMENT_LIGHTEN
	module_ALIGNMENT_DARKEN
	module_TIMING_LIGHTEN
	module_TIMING_DARKEN
	module_DARK
	module_RESERVED
	module_EMPTY
)

var logTable = make([]int, 256)
var antilogTable = make([]int, 256)

var ecInfo = map[string]qrEcInfo{
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

var fmtInfoCodes = map[rune]map[int]string{
	'L': {
		0: "111011111000100",
		1: "111001011110011",
		2: "111110110101010",
		3: "111100010011101",
		4: "110011000101111",
		5: "110001100011000",
		6: "110110001000001",
		7: "110100101110110",
	},
	'M': {
		0: "101010000010010",
		1: "101000100100101",
		2: "101111001111100",
		3: "101101101001011",
		4: "100010111111001",
		5: "100000011001110",
		6: "100111110010111",
		7: "100101010100000",
	},
	'Q': {
		0: "011010101011111",
		1: "011000001101000",
		2: "011111100110001",
		3: "011101000000110",
		4: "010010010110100",
		5: "010000110000011",
		6: "010111011011010",
		7: "010101111101101",
	},
	'H': {
		0: "001011010001001",
		1: "001001110111110",
		2: "001110011100111",
		3: "001100111010000",
		4: "000011101100010",
		5: "000001001010101",
		6: "000110100001100",
		7: "000100000111011",
	},
}

func isModuleLighten(module module) bool {
	return module == module_LIGHTEN || module == module_FINDER_LIGHTEN ||
		module == module_ALIGNMENT_LIGHTEN || module == module_TIMING_LIGHTEN ||
		module == module_SEPARATOR
}

func isModuleSkippedForFormat(module module) bool {
	return module == module_TIMING_LIGHTEN || module == module_TIMING_DARKEN || module == module_DARK
}

func getDataModule(value int) module {
	if value == 0 {
		return module_LIGHTEN
	}
	return module_DARKEN
}

// computeAlphaToPower computes the power of a to p
// in the Galois field of order 256.
func computeAlphaToPower(a, p int) int {
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

// computeLogAntilogTables computes the Log and Antilog tables
// in the Galois field of order 256.
func computeLogAntilogTables() {
	for i := 0; i < 256; i++ {
		logTable[i] = computeAlphaToPower(2, i)
		antilogTable[logTable[i]] = i
	}
	antilogTable[1] = 0
}

// convertValueToExponent converts the value n into the exponent
// from the Antilog table in the Galois field of order 256.
func convertValueToExponent(n int) int {
	if n < 0 || n > len(antilogTable)-1 {
		return 0
	}

	return antilogTable[n]
}

// convertExponentToValue converts the exponent n into the base
// from the Log table in the Galois field of order 256.
func convertExponentToValue(n int) int {
	if n < 0 || n > len(logTable)-1 {
		return 0
	}

	return logTable[n]
}

// convertIntListToBin converts a list of integers into a list
// of 8 bit binary strings.
func convertIntListToBin(list []int) []string {
	result := make([]string, len(list))

	for i, elem := range list {
		bin := strconv.FormatInt(int64(elem), 2)
		result[i] = padLeft(bin, "0", 8)
	}

	return result
}

// convertIntListToCodewords converts a list of integers into
// a binary string of codewords.
func convertIntListToCodewords(list []int) string {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return strings.Join(convertIntListToBin(list), "")
}

// getClosestMultiple computes the closest to n mutiple of m.
func getClosestMultiple(n int, m int) int {
	multiple := int(math.Round(float64(n) / float64(m)))
	return multiple * m
}

// max computes the maximum value between two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// padLeft applies a padding to the left with the character c
// such that the padded string has length n.
func padLeft(s string, c string, n int) string {
	return strings.Repeat(c, n-len(s)) + s
}

// padRight applies a padding to the right with the character c
// such that the padded string has length n.
func padRight(s string, c string, n int) string {
	return s + strings.Repeat(c, n-len(s))
}

// splitInGroups splits a string into groups of at least n characters.
func splitInGroups(s string, n int) []string {
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

func getECMappingKey(v int, lvl string) string {
	return strconv.Itoa(v) + "-" + lvl
}
