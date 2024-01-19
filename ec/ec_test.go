package ec

import (
	"qr/qr-gen/encoder"
	"qr/qr-gen/versioner"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagePolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	v := versioner.New()
	e := encoder.New()
	ec := New()

	var input string
	input = "HELLO WORLD"
	mode, _ := v.GetMode(input)
	version, _ := v.GetVersion(input, mode, versioner.QrEcMedium)
	encoded, _ := e.Encode(input, versioner.QrEcMedium)
	augmented := e.AugmentEncodedInput(encoded, version, versioner.QrEcMedium)
	assert.Equal("00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100000100011110110000010001",
		augmented, "Augmented encoded input should match binary representation")

	actual := ec.GetMessagePolynomial(augmented)
	var expected QrPolynomial = []int{17, 236, 17, 236, 17, 236, 64, 67, 77, 220, 114, 209, 120, 11, 91, 32}
	assert.Equal(expected, actual, "Message polynomial coefficients should match")
}

func TestGeneratorPolynomialGeneration(t *testing.T) {
	assert := assert.New(t)
	ec := New()

	actual := ec.GetGeneratorPolynomial(1, versioner.QrEcMedium)
	var expected QrPolynomial = []int{193, 157, 113, 95, 94, 199, 111, 159, 194, 216, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")

	actual = ec.GetGeneratorPolynomial(4, versioner.QrEcQuartile)
	expected = []int{94, 43, 77, 146, 144, 70, 68, 135, 42, 233, 117, 209, 40, 145, 24, 206, 56, 77, 152, 199, 98, 136, 4, 183, 51, 246, 1}
	assert.Equal(expected, actual, "Generator polynomial coefficients should match")
}

func TestErrorCorrectionCodewordsGenerator(t *testing.T) {
	assert := assert.New(t)
	v := versioner.New()
	e := encoder.New()
	ec := New()

	input := "HELLO WORLD"
	mode, _ := v.GetMode(input)
	version, _ := v.GetVersion(input, mode, versioner.QrEcMedium)
	encoded, _ := e.Encode(input, versioner.QrEcMedium)
	augmented := e.AugmentEncodedInput(encoded, version, versioner.QrEcMedium)

	actual := ec.GetErrorCorrectionCodewords(augmented, 1, versioner.QrEcMedium)
	var expected QrPolynomial = []int{23, 93, 226, 231, 215, 235, 119, 39, 35, 196}
	assert.Equal(expected, actual, "Error correction codewords should match")
}
