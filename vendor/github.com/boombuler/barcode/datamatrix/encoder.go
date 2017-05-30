// Package datamatrix can create Datamatrix barcodes
package datamatrix

import (
	"errors"

	"github.com/boombuler/barcode"
)

// Encode returns a Datamatrix barcode for the given content
func Encode(content string) (barcode.Barcode, error) {
	data := encodeText(content)

	var size *dmCodeSize
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(data) {
			size = s
			break
		}
	}
	if size == nil {
		return nil, errors.New("to much data to encode")
	}
	data = addPadding(data, size.DataCodewords())
	data = ec.calcECC(data, size)
	code := render(data, size)
	if code != nil {
		code.content = content
		return code, nil
	}
	return nil, errors.New("unable to render barcode")
}

func render(data []byte, size *dmCodeSize) *datamatrixCode {
	cl := newCodeLayout(size)

	cl.SetValues(data)

	return cl.Merge()
}

func encodeText(content string) []byte {
	var result []byte
	input := []byte(content)

	for i := 0; i < len(input); {
		c := input[i]
		i++

		if c >= '0' && c <= '9' && i < len(input) && input[i] >= '0' && input[i] <= '9' {
			// two numbers...
			c2 := input[i]
			i++
			cw := byte(((c-'0')*10 + (c2 - '0')) + 130)
			result = append(result, cw)
		} else if c > 127 {
			// not correct... needs to be redone later...
			result = append(result, 235, c-127)
		} else {
			result = append(result, c+1)
		}
	}
	return result
}

func addPadding(data []byte, toCount int) []byte {
	if len(data) < toCount {
		data = append(data, 129)
	}
	for len(data) < toCount {
		R := ((149 * (len(data) + 1)) % 253) + 1
		data = append(data, byte((129+R)%254))
	}
	return data
}
