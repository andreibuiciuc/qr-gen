package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string
	var err error

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
	assert.Equal(QrCodeMode(EMPTY_STRING), actual, "Input pattern should be invalid")

	input = "こんにちは"
	actual, err = GetMode(input)
	assert.Equal(QrCodeMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string
	var err error

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
	actual, err = GetMode(input)
	assert.Equal(QrCodeMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	var input string
	var err error

	input = "ABCabc123"
	actual, _ := GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = "https://www.google.com"
	actual, _ = GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = ""
	actual, err = GetMode(input)
	assert.Equal(QrCodeMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	var input string
	var err error

	input = "HELLO WORLD"
	actual, _ := GetSmallestVersion(input, ALPHA_NUMERIC, QUARTILE)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "HELLO THERE WORLD"
	actual, _ = GetSmallestVersion(input, ALPHA_NUMERIC, QUARTILE)
	assert.Equal(VERSION_2, actual, "Input should be version 2")

	input = "12345"
	actual, _ = GetSmallestVersion(input, NUMERIC, MEDIUM)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "https://www.google.com"
	actual, _ = GetSmallestVersion(input, BYTE, HIGH)
	assert.Equal(VERSION_3, actual, "Input should be version 3")

	input = "this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version"
	actual, err = GetSmallestVersion(input, BYTE, LOW)
	assert.Error(err)
}

func TestGetModeIndicator(t *testing.T) {
	assert := assert.New(t)
	var input string

	input = "1233"
	mode, _ := GetMode(input)
	assert.Equal(NUMERIC, mode, "Input should be numeric")
	actual := GetModeIndicator(mode)
	assert.Equal(NUMERIC_INDICATOR, actual, "Input should have a numeric indicator")

	input = "HELLO WORLD"
	mode, _ = GetMode(input)
	assert.Equal(ALPHA_NUMERIC, mode, "Input should be alpha numeric")
	actual = GetModeIndicator(mode)
	assert.Equal(ALPHA_NUMERIC_INDICATOR, actual, "Input should have a alpha numeric indicator")

	input = "Hello, world!"
	mode, _ = GetMode(input)
	assert.Equal(BYTE, mode, "Input should be byte")
	actual = GetModeIndicator(mode)
	assert.Equal(BYTE_INDICATOR, actual, "Input should have a byte indicator")
}
