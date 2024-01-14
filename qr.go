package qr

import (
	"fmt"
	"math"
	"qr/qr-gen/util"
	"regexp"
	"strconv"
	"strings"
)

type QrMode string
type QrErrCorrectionLvl string
type QrVersion int
type QrModeIndicator string

type QrEncoder interface {
	Encode(s string, lvl QrErrCorrectionLvl) (string, error)
}

type QrEncoderTest interface {
	// QR versioning
	GetMode(s string) (QrMode, error)
	GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error)
	GetModeIndicator(mode QrMode) QrModeIndicator
	GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error)
	// Data encoding
	EncodeNumericInput(s string) string
	EncodeAlphanumericInput(s string) string
	EncodeByteInput(s string) string
	Encode(s string, lvl QrErrCorrectionLvl) (string, error)
	AugmentEncodedInput(s string, version QrVersion, lvl QrErrCorrectionLvl) string
	// Error Correction encoding
	GetMessagePolynomial(encoded string) QrPolynomial
	GetGeneratorPolynomial(version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial
	GetErrorCorrectionCodewords(encoded string, version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial
	// Interleaving
	GetFinalMessage(encoded string, version QrVersion, lvl QrErrCorrectionLvl) string
}

type Encoder struct {
}

func NewEncoderTest() QrEncoderTest {
	return &Encoder{}
}

func NewEncoder() QrEncoder {
	return &Encoder{}
}

// QR versioning

func (e *Encoder) GetMode(s string) (QrMode, error) {
	if matched, _ := regexp.MatchString(PATTERNS[NUMERIC], s); matched {
		return NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(PATTERNS[ALPHA_NUMERIC], s); matched {
		return ALPHA_NUMERIC, nil
	}

	if matched, _ := regexp.MatchString(PATTERNS[BYTE], s); matched {
		return BYTE, nil
	}

	return QrMode(""), fmt.Errorf("Invalid input pattern")
}

func (e *Encoder) GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error) {
	version := 1

	for version <= len(CAPACITIES) {
		if len(s) <= CAPACITIES[QrVersion(version)][lvl][MODE_INDICES[mode]] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return QrVersion(-1), fmt.Errorf("Cannot compute QR version")
}

func (e *Encoder) GetModeIndicator(mode QrMode) QrModeIndicator {
	return MODE_INDICATORS[mode]
}

func (e *Encoder) GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error) {
	cntIndicatorLen, err := e.getCountIndicatorLen(version, mode)

	if err != nil {
		return "", err
	}

	sLenBinary := strconv.FormatInt(int64(len(s)), 2)
	return util.PadLeft(sLenBinary, "0", cntIndicatorLen), nil
}

// Data encoding

func (e *Encoder) EncodeNumericInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryString := strconv.FormatInt(int64(numericValue), 2)

		switch true {
		case numericValue <= 9:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[DIGIT])
		case 10 <= numericValue && numericValue <= 99:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[TEN])
		default:
			binaryString = util.PadLeft(binaryString, "0", QR_NUMERIC_MASKS[HUNDRED])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *Encoder) EncodeAlphanumericInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[ALPHA_NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		var binaryString string
		var firstCharValue int
		var secondCharValue int
		var groupValue int

		if len(group) == 2 {
			firstCharValue, secondCharValue = ALPHA_NUMERIC_VALUES[group[0]], ALPHA_NUMERIC_VALUES[group[1]]
			groupValue = QR_ALPHA_NUMERIC_FACTOR*firstCharValue + secondCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = util.PadLeft(binaryString, "0", QR_ALPHA_NUMERIC_MASKS[FULL_GROUP])
		} else {
			firstCharValue = ALPHA_NUMERIC_VALUES[group[0]]
			groupValue = firstCharValue
			binaryString = strconv.FormatInt(int64(groupValue), 2)
			binaryString = util.PadLeft(binaryString, "0", QR_ALPHA_NUMERIC_MASKS[ONE_ONLY])
		}

		result[index] = binaryString
	}

	return strings.Join(result, "")
}

func (e *Encoder) EncodeByteInput(s string) string {
	groups := util.SplitInGroups(s, SPLIT_VALUES[BYTE])
	result := make([]string, len(groups))

	for index, group := range groups {
		hex := strconv.FormatInt(int64(group[0]), 16)
		hex0, _ := strconv.ParseInt(string(hex[0]), 16, 64)
		hex1, _ := strconv.ParseInt(string(hex[1]), 16, 64)

		bin0 := util.PadLeft(strconv.FormatInt(hex0, 2), "0", QR_BYTE_MASKS[CHAR])
		bin1 := util.PadLeft(strconv.FormatInt(hex1, 2), "0", QR_BYTE_MASKS[CHAR])

		result[index] = bin0 + bin1
	}

	return strings.Join(result, "")
}

func (e *Encoder) EncodeInput(s string, mode QrMode) string {
	switch mode {
	case NUMERIC:
		return e.EncodeNumericInput(s)
	case ALPHA_NUMERIC:
		return e.EncodeAlphanumericInput(s)
	case BYTE:
		return e.EncodeByteInput(s)
	default:
		return ""
	}
}

func (e *Encoder) AugmentEncodedInput(s string, version QrVersion, lvl QrErrCorrectionLvl) string {
	requiredBitsCount := e.getNumberOfRequiredBits(version, lvl)

	s = e.augmentWithTerminatorBits(s, requiredBitsCount)
	remainingBitsCount := requiredBitsCount - len(s)
	if remainingBitsCount == 0 {
		return s
	}

	s = e.augmentWithZeroBits(s)
	remainingBitsCount = requiredBitsCount - len(s)
	if remainingBitsCount == 0 {
		return s
	}

	return e.augmentWithPaddingBits(s, requiredBitsCount)
}

func (e *Encoder) augmentWithTerminatorBits(s string, requiredBitsCount int) string {
	remainingBitsCount := requiredBitsCount - len(s)

	if remainingBitsCount >= 4 {
		return util.PadRight(s, "0", len(s)+4)
	}

	return util.PadRight(s, "0", len(s)+remainingBitsCount)
}

func (e *Encoder) augmentWithZeroBits(s string) string {
	multiple := e.getClosestMultiple(len(s), QR_CODEWORD_SIZE)
	return util.PadRight(s, "0", multiple)
}

func (e *Encoder) augmentWithPaddingBits(s string, requiredBitsCount int) string {
	numberOfPadBytes := (requiredBitsCount - len(s)) / QR_CODEWORD_SIZE
	paddingByteIndex := 0
	paddingSequence := ""

	for i := 0; i < numberOfPadBytes; i++ {
		if paddingByteIndex == 2 {
			paddingByteIndex = paddingByteIndex % 2
		}

		paddingSequence = paddingSequence + QR_PADDING_BYTES[QrPaddingByte(paddingByteIndex)]
		paddingByteIndex += 1
	}

	return s + paddingSequence
}

func (e *Encoder) getClosestMultiple(n int, multipleOf int) int {
	multiple := int(math.Round(float64(n) / float64(multipleOf)))
	return multiple * multipleOf
}

func (e *Encoder) Encode(s string, lvl QrErrCorrectionLvl) (string, error) {
	mode, err := e.GetMode(s)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding mode: %v", err)
	}

	modeIndicator := e.GetModeIndicator(mode)

	version, err := e.GetVersion(s, mode, lvl)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding version: %v", err)
	}

	countIndicator, err := e.GetCountIndicator(s, version, mode)
	if err != nil {
		return "", fmt.Errorf("Error on computing the encoding count indicator: %v", err)
	}

	encodedInput := e.EncodeInput(s, mode)

	return string(modeIndicator) + countIndicator + encodedInput, nil
}

func (e *Encoder) getCountIndicatorLen(version QrVersion, mode QrMode) (int, error) {
	// Extend this functionality for further versions support
	if VERSION_1 <= version && version <= VERSION_5 {
		switch mode {
		case NUMERIC:
			return 10, nil
		case ALPHA_NUMERIC:
			return 9, nil
		case BYTE:
			return 8, nil
		}
	}

	return 0, fmt.Errorf("Cannot compute QR Count Indicator length")
}

func (e *Encoder) getNumberOfRequiredBits(version QrVersion, lvl QrErrCorrectionLvl) int {
	key := e.getECMapKey(version, lvl)
	return QR_CODEWORD_SIZE * QR_EC_INFO[key].TotalDataCodewords
}

func (e *Encoder) getECMapKey(version QrVersion, lvl QrErrCorrectionLvl) string {
	return strconv.Itoa(int(version)) + "-" + string(lvl)
}

// Error Correction encoding

func (e *Encoder) GetMessagePolynomial(encoded string) QrPolynomial {
	codewords := util.SplitInGroups(encoded, QR_CODEWORD_SIZE)
	coefficients := make(QrPolynomial, len(codewords))

	for index, codeword := range codewords {
		decimalValue, _ := strconv.ParseInt(codeword, 2, 64)
		coefficients[len(coefficients)-index-1] = int(decimalValue)
	}

	return coefficients
}

func (e *Encoder) GetGeneratorPolynomial(version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial {
	degree := QR_EC_INFO[e.getECMapKey(version, lvl)].ECCodewordsPerBlock
	coefficients := make(QrPolynomial, 1)

	for i := 1; i <= degree; i++ {
		factorCoefficients := make(QrPolynomial, 2)
		factorCoefficients[0] = i - 1
		coefficients = e.multiplyPolynomials(coefficients, factorCoefficients)
	}

	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = util.ConvertExponentToValue(coefficients[i])
	}

	return coefficients
}

func (e *Encoder) GetErrorCorrectionCodewords(encoded string, version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial {
	messagePolynomial := e.GetMessagePolynomial(encoded)
	generatorPolynomial := e.GetGeneratorPolynomial(version, lvl)
	numErrCorrCodewords := QR_EC_INFO[e.getECMapKey(version, lvl)].ECCodewordsPerBlock

	divisionSteps := len(messagePolynomial)
	messagePolynomial = e.expandPolynomial(messagePolynomial, numErrCorrCodewords)
	generatorPolynomial = e.expandPolynomial(generatorPolynomial, len(messagePolynomial)-len(generatorPolynomial))

	errCorrCodewords := e.dividePolynomials(messagePolynomial, generatorPolynomial, divisionSteps, numErrCorrCodewords)
	return errCorrCodewords[0:numErrCorrCodewords]
}

func (e *Encoder) multiplyPolynomials(firstPoly, secondPoly QrPolynomial) QrPolynomial {
	degreeFirstPoly, degreeSecondPoly := len(firstPoly)-1, len(secondPoly)-1
	result := e.initializeResultAsExponents(degreeFirstPoly + degreeSecondPoly)

	for i := 0; i <= degreeFirstPoly; i++ {
		for j := 0; j <= degreeSecondPoly; j++ {
			exponent := firstPoly[i] + secondPoly[j]

			if exponent >= QR_GALOIS_ORDER {
				exponent = exponent % (QR_GALOIS_ORDER - 1)
			}

			currValue := util.ConvertExponentToValue(exponent)
			prevValue := e.getValueFromResultExponents(result, i+j)

			exponent = util.ConvertValueToExponent(currValue ^ prevValue)
			result[i+j] = exponent
		}
	}

	return result
}

func (e *Encoder) dividePolynomials(divident, divisor QrPolynomial, steps, remainderLen int) QrPolynomial {
	remainder := make(QrPolynomial, remainderLen)

	var currentDivisor QrPolynomial
	copy(currentDivisor, divisor)

	for i := 0; i < steps; i++ {
		dividentLeadTerm, leadTermIndex := e.getLeadingTerm(divident)

		currentDivisor = divisor
		currentDivisor = e.shiftPolynomial(currentDivisor, i)
		currentDivisor = e.multiplyPolynomialByScalar(currentDivisor, dividentLeadTerm, remainderLen)
		remainder = e.xorPolynomials(currentDivisor, divident, leadTermIndex)

		divident = remainder
	}

	return divident
}

func (e *Encoder) multiplyPolynomialByScalar(polynomial QrPolynomial, scalar, remainderLen int) QrPolynomial {
	scalarAlphaValue := util.ConvertValueToExponent(scalar)

	var termAlphaValue int
	result := make(QrPolynomial, len(polynomial))

	currIndex := 0
	isLeadTermFound := false

	for i := len(polynomial) - 1; i >= 0; i-- {
		termAlphaValue = util.ConvertValueToExponent(polynomial[i])

		if termAlphaValue == 0 {
			isLeadTermFound = true
			continue
		}

		if isLeadTermFound && currIndex < remainderLen {
			result[i] = util.ConvertExponentToValue((termAlphaValue + scalarAlphaValue) % (QR_GALOIS_ORDER - 1))
			currIndex++
		}
	}

	return result
}

func (e *Encoder) xorPolynomials(firstPolynomial, secondPolynomial QrPolynomial, leadTermIndex int) QrPolynomial {
	result := make(QrPolynomial, len(secondPolynomial))

	for i := len(secondPolynomial) - 1; i >= 0; i-- {
		if i < leadTermIndex {
			result[i] = firstPolynomial[i] ^ secondPolynomial[i]
		}
	}

	return result
}

func (e *Encoder) getLeadingTerm(polynomial QrPolynomial) (int, int) {
	for i := len(polynomial) - 1; i >= 0; i-- {
		if polynomial[i] != 0 {
			return polynomial[i], i
		}
	}
	return -1, -1
}

func (e *Encoder) shiftPolynomial(polynomial QrPolynomial, unit int) QrPolynomial {
	return append(polynomial[unit:], 0)
}

func (e *Encoder) expandPolynomial(polynomial QrPolynomial, n int) QrPolynomial {
	expandedPolynomial := make(QrPolynomial, len(polynomial)+n)
	copy(expandedPolynomial[n:], polynomial)
	return expandedPolynomial
}

func (e *Encoder) initializeResultAsExponents(degree int) QrPolynomial {
	exponents := make([]int, degree+1)

	for i := 0; i < degree; i++ {
		exponents[i] = -1
	}

	return exponents
}

func (e *Encoder) getValueFromResultExponents(result QrPolynomial, index int) int {
	if index < 0 || index > len(result)-1 || result[index] == -1 {
		return 0
	}
	return util.ConvertExponentToValue(result[index])
}

// Interleaving

func (e *Encoder) GetFinalMessage(inputCodewords string, version QrVersion, lvl QrErrCorrectionLvl) string {
	var codewords string

	if e.isInterleavingNecessary(version, lvl) {
		codewords = e.handleInterleaveProcess(version, lvl, inputCodewords)
	} else {
		errCorrCodewords := e.GetErrorCorrectionCodewords(inputCodewords, version, lvl)
		codewords = inputCodewords + util.ConvertIntListToCodewords(errCorrCodewords)
	}

	return util.PadRight(codewords, "0", len(codewords)+QR_REMAINDER_BITS[version])
}

func (e *Encoder) isInterleavingNecessary(version QrVersion, lvl QrErrCorrectionLvl) bool {
	return QR_EC_INFO[e.getECMapKey(version, lvl)].NumBlocksGroup2 != 0
}

func (e *Encoder) interleaveDataCodewords(version QrVersion, lvl QrErrCorrectionLvl, encoded string) ([]int, [][]int) {
	key := e.getECMapKey(version, lvl)

	group1Size := QR_EC_INFO[key].NumBlocksGroup1
	group1BlockSize := QR_EC_INFO[key].DataCodeworkdsInGroup1Block
	group1Codewords := encoded[:group1Size*group1BlockSize*QR_CODEWORD_SIZE]
	group1Blocks := e.getBlocksOfCodewords(group1Codewords, group1Size, group1BlockSize)

	group2Size := QR_EC_INFO[key].NumBlocksGroup2
	group2BlockSize := QR_EC_INFO[key].DataCodewordsInGroup2Block
	group2Codewords := encoded[group1Size*group1BlockSize*QR_CODEWORD_SIZE:]
	group2Blocks := e.getBlocksOfCodewords(group2Codewords, group2Size, group2BlockSize)

	dataBlocks := append(group1Blocks, group2Blocks...)
	return e.interleaveCodewords(dataBlocks, util.Max(group1BlockSize, group2BlockSize)), dataBlocks
}

func (e *Encoder) interleaveErrCorrCodewords(version QrVersion, lvl QrErrCorrectionLvl, dataBlocks [][]int) []int {
	blocks := make([][]int, len(dataBlocks))

	for i, block := range dataBlocks {
		encoded := strings.Join(util.ConvertIntListToBin(block), "")
		blocks[i] = e.GetErrorCorrectionCodewords(encoded, version, lvl)
	}

	return e.interleaveCodewords(blocks, QR_EC_INFO[e.getECMapKey(version, lvl)].ECCodewordsPerBlock)
}

func (e *Encoder) handleInterleaveProcess(version QrVersion, lvl QrErrCorrectionLvl, codewords string) string {
	interleavedDataCodewords, dataBlocks := e.interleaveDataCodewords(version, lvl, codewords)
	interleavedECCodewords := e.interleaveErrCorrCodewords(version, lvl, dataBlocks)

	interleavedDataBinary := util.ConvertIntListToBin(interleavedDataCodewords)
	interleavedECBinary := util.ConvertIntListToBin(interleavedECCodewords)

	return strings.Join(interleavedDataBinary, "") + strings.Join(interleavedECBinary, "")
}

func (e *Encoder) interleaveCodewords(blocks [][]int, length int) []int {
	var result []int

	for i := 0; i < length; i++ {
		for _, block := range blocks {
			if i < len(block) {
				result = append(result, block[i])
			}
		}
	}

	return result
}

func (e *Encoder) getBlocksOfCodewords(input string, blocksCount int, blockSize int) [][]int {
	blocks := make([][]int, blocksCount)

	for i := 0; i < blocksCount; i++ {
		currBlock := input[:blockSize*QR_CODEWORD_SIZE]
		currBlockSlice := make([]int, blockSize)

		for j := 0; j < blockSize; j++ {
			bin := currBlock[:QR_CODEWORD_SIZE]
			value, _ := strconv.ParseInt(bin, 2, 64)
			currBlockSlice[j] = int(value)
			currBlock = currBlock[QR_CODEWORD_SIZE:]
		}

		blocks[i] = currBlockSlice
		input = input[blockSize*QR_CODEWORD_SIZE:]
	}

	return blocks
}
