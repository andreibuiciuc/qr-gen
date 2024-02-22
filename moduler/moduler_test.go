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
	assert.Equal(version, versioner.QrVersion(2))

	e := encoder.New()
	encoded, _ := e.Encode(input, versioner.QrEcMedium)
	encoded = e.AugmentEncodedInput(encoded, version, versioner.QrEcMedium)

	i := interleaver.New()
	data := i.GetFinalMessage(encoded, version, versioner.QrEcMedium)

	m := NewModuler(version)
	matrix := m.CreateModuleMatrix(data)

	qi := img.New()
	qi.CreateImage("test2.png", matrix.GetMatrix())
}
