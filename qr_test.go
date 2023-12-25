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
	assert.Equal(INVALID, actual, "Input pattern should be invalid")
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
	assert.Equal(INVALID, actual, "Input pattern should be invalid")
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
	assert.Equal(INVALID, actual, "Input pattern should be invalud")
}
