package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadLeft(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		char     string
		len      int
		expected string
	}{
		{"hello", " ", 8, "   hello"},
		{"world", "*", 5, "world"},
		{"test", "-", 4, "test"},
		{"123", "0", 6, "000123"},
		{"", "*", 3, "***"},
		{"abcdef", "", 8, "abcdef"},
	}

	for _, test := range tests {
		actual := padLeft(test.input, test.char, test.len)
		assert.Equal(test.expected, actual, "Padded values should match")
	}
}

func TestPadRight(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		char     string
		len      int
		expected string
	}{
		{"hello", " ", 8, "hello   "},
		{"world", "*", 5, "world"},
		{"test", "-", 4, "test"},
		{"123", "0", 6, "123000"},
		{"", "*", 3, "***"},
		{"abcdef", "", 8, "abcdef"},
	}

	for _, test := range tests {
		actual := padRight(test.input, test.char, test.len)
		assert.Equal(test.expected, actual, "Padded values should match")
	}
}

func TestClosestMultiple(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input      int
		multipleOf int
		expected   int
	}{
		{10, 5, 10},
		{17, 4, 16},
		{25, 7, 28},
		{30, 10, 30},
		{13, 3, 12},
		{0, 5, 0},
		{100, 0, 0},
		{-15, 5, -15},
	}

	for _, test := range tests {
		actual := getClosestMultiple(test.input, test.multipleOf)
		assert.Equal(test.expected, actual, "Multiples should match")
	}
}

func TestSplitInGroups(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		n        int
		expected []string
	}{
		{"abcdefgh", 3, []string{"abc", "def", "gh"}},
		{"123456789", 2, []string{"12", "34", "56", "78", "9"}},
		{"abcdefgh", 5, []string{"abcde", "fgh"}},
		{"abcd", 2, []string{"ab", "cd"}},
		{"", 3, nil},
		{"abc", 1, []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		actual := splitInGroups(test.input, test.n)
		assert.Equal(test.expected, actual, "Array of groups should match")
	}
}

func TestConvertIntListToBin(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    []int
		expected []string
	}{
		{[]int{1, 2, 3}, []string{"00000001", "00000010", "00000011"}},
		{[]int{10, 20, 30}, []string{"00001010", "00010100", "00011110"}},
		{[]int{255, 0, 127}, []string{"11111111", "00000000", "01111111"}},
		{[]int{}, []string{}},
	}

	for _, test := range tests {
		actual := convertIntListToBin(test.input)
		assert.Equal(test.expected, actual, "Arrays of 8 bit binary representation should match")
	}
}

func TestConvertIntListToCodewords(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    []int
		expected string
	}{
		{[]int{1, 2, 3}, "000000110000001000000001"},
		{[]int{10, 20, 30}, "000111100001010000001010"},
		{[]int{255, 0, 127}, "011111110000000011111111"},
		{[]int{}, ""},
	}

	for _, test := range tests {
		actual := convertIntListToCodewords(test.input)
		assert.Equal(test.expected, actual, "string of codewords should match")
	}
}
