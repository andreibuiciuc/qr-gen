package moduler

import (
	"fmt"
	"qr/qr-gen/matrix"
	"qr/qr-gen/util"
	"qr/qr-gen/versioner"
	"strconv"
)

type ModulerInterface interface {
	CreateModuleMatrix(data string) *matrix.Matrix[util.Module]
}

type Moduler struct {
	version      versioner.QrVersion
	moduleMatrix *matrix.Matrix[util.Module]
}

type Coordintates struct {
	row int
	col int
}

type Boundary struct {
	lower Coordintates
	upper Coordintates
}

const finderPatternSize = 7

// These locations stand only for alignment patterns that do not overlap with finder patterns
// This can be improved by checking the overlap programatically
var allignmentPatternLocation = map[versioner.QrVersion][]Coordintates{
	1: {},
	2: {Coordintates{18, 18}},
	3: {Coordintates{22, 22}},
	4: {Coordintates{26, 26}},
	5: {Coordintates{30, 30}},
}

func NewModuler(version versioner.QrVersion) ModulerInterface {
	return &Moduler{
		version: version,
	}
}

func (m *Moduler) CreateModuleMatrix(data string) *matrix.Matrix[util.Module] {
	qrCodeSize := m.qrCodeSize()

	m.moduleMatrix = matrix.NewMatrix[util.Module](qrCodeSize, qrCodeSize)
	m.moduleMatrix.Init(util.Module_EMPTY)

	m.setTopLeftFinderPattern()
	m.setTopRightFinderPattern()
	m.setBottomLeftFinderPattern()
	m.setAlignmentPatterns()
	m.setTimingPatterns()
	m.setDarkModule()

	m.reserveFormatArea()

	m.placeDataBits(data)
	m.moduleMatrix.PrintMatrix()

	return m.moduleMatrix
}

func (m *Moduler) qrCodeSize() int {
	return (int(m.version)-1)*4 + 21
}

// MODULE PLACEMENT

// Sets the top left finder pattern in the module matrix
func (m *Moduler) setTopLeftFinderPattern() {
	boundary, _ := m.finderPatternBoundary(true, true)
	m.patchPattern(*boundary, util.Module_FINDER_LIGHTEN, util.Module_FINDER_DARKEN)

	for i := boundary.lower.row; i <= boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.upper.col, util.Module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.upper.row, i, util.Module_SEPARATOR)
	}
}

// Sets the top right finder pattern in the module matrix
func (m *Moduler) setTopRightFinderPattern() {
	boundary, _ := m.finderPatternBoundary(true, false)
	m.patchPattern(*boundary, util.Module_FINDER_LIGHTEN, util.Module_FINDER_DARKEN)

	for i := boundary.lower.row; i <= boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.lower.col-1, util.Module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.upper.row, i, util.Module_SEPARATOR)
	}
}

// Sets the bottom left finder pattern in the module matrix
func (m *Moduler) setBottomLeftFinderPattern() {
	boundary, _ := m.finderPatternBoundary(false, true)
	m.patchPattern(*boundary, util.Module_FINDER_LIGHTEN, util.Module_FINDER_DARKEN)

	for i := boundary.lower.row - 1; i < boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.upper.col, util.Module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.lower.row-1, i, util.Module_SEPARATOR)
	}
}

// Sets the alignment patterns in the module matrix
func (m *Moduler) setAlignmentPatterns() {
	for _, coordinates := range allignmentPatternLocation[m.version] {
		boundary := m.alignmentPatternBoundary(coordinates)
		m.patchPattern(boundary, util.Module_ALIGNMENT_LIGHTEN, util.Module_ALIGNMENT_DARKEN)
	}
}

// Sets the timing patterns in the module matrix
func (m *Moduler) setTimingPatterns() {
	topLeftFinderBoundary, _ := m.finderPatternBoundary(true, true)
	topRightFinderBoundary, _ := m.finderPatternBoundary(true, false)
	bottomLeftFinderBoundary, _ := m.finderPatternBoundary(false, true)

	val := 1
	for i := topLeftFinderBoundary.upper.col - 1; i < topRightFinderBoundary.lower.col; i++ {
		if val == 0 {
			m.moduleMatrix.Set(6, i, util.Module_TIMING_LIGHTEN)
		} else {
			m.moduleMatrix.Set(6, i, util.Module_TIMING_DARKEN)
		}
		val = (val + 1) % 2
	}

	val = 1
	for i := topLeftFinderBoundary.upper.row - 1; i < bottomLeftFinderBoundary.lower.row; i++ {
		if val == 0 {
			m.moduleMatrix.Set(i, 6, util.Module_TIMING_LIGHTEN)
		} else {
			m.moduleMatrix.Set(i, 6, util.Module_TIMING_DARKEN)
		}
		val = (val + 1) % 2
	}
}

// Sets the dark module in the module matrix
func (m *Moduler) setDarkModule() {
	m.moduleMatrix.Set(4*int(m.version)+9, 8, util.Module_DARK)
}

// Sets the reserved format information area in the module matrix
func (m *Moduler) reserveFormatArea() {
	boundary, _ := m.finderPatternBoundary(true, true)

	for i := boundary.lower.row; i < boundary.upper.row+2; i++ {
		if val, _ := m.moduleMatrix.At(i, boundary.upper.col+1); val == util.Module_EMPTY {
			m.moduleMatrix.Set(i, boundary.upper.col+1, util.Module_RESERVED)
		}
	}

	for i := boundary.lower.col; i < boundary.upper.col+2; i++ {
		if val, _ := m.moduleMatrix.At(boundary.upper.row+1, i); val == util.Module_EMPTY {
			m.moduleMatrix.Set(boundary.upper.row+1, i, util.Module_RESERVED)
		}
	}

	boundary, _ = m.finderPatternBoundary(true, false)

	for i := boundary.lower.col - 1; i < boundary.upper.col; i++ {
		if val, _ := m.moduleMatrix.At(boundary.upper.row+1, i); val == util.Module_EMPTY {
			m.moduleMatrix.Set(boundary.upper.row+1, i, util.Module_RESERVED)
		}
	}

	boundary, _ = m.finderPatternBoundary(false, true)

	for i := boundary.lower.row - 1; i < boundary.upper.row; i++ {
		if val, _ := m.moduleMatrix.At(i, boundary.upper.col+1); val == util.Module_EMPTY {
			m.moduleMatrix.Set(i, boundary.upper.col+1, util.Module_RESERVED)
		}
	}
}

// Places the encoded data bits in the module matrix
func (m *Moduler) placeDataBits(data string) {
	currentCellCoord := Coordintates{row: m.qrCodeSize() - 1, col: m.qrCodeSize() - 1}

	indexInBits := 0
	indexInModules := 0
	isUpwardMovement := -1
	module := util.Module_EMPTY

	for currentCellCoord.col >= 0 {

		if val, _ := m.moduleMatrix.At(currentCellCoord.row, currentCellCoord.col); val == util.Module_EMPTY {
			if bit, _ := strconv.ParseInt(string(data[indexInBits]), 2, 64); bit == 0 {
				module = util.Module_LIGHTEN
			} else {
				module = util.Module_DARKEN
			}
			m.moduleMatrix.Set(currentCellCoord.row, currentCellCoord.col, util.Module(module))
			indexInBits += 1
		}

		if indexInModules%2 == 1 {
			currentCellCoord.row += isUpwardMovement

			if currentCellCoord.row == -1 || currentCellCoord.row == m.qrCodeSize() {
				isUpwardMovement = -isUpwardMovement
				currentCellCoord.row += isUpwardMovement
				if currentCellCoord.col == 7 {
					currentCellCoord.col -= 2
				} else {
					currentCellCoord.col -= 1
				}
			} else {
				currentCellCoord.col += 1
			}
		} else {
			currentCellCoord.col -= 1
		}

		indexInModules += 1
	}
}

// Patches a squared shape pattern of modules alternatively (finder, alignment)
func (m *Moduler) patchPattern(boundary Boundary, ligthenModule util.Module, darkenModule util.Module) {
	for i := boundary.lower.row; i < boundary.upper.row; i++ {
		for j := boundary.lower.col; j < boundary.upper.col; j++ {
			if m.isPatternModuleDarken(i, j, boundary) {
				m.moduleMatrix.Set(i, j, darkenModule)
			} else if m.isPatternModuleLighten(i, j, boundary) {
				m.moduleMatrix.Set(i, j, ligthenModule)
			} else {
				m.moduleMatrix.Set(i, j, darkenModule)
			}
		}
	}
}

func (m *Moduler) isPatternModuleDarken(i, j int, boundary Boundary) bool {
	return i == boundary.lower.row || i == boundary.upper.row-1 || j == boundary.lower.col || j == boundary.upper.col-1
}

func (m *Moduler) isPatternModuleLighten(i, j int, boundary Boundary) bool {
	return i == boundary.lower.row+1 || i == boundary.upper.row-2 || j == boundary.lower.col+1 || j == boundary.upper.col-2
}

func (m *Moduler) finderPatternBoundary(top, left bool) (*Boundary, error) {
	if top && left {
		return &Boundary{
			lower: Coordintates{row: 0, col: 0},
			upper: Coordintates{row: finderPatternSize, col: finderPatternSize},
		}, nil
	}

	if top && !left {
		return &Boundary{
			lower: Coordintates{row: 0, col: m.qrCodeSize() - finderPatternSize},
			upper: Coordintates{row: finderPatternSize, col: m.qrCodeSize()},
		}, nil
	}

	if !top && left {
		return &Boundary{
			lower: Coordintates{row: m.qrCodeSize() - finderPatternSize, col: 0},
			upper: Coordintates{row: m.qrCodeSize(), col: finderPatternSize},
		}, nil
	}

	return nil, fmt.Errorf("invalid finder pattern location")
}

func (m *Moduler) alignmentPatternBoundary(coordinates Coordintates) Boundary {
	return Boundary{
		lower: Coordintates{row: coordinates.row - 2, col: coordinates.col - 2},
		upper: Coordintates{row: coordinates.row + 3, col: coordinates.col + 3},
	}
}
