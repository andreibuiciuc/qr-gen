package moduler

import (
	"fmt"
	"qr/qr-gen/matrix"
	"qr/qr-gen/versioner"
)

type Module int

type Moduler struct {
	version versioner.QrVersion
}
type Coordintates struct {
	row int
	col int
}

type Boundary struct {
	lowerRow int
	upperRow int
	lowerCol int
	upperCol int
}

const finderPatternSize = 7

const (
	moduleValue_EMPTY Module = 5
)

// Those locations stand only for alignment patterns that do not overlap with finder patterns
// This can be improved by checking the overlap programatically
var allignmentPatternLocation = map[versioner.QrVersion][]Coordintates{
	1: {},
	2: {Coordintates{18, 18}},
	3: {Coordintates{22, 22}},
	4: {Coordintates{26, 26}},
	5: {Coordintates{30, 30}},
}

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
	m.setAlignmentPatterns(moduleMatrix)
	m.setTimingPatterns(moduleMatrix)

	moduleMatrix.PrintMatrix()
	return *moduleMatrix
}

func (m *Moduler) qrCodeSize() int {
	return (int(m.version)-1)*4 + 21
}

// Sets the top left finder pattern in the module matrix
func (m *Moduler) setTopLeftFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary, _ := m.finderPatternBoundary(true, true)
	m.patchPatternInMatrix(moduleMatrix, *boundary)

	for i := boundary.lowerRow; i <= boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.upperCol, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.upperRow, i, 0)
	}
}

// Sets the top right finder pattern in the module matrix
func (m *Moduler) setTopRightFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary, _ := m.finderPatternBoundary(true, false)
	m.patchPatternInMatrix(moduleMatrix, *boundary)

	for i := boundary.lowerRow; i <= boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.lowerCol-1, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.upperRow, i, 0)
	}
}

// Sets the bottom left finder pattern in the module matrix
func (m *Moduler) setBottomLeftFinderPattern(moduleMatrix *matrix.Matrix[Module]) {
	boundary, _ := m.finderPatternBoundary(false, true)
	m.patchPatternInMatrix(moduleMatrix, *boundary)

	for i := boundary.lowerRow - 1; i < boundary.upperRow; i++ {
		moduleMatrix.Set(i, boundary.upperCol, 0)
	}

	for i := boundary.lowerCol; i < boundary.upperCol; i++ {
		moduleMatrix.Set(boundary.lowerRow-1, i, 0)
	}
}

// Sets the alignment patterns in the module matrix
func (m *Moduler) setAlignmentPatterns(moduleMatrix *matrix.Matrix[Module]) {
	for _, coordinates := range allignmentPatternLocation[m.version] {
		boundary := m.alignmentPatternBoundary(coordinates)
		m.patchPatternInMatrix(moduleMatrix, boundary)
	}
}

// Sets the timing patterns in the module matrix
func (m *Moduler) setTimingPatterns(moduleMatrix *matrix.Matrix[Module]) {
	topLeftFinderBoundary, _ := m.finderPatternBoundary(true, true)
	topRightFinderBoundary, _ := m.finderPatternBoundary(true, false)
	bottomLeftFinderBoundary, _ := m.finderPatternBoundary(false, true)

	val := 1
	for i := topLeftFinderBoundary.upperCol - 1; i < topRightFinderBoundary.lowerCol; i++ {
		moduleMatrix.Set(6, i, Module(val))
		val = (val + 1) % 2
	}

	val = 1
	for i := topLeftFinderBoundary.upperRow - 1; i < bottomLeftFinderBoundary.lowerRow; i++ {
		moduleMatrix.Set(i, 6, Module(val))
		val = (val + 1) % 2
	}
}

func (m *Moduler) patchPatternInMatrix(moduleMatrix *matrix.Matrix[Module], boundary Boundary) {
	for i := boundary.lowerRow; i < boundary.upperRow; i++ {
		for j := boundary.lowerCol; j < boundary.upperCol; j++ {
			if m.isPatternModuleDarken(i, j, boundary) {
				moduleMatrix.Set(i, j, 1)
			} else if m.isPatternModuleLighten(i, j, boundary) {
				moduleMatrix.Set(i, j, 0)
			} else {
				moduleMatrix.Set(i, j, 1)
			}
		}
	}
}

func (m *Moduler) isPatternModuleDarken(i, j int, boundary Boundary) bool {
	return i == boundary.lowerRow || i == boundary.upperRow-1 || j == boundary.lowerCol || j == boundary.upperCol-1
}

func (m *Moduler) isPatternModuleLighten(i, j int, boundary Boundary) bool {
	return i == boundary.lowerRow+1 || i == boundary.upperRow-2 || j == boundary.lowerCol+1 || j == boundary.upperCol-2
}

func (m *Moduler) finderPatternBoundary(top, left bool) (*Boundary, error) {
	if top && left {
		return &Boundary{
			lowerRow: 0,
			upperRow: finderPatternSize,
			lowerCol: 0,
			upperCol: finderPatternSize,
		}, nil
	}

	if top && !left {
		return &Boundary{
			lowerRow: 0,
			upperRow: finderPatternSize,
			lowerCol: m.qrCodeSize() - finderPatternSize,
			upperCol: m.qrCodeSize(),
		}, nil
	}

	if !top && left {
		return &Boundary{
			lowerRow: m.qrCodeSize() - finderPatternSize,
			upperRow: m.qrCodeSize(),
			lowerCol: 0,
			upperCol: finderPatternSize,
		}, nil
	}

	return nil, fmt.Errorf("invalid finder pattern location")
}

func (m *Moduler) alignmentPatternBoundary(coordinates Coordintates) Boundary {
	return Boundary{
		lowerRow: coordinates.row - 2,
		upperRow: coordinates.row + 3,
		lowerCol: coordinates.col - 2,
		upperCol: coordinates.col + 3,
	}
}
