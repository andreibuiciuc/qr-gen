package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string

	input = "1234"
	actual, _ := GetMode(input)
	assert.Equal(NUMERIC, actual, "Input pattern should match as numeric")

	input = "0"
	actual, _ = GetMode(input)
	assert.Equal(NUMERIC, actual, "Input pattern should match as numeric")

	input = "123asd"
	actual, _ = GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be invalid")

	input = " "
	actual, _ = GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be invalid")

	input = ""
	actual, _ = GetMode(input)
	assert.Equal(INVALID_MODE, actual, "Input pattern should be invalid")
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string

	input = "ABC"
	actual, _ := GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = "ABC123"
	actual, _ = GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = "+A-BC..."
	actual, _ = GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = " "
	actual, _ = GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = ""
	actual, _ = GetMode(input)
	assert.Equal(INVALID_MODE, actual, "Input pattern should be invalid")
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string

	input = "ABCabc123"
	actual, _ := GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = "https://www.google.com"
	actual, _ = GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = ""
	actual, _ = GetMode(input)
	assert.Equal(INVALID_MODE, actual, "Input pattern should be invalud")
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	var input string

	input = "HELLO WORLD"
	actual, _ := GetSmallestVersion(input, QUARTILE)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "HELLO THERE WORLD"
	actual, _ = GetSmallestVersion(input, QUARTILE)
	assert.Equal(VERSION_2, actual, "Input should be version 2")

	input = "12345"
	actual, _ = GetSmallestVersion(input, MEDIUM)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "https://www.google.com"
	actual, _ = GetSmallestVersion(input, HIGH)
	assert.Equal(VERSION_3, actual, "Input should be version 3")
}
