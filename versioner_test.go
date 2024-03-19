package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()

	var input string
	var err error

	input = "1234"
	actual, _ := v.getMode(input)
	assert.Equal(mode_NUMERIC, actual, "Input pattern should match as numeric")

	input = "0"
	actual, _ = v.getMode(input)
	assert.Equal(mode_NUMERIC, actual, "Input pattern should match as numeric")

	input = "123asd"
	actual, _ = v.getMode(input)
	assert.Equal(mode_BYTE, actual, "Input pattern should be invalid")

	input = " "
	actual, _ = v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, actual, "Input pattern should be invalid")

	input = ""
	actual, _ = v.getMode(input)
	assert.Equal(string(""), actual, "Input pattern should be invalid")

	input = "こんにちは"
	actual, err = v.getMode(input)
	assert.Equal(string(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	var input string
	var err error

	input = "ABC"
	actual, _ := v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, actual, "Input pattern should be alpha numeric")

	input = "ABC123"
	actual, _ = v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, actual, "Input pattern should be alpha numeric")

	input = "+A-BC..."
	actual, _ = v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, actual, "Input pattern should be alpha numeric")

	input = " "
	actual, _ = v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, actual, "Input pattern should be alpha numeric")

	input = ""
	actual, err = v.getMode(input)
	assert.Equal(string(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	var input string
	var err error

	input = "ABCabc123"
	actual, _ := v.getMode(input)
	assert.Equal(mode_BYTE, actual, "Input pattern should be byte")

	input = "https://www.google.com"
	actual, _ = v.getMode(input)
	assert.Equal(mode_BYTE, actual, "Input pattern should be byte")

	input = ""
	actual, err = v.getMode(input)
	assert.Equal(string(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	var input string
	var err error

	input = "HELLO WORLD"
	actual, _ := v.getVersion(input, mode_ALPHANUMERIC, ec_QUARTILE)
	assert.Equal(int(1), actual, "Input should be int 1")

	input = "HELLO THERE WORLD"
	actual, _ = v.getVersion(input, mode_ALPHANUMERIC, ec_QUARTILE)
	assert.Equal(int(2), actual, "Input should be int 2")

	input = "12345"
	actual, _ = v.getVersion(input, mode_NUMERIC, ec_MEDIUM)
	assert.Equal(int(1), actual, "Input should be int 1")

	input = "https://www.google.com"
	actual, _ = v.getVersion(input, mode_BYTE, ec_HIGH)
	assert.Equal(int(3), actual, "Input should be int 3")

	input = "this is a very long text fragment for which we cannot compute a compatible qr int" +
		"this is a very long text fragment for which we cannot compute a compatible qr int" +
		"this is a very long text fragment for which we cannot compute a compatible qr int"
	_, err = v.getVersion(input, mode_BYTE, ec_LOW)
	assert.Error(err)
}

func TestGetModeIndicator(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	var input string

	input = "1233"
	string, _ := v.getMode(input)
	assert.Equal(mode_NUMERIC, string, "Input should be numeric")
	actual := v.getModeIndicator(string)
	assert.Equal(indicator_NUMERIC, actual, "Input should have a numeric indicator")

	input = "HELLO WORLD"
	string, _ = v.getMode(input)
	assert.Equal(mode_ALPHANUMERIC, string, "Input should be alpha numeric")
	actual = v.getModeIndicator(string)
	assert.Equal(indicator_ALPHANUMERIC, actual, "Input should have a alpha numeric indicator")

	input = "Hello, world!"
	string, _ = v.getMode(input)
	assert.Equal(mode_BYTE, string, "Input should be byte")
	actual = v.getModeIndicator(string)
	assert.Equal(indicator_BYTE, actual, "Input should have a byte indicator")
}

func TestGetCountIndicator(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	var input string

	input = "HELLO WORLD"
	string, _ := v.getMode(input)
	int, _ := v.getVersion(input, string, ec_QUARTILE)
	actual, _ := v.getCountIndicator(input, int, string)
	assert.Equal("000001011", actual, "Input should match binary representation")

	input = "1234"
	string, _ = v.getMode(input)
	int, _ = v.getVersion(input, string, ec_HIGH)
	actual, _ = v.getCountIndicator(input, int, string)
	assert.Equal("0000000100", actual, "Input should match binary representation")

	input = "Hello and welcome!"
	string, _ = v.getMode(input)
	int, _ = v.getVersion(input, string, ec_QUARTILE)
	actual, _ = v.getCountIndicator(input, int, string)
	assert.Equal("00010010", actual, "Input should match binary representation")
}
