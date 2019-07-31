package code128

import (
	"image/color"
	"testing"
)

func testEncode(t *testing.T, txt, testResult string) {
	code, err := Encode(txt)
	if err != nil || code == nil {
		t.Error(err)
	} else {
		if code.Bounds().Max.X != len(testResult) {
			t.Errorf("%v: length missmatch. Got %d expected %d", txt, code.Bounds().Max.X, len(testResult))
		} else {
			encoded := ""
			failed := false
			for i, r := range testResult {
				if code.At(i, 0) == color.Black {
					encoded += "1"
				} else {
					encoded += "0"
				}

				if (code.At(i, 0) == color.Black) != (r == '1') {
					failed = true
					t.Errorf("%v: code missmatch on position %d", txt, i)
				}
			}
			if failed {
				t.Log("Encoded: ", encoded)
			}
		}
	}
}

func Test_EncodeFunctionChars(t *testing.T) {
	encFNC1 := "11110101110"
	encFNC2 := "11110101000"
	encFNC3 := "10111100010"
	encFNC4 := "10111101110"
	encStartB := "11010010000"
	encStop := "1100011101011"

	// Special Case FC1 can also be encoded to C Table therefor using 123 as suffix might have unexpected results.
	testEncode(t, string(FNC1)+"A23", encStartB+encFNC1+"10100011000"+"11001110010"+"11001011100"+"10100011110"+encStop)
	testEncode(t, string(FNC2)+"123", encStartB+encFNC2+"10011100110"+"11001110010"+"11001011100"+"11100010110"+encStop)
	testEncode(t, string(FNC3)+"123", encStartB+encFNC3+"10011100110"+"11001110010"+"11001011100"+"11101000110"+encStop)
	testEncode(t, string(FNC4)+"123", encStartB+encFNC4+"10011100110"+"11001110010"+"11001011100"+"11100011010"+encStop)
}

func Test_Unencodable(t *testing.T) {
	if _, err := Encode(""); err == nil {
		t.Fail()
	}
	if _, err := Encode("ä"); err == nil {
		t.Fail()
	}
}

func Test_EncodeCTable(t *testing.T) {
	testEncode(t, "HI345678H", "110100100001100010100011000100010101110111101000101100011100010110110000101001011110111011000101000111011000101100011101011")
	testEncode(t, "334455", "11010011100101000110001000110111011101000110100100111101100011101011")

	testEncode(t, string(FNC1)+"1234",
		"11010011100"+ // Start C
			"11110101110"+ // FNC1
			"10110011100"+ // 12
			"10001011000"+ // 34
			"11101001100"+ // CheckSum == 24
			"1100011101011") // Stop
}

func Test_shouldUseCTable(t *testing.T) {
	if !shouldUseCTable([]rune{FNC1, '1', '2'}, startCSymbol) {
		t.Error("[FNC1]12 failed")
	}
	if shouldUseCTable([]rune{FNC1, '1'}, startCSymbol) {
		t.Error("[FNC1]1 failed")
	}
	if shouldUseCTable([]rune{'0', FNC1, '1'}, startCSymbol) {
		t.Error("0[FNC1]1 failed")
	}
	if !shouldUseCTable([]rune{'0', '1', FNC1, '2', '3'}, startBSymbol) {
		t.Error("01[FNC1]23 failed")
	}
	if shouldUseCTable([]rune{'0', '1', FNC1}, startBSymbol) {
		t.Error("01[FNC1] failed")
	}
}

func Test_Issue16(t *testing.T) {
	if !shouldUseATable([]rune{'\r', 'A'}, 0) {
		t.Error("Code should start with A-Table if the text start with \\r")
	}
	if !shouldUseATable([]rune{FNC1, '\r'}, 0) {
		t.Error("Code should start with A-Table if the text start with <FNC1>\\r")
	}
	if shouldUseATable([]rune{FNC1, '1', '2', '3'}, 0) {
		t.Error("Code should not start with A-Table if the text start with <FNC1>123")
	}
	testEncode(t, string(FNC3)+"$P\rI", "110100001001011110001010010001100111011101101111011101011000100010110001010001100011101011")
}

func Test_Datalogic(t *testing.T) {
	// <Start A><FNC3>$P\r<checksum><STOP>
	testEncode(t, string(FNC3)+"$P\r",
		"11010000100"+ // <Start A>
			"10111100010"+ // <FNC3>
			"10010001100"+ // $
			"11101110110"+ // P
			"11110111010"+ // CR
			"11000100010"+ // checksum = 'I'
			"1100011101011") // STOP

	// <Start B><FNC3>$P,Ae,P<CR><checksum><STOP>
	testEncode(t, string(FNC3)+"$P,Ae,P\r",
		"11010010000"+ // <Start B>
			"10111100010"+ // <FNC3>
			"10010001100"+ // $
			"11101110110"+ // P
			"10110011100"+ // ,
			"10100011000"+ // A
			"10110010000"+ // e
			"10110011100"+ // ,
			"11101110110"+ // P
			"11101011110"+ // <Code A>
			"11110111010"+ // <CR>
			"10110001000"+ // checksum = 'D'
			"1100011101011") // STOP
}
