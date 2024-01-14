package versioner

import (
	"fmt"
	"qr/qr-gen/util"
	"regexp"
	"strconv"
)

type QrVersion int
type QrEcLevel string
type QrMode string
type QrModeIndicator string

type Versioner interface {
	GetMode(s string) (QrMode, error)
	GetVersion(s string, mode QrMode, lvl QrEcLevel) (QrVersion, error)
	GetModeIndicator(mode QrMode) string
	GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error)
}

type QrVersioner struct{}

func New() Versioner {
	return &QrVersioner{}
}

func (v *QrVersioner) GetMode(s string) (QrMode, error) {
	if matched, _ := regexp.MatchString(qrModeRegexex[QrNumericMode], s); matched {
		return QrNumericMode, nil
	}

	if matched, _ := regexp.MatchString(qrModeRegexex[QrAlphanumericMode], s); matched {
		return QrAlphanumericMode, nil
	}

	if matched, _ := regexp.MatchString(qrModeRegexex[QrByteMode], s); matched {
		return QrByteMode, nil
	}

	return QrMode(""), fmt.Errorf("Invalid input pattern")
}

func (v *QrVersioner) GetVersion(s string, mode QrMode, lvl QrEcLevel) (QrVersion, error) {
	version := 1

	for version <= len(qrCapacities) {
		if len(s) <= qrCapacities[QrVersion(version)][lvl][qrModeIndices[mode]] {
			return QrVersion(version), nil
		}
		version += 1
	}

	return QrVersion(-1), fmt.Errorf("Cannot compute QR version")
}

func (v *QrVersioner) GetModeIndicator(mode QrMode) string {
	return qrModeIndicators[mode]
}

func (v *QrVersioner) GetCountIndicator(s string, version QrVersion, mode QrMode) (string, error) {
	sLenBin := strconv.FormatInt(int64(len(s)), 2)
	return util.PadLeft(sLenBin, "0", qrCountIndLengths[mode]), nil
}

const (
	// Extend this to support also Kanji mode
	QrNumericMode      QrMode = "numeric"
	QrAlphanumericMode QrMode = "alphanumeric"
	QrByteMode         QrMode = "byte"
)

const (
	QrEcLow      QrEcLevel = "L"
	QrEcMedium   QrEcLevel = "M"
	QrEcQuartile QrEcLevel = "Q"
	QrECHigh     QrEcLevel = "H"
)

const (
	qrNumericInd      string = "0001"
	qrAlphanumericInd string = "0010"
	qrByteInd         string = "0100"
)

var qrModeRegexex = map[QrMode]string{
	QrNumericMode:      "^\\d+$",
	QrAlphanumericMode: "^[\\dA-Z $%*+\\-./:]+$",
	QrByteMode:         "^[\\x00-\\xff]+$",
}

var qrModeIndices = map[QrMode]int{
	QrNumericMode:      0,
	QrAlphanumericMode: 1,
	QrByteMode:         2,
}

var qrModeIndicators = map[QrMode]string{
	QrNumericMode:      qrNumericInd,
	QrAlphanumericMode: qrAlphanumericInd,
	QrByteMode:         qrByteInd,
}

var qrCountIndLengths = map[QrMode]int{
	QrNumericMode:      10,
	QrAlphanumericMode: 9,
	QrByteMode:         8,
}

var qrCapacities = map[QrVersion]map[QrEcLevel][]int{
	1: {
		QrEcLow:      {41, 25, 17},
		QrEcMedium:   {34, 20, 14},
		QrEcQuartile: {27, 16, 11},
		QrECHigh:     {17, 10, 7},
	},
	2: {
		QrEcLow:      {77, 47, 32},
		QrEcMedium:   {63, 38, 26},
		QrEcQuartile: {48, 29, 20},
		QrECHigh:     {34, 20, 14},
	},
	3: {
		QrEcLow:      {127, 77, 53},
		QrEcMedium:   {101, 61, 42},
		QrEcQuartile: {77, 47, 32},
		QrECHigh:     {58, 35, 24},
	},
	4: {
		QrEcLow:      {187, 114, 78},
		QrEcMedium:   {149, 90, 62},
		QrEcQuartile: {111, 67, 46},
		QrECHigh:     {82, 50, 34},
	},
	5: {
		QrEcLow:      {255, 154, 106},
		QrEcMedium:   {202, 122, 84},
		QrEcQuartile: {144, 87, 60},
		QrECHigh:     {106, 64, 44},
	},
}
