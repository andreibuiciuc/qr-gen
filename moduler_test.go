package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuler(t *testing.T) {
	assert := assert.New(t)
	input := "https://www.qrcode.com/"

	v := newVersioner()
	string, _ := v.getMode(input)
	vrs, _ := v.getVersion(input, string, ec_MEDIUM)

	e := newEncoder()
	encoded, _ := e.encode(input, ec_MEDIUM)
	encoded = e.augmentEncodedInput(encoded, vrs, ec_MEDIUM)

	i := newInterleaver()
	data := i.getFinalMessage(encoded, vrs, ec_MEDIUM)

	m := newModuler(vrs, ec_MEDIUM)
	matrix, penalty := m.createModuleMatrix(data)
	assert.Equal(415, penalty.total, "penalty score should match")

	qi := NewImage()
	qi.CreateImage("best.png", matrix.GetMatrix())
}
