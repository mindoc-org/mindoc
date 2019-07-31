package aztec

import (
	"testing"
)

func encodeTest(t *testing.T, data, wanted string) {
	result, err := Encode([]byte(data), DEFAULT_EC_PERCENT, DEFAULT_LAYERS)
	if err != nil {
		t.Error(err)
	} else {
		ac, ok := result.(*aztecCode)
		if !ok {
			t.Error("returned barcode is no aztec code...")
		} else if draw := ac.string(); draw != wanted {
			t.Errorf("Invalid Barcode returned:\n%s", draw)
		}

	}
}

func Test_Encode1(t *testing.T) {
	encodeTest(t, "This is an example Aztec symbol for Wikipedia.",
		"X     X X       X     X X     X     X         \n"+
			"X         X     X X     X   X X   X X       X \n"+
			"X X   X X X X X   X X X                 X     \n"+
			"X X                 X X   X       X X X X X X \n"+
			"    X X X   X   X     X X X X         X X     \n"+
			"  X X X   X X X X   X     X   X     X X   X   \n"+
			"        X X X X X     X X X X   X   X     X   \n"+
			"X       X   X X X X X X X X X X X     X   X X \n"+
			"X   X     X X X               X X X X   X X   \n"+
			"X     X X   X X   X X X X X   X X   X   X X X \n"+
			"X   X         X   X       X   X X X X       X \n"+
			"X       X     X   X   X   X   X   X X   X     \n"+
			"      X   X X X   X       X   X     X X X     \n"+
			"    X X X X X X   X X X X X   X X X X X X   X \n"+
			"  X X   X   X X               X X X   X X X X \n"+
			"  X   X       X X X X X X X X X X X X   X X   \n"+
			"  X X   X       X X X   X X X       X X       \n"+
			"  X               X   X X     X     X X X     \n"+
			"  X   X X X   X X   X   X X X X   X   X X X X \n"+
			"    X   X   X X X   X   X   X X X X     X     \n"+
			"        X               X                 X   \n"+
			"        X X     X   X X   X   X   X       X X \n"+
			"  X   X   X X       X   X         X X X     X \n")
}

func Test_Encode2(t *testing.T) {
	encodeTest(t, "Aztec Code is a public domain 2D matrix barcode symbology"+
		" of nominally square symbols built on a square grid with a "+
		"distinctive square bullseye pattern at their center.",
		"        X X     X X     X     X     X   X X X         X   X         X   X X       \n"+
			"  X       X X     X   X X   X X       X             X     X   X X   X           X \n"+
			"  X   X X X     X   X   X X     X X X   X   X X               X X       X X     X \n"+
			"X X X             X   X         X         X     X     X   X     X X       X   X   \n"+
			"X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X \n"+
			"    X X   X   X   X X X               X       X       X X     X X   X X       X   \n"+
			"X X     X       X       X X X X   X   X X       X   X X   X       X X   X X   X   \n"+
			"  X       X   X     X X   X   X X   X X   X X X X X X   X X           X   X   X X \n"+
			"X X   X X   X   X X X X   X X X X X X X X   X   X       X X   X X X X   X X X     \n"+
			"  X       X   X     X       X X     X X   X   X   X     X X   X X X   X     X X X \n"+
			"  X   X X X   X X       X X X         X X           X   X   X   X X X   X X     X \n"+
			"    X     X   X X     X X X X     X   X     X X X X   X X   X X   X X X     X   X \n"+
			"X X X   X             X         X X X X X   X   X X   X   X   X X   X   X   X   X \n"+
			"          X       X X X   X X     X   X           X   X X X X   X X               \n"+
			"  X     X X   X   X       X X X X X X X X X X X X X X X   X   X X   X   X X X     \n"+
			"    X X                 X   X                       X X   X       X         X X X \n"+
			"        X   X X   X X X X X X   X X X X X X X X X   X     X X           X X X X   \n"+
			"          X X X   X     X   X   X               X   X X     X X X   X X           \n"+
			"X X     X     X   X   X   X X   X   X X X X X   X   X X X X X X X       X   X X X \n"+
			"X X X X       X       X   X X   X   X       X   X   X     X X X     X X       X X \n"+
			"X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X \n"+
			"    X     X       X         X   X   X       X   X   X     X   X X                 \n"+
			"        X X     X X X X X   X   X   X X X X X   X   X X X     X X X X   X         \n"+
			"X     X   X   X         X   X   X               X   X X   X X   X X X     X   X   \n"+
			"  X   X X X   X   X X   X X X   X X X X X X X X X   X X         X X     X X X X   \n"+
			"    X X   X   X   X X X     X                       X X X   X X   X   X     X     \n"+
			"    X X X X   X         X   X X X X X X X X X X X X X X   X       X X   X X   X X \n"+
			"            X   X   X X       X X X X X     X X X       X       X X X         X   \n"+
			"X       X         X   X X X X   X     X X     X X     X X           X   X       X \n"+
			"X     X       X X X X X     X   X X X X   X X X     X       X X X X   X   X X   X \n"+
			"  X X X X X               X     X X X   X       X X   X X   X X X X     X X       \n"+
			"X             X         X   X X   X X     X     X     X   X   X X X X             \n"+
			"    X   X X       X     X       X   X X X X X X   X X   X X X X X X X X X   X   X \n"+
			"    X         X X   X       X     X   X   X       X     X X X     X       X X X X \n"+
			"X     X X     X X X X X X             X X X   X               X   X     X     X X \n"+
			"X   X X     X               X X X X X     X X     X X X X X X X X     X   X   X X \n"+
			"X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X   X \n"+
			"X           X     X X X X     X     X         X         X   X       X X   X X X   \n"+
			"X   X   X X   X X X   X         X X     X X X X     X X   X   X     X   X       X \n"+
			"      X     X     X     X X     X   X X   X X   X         X X       X       X   X \n"+
			"X       X           X   X   X     X X   X               X     X     X X X         \n")
}
