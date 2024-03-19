package qr

import "fmt"

type Qr struct{}

func New(*Qr) *Qr {
	return &Qr{}
}

func (qr *Qr) Encode(s string, lvl rune, filename string) (Image[module], error) {
	versioner := newVersioner()
	mode, err := versioner.getMode(s)

	if err != nil {
		return nil, fmt.Errorf("error")
	}

	version, err := versioner.getVersion(s, mode, lvl)

	if err != nil {
		return nil, fmt.Errorf("error")
	}

	encoder := newEncoder()
	encoded, err := encoder.encode(s, lvl)

	if err != nil {
		return nil, fmt.Errorf("error")
	}

	encoded = encoder.augmentEncodedInput(encoded, version, lvl)

	interleaver := newInterleaver()
	encoded = interleaver.getFinalMessage(encoded, version, lvl)

	moduler := newModuler(version, lvl)
	matrix, _ := moduler.createModuleMatrix(encoded)

	image := NewImage()
	image.CreateImage(filename, matrix.GetMatrix())

	return image, nil
}
