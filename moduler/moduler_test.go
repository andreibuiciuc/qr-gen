package moduler

import (
	"qr/qr-gen/encoder"
	"qr/qr-gen/img"
	"qr/qr-gen/interleaver"
	"qr/qr-gen/versioner"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuler(t *testing.T) {
	assert := assert.New(t)
	input := "https://www.qrcode.com/"

	v := versioner.New()
	mode, _ := v.GetMode(input)
	version, _ := v.GetVersion(input, mode, versioner.QrEcMedium)

	e := encoder.New()
	encoded, _ := e.Encode(input, versioner.QrEcMedium)
	encoded = e.AugmentEncodedInput(encoded, version, versioner.QrEcMedium)

	i := interleaver.New()
	data := i.GetFinalMessage(encoded, version, versioner.QrEcMedium)

	m := New(version, versioner.QrEcMedium)
	matrix, penalty := m.CreateModuleMatrix(data)
	assert.Equal(415, penalty.total, "penalty score should match")

	qi := img.New()
	qi.CreateImage("best.png", matrix.GetMatrix())
}
