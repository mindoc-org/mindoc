package twooffive

import (
	"image/color"
	"testing"
)

func Test_AddCheckSum(t *testing.T) {
	if sum, err := AddCheckSum("1234567"); err != nil || sum != "12345670" {
		t.Fail()
	}
	if _, err := AddCheckSum("1ABC"); err == nil {
		t.Fail()
	}
	if _, err := AddCheckSum(""); err == nil {
		t.Fail()
	}
}

func Test_Encode(t *testing.T) {
	_, err := Encode("FOOBAR", false)
	if err == nil {
		t.Error("\"FOOBAR\" should not be encodable")
	}

	testEncode := func(interleaved bool, txt, testResult string) {
		code, err := Encode(txt, interleaved)
		if err != nil || code == nil {
			t.Fail()
		} else {
			if code.Bounds().Max.X != len(testResult) {
				t.Errorf("%v: length missmatch! %v != %v", txt, code.Bounds().Max.X, len(testResult))
			} else {
				for i, r := range testResult {
					if (code.At(i, 0) == color.Black) != (r == '1') {
						t.Errorf("%v: code missmatch on position %d", txt, i)
					}
				}
			}
		}
	}

	testEncode(false, "12345670", "1101101011101010101110101110101011101110111010101010101110101110111010111010101011101110101010101011101110101011101110101101011")
	testEncode(true, "12345670", "10101110100010101110001110111010001010001110100011100010101010100011100011101101")
}
