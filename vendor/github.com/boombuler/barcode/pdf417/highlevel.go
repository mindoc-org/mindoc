package pdf417

import (
	"errors"
	"math/big"

	"github.com/boombuler/barcode/utils"
)

type encodingMode byte

type subMode byte

const (
	encText encodingMode = iota
	encNumeric
	encBinary

	subUpper subMode = iota
	subLower
	subMixed
	subPunct

	latch_to_text        = 900
	latch_to_byte_padded = 901
	latch_to_numeric     = 902
	latch_to_byte        = 924
	shift_to_byte        = 913

	min_numeric_count = 13
)

var (
	mixedMap map[rune]int
	punctMap map[rune]int
)

func init() {
	mixedMap = make(map[rune]int)
	mixedRaw := []rune{
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 38, 13, 9, 44, 58,
		35, 45, 46, 36, 47, 43, 37, 42, 61, 94, 0, 32, 0, 0, 0,
	}
	for idx, ch := range mixedRaw {
		if ch > 0 {
			mixedMap[ch] = idx
		}
	}

	punctMap = make(map[rune]int)
	punctRaw := []rune{
		59, 60, 62, 64, 91, 92, 93, 95, 96, 126, 33, 13, 9, 44, 58,
		10, 45, 46, 36, 47, 34, 124, 42, 40, 41, 63, 123, 125, 39, 0,
	}
	for idx, ch := range punctRaw {
		if ch > 0 {
			punctMap[ch] = idx
		}
	}
}

func determineConsecutiveDigitCount(data []rune) int {
	cnt := 0
	for _, r := range data {
		if utils.RuneToInt(r) == -1 {
			break
		}
		cnt++
	}
	return cnt
}

func encodeNumeric(digits []rune) ([]int, error) {
	digitCount := len(digits)
	chunkCount := digitCount / 44
	if digitCount%44 != 0 {
		chunkCount++
	}

	codeWords := []int{}

	for i := 0; i < chunkCount; i++ {
		start := i * 44
		end := start + 44
		if end > digitCount {
			end = digitCount
		}
		chunk := digits[start:end]

		chunkNum := big.NewInt(0)
		_, ok := chunkNum.SetString("1"+string(chunk), 10)

		if !ok {
			return nil, errors.New("Failed converting: " + string(chunk))
		}

		cws := []int{}

		for chunkNum.Cmp(big.NewInt(0)) > 0 {
			newChunk, cw := chunkNum.DivMod(chunkNum, big.NewInt(900), big.NewInt(0))
			chunkNum = newChunk
			cws = append([]int{int(cw.Int64())}, cws...)
		}

		codeWords = append(codeWords, cws...)
	}

	return codeWords, nil
}

func determineConsecutiveTextCount(msg []rune) int {
	result := 0

	isText := func(ch rune) bool {
		return ch == '\t' || ch == '\n' || ch == '\r' || (ch >= 32 && ch <= 126)
	}

	for i, ch := range msg {
		numericCount := determineConsecutiveDigitCount(msg[i:])
		if numericCount >= min_numeric_count || (numericCount == 0 && !isText(ch)) {
			break
		}

		result++
	}
	return result
}

func encodeText(text []rune, submode subMode) (subMode, []int) {
	isAlphaUpper := func(ch rune) bool {
		return ch == ' ' || (ch >= 'A' && ch <= 'Z')
	}
	isAlphaLower := func(ch rune) bool {
		return ch == ' ' || (ch >= 'a' && ch <= 'z')
	}
	isMixed := func(ch rune) bool {
		_, ok := mixedMap[ch]
		return ok
	}
	isPunctuation := func(ch rune) bool {
		_, ok := punctMap[ch]
		return ok
	}

	idx := 0
	var tmp []int
	for idx < len(text) {
		ch := text[idx]
		switch submode {
		case subUpper:
			if isAlphaUpper(ch) {
				if ch == ' ' {
					tmp = append(tmp, 26) //space
				} else {
					tmp = append(tmp, int(ch-'A'))
				}
			} else {
				if isAlphaLower(ch) {
					submode = subLower
					tmp = append(tmp, 27) // lower latch
					continue
				} else if isMixed(ch) {
					submode = subMixed
					tmp = append(tmp, 28) // mixed latch
					continue
				} else {
					tmp = append(tmp, 29) // punctuation switch
					tmp = append(tmp, punctMap[ch])
					break
				}
			}
			break
		case subLower:
			if isAlphaLower(ch) {
				if ch == ' ' {
					tmp = append(tmp, 26) //space
				} else {
					tmp = append(tmp, int(ch-'a'))
				}
			} else {
				if isAlphaUpper(ch) {
					tmp = append(tmp, 27) //upper switch
					tmp = append(tmp, int(ch-'A'))
					break
				} else if isMixed(ch) {
					submode = subMixed
					tmp = append(tmp, 28) //mixed latch
					continue
				} else {
					tmp = append(tmp, 29) //punctuation switch
					tmp = append(tmp, punctMap[ch])
					break
				}
			}
			break
		case subMixed:
			if isMixed(ch) {
				tmp = append(tmp, mixedMap[ch])
			} else {
				if isAlphaUpper(ch) {
					submode = subUpper
					tmp = append(tmp, 28) //upper latch
					continue
				} else if isAlphaLower(ch) {
					submode = subLower
					tmp = append(tmp, 27) //lower latch
					continue
				} else {
					if idx+1 < len(text) {
						next := text[idx+1]
						if isPunctuation(next) {
							submode = subPunct
							tmp = append(tmp, 25) //punctuation latch
							continue
						}
					}
					tmp = append(tmp, 29) //punctuation switch
					tmp = append(tmp, punctMap[ch])
				}
			}
			break
		default: //subPunct
			if isPunctuation(ch) {
				tmp = append(tmp, punctMap[ch])
			} else {
				submode = subUpper
				tmp = append(tmp, 29) //upper latch
				continue
			}
		}
		idx++
	}

	h := 0
	result := []int{}
	for i, val := range tmp {
		if i%2 != 0 {
			h = (h * 30) + val
			result = append(result, h)
		} else {
			h = val
		}
	}
	if len(tmp)%2 != 0 {
		result = append(result, (h*30)+29)
	}
	return submode, result
}

func determineConsecutiveBinaryCount(msg []byte) int {
	result := 0

	for i, _ := range msg {
		numericCount := determineConsecutiveDigitCount([]rune(string(msg[i:])))
		if numericCount >= min_numeric_count {
			break
		}
		textCount := determineConsecutiveTextCount([]rune(string(msg[i:])))
		if textCount > 5 {
			break
		}
		result++
	}
	return result
}

func encodeBinary(data []byte, startmode encodingMode) []int {
	result := []int{}

	count := len(data)
	if count == 1 && startmode == encText {
		result = append(result, shift_to_byte)
	} else if (count % 6) == 0 {
		result = append(result, latch_to_byte)
	} else {
		result = append(result, latch_to_byte_padded)
	}

	idx := 0
	// Encode sixpacks
	if count >= 6 {
		words := make([]int, 5)
		for (count - idx) >= 6 {
			var t int64 = 0
			for i := 0; i < 6; i++ {
				t = t << 8
				t += int64(data[idx+i])
			}
			for i := 0; i < 5; i++ {
				words[4-i] = int(t % 900)
				t = t / 900
			}
			result = append(result, words...)
			idx += 6
		}
	}
	//Encode rest (remaining n<5 bytes if any)
	for i := idx; i < count; i++ {
		result = append(result, int(data[i]&0xff))
	}
	return result
}

func highlevelEncode(dataStr string) ([]int, error) {
	encodingMode := encText
	textSubMode := subUpper

	result := []int{}

	data := []byte(dataStr)

	for len(data) > 0 {
		numericCount := determineConsecutiveDigitCount([]rune(string(data)))
		if numericCount >= min_numeric_count || numericCount == len(data) {
			result = append(result, latch_to_numeric)
			encodingMode = encNumeric
			textSubMode = subUpper
			numData, err := encodeNumeric([]rune(string(data[:numericCount])))
			if err != nil {
				return nil, err
			}
			result = append(result, numData...)
			data = data[numericCount:]
		} else {
			textCount := determineConsecutiveTextCount([]rune(string(data)))
			if textCount >= 5 || textCount == len(data) {
				if encodingMode != encText {
					result = append(result, latch_to_text)
					encodingMode = encText
					textSubMode = subUpper
				}
				var txtData []int
				textSubMode, txtData = encodeText([]rune(string(data[:textCount])), textSubMode)
				result = append(result, txtData...)
				data = data[textCount:]
			} else {
				binaryCount := determineConsecutiveBinaryCount(data)
				if binaryCount == 0 {
					binaryCount = 1
				}
				bytes := data[:binaryCount]
				if len(bytes) != 1 || encodingMode != encText {
					encodingMode = encBinary
					textSubMode = subUpper
				}
				byteData := encodeBinary(bytes, encodingMode)
				result = append(result, byteData...)
				data = data[binaryCount:]
			}
		}
	}

	return result, nil
}
