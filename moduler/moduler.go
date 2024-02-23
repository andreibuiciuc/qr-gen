package moduler

import (
	"fmt"
	"math"
	"qr/qr-gen/matrix"
	"qr/qr-gen/util"
	"qr/qr-gen/versioner"
	"strconv"
)

type ModulerInterface interface {
	CreateModuleMatrix(data string) (*matrix.Matrix[util.Module], []*matrix.Matrix[util.Module], Penalty)
}

type Moduler struct {
	version      versioner.QrVersion
	ecLevel      versioner.QrEcLevel
	moduleMatrix *matrix.Matrix[util.Module]
}

type Coordinates struct {
	row int
	col int
}

type Boundary struct {
	lower Coordinates
	upper Coordinates
}

type Penalty struct {
	score1 int
	score2 int
	score3 int
	score4 int
	total  int
}

const finderPatternSize = 7

var rulePattern = []util.Module{util.Module_DARKEN, util.Module_LIGHTEN, util.Module_DARKEN, util.Module_DARKEN, util.Module_DARKEN, util.Module_LIGHTEN, util.Module_DARKEN, util.Module_LIGHTEN, util.Module_LIGHTEN, util.Module_LIGHTEN, util.Module_LIGHTEN}
var reversedRulePattern = []util.Module{util.Module_LIGHTEN, util.Module_LIGHTEN, util.Module_LIGHTEN, util.Module_LIGHTEN, util.Module_DARKEN, util.Module_LIGHTEN, util.Module_DARKEN, util.Module_DARKEN, util.Module_DARKEN, util.Module_LIGHTEN, util.Module_DARKEN}

// These locations stand only for alignment patterns that do not overlap with finder patterns
// This can be improved by checking the overlap programatically
var allignmentPatternLocation = map[versioner.QrVersion][]Coordinates{
	1: {},
	2: {Coordinates{18, 18}},
	3: {Coordinates{22, 22}},
	4: {Coordinates{26, 26}},
	5: {Coordinates{30, 30}},
}

var maskFormula = map[int]func(Coordinates) bool{
	0: func(c Coordinates) bool {
		return (c.row+c.col)%2 == 0
	},
	1: func(c Coordinates) bool {
		return c.row%2 == 0
	},
	2: func(c Coordinates) bool {
		return c.col%3 == 0
	},
	3: func(c Coordinates) bool {
		return (c.row+c.col)%3 == 0
	},
	4: func(c Coordinates) bool {
		return (int(math.Floor(float64(c.row)/2))+int(math.Floor(float64(c.col)/3)))%2 == 0
	},
	5: func(c Coordinates) bool {
		return (c.row*c.col)%2+(c.row*c.col)%3 == 0
	},
	6: func(c Coordinates) bool {
		return ((c.row*c.col)%2+(c.row*c.col)%3)%2 == 0
	},
	7: func(c Coordinates) bool {
		return ((c.row+c.col)%2+(c.row*c.col%3))%2 == 0
	},
}

func New(version versioner.QrVersion, ecLevel versioner.QrEcLevel) ModulerInterface {
	return &Moduler{
		version: version,
		ecLevel: ecLevel,
	}
}

// TODO: Refactor this after defining better test cases
func (m *Moduler) CreateModuleMatrix(data string) (*matrix.Matrix[util.Module], []*matrix.Matrix[util.Module], Penalty) {
	m.prepareModuleMatrix(data)

	moduleCoords := m.placeDataBits(data)
	candidates := m.getModuleMatrixCandidates(moduleCoords)
	matrix, penalty := m.getBestMaskedMatrix(candidates)

	return matrix, candidates, penalty
}

func (m *Moduler) prepareModuleMatrix(data string) {
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
}

func (m *Moduler) qrCodeSize() int {
	return (int(m.version)-1)*4 + 21
}

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
	for _, c := range allignmentPatternLocation[m.version] {
		boundary := m.alignmentPatternBoundary(c)
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

	for i := boundary.lower.col; i < boundary.upper.col+2; i++ {
		if val, _ := m.moduleMatrix.At(boundary.upper.row+1, i); val == util.Module_EMPTY {
			m.moduleMatrix.Set(boundary.upper.row+1, i, util.Module_RESERVED)
		}
	}

	for i := boundary.lower.row; i < boundary.upper.row+2; i++ {
		if val, _ := m.moduleMatrix.At(i, boundary.upper.col+1); val == util.Module_EMPTY {
			m.moduleMatrix.Set(i, boundary.upper.col+1, util.Module_RESERVED)
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
func (m *Moduler) placeDataBits(data string) []Coordinates {
	var moduleCoords []Coordinates
	currentCellCoord := Coordinates{row: m.qrCodeSize() - 1, col: m.qrCodeSize() - 1}

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
			moduleCoords = append(moduleCoords, currentCellCoord)
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

	return moduleCoords
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
			lower: Coordinates{row: 0, col: 0},
			upper: Coordinates{row: finderPatternSize, col: finderPatternSize},
		}, nil
	}

	if top && !left {
		return &Boundary{
			lower: Coordinates{row: 0, col: m.qrCodeSize() - finderPatternSize},
			upper: Coordinates{row: finderPatternSize, col: m.qrCodeSize()},
		}, nil
	}

	if !top && left {
		return &Boundary{
			lower: Coordinates{row: m.qrCodeSize() - finderPatternSize, col: 0},
			upper: Coordinates{row: m.qrCodeSize(), col: finderPatternSize},
		}, nil
	}

	return nil, fmt.Errorf("invalid finder pattern location")
}

func (m *Moduler) alignmentPatternBoundary(c Coordinates) Boundary {
	return Boundary{
		lower: Coordinates{row: c.row - 2, col: c.col - 2},
		upper: Coordinates{row: c.row + 3, col: c.col + 3},
	}
}

// Gets the masked matrix candidates based on masking formulas
func (m *Moduler) getModuleMatrixCandidates(moduleCoords []Coordinates) []*matrix.Matrix[util.Module] {
	matrixCandidates := make([]*matrix.Matrix[util.Module], len(maskFormula))

	for i := range maskFormula {
		matrixCandidates[i] = m.maskModuleMatrix(moduleCoords, i)
	}

	return matrixCandidates
}

// Masks a module matrix based on the given rule
func (m *Moduler) maskModuleMatrix(moduleCoords []Coordinates, rule int) *matrix.Matrix[util.Module] {
	matrixCandidate := matrix.NewMatrix[util.Module](m.qrCodeSize(), m.qrCodeSize())
	matrixCandidate.SetMatrix(m.moduleMatrix.GetMatrix())
	m.setFormatInformationModules(matrixCandidate, rule)

	for _, c := range moduleCoords {
		module, _ := matrixCandidate.At(c.row, c.col)

		if maskFormula[rule](c) {
			module = m.toggleModule(c)
		}

		matrixCandidate.Set(c.row, c.col, module)
	}

	return matrixCandidate
}

// Toggles the value of a data module, necessary as the modules are expressive and ARE NOT values of 0 or 1
func (m *Moduler) toggleModule(c Coordinates) util.Module {
	if val, _ := m.moduleMatrix.At(c.row, c.col); val == util.Module_LIGHTEN {
		return util.Module_DARKEN
	}
	return util.Module_LIGHTEN
}

func (m *Moduler) setFormatInformationModules(matrix *matrix.Matrix[util.Module], rule int) {
	boundary, _ := m.finderPatternBoundary(true, true)
	format := util.FormatInformationStrings[rune(m.ecLevel)][rule]
	format += format
	index := 0

	for i := boundary.lower.col; i < boundary.upper.col+1; i++ {
		if val, _ := matrix.At(boundary.upper.row+1, i); !util.IsModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(boundary.upper.row+1, i, util.GetDataModule(int(bit)))
			index += 1
		}
	}

	for i := boundary.upper.row + 1; i >= 0; i-- {
		if val, _ := matrix.At(i, boundary.upper.col+1); !util.IsModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(i, boundary.upper.col+1, util.GetDataModule(int(bit)))
			index += 1
		}
	}

	boundary, _ = m.finderPatternBoundary(false, true)

	for i := boundary.upper.row - 1; i >= boundary.lower.row-1; i-- {
		if val, _ := matrix.At(i, boundary.upper.col+1); !util.IsModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(i, boundary.upper.col+1, util.GetDataModule(int(bit)))
			index += 1
		}
	}

	boundary, _ = m.finderPatternBoundary(true, false)

	for i := boundary.lower.col - 1; i < boundary.upper.col; i++ {
		if val, _ := matrix.At(boundary.upper.row+1, i); !util.IsModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(boundary.upper.row+1, i, util.GetDataModule(int(bit)))
			index += 1
		}
	}
}

func (m *Moduler) getBestMaskedMatrix(candidates []*matrix.Matrix[util.Module]) (*matrix.Matrix[util.Module], Penalty) {
	penalty := m.evaluateMatrixCandidate(candidates[0])
	matrix := candidates[0]

	for i := 1; i < len(candidates); i++ {
		currentPenalty := m.evaluateMatrixCandidate(candidates[i])
		if currentPenalty.total < penalty.total {
			penalty = currentPenalty
			matrix = candidates[i]
		}
	}

	return matrix, penalty
}

func (m *Moduler) evaluateMatrixCandidate(matrix *matrix.Matrix[util.Module]) Penalty {
	penalty := Penalty{}
	penalty.score1 = m.computeFirstPenalty(matrix)
	penalty.score2 = m.computeSecondPenalty(matrix)
	penalty.score3 = m.computeThirdPenalty(matrix)
	penalty.score4 = m.computeFourthPenalty(matrix)
	penalty.total = penalty.score1 + penalty.score2 + penalty.score3 + penalty.score4
	return penalty
}

// Implements the first penalty score strategy
func (m *Moduler) computeFirstPenalty(matrix *matrix.Matrix[util.Module]) int {
	rowPenalty := 0
	colPenalty := 0

	for i := 0; i < m.qrCodeSize(); i++ {
		row, _ := matrix.RowAt(i)
		rowPenalty += m.computeModulesLinePenalty(row)
		col, _ := matrix.ColumnAt(i)
		colPenalty += m.computeModulesLinePenalty(col)
	}

	return rowPenalty + colPenalty
}

// Implements the second penalty score strategy
func (m *Moduler) computeSecondPenalty(matrix *matrix.Matrix[util.Module]) int {
	count := 0

	for i := 0; i < matrix.Width()-1; i++ {
		for j := 0; j < matrix.Height()-1; j++ {
			module, _ := matrix.At(i, j)
			moduleRight, _ := matrix.At(i, j+1)
			moduleBottom, _ := matrix.At(i+1, j)
			moduleOpposite, _ := matrix.At(i+1, j+1)

			if util.IsModuleLighten(module) == util.IsModuleLighten(moduleRight) &&
				util.IsModuleLighten(module) == util.IsModuleLighten(moduleBottom) &&
				util.IsModuleLighten(module) == util.IsModuleLighten(moduleOpposite) {
				count += 1
			}
		}
	}

	return count * 3
}

// Implements the third penalty score strategy
func (m *Moduler) computeThirdPenalty(matrix *matrix.Matrix[util.Module]) int {
	count := 0

	for i := 0; i < m.qrCodeSize(); i++ {
		row, _ := matrix.RowAt(i)
		for col := 0; col < m.qrCodeSize()-11+1; col++ {
			if m.isRulePattern(rulePattern, row[col:col+11]) {
				count += 1
			}
			if m.isRulePattern(reversedRulePattern, row[col:col+11]) {
				count += 1
			}
		}
	}

	for i := 0; i < m.qrCodeSize(); i++ {
		col, _ := matrix.ColumnAt(i)
		for row := 0; row < m.qrCodeSize()-11+1; row++ {
			if m.isRulePattern(rulePattern, col[row:row+11]) {
				count += 1
			}
			if m.isRulePattern(reversedRulePattern, col[row:row+11]) {
				count += 1
			}
		}
	}

	return count * 40
}

// Implements the fourth penalty score strategy
func (m *Moduler) computeFourthPenalty(matrix *matrix.Matrix[util.Module]) int {
	total := m.qrCodeSize() * m.qrCodeSize()
	darkModules := 0

	for _, row := range matrix.GetMatrix() {
		for _, module := range row {
			if !util.IsModuleLighten(module) {
				darkModules += 1
			}
		}
	}

	percentage := darkModules * 100 / total

	if percentage > 50 {
		percentage = int(math.Floor(float64(percentage)/5)) * 5
	} else {
		percentage = int(math.Ceil(float64(percentage)/5)) * 5
	}

	return int(math.Abs(float64(percentage)-50)) * 2
}

func (m *Moduler) computeModulesLinePenalty(modules []util.Module) int {
	count := 0
	score := 0
	isLighten := false

	for _, module := range modules {
		if x := util.IsModuleLighten(module); x != isLighten {
			isLighten = x
			count = 1
		} else {
			count++
			if count == 5 {
				score += 3
			} else if count > 5 {
				score += 1
			}
		}
	}

	return score
}

func (m *Moduler) isRulePattern(pattern, seq []util.Module) bool {
	for i := range seq {
		val := util.Module_DARKEN
		if util.IsModuleLighten(seq[i]) {
			val = util.Module_LIGHTEN
		}

		if pattern[i] != val {
			return false
		}
	}
	return true
}
