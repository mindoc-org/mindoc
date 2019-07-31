package datamatrix

import (
	"bytes"
	"testing"
)

func Test_CalcECC(t *testing.T) {
	data := []byte{142, 164, 186}
	var size *dmCodeSize = nil
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(data) {
			size = s
			break
		}
	}
	if size == nil {
		t.Error("size not found")
	}

	if bytes.Compare(ec.calcECC(data, size), []byte{142, 164, 186, 114, 25, 5, 88, 102}) != 0 {
		t.Error("ECC Test 1 failed")
	}
	data = []byte{66, 129, 70}
	if bytes.Compare(ec.calcECC(data, size), []byte{66, 129, 70, 138, 234, 82, 82, 95}) != 0 {
		t.Error("ECC Test 2 failed")
	}
}
