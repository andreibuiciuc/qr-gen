package moduler

import (
	"qr/qr-gen/matrix"
	"qr/qr-gen/versioner"
)

type Module int

type Moduler struct {
	version versioner.QrVersion
}

type Boundary struct {
	lowerRow int
	upperRow int
	lowerCol int
	upperCol int
}

const (
	moduleValue_EMPTY Module = 5
)

const finderPatternSize = 7

func NewModuler(version int) *Moduler {
	return &Moduler{
		version: versioner.QrVersion(version),
	}
}

func (m *Moduler) CreateModuleMatrix() matrix.Matrix[Module] {
	qrCodeSize := m.qrCodeSize()

	moduleMatrix := matrix.NewMatrix[Module](qrCodeSize, qrCodeSize)
	moduleMatrix.Init(moduleValue_EMPTY)

	m.setTopLeftFinderPattern(moduleMatrix)
	m.setTopRightFinderPattern(moduleMatrix)
	m.setBottomLeftFinderPattern(moduleMatrix)

	moduleMatrix.PrintMatrix()
	return *moduleMatrix
}

func (m *Moduler) qrCodeSize() int {
	return (int(m.version)-1)*4 + 21
}

func (m *Moduler) setTopLeftFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary := Boundary{
		lowerRow: 0,
		upperRow: finderPatternSize,
		lowerCol: 0,
		upperCol: finderPatternSize,
	}

	m.patchFinderPattern(moduleMatrix, boundary)

	for i := boundary.lowerRow; i <= boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.upperCol, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.upperRow, i, 0)
	}
}

func (m *Moduler) setTopRightFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary := Boundary{
		lowerRow: 0,
		upperRow: finderPatternSize,
		lowerCol: m.qrCodeSize() - finderPatternSize,
		upperCol: m.qrCodeSize(),
	}
	m.patchFinderPattern(moduleMatrix, boundary)

	for i := boundary.lowerRow; i <= boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.lowerCol-1, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.upperRow, i, 0)
	}
}

func (m *Moduler) setBottomLeftFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary := Boundary{
		lowerRow: m.qrCodeSize() - finderPatternSize,
		upperRow: m.qrCodeSize(),
		lowerCol: 0,
		upperCol: finderPatternSize,
	}
	m.patchFinderPattern(moduleMatrix, boundary)

	for i := boundary.lowerRow - 1; i < boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.upperCol, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.lowerRow-1, i, 0)
	}
}

func (m *Moduler) patchFinderPattern(moduleMatrix *matrix.Matrix[Module], boundary Boundary) {
	for i := boundary.lowerRow; i < boundary.upperRow; i++ {
		for j := boundary.lowerCol; j < boundary.upperCol; j++ {
			if m.isFinderModuleDarken(i, j, boundary) {
				moduleMatrix.Set(i, j, 1)
			} else if m.isFinderModuleLighten(i, j, boundary) {
				moduleMatrix.Set(i, j, 0)
			} else {
				moduleMatrix.Set(i, j, 1)
			}
		}
	}
}

func (m *Moduler) isFinderModuleDarken(i, j int, boundary Boundary) bool {
	return i == boundary.lowerRow || i == boundary.upperRow-1 || j == boundary.lowerCol || j == boundary.upperCol-1
}

func (m *Moduler) isFinderModuleLighten(i, j int, boundary Boundary) bool {
	return i == boundary.lowerRow+1 || i == boundary.upperRow-2 || j == boundary.lowerCol+1 || j == boundary.upperCol-2
}
