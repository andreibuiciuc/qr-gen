package versioner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := New()

	var input string
	var err error

	input = "1234"
	actual, _ := v.GetMode(input)
	assert.Equal(QrNumericMode, actual, "Input pattern should match as numeric")

	input = "0"
	actual, _ = v.GetMode(input)
	assert.Equal(QrNumericMode, actual, "Input pattern should match as numeric")

	input = "123asd"
	actual, _ = v.GetMode(input)
	assert.Equal(QrByteMode, actual, "Input pattern should be invalid")

	input = " "
	actual, _ = v.GetMode(input)
	assert.Equal(QrAlphanumericMode, actual, "Input pattern should be invalid")

	input = ""
	actual, _ = v.GetMode(input)
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")

	input = "こんにちは"
	actual, err = v.GetMode(input)
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestAlphaNumericInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := New()
	var input string
	var err error

	input = "ABC"
	actual, _ := v.GetMode(input)
	assert.Equal(QrAlphanumericMode, actual, "Input pattern should be alpha numeric")

	input = "ABC123"
	actual, _ = v.GetMode(input)
	assert.Equal(QrAlphanumericMode, actual, "Input pattern should be alpha numeric")

	input = "+A-BC..."
	actual, _ = v.GetMode(input)
	assert.Equal(QrAlphanumericMode, actual, "Input pattern should be alpha numeric")

	input = " "
	actual, _ = v.GetMode(input)
	assert.Equal(QrAlphanumericMode, actual, "Input pattern should be alpha numeric")

	input = ""
	actual, err = v.GetMode(input)
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestByteInputOnly(t *testing.T) {
	assert := assert.New(t)
	v := New()
	var input string
	var err error

	input = "ABCabc123"
	actual, _ := v.GetMode(input)
	assert.Equal(QrByteMode, actual, "Input pattern should be byte")

	input = "https://www.google.com"
	actual, _ = v.GetMode(input)
	assert.Equal(QrByteMode, actual, "Input pattern should be byte")

	input = ""
	actual, err = v.GetMode(input)
	assert.Equal(QrMode(""), actual, "Input pattern should be invalid")
	assert.Error(err)
}

func TestComputeQrVersion(t *testing.T) {
	assert := assert.New(t)
	v := New()
	var input string
	var err error

	input = "HELLO WORLD"
	actual, _ := v.GetVersion(input, QrAlphanumericMode, QrEcQuartile)
	assert.Equal(QrVersion(1), actual, "Input should be version 1")

	input = "HELLO THERE WORLD"
	actual, _ = v.GetVersion(input, QrAlphanumericMode, QrEcQuartile)
	assert.Equal(QrVersion(2), actual, "Input should be version 2")

	input = "12345"
	actual, _ = v.GetVersion(input, QrNumericMode, QrEcMedium)
	assert.Equal(QrVersion(1), actual, "Input should be version 1")

	input = "https://www.google.com"
	actual, _ = v.GetVersion(input, QrByteMode, QrECHigh)
	assert.Equal(QrVersion(3), actual, "Input should be version 3")

	input = "this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version" +
		"this is a very long text fragment for which we cannot compute a compatible qr version"
	actual, err = v.GetVersion(input, QrByteMode, QrEcLow)
	assert.Error(err)
}

func TestGetModeIndicator(t *testing.T) {
	assert := assert.New(t)
	v := New()
	var input string

	input = "1233"
	mode, _ := v.GetMode(input)
	assert.Equal(QrNumericMode, mode, "Input should be numeric")
	actual := v.GetModeIndicator(mode)
	assert.Equal(qrNumericInd, actual, "Input should have a numeric indicator")

	input = "HELLO WORLD"
	mode, _ = v.GetMode(input)
	assert.Equal(QrAlphanumericMode, mode, "Input should be alpha numeric")
	actual = v.GetModeIndicator(mode)
	assert.Equal(qrAlphanumericInd, actual, "Input should have a alpha numeric indicator")

	input = "Hello, world!"
	mode, _ = v.GetMode(input)
	assert.Equal(QrByteMode, mode, "Input should be byte")
	actual = v.GetModeIndicator(mode)
	assert.Equal(qrByteInd, actual, "Input should have a byte indicator")
}

func TestGetCountIndicator(t *testing.T) {
	assert := assert.New(t)
	v := New()
	var input string

	input = "HELLO WORLD"
	mode, _ := v.GetMode(input)
	version, _ := v.GetVersion(input, mode, QrEcQuartile)
	actual, _ := v.GetCountIndicator(input, version, mode)
	assert.Equal("000001011", actual, "Input should match binary representation")

	input = "1234"
	mode, _ = v.GetMode(input)
	version, _ = v.GetVersion(input, mode, QrECHigh)
	actual, _ = v.GetCountIndicator(input, version, mode)
	assert.Equal("0000000100", actual, "Input should match binary representation")

	input = "Hello and welcome!"
	mode, _ = v.GetMode(input)
	version, _ = v.GetVersion(input, mode, QrEcQuartile)
	actual, _ = v.GetCountIndicator(input, version, mode)
	assert.Equal("00010010", actual, "Input should match binary representation")
}
