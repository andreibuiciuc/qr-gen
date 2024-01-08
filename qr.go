package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type QrMode string
type QrErrCorrectionLvl string
type QrVersion int
type QrModeIndicator string

type QrEncoder interface {
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
	splitInGroups(s string, n int) []string
	// Error Correction encoding
	GetMessagePolynomial(encoded string) QrPolynomial
	GetGeneratorPolynomial(version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial
	GetErrorCorrectionCodewords(encoded string, version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial
}

type Encoder struct {
	alpha int
}

func NewEncoder() QrEncoder {
	return &Encoder{
		alpha: 2,
	}
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

	return QrMode(EMPTY_STRING), fmt.Errorf("Invalid input pattern")
}

func (e *Encoder) GetVersion(s string, mode QrMode, lvl QrErrCorrectionLvl) (QrVersion, error) {
	version := 1

	for version <= len(CAPACITIES) {
		if len(s) <= CAPACITIES[QrVersion(version)][lvl][MODE_INDICES[mode]] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return QrVersion(INVALID_IDX), fmt.Errorf("Cannot compute QR version")
}

func (e *Encoder) GetModeIndicator(mode QrMode) QrModeIndicator {
	return MODE_INDICATORS[mode]
}

func (e *Encoder) GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error) {
	cntIndicatorLen, err := e.getCountIndicatorLen(version, mode)

	if err != nil {
		return EMPTY_STRING, err
	}

	sLenBinary := strconv.FormatInt(int64(len(s)), BINARY_RADIX)
	return e.padStart(sLenBinary, DEFAULT_PAD_CHAR, cntIndicatorLen), nil
}

// Data encoding

func (e *Encoder) EncodeNumericInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		numericValue, _ := strconv.Atoi(group)
		binaryString := strconv.FormatInt(int64(numericValue), BINARY_RADIX)

		switch true {
		case numericValue <= 9:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[DIGIT])
		case 10 <= numericValue && numericValue <= 99:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[TEN])
		default:
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, NUMERIC_MASKS[HUNDRED])
		}

		result[index] = binaryString
	}

	return strings.Join(result, EMPTY_STRING)
}

func (e *Encoder) EncodeAlphanumericInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[ALPHA_NUMERIC])
	result := make([]string, len(groups))

	for index, group := range groups {
		var binaryString string
		var firstCharValue int
		var secondCharValue int
		var groupValue int

		if len(group) == 2 {
			firstCharValue, secondCharValue = ALPHA_NUMERIC_VALUES[group[0]], ALPHA_NUMERIC_VALUES[group[1]]
			groupValue = 45*firstCharValue + secondCharValue
			binaryString = strconv.FormatInt(int64(groupValue), BINARY_RADIX)
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, ALPHA_NUMERIC_MASKS[FULL_GROUP])
		} else {
			firstCharValue = ALPHA_NUMERIC_VALUES[group[0]]
			groupValue = firstCharValue
			binaryString = strconv.FormatInt(int64(groupValue), BINARY_RADIX)
			binaryString = e.padStart(binaryString, DEFAULT_PAD_CHAR, ALPHA_NUMERIC_MASKS[ONE_ONLY])
		}

		result[index] = binaryString
	}

	return strings.Join(result, EMPTY_STRING)
}

func (e *Encoder) EncodeByteInput(s string) string {
	groups := e.splitInGroups(s, SPLIT_VALUES[BYTE])
	result := make([]string, len(groups))

	for index, group := range groups {
		hex := strconv.FormatInt(int64(group[0]), HEXADECIMAL_RADIX)
		hex0, _ := strconv.ParseInt(string(hex[0]), HEXADECIMAL_RADIX, INTEGER_RADIX)
		hex1, _ := strconv.ParseInt(string(hex[1]), HEXADECIMAL_RADIX, INTEGER_RADIX)

		bin0 := e.padStart(strconv.FormatInt(hex0, BINARY_RADIX), DEFAULT_PAD_CHAR, BYTE_MASKS[CHAR])
		bin1 := e.padStart(strconv.FormatInt(hex1, BINARY_RADIX), DEFAULT_PAD_CHAR, BYTE_MASKS[CHAR])

		result[index] = bin0 + bin1
	}

	return strings.Join(result, EMPTY_STRING)
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
		return EMPTY_STRING
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
		return e.padRight(s, DEFAULT_PAD_CHAR, len(s)+4)
	}

	return e.padRight(s, DEFAULT_PAD_CHAR, len(s)+remainingBitsCount)
}

func (e *Encoder) augmentWithZeroBits(s string) string {
	multiple := e.getClosestMultiple(len(s), CODEWORD_BITS)
	return e.padRight(s, DEFAULT_PAD_CHAR, multiple)
}

func (e *Encoder) augmentWithPaddingBits(s string, requiredBitsCount int) string {
	numberOfPadBytes := (requiredBitsCount - len(s)) / 8
	paddingByteIndex := 0
	paddingSequence := EMPTY_STRING

	for i := 0; i < numberOfPadBytes; i++ {
		if paddingByteIndex == 2 {
			paddingByteIndex = paddingByteIndex % 2
		}

		paddingSequence = paddingSequence + PADDING_BYTES[QrPaddingByte(paddingByteIndex)]
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
		return EMPTY_STRING, fmt.Errorf("Error on computing the encoding mode: %v", err)
	}

	modeIndicator := e.GetModeIndicator(mode)

	version, err := e.GetVersion(s, mode, lvl)
	if err != nil {
		return EMPTY_STRING, fmt.Errorf("Error on computing the encoding version: %v", err)
	}

	countIndicator, err := e.GetCountIndicator(s, version, mode)
	if err != nil {
		return EMPTY_STRING, fmt.Errorf("Error on computing the encoding count indicator: %v", err)
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
	return CODEWORD_BITS * ERROR_CORRECTION_INFO[key].TotalDataCodewords
}

func (e *Encoder) padStart(s string, padChar string, n int) string {
	return strings.Repeat(padChar, n-len(s)) + s
}

func (e *Encoder) padRight(s string, padChar string, n int) string {
	return s + strings.Repeat(padChar, n-len(s))
}

func (e *Encoder) splitInGroups(s string, n int) []string {
	var result []string
	count := 0
	start := 0

	for i := 0; i < len(s); i++ {
		count += 1

		if count == n {
			result = append(result, s[start:i+1])
			start = i + 1
			count = 0
		}
	}

	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}

func (e *Encoder) getECMapKey(version QrVersion, lvl QrErrCorrectionLvl) string {
	return strconv.Itoa(int(version)) + "-" + string(lvl)
}

// Error Correction encoding

func (e *Encoder) GetMessagePolynomial(encoded string) QrPolynomial {
	codewords := e.splitInGroups(encoded, CODEWORD_BITS)
	coefficients := make(QrPolynomial, len(codewords))

	for index, codeword := range codewords {
		decimalValue, _ := strconv.ParseInt(codeword, BINARY_RADIX, INTEGER_RADIX)
		coefficients[len(coefficients)-index-1] = int(decimalValue)
	}

	return coefficients
}

func (e *Encoder) GetGeneratorPolynomial(version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial {
	degree := ERROR_CORRECTION_INFO[e.getECMapKey(version, lvl)].ECCodewordsPerBlock
	coefficients := make(QrPolynomial, 1)

	for i := 1; i <= degree; i++ {
		factorCoefficients := make(QrPolynomial, 2)
		factorCoefficients[0] = i - 1
		coefficients = e.multiplyPolynomials(coefficients, factorCoefficients)
	}

	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = convertExponentToValue(coefficients[i])
	}

	return coefficients
}

func (e *Encoder) GetErrorCorrectionCodewords(encoded string, version QrVersion, lvl QrErrCorrectionLvl) QrPolynomial {
	messagePolynomial := e.GetMessagePolynomial(encoded)
	generatorPolynomial := e.GetGeneratorPolynomial(version, lvl)
	numErrCorrCodewords := ERROR_CORRECTION_INFO[e.getECMapKey(version, lvl)].ECCodewordsPerBlock

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

			currValue := convertExponentToValue(exponent)
			prevValue := e.getValueFromResultExponents(result, i+j)

			exponent = convertValueToExponent(currValue ^ prevValue)
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
		currentDivisor = e.multiplyPolynomialByScalar(currentDivisor, dividentLeadTerm)
		remainder = e.xorPolynomials(currentDivisor, divident, leadTermIndex)

		divident = remainder
	}

	return divident
}

func (e *Encoder) multiplyPolynomialByScalar(polynomial QrPolynomial, scalar int) QrPolynomial {
	scalarAlphaValue := convertValueToExponent(scalar)

	var termAlphaValue int
	result := make(QrPolynomial, len(polynomial))

	currIndex := 0
	isLeadTermFound := false

	for i := len(polynomial) - 1; i >= 0; i-- {
		termAlphaValue = convertValueToExponent(polynomial[i])

		if termAlphaValue == 0 {
			isLeadTermFound = true
			continue
		}

		if isLeadTermFound && currIndex < 10 {
			result[i] = convertExponentToValue((termAlphaValue + scalarAlphaValue) % (QR_GALOIS_ORDER - 1))
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
	return convertExponentToValue(result[index])
}

// Utilities

func computeLogAntilogTables() {
	for i := 0; i < QR_GALOIS_ORDER; i++ {
		QR_GALOIS_LOG_TABLE[i] = computeAlphaToPower(2, i)
		QR_GALOIS_ANTILOG_TABLE[QR_GALOIS_LOG_TABLE[i]] = i
	}
	QR_GALOIS_ANTILOG_TABLE[1] = 0
}

func convertValueToExponent(value int) int {
	if value < 0 || value > len(QR_GALOIS_ANTILOG_TABLE)-1 {
		return 0
	}
	return QR_GALOIS_ANTILOG_TABLE[value]
}

func convertExponentToValue(exponent int) int {
	if exponent < 0 || exponent > len(QR_GALOIS_LOG_TABLE)-1 {
		return 0
	}
	return QR_GALOIS_LOG_TABLE[exponent]
}

func computeAlphaToPower(alpha, power int) int {
	result := 1

	for i := 0; i < power; i++ {
		prod := result * alpha

		if prod >= QR_GALOIS_ORDER {
			prod = prod ^ QR_GALOIS_MOD_VALUE
		}

		result = prod
	}
	return result
}

func main() {
	computeLogAntilogTables()
	e := NewEncoder()
	input := "HELLO WORLD"
	mode, _ := e.GetMode(input)
	version, _ := e.GetVersion(input, mode, MEDIUM)
	encoded, _ := e.Encode(input, MEDIUM)
	augmented := e.AugmentEncodedInput(encoded, version, MEDIUM)

	e.GetErrorCorrectionCodewords(augmented, VERSION_1, MEDIUM)
}
