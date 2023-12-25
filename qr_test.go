package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()

	var input string
	var err error

	input = "1234"
	actual, _ := e.GetMode(input)
	assert.Equal(NUMERIC, actual, "Input pattern should match as numeric")

	input = "0"
	actual, _ = e.GetMode(input)
	assert.Equal(NUMERIC, actual, "Input pattern should match as numeric")

	input = "123asd"
	actual, _ = e.GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be invalid")

	input = " "
	actual, _ = e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be invalid")

	input = ""
	actual, _ = e.GetMode(input)
	assert.Equal(QrMode(EMPTY_STRING), actual, "Input pattern should be invalid")

	input = "こんにちは"
	actual, err = e.GetMode(input)
	assert.Equal(QrMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string
	var err error

	input = "ABC"
	actual, _ := e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = "ABC123"
	actual, _ = e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = "+A-BC..."
	actual, _ = e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = " "
	actual, _ = e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, actual, "Input pattern should be alpha numeric")

	input = ""
	actual, err = e.GetMode(input)
	assert.Equal(QrMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string
	var err error

	input = "ABCabc123"
	actual, _ := e.GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = "https://www.google.com"
	actual, _ = e.GetMode(input)
	assert.Equal(BYTE, actual, "Input pattern should be byte")

	input = ""
	actual, err = e.GetMode(input)
	assert.Equal(QrMode(EMPTY_STRING), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string
	var err error

	input = "HELLO WORLD"
	actual, _ := e.GetVersion(input, ALPHA_NUMERIC, QUARTILE)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "HELLO THERE WORLD"
	actual, _ = e.GetVersion(input, ALPHA_NUMERIC, QUARTILE)
	assert.Equal(VERSION_2, actual, "Input should be version 2")

	input = "12345"
	actual, _ = e.GetVersion(input, NUMERIC, MEDIUM)
	assert.Equal(VERSION_1, actual, "Input should be version 1")

	input = "https://www.google.com"
	actual, _ = e.GetVersion(input, BYTE, HIGH)
	assert.Equal(VERSION_3, actual, "Input should be version 3")

	input = "this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version"
	actual, err = e.GetVersion(input, BYTE, LOW)
	assert.Error(err)
}

func TestGetModeIndicator(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "1233"
	mode, _ := e.GetMode(input)
	assert.Equal(NUMERIC, mode, "Input should be numeric")
	actual := e.GetModeIndicator(mode)
	assert.Equal(NUMERIC_INDICATOR, actual, "Input should have a numeric indicator")

	input = "HELLO WORLD"
	mode, _ = e.GetMode(input)
	assert.Equal(ALPHA_NUMERIC, mode, "Input should be alpha numeric")
	actual = e.GetModeIndicator(mode)
	assert.Equal(ALPHA_NUMERIC_INDICATOR, actual, "Input should have a alpha numeric indicator")

	input = "Hello, world!"
	mode, _ = e.GetMode(input)
	assert.Equal(BYTE, mode, "Input should be byte")
	actual = e.GetModeIndicator(mode)
	assert.Equal(BYTE_INDICATOR, actual, "Input should have a byte indicator")
}
