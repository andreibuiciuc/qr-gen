package util

import (
	"strconv"
	"strings"
)

var logTable = make([]int, 256)
var antilogTable = make([]int, 256)

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
