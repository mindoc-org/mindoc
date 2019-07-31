// Package twooffive can create interleaved and standard "2 of 5" barcodes.
package twooffive

import (
	"errors"
	"fmt"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

const patternWidth = 5

type pattern [patternWidth]bool
type encodeInfo struct {
	start  []bool
	end    []bool
	widths map[bool]int
}

var (
	encodingTable = map[rune]pattern{
		'0': pattern{false, false, true, true, false},
		'1': pattern{true, false, false, false, true},
		'2': pattern{false, true, false, false, true},
		'3': pattern{true, true, false, false, false},
		'4': pattern{false, false, true, false, true},
		'5': pattern{true, false, true, false, false},
		'6': pattern{false, true, true, false, false},
		'7': pattern{false, false, false, true, true},
		'8': pattern{true, false, false, true, false},
		'9': pattern{false, true, false, true, false},
	}

	modes = map[bool]encodeInfo{
		false: encodeInfo{ // non-interleaved
			start: []bool{true, true, false, true, true, false, true, false},
			end:   []bool{true, true, false, true, false, true, true},
			widths: map[bool]int{
				true:  3,
				false: 1,
			},
		},
		true: encodeInfo{ // interleaved
			start: []bool{true, false, true, false},
			end:   []bool{true, true, false, true},
			widths: map[bool]int{
				true:  3,
				false: 1,
			},
		},
	}
	nonInterleavedSpace = pattern{false, false, false, false, false}
)

// AddCheckSum calculates the correct check-digit and appends it to the given content.
func AddCheckSum(content string) (string, error) {
	if content == "" {
		return "", errors.New("content is empty")
	}

	even := len(content)%2 == 1
	sum := 0
	for _, r := range content {
		if _, ok := encodingTable[r]; ok {
			value := utils.RuneToInt(r)
			if even {
				sum += value * 3
			} else {
				sum += value
			}
			even = !even
		} else {
			return "", fmt.Errorf("can not encode \"%s\"", content)
		}
	}

	return content + string(utils.IntToRune(sum%10)), nil
}

// Encode creates a codabar barcode for the given content
func Encode(content string, interleaved bool) (barcode.Barcode, error) {
	if content == "" {
		return nil, errors.New("content is empty")
	}

	if interleaved && len(content)%2 == 1 {
		return nil, errors.New("can only encode even number of digits in interleaved mode")
	}

	mode := modes[interleaved]
	resBits := new(utils.BitList)
	resBits.AddBit(mode.start...)

	var lastRune *rune
	for _, r := range content {
		var a, b pattern
		if interleaved {
			if lastRune == nil {
				lastRune = new(rune)
				*lastRune = r
				continue
			} else {
				var o1, o2 bool
				a, o1 = encodingTable[*lastRune]
				b, o2 = encodingTable[r]
				if !o1 || !o2 {
					return nil, fmt.Errorf("can not encode \"%s\"", content)
				}
				lastRune = nil
			}
		} else {
			var ok bool
			a, ok = encodingTable[r]
			if !ok {
				return nil, fmt.Errorf("can not encode \"%s\"", content)
			}
			b = nonInterleavedSpace
		}

		for i := 0; i < patternWidth; i++ {
			for x := 0; x < mode.widths[a[i]]; x++ {
				resBits.AddBit(true)
			}
			for x := 0; x < mode.widths[b[i]]; x++ {
				resBits.AddBit(false)
			}
		}
	}

	resBits.AddBit(mode.end...)

	if interleaved {
		return utils.New1DCode(barcode.Type2of5Interleaved, content, resBits), nil
	} else {
		return utils.New1DCode(barcode.Type2of5, content, resBits), nil
	}
}
