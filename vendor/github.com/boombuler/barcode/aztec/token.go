package aztec

import (
	"fmt"

	"github.com/boombuler/barcode/utils"
)

type token interface {
	fmt.Stringer
	prev() token
	appendTo(bits *utils.BitList, text []byte)
}

type simpleToken struct {
	token
	value    int
	bitCount byte
}

type binaryShiftToken struct {
	token
	bShiftStart   int
	bShiftByteCnt int
}

func newSimpleToken(prev token, value int, bitCount byte) token {
	return &simpleToken{prev, value, bitCount}
}
func newShiftToken(prev token, bShiftStart int, bShiftCnt int) token {
	return &binaryShiftToken{prev, bShiftStart, bShiftCnt}
}

func (st *simpleToken) prev() token {
	return st.token
}
func (st *simpleToken) appendTo(bits *utils.BitList, text []byte) {
	bits.AddBits(st.value, st.bitCount)
}
func (st *simpleToken) String() string {
	value := st.value & ((1 << st.bitCount) - 1)
	value |= 1 << st.bitCount
	return "<" + fmt.Sprintf("%b", value)[1:] + ">"
}

func (bst *binaryShiftToken) prev() token {
	return bst.token
}
func (bst *binaryShiftToken) appendTo(bits *utils.BitList, text []byte) {
	for i := 0; i < bst.bShiftByteCnt; i++ {
		if i == 0 || (i == 31 && bst.bShiftByteCnt <= 62) {
			// We need a header before the first character, and before
			// character 31 when the total byte code is <= 62
			bits.AddBits(31, 5) // BINARY_SHIFT
			if bst.bShiftByteCnt > 62 {
				bits.AddBits(bst.bShiftByteCnt-31, 16)
			} else if i == 0 {
				// 1 <= binaryShiftByteCode <= 62
				if bst.bShiftByteCnt < 31 {
					bits.AddBits(bst.bShiftByteCnt, 5)
				} else {
					bits.AddBits(31, 5)
				}
			} else {
				// 32 <= binaryShiftCount <= 62 and i == 31
				bits.AddBits(bst.bShiftByteCnt-31, 5)
			}
		}
		bits.AddByte(text[bst.bShiftStart+i])
	}
}

func (bst *binaryShiftToken) String() string {
	return fmt.Sprintf("<%d::%d>", bst.bShiftStart, (bst.bShiftStart + bst.bShiftByteCnt - 1))
}
