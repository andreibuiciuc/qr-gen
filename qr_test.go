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

func TestGetCountIndicator(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, QUARTILE)
	actual, _ := e.GetCountIndicator(input, version, mode)
	assert.Equal("000001011", actual, "Input should match binary representation")

	input = "1234"
	mode, _ = e.GetMode(input)
	version, _ = e.GetVersion(input, mode, HIGH)
	actual, _ = e.GetCountIndicator(input, version, mode)
	assert.Equal("0000000100", actual, "Input should match binary representation")

	input = "Hello and welcome!"
	mode, _ = e.GetMode(input)
	version, _ = e.GetVersion(input, mode, QUARTILE)
	actual, _ = e.GetCountIndicator(input, version, mode)
	assert.Equal("00010010", actual, "Input should match binary representation")
}

func TestNumericEncoding(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "8675309"
	actual := e.EncodeNumericInput(input)
	assert.Equal("110110001110000100101001", actual, "Input should match binary representation")

	input = "1234"
	actual = e.EncodeNumericInput(input)
	assert.Equal("00011110110100", actual, "Input should match binary representation")
}

func TestAlphaNumericEncoding(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "HE"
	actual := e.EncodeAlphanumericInput(input)
	assert.Equal("01100001011", actual, "Input should match binary representation")

	input = "HED"
	actual = e.EncodeAlphanumericInput(input)
	assert.Equal("01100001011001101", actual, "Input should match binary representation")

	input = "HELLO WORLD"
	actual = e.EncodeAlphanumericInput(input)
	assert.Equal("0110000101101111000110100010111001011011100010011010100001101", actual, "Input should mathc binary representation")
}

func TestByteNumericInput(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "Hello"
	actual := e.EncodeByteInput(input)
	assert.Equal("0100100001100101011011000110110001101111",
		actual, "Input should match binary representation")

	input = "Hello, world!"
	actual = e.EncodeByteInput(input)
	assert.Equal("01001000011001010110110001101100011011110010110000100000011101110110111101110010011011000110010000100001",
		actual, "Input should match binary representation")
}

func TestEncodedInputAugmentation(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, QUARTILE)
	encoded, _ := e.Encode(input, QUARTILE)
	actual := e.AugmentEncodedInput(encoded, version, QUARTILE)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100",
		actual, "Augmented encoded input should match binary representation")
}

func TestMessagePolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoder()
	var input string

	input = "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, MEDIUM)
	encoded, _ := e.Encode(input, MEDIUM)
	augmented := e.AugmentEncodedInput(encoded, version, MEDIUM)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100000100011110110000010001",
		augmented, "Augmented encoded input should match binary representation")

	actual := e.GetMessagePolynomial(augmented)
	var expected QrPolynomial = []int{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236, 17, 236, 17}
	assert.Equal(expected, actual, "Message polynomial coefficients should match")
}

func TestGeneratorPolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	computeLogAntilogTables()
	e := NewEncoder()

	actual := e.GetGeneratorPolynomial(VERSION_1, MEDIUM)
	var expected QrPolynomial = []int{45, 32, 94, 64, 70, 118, 61, 46, 67, 251, 0}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")

	actual = e.GetGeneratorPolynomial(VERSION_4, QUARTILE)
	expected = []int{70, 218, 145, 153, 227, 48, 102, 13, 142, 245, 21, 161, 53, 165, 28, 111, 201, 145, 17, 118, 182, 103, 2, 158, 125, 173, 0}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")
}
