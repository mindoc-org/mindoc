package aztec

import (
	"strings"
	"testing"

	"github.com/boombuler/barcode/utils"
)

func Test_StuffBits(t *testing.T) {
	testStuffBits := func(wordSize int, bits string, expected string) {
		bl := new(utils.BitList)
		for _, r := range bits {
			if r == 'X' {
				bl.AddBit(true)
			} else if r == '.' {
				bl.AddBit(false)
			}
		}
		stuffed := stuffBits(bl, wordSize)
		expectedBits := strings.Replace(expected, " ", "", -1)
		result := bitStr(stuffed)

		if result != expectedBits {
			t.Errorf("stuffBits failed for %q\nGot: %q", bits, result)
		}
	}

	testStuffBits(5, ".X.X. X.X.X .X.X.",
		".X.X. X.X.X .X.X.")
	testStuffBits(5, ".X.X. ..... .X.X",
		".X.X. ....X ..X.X")
	testStuffBits(3, "XX. ... ... ..X XXX .X. ..",
		"XX. ..X ..X ..X ..X .XX XX. .X. ..X")
	testStuffBits(6, ".X.X.. ...... ..X.XX",
		".X.X.. .....X. ..X.XX XXXX.")
	testStuffBits(6, ".X.X.. ...... ...... ..X.X.",
		".X.X.. .....X .....X ....X. X.XXXX")
	testStuffBits(6, ".X.X.. XXXXXX ...... ..X.XX",
		".X.X.. XXXXX. X..... ...X.X XXXXX.")
	testStuffBits(6,
		"...... ..XXXX X..XX. .X.... .X.X.X .....X .X.... ...X.X .....X ....XX ..X... ....X. X..XXX X.XX.X",
		".....X ...XXX XX..XX ..X... ..X.X. X..... X.X... ....X. X..... X....X X..X.. .....X X.X..X XXX.XX .XXXXX")
}

func Test_ModeMessage(t *testing.T) {
	testModeMessage := func(compact bool, layers, words int, expected string) {
		result := bitStr(generateModeMessage(compact, layers, words))
		expectedBits := strings.Replace(expected, " ", "", -1)
		if result != expectedBits {
			t.Errorf("generateModeMessage(%v, %d, %d) failed.\nGot:%s", compact, layers, words, result)
		}
	}
	testModeMessage(true, 2, 29, ".X .XXX.. ...X XX.. ..X .XX. .XX.X")
	testModeMessage(true, 4, 64, "XX XXXXXX .X.. ...X ..XX .X.. XX..")
	testModeMessage(false, 21, 660, "X.X.. .X.X..X..XX .XXX ..X.. .XXX. .X... ..XXX")
	testModeMessage(false, 32, 4096, "XXXXX XXXXXXXXXXX X.X. ..... XXX.X ..X.. X.XXX")
}
