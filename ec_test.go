package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagePolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	e := newEncoder()
	ec := NewErrorCorrector()

	var input string = "HELLO WORLD"
	string, _ := v.getMode(input)
	ver, _ := v.getVersion(input, string, ec_MEDIUM)
	encoded, _ := e.encode(input, ec_MEDIUM)
	augmented := e.augmentEncodedInput(encoded, ver, ec_MEDIUM)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100000100011110110000010001",
		augmented, "Augmented encoded input should match binary representation")

	actual := ec.getMessagePolynomial(augmented)
	var expected QrPolynomial = []int{17, 236, 17, 236, 17, 236, 64, 67, 77, 220, 114, 209, 120, 11, 91, 32}
	assert.Equal(expected, actual, "Message polynomial coefficients should match")
}

func TestGeneratorPolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	ec := NewErrorCorrector()

	actual := ec.getGeneratorPolynomial(1, ec_MEDIUM)
	var expected QrPolynomial = []int{193, 157, 113, 95, 94, 199, 111, 159, 194, 216, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")

	actual = ec.getGeneratorPolynomial(4, ec_QUARTILE)
	expected = []int{94, 43, 77, 146, 144, 70, 68, 135, 42, 233, 117, 209, 40, 145, 24, 206, 56, 77, 152, 199, 98, 136, 4, 183, 51, 246, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")
}

func TestErrorCorrectionCodewordsGenerator(t *testing.T) {
	assert := assert.New(t)
	v := newVersioner()
	e := newEncoder()
	ec := NewErrorCorrector()

	input := "HELLO WORLD"
	string, _ := v.getMode(input)
	ver, _ := v.getVersion(input, string, ec_MEDIUM)
	encoded, _ := e.encode(input, ec_MEDIUM)
	augmented := e.augmentEncodedInput(encoded, ver, ec_MEDIUM)

	actual := ec.getErrorCorrectionCodewords(augmented, 1, ec_MEDIUM)
	var expected QrPolynomial = []int{23, 93, 226, 231, 215, 235, 119, 39, 35, 196}
	assert.Equal(expected, actual, "Error correction codewords should match")
}
