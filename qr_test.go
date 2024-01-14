package qr

import (
	"qr/qr-gen/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoderTest()

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
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")

	input = "こんにちは"
	actual, err = e.GetMode(input)
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoderTest()
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
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoderTest()
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
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
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
	e := NewEncoderTest()
	var input string

	input = "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, MEDIUM)
	encoded, _ := e.Encode(input, MEDIUM)
	augmented := e.AugmentEncodedInput(encoded, version, MEDIUM)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100000100011110110000010001",
		augmented, "Augmented encoded input should match binary representation")

	actual := e.GetMessagePolynomial(augmented)
	var expected QrPolynomial = []int{17, 236, 17, 236, 17, 236, 64, 67, 77, 220, 114, 209, 120, 11, 91, 32}
	assert.Equal(expected, actual, "Message polynomial coefficients should match")
}

func TestGeneratorPolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	util.ComputeLogAntilogTables()
	e := NewEncoderTest()

	actual := e.GetGeneratorPolynomial(VERSION_1, MEDIUM)
	var expected QrPolynomial = []int{193, 157, 113, 95, 94, 199, 111, 159, 194, 216, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")

	actual = e.GetGeneratorPolynomial(VERSION_4, QUARTILE)
	expected = []int{94, 43, 77, 146, 144, 70, 68, 135, 42, 233, 117, 209, 40, 145, 24, 206, 56, 77, 152, 199, 98, 136, 4, 183, 51, 246, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")
}

func TestErrorCorrectionCodewordsGenerator(t *testing.T) {
	assert := assert.New(t)
	util.ComputeLogAntilogTables()
	e := NewEncoderTest()

	input := "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, MEDIUM)
	encoded, _ := e.Encode(input, MEDIUM)
	augmented := e.AugmentEncodedInput(encoded, version, MEDIUM)

	actual := e.GetErrorCorrectionCodewords(augmented, VERSION_1, MEDIUM)
	var expected QrPolynomial = []int{23, 93, 226, 231, 215, 235, 119, 39, 35, 196}
	assert.Equal(expected, actual, "Error correction codewords should match")
}

func TestInterleavingProcess(t *testing.T) {
	assert := assert.New(t)
	util.ComputeLogAntilogTables()
	e := NewEncoderTest()

	input := "0100001101010101010001101000011001010111001001100101010111000010011101110011001000000110000100100000011001100111001001101111011011110110010000100000011101110110100001101111001000000111001001100101011000010110110001101100011110010010000001101011011011100110111101110111011100110010000001110111011010000110010101110010011001010010000001101000011010010111001100100000011101000110111101110111011001010110110000100000011010010111001100101110000011101100000100011110110000010001111011000001000111101100"
	actual := e.GetFinalMessage(input, VERSION_5, QUARTILE)
	expected := "01000011111101101011011001000110010101011111011011100110111101110100011001000010111101110111011010000110000001110111011101010110010101110111011000110010110000100010011010000110000001110000011001010101111100100111011010010111110000100000011110000110001100100111011100100110010101111110000000110010010101100010011011101100000001100001011001010010000100010001001011000110000001101110110000000110110001111000011000010001011001111001001010010111111011000010011000000110001100100001000100000111111011000010011110000101100011010100101001101111011110001100000001011001101000011010001111110000110011010101011010100011011011000000011001101111000100010000101001010011100110101101000110111101110001010111010110000001111001101110101110011010000110111100001101101111111110001000011001001100011100011110010111001000111011101111110111011111100111011111001001101000111100010111110001001011001001011111011110110110100001011000001101110011110010100100110001101100001011010011110011010100111101110000101101100000101100011111101011000111110011000111010001100100110101010101011110010100100011000000000"
	assert.Equal(expected, actual, "Interleaved codewords should match")
}
