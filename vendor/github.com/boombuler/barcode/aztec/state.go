package aztec

import (
	"fmt"

	"github.com/boombuler/barcode/utils"
)

type encodingMode byte

const (
	mode_upper encodingMode = iota // 5 bits
	mode_lower                     // 5 bits
	mode_digit                     // 4 bits
	mode_mixed                     // 5 bits
	mode_punct                     // 5 bits
)

var (
	// The Latch Table shows, for each pair of Modes, the optimal method for
	// getting from one mode to another.  In the worst possible case, this can
	// be up to 14 bits.  In the best possible case, we are already there!
	// The high half-word of each entry gives the number of bits.
	// The low half-word of each entry are the actual bits necessary to change
	latchTable = map[encodingMode]map[encodingMode]int{
		mode_upper: {
			mode_upper: 0,
			mode_lower: (5 << 16) + 28,
			mode_digit: (5 << 16) + 30,
			mode_mixed: (5 << 16) + 29,
			mode_punct: (10 << 16) + (29 << 5) + 30,
		},
		mode_lower: {
			mode_upper: (9 << 16) + (30 << 4) + 14,
			mode_lower: 0,
			mode_digit: (5 << 16) + 30,
			mode_mixed: (5 << 16) + 29,
			mode_punct: (10 << 16) + (29 << 5) + 30,
		},
		mode_digit: {
			mode_upper: (4 << 16) + 14,
			mode_lower: (9 << 16) + (14 << 5) + 28,
			mode_digit: 0,
			mode_mixed: (9 << 16) + (14 << 5) + 29,
			mode_punct: (14 << 16) + (14 << 10) + (29 << 5) + 30,
		},
		mode_mixed: {
			mode_upper: (5 << 16) + 29,
			mode_lower: (5 << 16) + 28,
			mode_digit: (10 << 16) + (29 << 5) + 30,
			mode_mixed: 0,
			mode_punct: (5 << 16) + 30,
		},
		mode_punct: {
			mode_upper: (5 << 16) + 31,
			mode_lower: (10 << 16) + (31 << 5) + 28,
			mode_digit: (10 << 16) + (31 << 5) + 30,
			mode_mixed: (10 << 16) + (31 << 5) + 29,
			mode_punct: 0,
		},
	}
	// A map showing the available shift codes.  (The shifts to BINARY are not shown)
	shiftTable = map[encodingMode]map[encodingMode]int{
		mode_upper: {
			mode_punct: 0,
		},
		mode_lower: {
			mode_punct: 0,
			mode_upper: 28,
		},
		mode_mixed: {
			mode_punct: 0,
		},
		mode_digit: {
			mode_punct: 0,
			mode_upper: 15,
		},
	}
	charMap map[encodingMode][]int
)

type state struct {
	mode            encodingMode
	tokens          token
	bShiftByteCount int
	bitCount        int
}
type stateSlice []*state

var initialState *state = &state{
	mode:            mode_upper,
	tokens:          nil,
	bShiftByteCount: 0,
	bitCount:        0,
}

func init() {
	charMap = make(map[encodingMode][]int)
	charMap[mode_upper] = make([]int, 256)
	charMap[mode_lower] = make([]int, 256)
	charMap[mode_digit] = make([]int, 256)
	charMap[mode_mixed] = make([]int, 256)
	charMap[mode_punct] = make([]int, 256)

	charMap[mode_upper][' '] = 1
	for c := 'A'; c <= 'Z'; c++ {
		charMap[mode_upper][int(c)] = int(c - 'A' + 2)
	}

	charMap[mode_lower][' '] = 1
	for c := 'a'; c <= 'z'; c++ {
		charMap[mode_lower][c] = int(c - 'a' + 2)
	}
	charMap[mode_digit][' '] = 1
	for c := '0'; c <= '9'; c++ {
		charMap[mode_digit][c] = int(c - '0' + 2)
	}
	charMap[mode_digit][','] = 12
	charMap[mode_digit]['.'] = 13

	mixedTable := []int{
		0, ' ', 1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 27, 28, 29, 30, 31, '@', '\\', '^',
		'_', '`', '|', '~', 127,
	}
	for i, v := range mixedTable {
		charMap[mode_mixed][v] = i
	}

	punctTable := []int{
		0, '\r', 0, 0, 0, 0, '!', '\'', '#', '$', '%', '&', '\'',
		'(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>', '?',
		'[', ']', '{', '}',
	}
	for i, v := range punctTable {
		if v > 0 {
			charMap[mode_punct][v] = i
		}
	}
}

func (em encodingMode) BitCount() byte {
	if em == mode_digit {
		return 4
	}
	return 5
}

// Create a new state representing this state with a latch to a (not
// necessary different) mode, and then a code.
func (s *state) latchAndAppend(mode encodingMode, value int) *state {
	bitCount := s.bitCount
	tokens := s.tokens

	if mode != s.mode {
		latch := latchTable[s.mode][mode]
		tokens = newSimpleToken(tokens, latch&0xFFFF, byte(latch>>16))
		bitCount += latch >> 16
	}
	tokens = newSimpleToken(tokens, value, mode.BitCount())
	return &state{
		mode:            mode,
		tokens:          tokens,
		bShiftByteCount: 0,
		bitCount:        bitCount + int(mode.BitCount()),
	}
}

// Create a new state representing this state, with a temporary shift
// to a different mode to output a single value.
func (s *state) shiftAndAppend(mode encodingMode, value int) *state {
	tokens := s.tokens

	// Shifts exist only to UPPER and PUNCT, both with tokens size 5.
	tokens = newSimpleToken(tokens, shiftTable[s.mode][mode], s.mode.BitCount())
	tokens = newSimpleToken(tokens, value, 5)

	return &state{
		mode:            s.mode,
		tokens:          tokens,
		bShiftByteCount: 0,
		bitCount:        s.bitCount + int(s.mode.BitCount()) + 5,
	}
}

// Create a new state representing this state, but an additional character
// output in Binary Shift mode.
func (s *state) addBinaryShiftChar(index int) *state {
	tokens := s.tokens
	mode := s.mode
	bitCnt := s.bitCount
	if s.mode == mode_punct || s.mode == mode_digit {
		latch := latchTable[s.mode][mode_upper]
		tokens = newSimpleToken(tokens, latch&0xFFFF, byte(latch>>16))
		bitCnt += latch >> 16
		mode = mode_upper
	}
	deltaBitCount := 8
	if s.bShiftByteCount == 0 || s.bShiftByteCount == 31 {
		deltaBitCount = 18
	} else if s.bShiftByteCount == 62 {
		deltaBitCount = 9
	}
	result := &state{
		mode:            mode,
		tokens:          tokens,
		bShiftByteCount: s.bShiftByteCount + 1,
		bitCount:        bitCnt + deltaBitCount,
	}
	if result.bShiftByteCount == 2047+31 {
		// The string is as long as it's allowed to be.  We should end it.
		result = result.endBinaryShift(index + 1)
	}

	return result
}

// Create the state identical to this one, but we are no longer in
// Binary Shift mode.
func (s *state) endBinaryShift(index int) *state {
	if s.bShiftByteCount == 0 {
		return s
	}
	tokens := newShiftToken(s.tokens, index-s.bShiftByteCount, s.bShiftByteCount)
	return &state{
		mode:            s.mode,
		tokens:          tokens,
		bShiftByteCount: 0,
		bitCount:        s.bitCount,
	}
}

// Returns true if "this" state is better (or equal) to be in than "that"
// state under all possible circumstances.
func (this *state) isBetterThanOrEqualTo(other *state) bool {
	mySize := this.bitCount + (latchTable[this.mode][other.mode] >> 16)

	if other.bShiftByteCount > 0 && (this.bShiftByteCount == 0 || this.bShiftByteCount > other.bShiftByteCount) {
		mySize += 10 // Cost of entering Binary Shift mode.
	}
	return mySize <= other.bitCount
}

func (s *state) toBitList(text []byte) *utils.BitList {
	tokens := make([]token, 0)
	se := s.endBinaryShift(len(text))

	for t := se.tokens; t != nil; t = t.prev() {
		tokens = append(tokens, t)
	}
	res := new(utils.BitList)
	for i := len(tokens) - 1; i >= 0; i-- {
		tokens[i].appendTo(res, text)
	}
	return res
}

func (s *state) String() string {
	tokens := make([]token, 0)
	for t := s.tokens; t != nil; t = t.prev() {
		tokens = append([]token{t}, tokens...)
	}
	return fmt.Sprintf("M:%d bits=%d bytes=%d: %v", s.mode, s.bitCount, s.bShiftByteCount, tokens)
}
