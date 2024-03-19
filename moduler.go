package qr

import (
	"fmt"
	"math"
	"strconv"
)

type moduler struct {
	version      int
	lvl          rune
	moduleMatrix *matrix[module]
}

type coordinates struct {
	row int
	col int
}

type boundary struct {
	lower coordinates
	upper coordinates
}

type penalty struct {
	score1 int
	score2 int
	score3 int
	score4 int
	total  int
}

const finderPatternSize = 7

var rulePattern = []module{module_DARKEN, module_LIGHTEN, module_DARKEN, module_DARKEN, module_DARKEN, module_LIGHTEN, module_DARKEN, module_LIGHTEN, module_LIGHTEN, module_LIGHTEN, module_LIGHTEN}
var reversedRulePattern = []module{module_LIGHTEN, module_LIGHTEN, module_LIGHTEN, module_LIGHTEN, module_DARKEN, module_LIGHTEN, module_DARKEN, module_DARKEN, module_DARKEN, module_LIGHTEN, module_DARKEN}

var allignmentPatternLocation = map[int][]coordinates{
	1: {},
	2: {coordinates{18, 18}},
	3: {coordinates{22, 22}},
	4: {coordinates{26, 26}},
	5: {coordinates{30, 30}},
}

var maskFormula = map[int]func(coordinates) bool{
	0: func(c coordinates) bool {
		return (c.row+c.col)%2 == 0
	},
	1: func(c coordinates) bool {
		return c.row%2 == 0
	},
	2: func(c coordinates) bool {
		return c.col%3 == 0
	},
	3: func(c coordinates) bool {
		return (c.row+c.col)%3 == 0
	},
	4: func(c coordinates) bool {
		return (int(math.Floor(float64(c.row)/2))+int(math.Floor(float64(c.col)/3)))%2 == 0
	},
	5: func(c coordinates) bool {
		return (c.row*c.col)%2+(c.row*c.col)%3 == 0
	},
	6: func(c coordinates) bool {
		return ((c.row*c.col)%2+(c.row*c.col)%3)%2 == 0
	},
	7: func(c coordinates) bool {
		return ((c.row+c.col)%2+(c.row*c.col%3))%2 == 0
	},
}

func newModuler(v int, lvl rune) *moduler {
	return &moduler{
		version: v,
		lvl:     lvl,
	}
}

func (m *moduler) createModuleMatrix(data string) (*matrix[module], penalty) {
	m.prepareModuleMatrix(data)

	moduleCoords := m.placeDataBits(data)
	candidates := m.getModuleMatrixCandidates(moduleCoords)
	matrix, penalty := m.getBestMaskedMatrix(candidates)
	matrix.Expand(4)

	return matrix, penalty
}

func (m *moduler) prepareModuleMatrix(data string) {
	qrCodeSize := m.qrCodeSize()
	m.moduleMatrix = newMatrix[module](qrCodeSize, qrCodeSize)
	m.moduleMatrix.Init(module_EMPTY)

	m.setTopLeftFinderPattern()
	m.setTopRightFinderPattern()
	m.setBottomLeftFinderPattern()
	m.setAlignmentPatterns()
	m.setTimingPatterns()
	m.setDarkModule()
	m.reserveFormatArea()
}

func (m *moduler) qrCodeSize() int {
	return (int(m.version)-1)*4 + 21
}

// Sets the top left finder pattern in the module matrix
func (m *moduler) setTopLeftFinderPattern() {
	boundary, _ := m.finderPatternBoundary(true, true)
	m.patchPattern(*boundary, module_FINDER_LIGHTEN, module_FINDER_DARKEN)

	for i := boundary.lower.row; i <= boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.upper.col, module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.upper.row, i, module_SEPARATOR)
	}
}

// Sets the top right finder pattern in the module matrix
func (m *moduler) setTopRightFinderPattern() {
	boundary, _ := m.finderPatternBoundary(true, false)
	m.patchPattern(*boundary, module_FINDER_LIGHTEN, module_FINDER_DARKEN)

	for i := boundary.lower.row; i <= boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.lower.col-1, module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.upper.row, i, module_SEPARATOR)
	}
}

// Sets the bottom left finder pattern in the module matrix
func (m *moduler) setBottomLeftFinderPattern() {
	boundary, _ := m.finderPatternBoundary(false, true)
	m.patchPattern(*boundary, module_FINDER_LIGHTEN, module_FINDER_DARKEN)

	for i := boundary.lower.row - 1; i < boundary.upper.row; i++ {
		m.moduleMatrix.Set(i, boundary.upper.col, module_SEPARATOR)
	}

	for i := boundary.lower.col; i < boundary.upper.col; i++ {
		m.moduleMatrix.Set(boundary.lower.row-1, i, module_SEPARATOR)
	}
}

// Sets the alignment patterns in the module matrix
func (m *moduler) setAlignmentPatterns() {
	for _, c := range allignmentPatternLocation[m.version] {
		boundary := m.alignmentPatternBoundary(c)
		m.patchPattern(boundary, module_ALIGNMENT_LIGHTEN, module_ALIGNMENT_DARKEN)
	}
}

// Sets the timing patterns in the module matrix
func (m *moduler) setTimingPatterns() {
	topLeftFinderBoundary, _ := m.finderPatternBoundary(true, true)
	topRightFinderBoundary, _ := m.finderPatternBoundary(true, false)
	bottomLeftFinderBoundary, _ := m.finderPatternBoundary(false, true)

	val := 1
	for i := topLeftFinderBoundary.upper.col - 1; i < topRightFinderBoundary.lower.col; i++ {
		if val == 0 {
			m.moduleMatrix.Set(6, i, module_TIMING_LIGHTEN)
		} else {
			m.moduleMatrix.Set(6, i, module_TIMING_DARKEN)
		}
		val = (val + 1) % 2
	}

	val = 1
	for i := topLeftFinderBoundary.upper.row - 1; i < bottomLeftFinderBoundary.lower.row; i++ {
		if val == 0 {
			m.moduleMatrix.Set(i, 6, module_TIMING_LIGHTEN)
		} else {
			m.moduleMatrix.Set(i, 6, module_TIMING_DARKEN)
		}
		val = (val + 1) % 2
	}
}

// Sets the dark module in the module matrix
func (m *moduler) setDarkModule() {
	m.moduleMatrix.Set(4*int(m.version)+9, 8, module_DARK)
}

// Sets the reserved format information area in the module matrix
func (m *moduler) reserveFormatArea() {
	boundary, _ := m.finderPatternBoundary(true, true)

	for i := boundary.lower.col; i < boundary.upper.col+2; i++ {
		if val, _ := m.moduleMatrix.At(boundary.upper.row+1, i); val == module_EMPTY {
			m.moduleMatrix.Set(boundary.upper.row+1, i, module_RESERVED)
		}
	}

	for i := boundary.lower.row; i < boundary.upper.row+2; i++ {
		if val, _ := m.moduleMatrix.At(i, boundary.upper.col+1); val == module_EMPTY {
			m.moduleMatrix.Set(i, boundary.upper.col+1, module_RESERVED)
		}
	}

	boundary, _ = m.finderPatternBoundary(true, false)

	for i := boundary.lower.col - 1; i < boundary.upper.col; i++ {
		if val, _ := m.moduleMatrix.At(boundary.upper.row+1, i); val == module_EMPTY {
			m.moduleMatrix.Set(boundary.upper.row+1, i, module_RESERVED)
		}
	}

	boundary, _ = m.finderPatternBoundary(false, true)

	for i := boundary.lower.row - 1; i < boundary.upper.row; i++ {
		if val, _ := m.moduleMatrix.At(i, boundary.upper.col+1); val == module_EMPTY {
			m.moduleMatrix.Set(i, boundary.upper.col+1, module_RESERVED)
		}
	}
}

// Places the encoded data bits in the module matrix
func (m *moduler) placeDataBits(data string) []coordinates {
	var moduleCoords []coordinates
	currentCellCoord := coordinates{row: m.qrCodeSize() - 1, col: m.qrCodeSize() - 1}

	indexInBits := 0
	indexInModules := 0
	isUpwardMovement := -1
	mod := module_EMPTY

	for currentCellCoord.col >= 0 {

		if val, _ := m.moduleMatrix.At(currentCellCoord.row, currentCellCoord.col); val == module_EMPTY {
			if bit, _ := strconv.ParseInt(string(data[indexInBits]), 2, 64); bit == 0 {
				mod = module_LIGHTEN
			} else {
				mod = module_DARKEN
			}
			m.moduleMatrix.Set(currentCellCoord.row, currentCellCoord.col, module(mod))
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
func (m *moduler) patchPattern(boundary boundary, ligthenModule module, darkenModule module) {
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

func (m *moduler) isPatternModuleDarken(i, j int, boundary boundary) bool {
	return i == boundary.lower.row || i == boundary.upper.row-1 || j == boundary.lower.col || j == boundary.upper.col-1
}

func (m *moduler) isPatternModuleLighten(i, j int, boundary boundary) bool {
	return i == boundary.lower.row+1 || i == boundary.upper.row-2 || j == boundary.lower.col+1 || j == boundary.upper.col-2
}

func (m *moduler) finderPatternBoundary(top, left bool) (*boundary, error) {
	if top && left {
		return &boundary{
			lower: coordinates{row: 0, col: 0},
			upper: coordinates{row: finderPatternSize, col: finderPatternSize},
		}, nil
	}

	if top && !left {
		return &boundary{
			lower: coordinates{row: 0, col: m.qrCodeSize() - finderPatternSize},
			upper: coordinates{row: finderPatternSize, col: m.qrCodeSize()},
		}, nil
	}

	if !top && left {
		return &boundary{
			lower: coordinates{row: m.qrCodeSize() - finderPatternSize, col: 0},
			upper: coordinates{row: m.qrCodeSize(), col: finderPatternSize},
		}, nil
	}

	return nil, fmt.Errorf("invalid finder pattern location")
}

func (m *moduler) alignmentPatternBoundary(c coordinates) boundary {
	return boundary{
		lower: coordinates{row: c.row - 2, col: c.col - 2},
		upper: coordinates{row: c.row + 3, col: c.col + 3},
	}
}

// Gets the masked matrix candidates based on masking formulas
func (m *moduler) getModuleMatrixCandidates(moduleCoords []coordinates) []*matrix[module] {
	matrixCandidates := make([]*matrix[module], len(maskFormula))

	for i := range maskFormula {
		matrixCandidates[i] = m.maskModuleMatrix(moduleCoords, i)
	}

	return matrixCandidates
}

// Masks a module matrix based on the given rule
func (m *moduler) maskModuleMatrix(moduleCoords []coordinates, rule int) *matrix[module] {
	matrixCandidate := newMatrix[module](m.qrCodeSize(), m.qrCodeSize())
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
func (m *moduler) toggleModule(c coordinates) module {
	if val, _ := m.moduleMatrix.At(c.row, c.col); val == module_LIGHTEN {
		return module_DARKEN
	}
	return module_LIGHTEN
}

func (m *moduler) setFormatInformationModules(matrix *matrix[module], rule int) {
	boundary, _ := m.finderPatternBoundary(true, true)
	format := fmtInfoCodes[rune(m.lvl)][rule]
	format += format
	index := 0

	for i := boundary.lower.col; i < boundary.upper.col+1; i++ {
		if val, _ := matrix.At(boundary.upper.row+1, i); !isModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(boundary.upper.row+1, i, getDataModule(int(bit)))
			index += 1
		}
	}

	for i := boundary.upper.row + 1; i >= 0; i-- {
		if val, _ := matrix.At(i, boundary.upper.col+1); !isModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(i, boundary.upper.col+1, getDataModule(int(bit)))
			index += 1
		}
	}

	boundary, _ = m.finderPatternBoundary(false, true)

	for i := boundary.upper.row - 1; i >= boundary.lower.row-1; i-- {
		if val, _ := matrix.At(i, boundary.upper.col+1); !isModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(i, boundary.upper.col+1, getDataModule(int(bit)))
			index += 1
		}
	}

	boundary, _ = m.finderPatternBoundary(true, false)

	for i := boundary.lower.col - 1; i < boundary.upper.col; i++ {
		if val, _ := matrix.At(boundary.upper.row+1, i); !isModuleSkippedForFormat(val) {
			bit, _ := strconv.ParseInt(string(format[index]), 2, 64)
			matrix.Set(boundary.upper.row+1, i, getDataModule(int(bit)))
			index += 1
		}
	}
}

func (m *moduler) getBestMaskedMatrix(candidates []*matrix[module]) (*matrix[module], penalty) {
	penalty := m.evaluateMatrixCandidate(candidates[0])
	scores := make([]int, len(candidates))
	scores[0] = penalty.total
	matrix := candidates[0]

	for i := 1; i < len(candidates); i++ {
		currentPenalty := m.evaluateMatrixCandidate(candidates[i])
		scores[i] = currentPenalty.total
		if currentPenalty.total < penalty.total {
			penalty = currentPenalty
			matrix = candidates[i]
		}
	}

	return matrix, penalty
}

func (m *moduler) evaluateMatrixCandidate(matrix *matrix[module]) penalty {
	penalty := penalty{}
	penalty.score1 = m.computeFirstPenalty(matrix)
	penalty.score2 = m.computeSecondPenalty(matrix)
	penalty.score3 = m.computeThirdPenalty(matrix)
	penalty.score4 = m.computeFourthPenalty(matrix)
	penalty.total = penalty.score1 + penalty.score2 + penalty.score3 + penalty.score4
	return penalty
}

// Implements the first penalty score strategy
func (m *moduler) computeFirstPenalty(matrix *matrix[module]) int {
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
func (m *moduler) computeSecondPenalty(matrix *matrix[module]) int {
	count := 0

	for i := 0; i < matrix.Width()-1; i++ {
		for j := 0; j < matrix.Height()-1; j++ {
			module, _ := matrix.At(i, j)
			moduleRight, _ := matrix.At(i, j+1)
			moduleBottom, _ := matrix.At(i+1, j)
			moduleOpposite, _ := matrix.At(i+1, j+1)

			if isModuleLighten(module) == isModuleLighten(moduleRight) &&
				isModuleLighten(module) == isModuleLighten(moduleBottom) &&
				isModuleLighten(module) == isModuleLighten(moduleOpposite) {
				count += 1
			}
		}
	}

	return count * 3
}

// Implements the third penalty score strategy
func (m *moduler) computeThirdPenalty(matrix *matrix[module]) int {
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
func (m *moduler) computeFourthPenalty(matrix *matrix[module]) int {
	total := m.qrCodeSize() * m.qrCodeSize()
	darkModules := 0

	for _, row := range matrix.GetMatrix() {
		for _, module := range row {
			if !isModuleLighten(module) {
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

func (m *moduler) computeModulesLinePenalty(modules []module) int {
	count := 0
	score := 0
	isLighten := false

	for _, module := range modules {
		if x := isModuleLighten(module); x != isLighten {
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

func (m *moduler) isRulePattern(pattern, seq []module) bool {
	for i := range seq {
		val := module_DARKEN
		if isModuleLighten(seq[i]) {
			val = module_LIGHTEN
		}

		if pattern[i] != val {
			return false
		}
	}
	return true
}
