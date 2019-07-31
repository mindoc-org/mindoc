// Package code39 can create Code39 barcodes
package code39

import (
	"errors"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type encodeInfo struct {
	value int
	data  []bool
}

var encodeTable = map[rune]encodeInfo{
	'0': encodeInfo{0, []bool{true, false, true, false, false, true, true, false, true, true, false, true}},
	'1': encodeInfo{1, []bool{true, true, false, true, false, false, true, false, true, false, true, true}},
	'2': encodeInfo{2, []bool{true, false, true, true, false, false, true, false, true, false, true, true}},
	'3': encodeInfo{3, []bool{true, true, false, true, true, false, false, true, false, true, false, true}},
	'4': encodeInfo{4, []bool{true, false, true, false, false, true, true, false, true, false, true, true}},
	'5': encodeInfo{5, []bool{true, true, false, true, false, false, true, true, false, true, false, true}},
	'6': encodeInfo{6, []bool{true, false, true, true, false, false, true, true, false, true, false, true}},
	'7': encodeInfo{7, []bool{true, false, true, false, false, true, false, true, true, false, true, true}},
	'8': encodeInfo{8, []bool{true, true, false, true, false, false, true, false, true, true, false, true}},
	'9': encodeInfo{9, []bool{true, false, true, true, false, false, true, false, true, true, false, true}},
	'A': encodeInfo{10, []bool{true, true, false, true, false, true, false, false, true, false, true, true}},
	'B': encodeInfo{11, []bool{true, false, true, true, false, true, false, false, true, false, true, true}},
	'C': encodeInfo{12, []bool{true, true, false, true, true, false, true, false, false, true, false, true}},
	'D': encodeInfo{13, []bool{true, false, true, false, true, true, false, false, true, false, true, true}},
	'E': encodeInfo{14, []bool{true, true, false, true, false, true, true, false, false, true, false, true}},
	'F': encodeInfo{15, []bool{true, false, true, true, false, true, true, false, false, true, false, true}},
	'G': encodeInfo{16, []bool{true, false, true, false, true, false, false, true, true, false, true, true}},
	'H': encodeInfo{17, []bool{true, true, false, true, false, true, false, false, true, true, false, true}},
	'I': encodeInfo{18, []bool{true, false, true, true, false, true, false, false, true, true, false, true}},
	'J': encodeInfo{19, []bool{true, false, true, false, true, true, false, false, true, true, false, true}},
	'K': encodeInfo{20, []bool{true, true, false, true, false, true, false, true, false, false, true, true}},
	'L': encodeInfo{21, []bool{true, false, true, true, false, true, false, true, false, false, true, true}},
	'M': encodeInfo{22, []bool{true, true, false, true, true, false, true, false, true, false, false, true}},
	'N': encodeInfo{23, []bool{true, false, true, false, true, true, false, true, false, false, true, true}},
	'O': encodeInfo{24, []bool{true, true, false, true, false, true, true, false, true, false, false, true}},
	'P': encodeInfo{25, []bool{true, false, true, true, false, true, true, false, true, false, false, true}},
	'Q': encodeInfo{26, []bool{true, false, true, false, true, false, true, true, false, false, true, true}},
	'R': encodeInfo{27, []bool{true, true, false, true, false, true, false, true, true, false, false, true}},
	'S': encodeInfo{28, []bool{true, false, true, true, false, true, false, true, true, false, false, true}},
	'T': encodeInfo{29, []bool{true, false, true, false, true, true, false, true, true, false, false, true}},
	'U': encodeInfo{30, []bool{true, true, false, false, true, false, true, false, true, false, true, true}},
	'V': encodeInfo{31, []bool{true, false, false, true, true, false, true, false, true, false, true, true}},
	'W': encodeInfo{32, []bool{true, true, false, false, true, true, false, true, false, true, false, true}},
	'X': encodeInfo{33, []bool{true, false, false, true, false, true, true, false, true, false, true, true}},
	'Y': encodeInfo{34, []bool{true, true, false, false, true, false, true, true, false, true, false, true}},
	'Z': encodeInfo{35, []bool{true, false, false, true, true, false, true, true, false, true, false, true}},
	'-': encodeInfo{36, []bool{true, false, false, true, false, true, false, true, true, false, true, true}},
	'.': encodeInfo{37, []bool{true, true, false, false, true, false, true, false, true, true, false, true}},
	' ': encodeInfo{38, []bool{true, false, false, true, true, false, true, false, true, true, false, true}},
	'$': encodeInfo{39, []bool{true, false, false, true, false, false, true, false, false, true, false, true}},
	'/': encodeInfo{40, []bool{true, false, false, true, false, false, true, false, true, false, false, true}},
	'+': encodeInfo{41, []bool{true, false, false, true, false, true, false, false, true, false, false, true}},
	'%': encodeInfo{42, []bool{true, false, true, false, false, true, false, false, true, false, false, true}},
	'*': encodeInfo{-1, []bool{true, false, false, true, false, true, true, false, true, true, false, true}},
}

var extendedTable = map[rune]string{
	0: `%U`, 1: `$A`, 2: `$B`, 3: `$C`, 4: `$D`, 5: `$E`, 6: `$F`, 7: `$G`, 8: `$H`, 9: `$I`, 10: `$J`,
	11: `$K`, 12: `$L`, 13: `$M`, 14: `$N`, 15: `$O`, 16: `$P`, 17: `$Q`, 18: `$R`, 19: `$S`, 20: `$T`,
	21: `$U`, 22: `$V`, 23: `$W`, 24: `$X`, 25: `$Y`, 26: `$Z`, 27: `%A`, 28: `%B`, 29: `%C`, 30: `%D`,
	31: `%E`, 33: `/A`, 34: `/B`, 35: `/C`, 36: `/D`, 37: `/E`, 38: `/F`, 39: `/G`, 40: `/H`, 41: `/I`,
	42: `/J`, 43: `/K`, 44: `/L`, 47: `/O`, 58: `/Z`, 59: `%F`, 60: `%G`, 61: `%H`, 62: `%I`, 63: `%J`,
	64: `%V`, 91: `%K`, 92: `%L`, 93: `%M`, 94: `%N`, 95: `%O`, 96: `%W`, 97: `+A`, 98: `+B`, 99: `+C`,
	100: `+D`, 101: `+E`, 102: `+F`, 103: `+G`, 104: `+H`, 105: `+I`, 106: `+J`, 107: `+K`, 108: `+L`,
	109: `+M`, 110: `+N`, 111: `+O`, 112: `+P`, 113: `+Q`, 114: `+R`, 115: `+S`, 116: `+T`, 117: `+U`,
	118: `+V`, 119: `+W`, 120: `+X`, 121: `+Y`, 122: `+Z`, 123: `%P`, 124: `%Q`, 125: `%R`, 126: `%S`,
	127: `%T`,
}

func getChecksum(content string) string {
	sum := 0
	for _, r := range content {
		info, ok := encodeTable[r]
		if !ok || info.value < 0 {
			return "#"
		}

		sum += info.value
	}

	sum = sum % 43
	for r, v := range encodeTable {
		if v.value == sum {
			return string(r)
		}
	}
	return "#"
}

func prepare(content string) (string, error) {
	result := ""
	for _, r := range content {
		if r > 127 {
			return "", errors.New("Only ASCII strings can be encoded")
		}
		val, ok := extendedTable[r]
		if ok {
			result += val
		} else {
			result += string([]rune{r})
		}
	}
	return result, nil
}

// Encode returns a code39 barcode for the given content
// if includeChecksum is set to true, a checksum character is calculated and added to the content
func Encode(content string, includeChecksum bool, fullASCIIMode bool) (barcode.BarcodeIntCS, error) {
	if fullASCIIMode {
		var err error
		content, err = prepare(content)
		if err != nil {
			return nil, err
		}
	} else if strings.ContainsRune(content, '*') {
		return nil, errors.New("invalid data! try full ascii mode")
	}

	data := "*" + content
	if includeChecksum {
		data += getChecksum(content)
	}
	data += "*"

	result := new(utils.BitList)

	for i, r := range data {
		if i != 0 {
			result.AddBit(false)
		}

		info, ok := encodeTable[r]
		if !ok {
			return nil, errors.New("invalid data! try full ascii mode")
		}
		result.AddBit(info.data...)
	}

	checkSum, err := strconv.ParseInt(getChecksum(content), 10, 64)
	if err != nil {
		checkSum = 0
	}
	return utils.New1DCodeIntCheckSum(barcode.TypeCode39, content, result, int(checkSum)), nil
}
