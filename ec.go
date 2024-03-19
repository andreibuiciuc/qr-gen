package qr

import (
	"strconv"
)

const qrGaloisOrder = 256

type errorCorrector struct{}
type QrPolynomial []int

func NewErrorCorrector() *errorCorrector {
	computeLogAntilogTables()
	return &errorCorrector{}
}

func (ec *errorCorrector) getMessagePolynomial(encoded string) QrPolynomial {
	codewords := splitInGroups(encoded, codewordSize)
	coefficients := make(QrPolynomial, len(codewords))

	for index, codeword := range codewords {
		decimalValue, _ := strconv.ParseInt(codeword, 2, 64)
		coefficients[len(coefficients)-index-1] = int(decimalValue)
	}

	return coefficients
}

func (ec *errorCorrector) getGeneratorPolynomial(v int, lvl rune) QrPolynomial {
	degree := ecInfo[getECMappingKey(v, string(lvl))].ECCodewordsPerBlock
	coefficients := make(QrPolynomial, 1)

	for i := 1; i <= degree; i++ {
		factorCoefficients := make(QrPolynomial, 2)
		factorCoefficients[0] = i - 1
		coefficients = ec.multiplyPolynomials(coefficients, factorCoefficients)
	}

	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = convertExponentToValue(coefficients[i])
	}

	return coefficients
}

func (ec *errorCorrector) getErrorCorrectionCodewords(encoded string, v int, lvl rune) QrPolynomial {
	messagePolynomial := ec.getMessagePolynomial(encoded)
	generatorPolynomial := ec.getGeneratorPolynomial(v, lvl)
	numErrCorrCodewords := ecInfo[getECMappingKey(v, string(lvl))].ECCodewordsPerBlock

	divisionSteps := len(messagePolynomial)
	messagePolynomial = ec.expandPolynomial(messagePolynomial, numErrCorrCodewords)
	generatorPolynomial = ec.expandPolynomial(generatorPolynomial, len(messagePolynomial)-len(generatorPolynomial))

	errCorrCodewords := ec.dividePolynomials(messagePolynomial, generatorPolynomial, divisionSteps, numErrCorrCodewords)
	return errCorrCodewords[0:numErrCorrCodewords]
}

func (ec *errorCorrector) multiplyPolynomials(firstPoly, secondPoly QrPolynomial) QrPolynomial {
	degreeFirstPoly, degreeSecondPoly := len(firstPoly)-1, len(secondPoly)-1
	result := ec.initializeResultAsExponents(degreeFirstPoly + degreeSecondPoly)

	for i := 0; i <= degreeFirstPoly; i++ {
		for j := 0; j <= degreeSecondPoly; j++ {
			exponent := firstPoly[i] + secondPoly[j]

			if exponent >= qrGaloisOrder {
				exponent = exponent % (qrGaloisOrder - 1)
			}

			currValue := convertExponentToValue(exponent)
			prevValue := ec.getValueFromResultExponents(result, i+j)

			exponent = convertValueToExponent(currValue ^ prevValue)
			result[i+j] = exponent
		}
	}

	return result
}

func (ec *errorCorrector) dividePolynomials(divident, divisor QrPolynomial, steps, remainderLen int) QrPolynomial {
	var remainder QrPolynomial
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

func (e *errorCorrector) xorPolynomials(firstPolynomial, secondPolynomial QrPolynomial, leadTermIndex int) QrPolynomial {
	result := make(QrPolynomial, len(secondPolynomial))

	for i := len(secondPolynomial) - 1; i >= 0; i-- {
		if i < leadTermIndex {
			result[i] = firstPolynomial[i] ^ secondPolynomial[i]
		}
	}

	return result
}

func (e *errorCorrector) multiplyPolynomialByScalar(polynomial QrPolynomial, scalar, remainderLen int) QrPolynomial {
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

		if isLeadTermFound && currIndex < remainderLen {
			result[i] = convertExponentToValue((termAlphaValue + scalarAlphaValue) % (qrGaloisOrder - 1))
			currIndex++
		}
	}

	return result
}

func (e *errorCorrector) getLeadingTerm(polynomial QrPolynomial) (int, int) {
	for i := len(polynomial) - 1; i >= 0; i-- {
		if polynomial[i] != 0 {
			return polynomial[i], i
		}
	}
	return -1, -1
}

func (ec *errorCorrector) expandPolynomial(polynomial QrPolynomial, n int) QrPolynomial {
	expandedPolynomial := make(QrPolynomial, len(polynomial)+n)
	copy(expandedPolynomial[n:], polynomial)
	return expandedPolynomial
}

func (e *errorCorrector) shiftPolynomial(polynomial QrPolynomial, unit int) QrPolynomial {
	return append(polynomial[unit:], 0)
}

func (ec *errorCorrector) initializeResultAsExponents(degree int) QrPolynomial {
	exponents := make([]int, degree+1)

	for i := 0; i < degree; i++ {
		exponents[i] = -1
	}

	return exponents
}

func (e *errorCorrector) getValueFromResultExponents(result QrPolynomial, index int) int {
	if index < 0 || index > len(result)-1 || result[index] == -1 {
		return 0
	}
	return convertExponentToValue(result[index])
}
