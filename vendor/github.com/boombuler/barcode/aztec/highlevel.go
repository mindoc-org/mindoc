package aztec

import (
	"github.com/boombuler/barcode/utils"
)

func highlevelEncode(data []byte) *utils.BitList {
	states := stateSlice{initialState}

	for index := 0; index < len(data); index++ {
		pairCode := 0
		nextChar := byte(0)
		if index+1 < len(data) {
			nextChar = data[index+1]
		}

		switch cur := data[index]; {
		case cur == '\r' && nextChar == '\n':
			pairCode = 2
		case cur == '.' && nextChar == ' ':
			pairCode = 3
		case cur == ',' && nextChar == ' ':
			pairCode = 4
		case cur == ':' && nextChar == ' ':
			pairCode = 5
		}
		if pairCode > 0 {
			// We have one of the four special PUNCT pairs.  Treat them specially.
			// Get a new set of states for the two new characters.
			states = updateStateListForPair(states, data, index, pairCode)
			index++
		} else {
			// Get a new set of states for the new character.
			states = updateStateListForChar(states, data, index)
		}
	}
	minBitCnt := int((^uint(0)) >> 1)
	var result *state = nil
	for _, s := range states {
		if s.bitCount < minBitCnt {
			minBitCnt = s.bitCount
			result = s
		}
	}
	if result != nil {
		return result.toBitList(data)
	} else {
		return new(utils.BitList)
	}
}

func simplifyStates(states stateSlice) stateSlice {
	var result stateSlice = nil
	for _, newState := range states {
		add := true
		var newResult stateSlice = nil

		for _, oldState := range result {
			if add && oldState.isBetterThanOrEqualTo(newState) {
				add = false
			}
			if !(add && newState.isBetterThanOrEqualTo(oldState)) {
				newResult = append(newResult, oldState)
			}
		}

		if add {
			result = append(newResult, newState)
		} else {
			result = newResult
		}

	}

	return result
}

// We update a set of states for a new character by updating each state
// for the new character, merging the results, and then removing the
// non-optimal states.
func updateStateListForChar(states stateSlice, data []byte, index int) stateSlice {
	var result stateSlice = nil
	for _, s := range states {
		if r := updateStateForChar(s, data, index); len(r) > 0 {
			result = append(result, r...)
		}
	}
	return simplifyStates(result)
}

// Return a set of states that represent the possible ways of updating this
// state for the next character.  The resulting set of states are added to
// the "result" list.
func updateStateForChar(s *state, data []byte, index int) stateSlice {
	var result stateSlice = nil
	ch := data[index]
	charInCurrentTable := charMap[s.mode][ch] > 0

	var stateNoBinary *state = nil
	for mode := mode_upper; mode <= mode_punct; mode++ {
		charInMode := charMap[mode][ch]
		if charInMode > 0 {
			if stateNoBinary == nil {
				// Only create stateNoBinary the first time it's required.
				stateNoBinary = s.endBinaryShift(index)
			}
			// Try generating the character by latching to its mode
			if !charInCurrentTable || mode == s.mode || mode == mode_digit {
				// If the character is in the current table, we don't want to latch to
				// any other mode except possibly digit (which uses only 4 bits).  Any
				// other latch would be equally successful *after* this character, and
				// so wouldn't save any bits.
				res := stateNoBinary.latchAndAppend(mode, charInMode)
				result = append(result, res)
			}
			// Try generating the character by switching to its mode.
			if _, ok := shiftTable[s.mode][mode]; !charInCurrentTable && ok {
				// It never makes sense to temporarily shift to another mode if the
				// character exists in the current mode.  That can never save bits.
				res := stateNoBinary.shiftAndAppend(mode, charInMode)
				result = append(result, res)
			}
		}
	}
	if s.bShiftByteCount > 0 || charMap[s.mode][ch] == 0 {
		// It's never worthwhile to go into binary shift mode if you're not already
		// in binary shift mode, and the character exists in your current mode.
		// That can never save bits over just outputting the char in the current mode.
		res := s.addBinaryShiftChar(index)
		result = append(result, res)
	}
	return result
}

// We update a set of states for a new character by updating each state
// for the new character, merging the results, and then removing the
// non-optimal states.
func updateStateListForPair(states stateSlice, data []byte, index int, pairCode int) stateSlice {
	var result stateSlice = nil
	for _, s := range states {
		if r := updateStateForPair(s, data, index, pairCode); len(r) > 0 {
			result = append(result, r...)
		}
	}
	return simplifyStates(result)
}

func updateStateForPair(s *state, data []byte, index int, pairCode int) stateSlice {
	var result stateSlice
	stateNoBinary := s.endBinaryShift(index)
	// Possibility 1.  Latch to MODE_PUNCT, and then append this code
	result = append(result, stateNoBinary.latchAndAppend(mode_punct, pairCode))
	if s.mode != mode_punct {
		// Possibility 2.  Shift to MODE_PUNCT, and then append this code.
		// Every state except MODE_PUNCT (handled above) can shift
		result = append(result, stateNoBinary.shiftAndAppend(mode_punct, pairCode))
	}
	if pairCode == 3 || pairCode == 4 {
		// both characters are in DIGITS.  Sometimes better to just add two digits
		digitState := stateNoBinary.
			latchAndAppend(mode_digit, 16-pairCode). // period or comma in DIGIT
			latchAndAppend(mode_digit, 1)            // space in DIGIT
		result = append(result, digitState)
	}
	if s.bShiftByteCount > 0 {
		// It only makes sense to do the characters as binary if we're already
		// in binary mode.
		result = append(result, s.addBinaryShiftChar(index).addBinaryShiftChar(index+1))
	}
	return result
}
