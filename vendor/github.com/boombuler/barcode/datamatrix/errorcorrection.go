package datamatrix

import (
	"github.com/boombuler/barcode/utils"
)

type errorCorrection struct {
	rs *utils.ReedSolomonEncoder
}

var ec *errorCorrection = newErrorCorrection()

func newErrorCorrection() *errorCorrection {
	gf := utils.NewGaloisField(301, 256, 1)

	return &errorCorrection{utils.NewReedSolomonEncoder(gf)}
}

func (ec *errorCorrection) calcECC(data []byte, size *dmCodeSize) []byte {
	dataSize := len(data)
	// make some space for error correction codes
	data = append(data, make([]byte, size.ECCCount)...)

	for block := 0; block < size.BlockCount; block++ {
		dataCnt := size.DataCodewordsForBlock(block)

		buff := make([]int, dataCnt)
		// copy the data for the current block to buff
		j := 0
		for i := block; i < dataSize; i += size.BlockCount {
			buff[j] = int(data[i])
			j++
		}
		// calc the error correction codes
		ecc := ec.rs.Encode(buff, size.ErrorCorrectionCodewordsPerBlock())
		// and append them to the result
		j = 0
		for i := block; i < size.ErrorCorrectionCodewordsPerBlock()*size.BlockCount; i += size.BlockCount {
			data[dataSize+i] = byte(ecc[j])
			j++
		}
	}

	return data
}
