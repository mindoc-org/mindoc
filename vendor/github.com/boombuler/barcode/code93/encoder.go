// Package code93 can create Code93 barcodes
package code93

import (
	"errors"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type encodeInfo struct {
	value int
	data  int
}

const (
	// Special Function 1 ($)
	FNC1 = '\u00f1'
	// Special Function 2 (%)
	FNC2 = '\u00f2'
	// Special Function 3 (/)
	FNC3 = '\u00f3'
	// Special Function 4 (+)
	FNC4 = '\u00f4'
)

var encodeTable = map[rune]encodeInfo{
	'0': encodeInfo{0, 0x114}, '1': encodeInfo{1, 0x148}, '2': encodeInfo{2, 0x144},
	'3': encodeInfo{3, 0x142}, '4': encodeInfo{4, 0x128}, '5': encodeInfo{5, 0x124},
	'6': encodeInfo{6, 0x122}, '7': encodeInfo{7, 0x150}, '8': encodeInfo{8, 0x112},
	'9': encodeInfo{9, 0x10A}, 'A': encodeInfo{10, 0x1A8}, 'B': encodeInfo{11, 0x1A4},
	'C': encodeInfo{12, 0x1A2}, 'D': encodeInfo{13, 0x194}, 'E': encodeInfo{14, 0x192},
	'F': encodeInfo{15, 0x18A}, 'G': encodeInfo{16, 0x168}, 'H': encodeInfo{17, 0x164},
	'I': encodeInfo{18, 0x162}, 'J': encodeInfo{19, 0x134}, 'K': encodeInfo{20, 0x11A},
	'L': encodeInfo{21, 0x158}, 'M': encodeInfo{22, 0x14C}, 'N': encodeInfo{23, 0x146},
	'O': encodeInfo{24, 0x12C}, 'P': encodeInfo{25, 0x116}, 'Q': encodeInfo{26, 0x1B4},
	'R': encodeInfo{27, 0x1B2}, 'S': encodeInfo{28, 0x1AC}, 'T': encodeInfo{29, 0x1A6},
	'U': encodeInfo{30, 0x196}, 'V': encodeInfo{31, 0x19A}, 'W': encodeInfo{32, 0x16C},
	'X': encodeInfo{33, 0x166}, 'Y': encodeInfo{34, 0x136}, 'Z': encodeInfo{35, 0x13A},
	'-': encodeInfo{36, 0x12E}, '.': encodeInfo{37, 0x1D4}, ' ': encodeInfo{38, 0x1D2},
	'$': encodeInfo{39, 0x1CA}, '/': encodeInfo{40, 0x16E}, '+': encodeInfo{41, 0x176},
	'%': encodeInfo{42, 0x1AE}, FNC1: encodeInfo{43, 0x126}, FNC2: encodeInfo{44, 0x1DA},
	FNC3: encodeInfo{45, 0x1D6}, FNC4: encodeInfo{46, 0x132}, '*': encodeInfo{47, 0x15E},
}

var extendedTable = []string{
	"\u00f2U", "\u00f1A", "\u00f1B", "\u00f1C", "\u00f1D", "\u00f1E", "\u00f1F", "\u00f1G",
	"\u00f1H", "\u00f1I", "\u00f1J", "\u00f1K", "\u00f1L", "\u00f1M", "\u00f1N", "\u00f1O",
	"\u00f1P", "\u00f1Q", "\u00f1R", "\u00f1S", "\u00f1T", "\u00f1U", "\u00f1V", "\u00f1W",
	"\u00f1X", "\u00f1Y", "\u00f1Z", "\u00f2A", "\u00f2B", "\u00f2C", "\u00f2D", "\u00f2E",
	" ", "\u00f3A", "\u00f3B", "\u00f3C", "\u00f3D", "\u00f3E", "\u00f3F", "\u00f3G",
	"\u00f3H", "\u00f3I", "\u00f3J", "\u00f3K", "\u00f3L", "-", ".", "\u00f3O",
	"0", "1", "2", "3", "4", "5", "6", "7",
	"8", "9", "\u00f3Z", "\u00f2F", "\u00f2G", "\u00f2H", "\u00f2I", "\u00f2J",
	"\u00f2V", "A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W",
	"X", "Y", "Z", "\u00f2K", "\u00f2L", "\u00f2M", "\u00f2N", "\u00f2O",
	"\u00f2W", "\u00f4A", "\u00f4B", "\u00f4C", "\u00f4D", "\u00f4E", "\u00f4F", "\u00f4G",
	"\u00f4H", "\u00f4I", "\u00f4J", "\u00f4K", "\u00f4L", "\u00f4M", "\u00f4N", "\u00f4O",
	"\u00f4P", "\u00f4Q", "\u00f4R", "\u00f4S", "\u00f4T", "\u00f4U", "\u00f4V", "\u00f4W",
	"\u00f4X", "\u00f4Y", "\u00f4Z", "\u00f2P", "\u00f2Q", "\u00f2R", "\u00f2S", "\u00f2T",
}

func prepare(content string) (string, error) {
	result := ""
	for _, r := range content {
		if r > 127 {
			return "", errors.New("Only ASCII strings can be encoded")
		}
		result += extendedTable[int(r)]
	}
	return result, nil
}

// Encode returns a code93 barcode for the given content
// if includeChecksum is set to true, two checksum characters are calculated and added to the content
func Encode(content string, includeChecksum bool, fullASCIIMode bool) (barcode.Barcode, error) {
	if fullASCIIMode {
		var err error
		content, err = prepare(content)
		if err != nil {
			return nil, err
		}
	} else if strings.ContainsRune(content, '*') {
		return nil, errors.New("invalid data! content may not contain '*'")
	}

	data := content + string(getChecksum(content, 20))
	data += string(getChecksum(data, 15))

	data = "*" + data + "*"
	result := new(utils.BitList)

	for _, r := range data {
		info, ok := encodeTable[r]
		if !ok {
			return nil, errors.New("invalid data!")
		}
		result.AddBits(info.data, 9)
	}
	result.AddBit(true)

	return utils.New1DCode(barcode.TypeCode93, content, result), nil
}

func getChecksum(content string, maxWeight int) rune {
	weight := 1
	total := 0

	data := []rune(content)
	for i := len(data) - 1; i >= 0; i-- {
		r := data[i]
		info, ok := encodeTable[r]
		if !ok {
			return ' '
		}
		total += info.value * weight
		if weight++; weight > maxWeight {
			weight = 1
		}
	}
	total = total % 47
	for r, info := range encodeTable {
		if info.value == total {
			return r
		}
	}
	return ' '
}
