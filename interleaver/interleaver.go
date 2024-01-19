package interleaver

import (
	"qr/qr-gen/ec"
	"qr/qr-gen/util"
	"qr/qr-gen/versioner"
	"strconv"
	"strings"
)

type Interleaver interface {
	GetFinalMessage(encoded string, version versioner.QrVersion, lvl versioner.QrEcLevel) string
}

type QrInterleaver struct{}

func New() Interleaver {
	// TODO: Extract this in another layer, to avoid redundant computation
	util.ComputeLogAntilogTables()
	return &QrInterleaver{}
}

func (i *QrInterleaver) GetFinalMessage(inputCodewords string, version versioner.QrVersion, lvl versioner.QrEcLevel) string {
	var codewords string
	ec := ec.New()

	if i.isInterleavingNecessary(version, lvl) {
		codewords = i.handleInterleaveProcess(version, lvl, inputCodewords)
	} else {
		errCorrCodewords := ec.GetErrorCorrectionCodewords(inputCodewords, version, lvl)
		codewords = inputCodewords + util.ConvertIntListToCodewords(errCorrCodewords)
	}

	return util.PadRight(codewords, "0", len(codewords)+QR_REMAINDER_BITS[version])
}

func (i *QrInterleaver) isInterleavingNecessary(version versioner.QrVersion, lvl versioner.QrEcLevel) bool {
	return util.QrEcInfo[util.GetECMappingKey(int(version), string(lvl))].NumBlocksGroup2 != 0
}

func (i *QrInterleaver) handleInterleaveProcess(version versioner.QrVersion, lvl versioner.QrEcLevel, codewords string) string {
	interleavedDataCodewords, dataBlocks := i.interleaveDataCodewords(version, lvl, codewords)
	interleavedECCodewords := i.interleaveErrCorrCodewords(version, lvl, dataBlocks)

	interleavedDataBinary := util.ConvertIntListToBin(interleavedDataCodewords)
	interleavedECBinary := util.ConvertIntListToBin(interleavedECCodewords)

	return strings.Join(interleavedDataBinary, "") + strings.Join(interleavedECBinary, "")
}

func (i *QrInterleaver) interleaveDataCodewords(version versioner.QrVersion, lvl versioner.QrEcLevel, encoded string) ([]int, [][]int) {
	key := util.GetECMappingKey(int(version), string(lvl))

	group1Size := util.QrEcInfo[key].NumBlocksGroup1
	group1BlockSize := util.QrEcInfo[key].DataCodeworkdsInGroup1Block
	group1Codewords := encoded[:group1Size*group1BlockSize*util.QrCodewordSize]
	group1Blocks := i.getBlocksOfCodewords(group1Codewords, group1Size, group1BlockSize)

	group2Size := util.QrEcInfo[key].NumBlocksGroup2
	group2BlockSize := util.QrEcInfo[key].DataCodewordsInGroup2Block
	group2Codewords := encoded[group1Size*group1BlockSize*util.QrCodewordSize:]
	group2Blocks := i.getBlocksOfCodewords(group2Codewords, group2Size, group2BlockSize)

	dataBlocks := append(group1Blocks, group2Blocks...)
	return i.interleaveCodewords(dataBlocks, util.Max(group1BlockSize, group2BlockSize)), dataBlocks
}

func (i *QrInterleaver) interleaveErrCorrCodewords(version versioner.QrVersion, lvl versioner.QrEcLevel, dataBlocks [][]int) []int {
	blocks := make([][]int, len(dataBlocks))
	ec := ec.New()

	for i, block := range dataBlocks {
		encoded := strings.Join(util.ConvertIntListToBin(block), "")
		blocks[i] = ec.GetErrorCorrectionCodewords(encoded, version, lvl)
	}

	return i.interleaveCodewords(blocks, util.QrEcInfo[util.GetECMappingKey(int(version), string(lvl))].ECCodewordsPerBlock)
}

func (i *QrInterleaver) getBlocksOfCodewords(input string, blocksCount int, blockSize int) [][]int {
	blocks := make([][]int, blocksCount)

	for i := 0; i < blocksCount; i++ {
		currBlock := input[:blockSize*util.QrCodewordSize]
		currBlockSlice := make([]int, blockSize)

		for j := 0; j < blockSize; j++ {
			bin := currBlock[:util.QrCodewordSize]
			value, _ := strconv.ParseInt(bin, 2, 64)
			currBlockSlice[j] = int(value)
			currBlock = currBlock[util.QrCodewordSize:]
		}

		blocks[i] = currBlockSlice
		input = input[blockSize*util.QrCodewordSize:]
	}

	return blocks
}

func (i *QrInterleaver) interleaveCodewords(blocks [][]int, length int) []int {
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

var QR_REMAINDER_BITS = map[versioner.QrVersion]int{
	1:  0,
	2:  7,
	3:  7,
	4:  7,
	5:  7,
	6:  7,
	7:  0,
	8:  0,
	9:  0,
	10: 0,
	11: 0,
	12: 0,
	13: 0,
	14: 3,
	15: 3,
	16: 3,
	17: 3,
	18: 3,
	19: 3,
	20: 3,
	21: 4,
	22: 4,
	23: 4,
	24: 4,
	25: 4,
	26: 4,
	27: 4,
	28: 3,
	29: 3,
	30: 3,
	31: 3,
	32: 3,
	33: 3,
	34: 3,
	35: 0,
	36: 0,
	37: 0,
	38: 0,
	39: 0,
	40: 0,
}
