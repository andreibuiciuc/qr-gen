package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericEncoding(t *testing.T) {
	assert := assert.New(t)
	e := newEncoder()
	var input string

	input = "8675309"
	actual := e.encodeNumericInput(input)
	assert.Equal("110110001110000100101001", actual, "Input should match binary representation")

	input = "1234"
	actual = e.encodeNumericInput(input)
	assert.Equal("00011110110100", actual, "Input should match binary representation")
}

func TestAlphaNumericEncoding(t *testing.T) {
	assert := assert.New(t)
	e := newEncoder()
	var input string

	input = "HE"
	actual := e.encodeAlphanumericInput(input)
	assert.Equal("01100001011", actual, "Input should match binary representation")

	input = "HED"
	actual = e.encodeAlphanumericInput(input)
	assert.Equal("01100001011001101", actual, "Input should match binary representation")

	input = "HELLO WORLD"
	actual = e.encodeAlphanumericInput(input)
	assert.Equal("0110000101101111000110100010111001011011100010011010100001101", actual, "Input should mathc binary representation")
}

func TestByteNumericInput(t *testing.T) {
	assert := assert.New(t)
	e := newEncoder()
	var input string

	input = "Hello"
	actual := e.encodeByteInput(input)
	assert.Equal("0100100001100101011011000110110001101111",
		actual, "Input should match binary representation")

	input = "Hello, world!"
	actual = e.encodeByteInput(input)
	assert.Equal("01001000011001010110110001101100011011110010110000100000011101110110111101110010011011000110010000100001",
		actual, "Input should match binary representation")
}

func TestEncodedInputAugmentation(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	e := newEncoder()
	var input = "HELLO WORLD"
	string, _ := v.getMode(input)
	int, _ := v.getVersion(input, string, ec_QUARTILE)
	encoded, _ := e.encode(input, ec_QUARTILE)
	actual := e.augmentEncodedInput(encoded, int, ec_QUARTILE)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100",
		actual, "Augmented encoded input should match binary representation")
}
