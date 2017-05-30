package aztec

import (
	"github.com/boombuler/barcode/utils"
)

func bitsToWords(stuffedBits *utils.BitList, wordSize int, wordCount int) []int {
	message := make([]int, wordCount)

	for i := 0; i < wordCount; i++ {
		value := 0
		for j := 0; j < wordSize; j++ {
			if stuffedBits.GetBit(i*wordSize + j) {
				value |= (1 << uint(wordSize-j-1))
			}
		}
		message[i] = value
	}
	return message
}

func generateCheckWords(bits *utils.BitList, totalBits, wordSize int) *utils.BitList {
	rs := utils.NewReedSolomonEncoder(getGF(wordSize))

	// bits is guaranteed to be a multiple of the wordSize, so no padding needed
	messageWordCount := bits.Len() / wordSize
	totalWordCount := totalBits / wordSize
	eccWordCount := totalWordCount - messageWordCount

	messageWords := bitsToWords(bits, wordSize, messageWordCount)
	eccWords := rs.Encode(messageWords, eccWordCount)
	startPad := totalBits % wordSize

	messageBits := new(utils.BitList)
	messageBits.AddBits(0, byte(startPad))

	for _, messageWord := range messageWords {
		messageBits.AddBits(messageWord, byte(wordSize))
	}
	for _, eccWord := range eccWords {
		messageBits.AddBits(eccWord, byte(wordSize))
	}
	return messageBits
}

func getGF(wordSize int) *utils.GaloisField {
	switch wordSize {
	case 4:
		return utils.NewGaloisField(0x13, 16, 1)
	case 6:
		return utils.NewGaloisField(0x43, 64, 1)
	case 8:
		return utils.NewGaloisField(0x012D, 256, 1)
	case 10:
		return utils.NewGaloisField(0x409, 1024, 1)
	case 12:
		return utils.NewGaloisField(0x1069, 4096, 1)
	default:
		return nil
	}
}
