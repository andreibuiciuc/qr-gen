package qr

import (
	"strconv"
	"strings"
)

type interleaver struct{}

var remainderBits = map[int]int{
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

func newInterleaver() *interleaver {
	computeLogAntilogTables()
	return &interleaver{}
}

func (i *interleaver) getFinalMessage(inputCodewords string, v int, lvl rune) string {
	var codewords string
	ec := NewErrorCorrector()

	if i.isInterleavingNecessary(v, lvl) {
		codewords = i.handleInterleaveProcess(v, lvl, inputCodewords)
	} else {
		errCorrCodewords := ec.getErrorCorrectionCodewords(inputCodewords, v, lvl)
		codewords = inputCodewords + convertIntListToCodewords(errCorrCodewords)
	}

	return padRight(codewords, "0", len(codewords)+remainderBits[v])
}

func (i *interleaver) isInterleavingNecessary(v int, lvl rune) bool {
	return ecInfo[getECMappingKey(v, string(lvl))].NumBlocksGroup2 != 0
}

func (i *interleaver) handleInterleaveProcess(v int, lvl rune, codewords string) string {
	interleavedDataCodewords, dataBlocks := i.interleaveDataCodewords(v, lvl, codewords)
	interleavedECCodewords := i.interleaveErrCorrCodewords(v, lvl, dataBlocks)

	interleavedDataBinary := convertIntListToBin(interleavedDataCodewords)
	interleavedECBinary := convertIntListToBin(interleavedECCodewords)

	return strings.Join(interleavedDataBinary, "") + strings.Join(interleavedECBinary, "")
}

func (i *interleaver) interleaveDataCodewords(v int, lvl rune, encoded string) ([]int, [][]int) {
	key := getECMappingKey(v, string(lvl))

	group1Size := ecInfo[key].NumBlocksGroup1
	group1BlockSize := ecInfo[key].DataCodeworkdsInGroup1Block
	group1Codewords := encoded[:group1Size*group1BlockSize*codewordSize]
	group1Blocks := i.getBlocksOfCodewords(group1Codewords, group1Size, group1BlockSize)

	group2Size := ecInfo[key].NumBlocksGroup2
	group2BlockSize := ecInfo[key].DataCodewordsInGroup2Block
	group2Codewords := encoded[group1Size*group1BlockSize*codewordSize:]
	group2Blocks := i.getBlocksOfCodewords(group2Codewords, group2Size, group2BlockSize)

	dataBlocks := append(group1Blocks, group2Blocks...)
	return i.interleaveCodewords(dataBlocks, max(group1BlockSize, group2BlockSize)), dataBlocks
}

func (i *interleaver) interleaveErrCorrCodewords(v int, lvl rune, dataBlocks [][]int) []int {
	blocks := make([][]int, len(dataBlocks))
	ec := NewErrorCorrector()

	for i, block := range dataBlocks {
		encoded := strings.Join(convertIntListToBin(block), "")
		blocks[i] = ec.getErrorCorrectionCodewords(encoded, v, lvl)
	}

	return i.interleaveCodewords(blocks, ecInfo[getECMappingKey(v, string(lvl))].ECCodewordsPerBlock)
}

func (i *interleaver) getBlocksOfCodewords(input string, blocksCount int, blockSize int) [][]int {
	blocks := make([][]int, blocksCount)

	for i := 0; i < blocksCount; i++ {
		currBlock := input[:blockSize*codewordSize]
		currBlockSlice := make([]int, blockSize)

		for j := 0; j < blockSize; j++ {
			bin := currBlock[:codewordSize]
			value, _ := strconv.ParseInt(bin, 2, 64)
			currBlockSlice[j] = int(value)
			currBlock = currBlock[codewordSize:]
		}

		blocks[i] = currBlockSlice
		input = input[blockSize*codewordSize:]
	}

	return blocks
}

func (i *interleaver) interleaveCodewords(blocks [][]int, length int) []int {
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
