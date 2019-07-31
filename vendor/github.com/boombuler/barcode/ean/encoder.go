// Package ean can create EAN 8 and EAN 13 barcodes.
package ean

import (
	"errors"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type encodedNumber struct {
	LeftOdd  []bool
	LeftEven []bool
	Right    []bool
	CheckSum []bool
}

var encoderTable = map[rune]encodedNumber{
	'0': encodedNumber{
		[]bool{false, false, false, true, true, false, true},
		[]bool{false, true, false, false, true, true, true},
		[]bool{true, true, true, false, false, true, false},
		[]bool{false, false, false, false, false, false},
	},
	'1': encodedNumber{
		[]bool{false, false, true, true, false, false, true},
		[]bool{false, true, true, false, false, true, true},
		[]bool{true, true, false, false, true, true, false},
		[]bool{false, false, true, false, true, true},
	},
	'2': encodedNumber{
		[]bool{false, false, true, false, false, true, true},
		[]bool{false, false, true, true, false, true, true},
		[]bool{true, true, false, true, true, false, false},
		[]bool{false, false, true, true, false, true},
	},
	'3': encodedNumber{
		[]bool{false, true, true, true, true, false, true},
		[]bool{false, true, false, false, false, false, true},
		[]bool{true, false, false, false, false, true, false},
		[]bool{false, false, true, true, true, false},
	},
	'4': encodedNumber{
		[]bool{false, true, false, false, false, true, true},
		[]bool{false, false, true, true, true, false, true},
		[]bool{true, false, true, true, true, false, false},
		[]bool{false, true, false, false, true, true},
	},
	'5': encodedNumber{
		[]bool{false, true, true, false, false, false, true},
		[]bool{false, true, true, true, false, false, true},
		[]bool{true, false, false, true, true, true, false},
		[]bool{false, true, true, false, false, true},
	},
	'6': encodedNumber{
		[]bool{false, true, false, true, true, true, true},
		[]bool{false, false, false, false, true, false, true},
		[]bool{true, false, true, false, false, false, false},
		[]bool{false, true, true, true, false, false},
	},
	'7': encodedNumber{
		[]bool{false, true, true, true, false, true, true},
		[]bool{false, false, true, false, false, false, true},
		[]bool{true, false, false, false, true, false, false},
		[]bool{false, true, false, true, false, true},
	},
	'8': encodedNumber{
		[]bool{false, true, true, false, true, true, true},
		[]bool{false, false, false, true, false, false, true},
		[]bool{true, false, false, true, false, false, false},
		[]bool{false, true, false, true, true, false},
	},
	'9': encodedNumber{
		[]bool{false, false, false, true, false, true, true},
		[]bool{false, false, true, false, true, true, true},
		[]bool{true, true, true, false, true, false, false},
		[]bool{false, true, true, false, true, false},
	},
}

func calcCheckNum(code string) rune {
	x3 := len(code) == 7
	sum := 0
	for _, r := range code {
		curNum := utils.RuneToInt(r)
		if curNum < 0 || curNum > 9 {
			return 'B'
		}
		if x3 {
			curNum = curNum * 3
		}
		x3 = !x3
		sum += curNum
	}

	return utils.IntToRune((10 - (sum % 10)) % 10)
}

func encodeEAN8(code string) *utils.BitList {
	result := new(utils.BitList)
	result.AddBit(true, false, true)

	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return nil
		}
		var data []bool
		if cpos < 4 {
			data = num.LeftOdd
		} else {
			data = num.Right
		}

		if cpos == 4 {
			result.AddBit(false, true, false, true, false)
		}
		result.AddBit(data...)
	}
	result.AddBit(true, false, true)

	return result
}

func encodeEAN13(code string) *utils.BitList {
	result := new(utils.BitList)
	result.AddBit(true, false, true)

	var firstNum []bool
	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return nil
		}
		if cpos == 0 {
			firstNum = num.CheckSum
			continue
		}

		var data []bool
		if cpos < 7 { // Left
			if firstNum[cpos-1] {
				data = num.LeftEven
			} else {
				data = num.LeftOdd
			}
		} else {
			data = num.Right
		}

		if cpos == 7 {
			result.AddBit(false, true, false, true, false)
		}
		result.AddBit(data...)
	}
	result.AddBit(true, false, true)
	return result
}

// Encode returns a EAN 8 or EAN 13 barcode for the given code
func Encode(code string) (barcode.BarcodeIntCS, error) {
	var checkSum int
	if len(code) == 7 || len(code) == 12 {
		code += string(calcCheckNum(code))
		checkSum = utils.RuneToInt(calcCheckNum(code))
	} else if len(code) == 8 || len(code) == 13 {
		check := code[0 : len(code)-1]
		check += string(calcCheckNum(check))
		if check != code {
			return nil, errors.New("checksum missmatch")
		}
		checkSum = utils.RuneToInt(rune(code[len(code)-1]))
	}

	if len(code) == 8 {
		result := encodeEAN8(code)
		if result != nil {
			return utils.New1DCodeIntCheckSum(barcode.TypeEAN8, code, result, checkSum), nil
		}
	} else if len(code) == 13 {
		result := encodeEAN13(code)
		if result != nil {
			return utils.New1DCodeIntCheckSum(barcode.TypeEAN13, code, result, checkSum), nil
		}
	}
	return nil, errors.New("invalid ean code data")
}
