// Package aztec can create Aztec Code barcodes
package aztec

import (
	"fmt"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

const (
	DEFAULT_EC_PERCENT  = 33
	DEFAULT_LAYERS      = 0
	max_nb_bits         = 32
	max_nb_bits_compact = 4
)

var (
	word_size = []int{
		4, 6, 6, 8, 8, 8, 8, 8, 8, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
		12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	}
)

func totalBitsInLayer(layers int, compact bool) int {
	tmp := 112
	if compact {
		tmp = 88
	}
	return (tmp + 16*layers) * layers
}

func stuffBits(bits *utils.BitList, wordSize int) *utils.BitList {
	out := new(utils.BitList)
	n := bits.Len()
	mask := (1 << uint(wordSize)) - 2
	for i := 0; i < n; i += wordSize {
		word := 0
		for j := 0; j < wordSize; j++ {
			if i+j >= n || bits.GetBit(i+j) {
				word |= 1 << uint(wordSize-1-j)
			}
		}
		if (word & mask) == mask {
			out.AddBits(word&mask, byte(wordSize))
			i--
		} else if (word & mask) == 0 {
			out.AddBits(word|1, byte(wordSize))
			i--
		} else {
			out.AddBits(word, byte(wordSize))
		}
	}
	return out
}

func generateModeMessage(compact bool, layers, messageSizeInWords int) *utils.BitList {
	modeMessage := new(utils.BitList)
	if compact {
		modeMessage.AddBits(layers-1, 2)
		modeMessage.AddBits(messageSizeInWords-1, 6)
		modeMessage = generateCheckWords(modeMessage, 28, 4)
	} else {
		modeMessage.AddBits(layers-1, 5)
		modeMessage.AddBits(messageSizeInWords-1, 11)
		modeMessage = generateCheckWords(modeMessage, 40, 4)
	}
	return modeMessage
}

func drawModeMessage(matrix *aztecCode, compact bool, matrixSize int, modeMessage *utils.BitList) {
	center := matrixSize / 2
	if compact {
		for i := 0; i < 7; i++ {
			offset := center - 3 + i
			if modeMessage.GetBit(i) {
				matrix.set(offset, center-5)
			}
			if modeMessage.GetBit(i + 7) {
				matrix.set(center+5, offset)
			}
			if modeMessage.GetBit(20 - i) {
				matrix.set(offset, center+5)
			}
			if modeMessage.GetBit(27 - i) {
				matrix.set(center-5, offset)
			}
		}
	} else {
		for i := 0; i < 10; i++ {
			offset := center - 5 + i + i/5
			if modeMessage.GetBit(i) {
				matrix.set(offset, center-7)
			}
			if modeMessage.GetBit(i + 10) {
				matrix.set(center+7, offset)
			}
			if modeMessage.GetBit(29 - i) {
				matrix.set(offset, center+7)
			}
			if modeMessage.GetBit(39 - i) {
				matrix.set(center-7, offset)
			}
		}
	}
}

func drawBullsEye(matrix *aztecCode, center, size int) {
	for i := 0; i < size; i += 2 {
		for j := center - i; j <= center+i; j++ {
			matrix.set(j, center-i)
			matrix.set(j, center+i)
			matrix.set(center-i, j)
			matrix.set(center+i, j)
		}
	}
	matrix.set(center-size, center-size)
	matrix.set(center-size+1, center-size)
	matrix.set(center-size, center-size+1)
	matrix.set(center+size, center-size)
	matrix.set(center+size, center-size+1)
	matrix.set(center+size, center+size-1)
}

// Encode returns an aztec barcode with the given content
func Encode(data []byte, minECCPercent int, userSpecifiedLayers int) (barcode.Barcode, error) {
	bits := highlevelEncode(data)
	eccBits := ((bits.Len() * minECCPercent) / 100) + 11
	totalSizeBits := bits.Len() + eccBits
	var layers, TotalBitsInLayer, wordSize int
	var compact bool
	var stuffedBits *utils.BitList
	if userSpecifiedLayers != DEFAULT_LAYERS {
		compact = userSpecifiedLayers < 0
		if compact {
			layers = -userSpecifiedLayers
		} else {
			layers = userSpecifiedLayers
		}
		if (compact && layers > max_nb_bits_compact) || (!compact && layers > max_nb_bits) {
			return nil, fmt.Errorf("Illegal value %d for layers", userSpecifiedLayers)
		}
		TotalBitsInLayer = totalBitsInLayer(layers, compact)
		wordSize = word_size[layers]
		usableBitsInLayers := TotalBitsInLayer - (TotalBitsInLayer % wordSize)
		stuffedBits = stuffBits(bits, wordSize)
		if stuffedBits.Len()+eccBits > usableBitsInLayers {
			return nil, fmt.Errorf("Data to large for user specified layer")
		}
		if compact && stuffedBits.Len() > wordSize*64 {
			return nil, fmt.Errorf("Data to large for user specified layer")
		}
	} else {
		wordSize = 0
		stuffedBits = nil
		// We look at the possible table sizes in the order Compact1, Compact2, Compact3,
		// Compact4, Normal4,...  Normal(i) for i < 4 isn't typically used since Compact(i+1)
		// is the same size, but has more data.
		for i := 0; ; i++ {
			if i > max_nb_bits {
				return nil, fmt.Errorf("Data too large for an aztec code")
			}
			compact = i <= 3
			layers = i
			if compact {
				layers = i + 1
			}
			TotalBitsInLayer = totalBitsInLayer(layers, compact)
			if totalSizeBits > TotalBitsInLayer {
				continue
			}
			// [Re]stuff the bits if this is the first opportunity, or if the
			// wordSize has changed
			if wordSize != word_size[layers] {
				wordSize = word_size[layers]
				stuffedBits = stuffBits(bits, wordSize)
			}
			usableBitsInLayers := TotalBitsInLayer - (TotalBitsInLayer % wordSize)
			if compact && stuffedBits.Len() > wordSize*64 {
				// Compact format only allows 64 data words, though C4 can hold more words than that
				continue
			}
			if stuffedBits.Len()+eccBits <= usableBitsInLayers {
				break
			}
		}
	}
	messageBits := generateCheckWords(stuffedBits, TotalBitsInLayer, wordSize)
	messageSizeInWords := stuffedBits.Len() / wordSize
	modeMessage := generateModeMessage(compact, layers, messageSizeInWords)

	// allocate symbol
	var baseMatrixSize int
	if compact {
		baseMatrixSize = 11 + layers*4
	} else {
		baseMatrixSize = 14 + layers*4
	}
	alignmentMap := make([]int, baseMatrixSize)
	var matrixSize int

	if compact {
		// no alignment marks in compact mode, alignmentMap is a no-op
		matrixSize = baseMatrixSize
		for i := 0; i < len(alignmentMap); i++ {
			alignmentMap[i] = i
		}
	} else {
		matrixSize = baseMatrixSize + 1 + 2*((baseMatrixSize/2-1)/15)
		origCenter := baseMatrixSize / 2
		center := matrixSize / 2
		for i := 0; i < origCenter; i++ {
			newOffset := i + i/15
			alignmentMap[origCenter-i-1] = center - newOffset - 1
			alignmentMap[origCenter+i] = center + newOffset + 1
		}
	}
	code := newAztecCode(matrixSize)
	code.content = data

	// draw data bits
	for i, rowOffset := 0, 0; i < layers; i++ {
		rowSize := (layers - i) * 4
		if compact {
			rowSize += 9
		} else {
			rowSize += 12
		}

		for j := 0; j < rowSize; j++ {
			columnOffset := j * 2
			for k := 0; k < 2; k++ {
				if messageBits.GetBit(rowOffset + columnOffset + k) {
					code.set(alignmentMap[i*2+k], alignmentMap[i*2+j])
				}
				if messageBits.GetBit(rowOffset + rowSize*2 + columnOffset + k) {
					code.set(alignmentMap[i*2+j], alignmentMap[baseMatrixSize-1-i*2-k])
				}
				if messageBits.GetBit(rowOffset + rowSize*4 + columnOffset + k) {
					code.set(alignmentMap[baseMatrixSize-1-i*2-k], alignmentMap[baseMatrixSize-1-i*2-j])
				}
				if messageBits.GetBit(rowOffset + rowSize*6 + columnOffset + k) {
					code.set(alignmentMap[baseMatrixSize-1-i*2-j], alignmentMap[i*2+k])
				}
			}
		}
		rowOffset += rowSize * 8
	}

	// draw mode message
	drawModeMessage(code, compact, matrixSize, modeMessage)

	// draw alignment marks
	if compact {
		drawBullsEye(code, matrixSize/2, 5)
	} else {
		drawBullsEye(code, matrixSize/2, 7)
		for i, j := 0, 0; i < baseMatrixSize/2-1; i, j = i+15, j+16 {
			for k := (matrixSize / 2) & 1; k < matrixSize; k += 2 {
				code.set(matrixSize/2-j, k)
				code.set(matrixSize/2+j, k)
				code.set(k, matrixSize/2-j)
				code.set(k, matrixSize/2+j)
			}
		}
	}
	return code, nil
}
