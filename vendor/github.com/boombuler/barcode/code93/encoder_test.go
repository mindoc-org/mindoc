package code93

import (
	"image/color"
	"testing"
)

func doTest(t *testing.T, data, testResult string) {
	code, err := Encode(data, true, false)
	if err != nil {
		t.Error(err)
	}
	if len(testResult) != code.Bounds().Max.X {
		t.Errorf("Invalid code size. Expected %d got %d", len(testResult), code.Bounds().Max.X)
	}
	for i, r := range testResult {
		if (code.At(i, 0) == color.Black) != (r == '1') {
			t.Errorf("Failed at position %d", i)
		}
	}
}

func Test_CheckSum(t *testing.T) {
	if r := getChecksum("TEST93", 20); r != '+' {
		t.Errorf("Checksum C-Failed. Got %s", string(r))
	}
	if r := getChecksum("TEST93+", 15); r != '6' {
		t.Errorf("Checksum K-Failed. Got %s", string(r))
	}
}

func Test_Encode(t *testing.T) {
	doTest(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"1010111101101010001101001001101000101100101001100100101100010101011010001011001"+
			"001011000101001101001000110101010110001010011001010001101001011001000101101101101001"+
			"101100101101011001101001101100101101100110101011011001011001101001101101001110101000"+
			"101001010010001010001001010000101001010001001001001001000101010100001000100101000010"+
			"101001110101010000101010111101")
}
