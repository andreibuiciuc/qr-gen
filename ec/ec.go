package ec

import (
	"qr/qr-gen/util"
	"qr/qr-gen/versioner"
	"strconv"
)

type ErrorCorrector interface {
	GetMessagePolynomial(encoded string) QrPolynomial
	GetGeneratorPolynomial(version versioner.QrVersion, lvl versioner.QrEcLevel) QrPolynomial
	GetErrorCorrectionCodewords(encoded string, version versioner.QrVersion, lvl versioner.QrEcLevel) QrPolynomial
}

type QrErrorCorrector struct{}

type QrPolynomial []int

func New() ErrorCorrector {
	util.ComputeLogAntilogTables()
	return &QrErrorCorrector{}
}

func (ec *QrErrorCorrector) GetMessagePolynomial(encoded string) QrPolynomial {
	codewords := util.SplitInGroups(encoded, util.QrCodewordSize)
	coefficients := make(QrPolynomial, len(codewords))

	for index, codeword := range codewords {
		decimalValue, _ := strconv.ParseInt(codeword, 2, 64)
		coefficients[len(coefficients)-index-1] = int(decimalValue)
	}

	return coefficients
}

func (ec *QrErrorCorrector) GetGeneratorPolynomial(version versioner.QrVersion, lvl versioner.QrEcLevel) QrPolynomial {
	degree := util.QrEcInfo[util.GetECMappingKey(int(version), string(lvl))].ECCodewordsPerBlock
	coefficients := make(QrPolynomial, 1)

	for i := 1; i <= degree; i++ {
		factorCoefficients := make(QrPolynomial, 2)
		factorCoefficients[0] = i - 1
		coefficients = ec.multiplyPolynomials(coefficients, factorCoefficients)
	}

	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = util.ConvertExponentToValue(coefficients[i])
	}

	return coefficients
}

func (ec *QrErrorCorrector) GetErrorCorrectionCodewords(encoded string, version versioner.QrVersion, lvl versioner.QrEcLevel) QrPolynomial {
	messagePolynomial := ec.GetMessagePolynomial(encoded)
	generatorPolynomial := ec.GetGeneratorPolynomial(version, lvl)
	numErrCorrCodewords := util.QrEcInfo[util.GetECMappingKey(int(version), string(lvl))].ECCodewordsPerBlock

	divisionSteps := len(messagePolynomial)
	messagePolynomial = ec.expandPolynomial(messagePolynomial, numErrCorrCodewords)
	generatorPolynomial = ec.expandPolynomial(generatorPolynomial, len(messagePolynomial)-len(generatorPolynomial))

	errCorrCodewords := ec.dividePolynomials(messagePolynomial, generatorPolynomial, divisionSteps, numErrCorrCodewords)
	return errCorrCodewords[0:numErrCorrCodewords]
}

func (ec *QrErrorCorrector) multiplyPolynomials(firstPoly, secondPoly QrPolynomial) QrPolynomial {
	degreeFirstPoly, degreeSecondPoly := len(firstPoly)-1, len(secondPoly)-1
	result := ec.initializeResultAsExponents(degreeFirstPoly + degreeSecondPoly)

	for i := 0; i <= degreeFirstPoly; i++ {
		for j := 0; j <= degreeSecondPoly; j++ {
			exponent := firstPoly[i] + secondPoly[j]

			if exponent >= qrGaloisOrder {
				exponent = exponent % (qrGaloisOrder - 1)
			}

			currValue := util.ConvertExponentToValue(exponent)
			prevValue := ec.getValueFromResultExponents(result, i+j)

			exponent = util.ConvertValueToExponent(currValue ^ prevValue)
			result[i+j] = exponent
		}
	}

	return result
}

func (ec *QrErrorCorrector) dividePolynomials(divident, divisor QrPolynomial, steps, remainderLen int) QrPolynomial {
	remainder := make(QrPolynomial, remainderLen)

	var currentDivisor QrPolynomial
	copy(currentDivisor, divisor)

	for i := 0; i < steps; i++ {
		dividentLeadTerm, leadTermIndex := ec.getLeadingTerm(divident)

		currentDivisor = divisor
		currentDivisor = ec.shiftPolynomial(currentDivisor, i)
		currentDivisor = ec.multiplyPolynomialByScalar(currentDivisor, dividentLeadTerm, remainderLen)
		remainder = ec.xorPolynomials(currentDivisor, divident, leadTermIndex)

		divident = remainder
	}

	return divident
}

func (e *QrErrorCorrector) xorPolynomials(firstPolynomial, secondPolynomial QrPolynomial, leadTermIndex int) QrPolynomial {
	result := make(QrPolynomial, len(secondPolynomial))

	for i := len(secondPolynomial) - 1; i >= 0; i-- {
		if i < leadTermIndex {
			result[i] = firstPolynomial[i] ^ secondPolynomial[i]
		}
	}

	return result
}

func (e *QrErrorCorrector) multiplyPolynomialByScalar(polynomial QrPolynomial, scalar, remainderLen int) QrPolynomial {
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
			result[i] = util.ConvertExponentToValue((termAlphaValue + scalarAlphaValue) % (qrGaloisOrder - 1))
			currIndex++
		}
	}

	return result
}

func (e *QrErrorCorrector) getLeadingTerm(polynomial QrPolynomial) (int, int) {
	for i := len(polynomial) - 1; i >= 0; i-- {
		if polynomial[i] != 0 {
			return polynomial[i], i
		}
	}
	return -1, -1
}

func (ec *QrErrorCorrector) expandPolynomial(polynomial QrPolynomial, n int) QrPolynomial {
	expandedPolynomial := make(QrPolynomial, len(polynomial)+n)
	copy(expandedPolynomial[n:], polynomial)
	return expandedPolynomial
}

func (e *QrErrorCorrector) shiftPolynomial(polynomial QrPolynomial, unit int) QrPolynomial {
	return append(polynomial[unit:], 0)
}

func (ec *QrErrorCorrector) initializeResultAsExponents(degree int) QrPolynomial {
	exponents := make([]int, degree+1)

	for i := 0; i < degree; i++ {
		exponents[i] = -1
	}

	return exponents
}

func (e *QrErrorCorrector) getValueFromResultExponents(result QrPolynomial, index int) int {
	if index < 0 || index > len(result)-1 || result[index] == -1 {
		return 0
	}
	return util.ConvertExponentToValue(result[index])
}

const qrGaloisOrder = 256
const qrGaloisModTerm = 285
