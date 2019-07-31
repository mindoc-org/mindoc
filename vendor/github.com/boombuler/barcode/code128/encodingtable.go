package code128

var encodingTable = [107][]bool{
	[]bool{true, true, false, true, true, false, false, true, true, false, false},
	[]bool{true, true, false, false, true, true, false, true, true, false, false},
	[]bool{true, true, false, false, true, true, false, false, true, true, false},
	[]bool{true, false, false, true, false, false, true, true, false, false, false},
	[]bool{true, false, false, true, false, false, false, true, true, false, false},
	[]bool{true, false, false, false, true, false, false, true, true, false, false},
	[]bool{true, false, false, true, true, false, false, true, false, false, false},
	[]bool{true, false, false, true, true, false, false, false, true, false, false},
	[]bool{true, false, false, false, true, true, false, false, true, false, false},
	[]bool{true, true, false, false, true, false, false, true, false, false, false},
	[]bool{true, true, false, false, true, false, false, false, true, false, false},
	[]bool{true, true, false, false, false, true, false, false, true, false, false},
	[]bool{true, false, true, true, false, false, true, true, true, false, false},
	[]bool{true, false, false, true, true, false, true, true, true, false, false},
	[]bool{true, false, false, true, true, false, false, true, true, true, false},
	[]bool{true, false, true, true, true, false, false, true, true, false, false},
	[]bool{true, false, false, true, true, true, false, true, true, false, false},
	[]bool{true, false, false, true, true, true, false, false, true, true, false},
	[]bool{true, true, false, false, true, true, true, false, false, true, false},
	[]bool{true, true, false, false, true, false, true, true, true, false, false},
	[]bool{true, true, false, false, true, false, false, true, true, true, false},
	[]bool{true, true, false, true, true, true, false, false, true, false, false},
	[]bool{true, true, false, false, true, true, true, false, true, false, false},
	[]bool{true, true, true, false, true, true, false, true, true, true, false},
	[]bool{true, true, true, false, true, false, false, true, true, false, false},
	[]bool{true, true, true, false, false, true, false, true, true, false, false},
	[]bool{true, true, true, false, false, true, false, false, true, true, false},
	[]bool{true, true, true, false, true, true, false, false, true, false, false},
	[]bool{true, true, true, false, false, true, true, false, true, false, false},
	[]bool{true, true, true, false, false, true, true, false, false, true, false},
	[]bool{true, true, false, true, true, false, true, true, false, false, false},
	[]bool{true, true, false, true, true, false, false, false, true, true, false},
	[]bool{true, true, false, false, false, true, true, false, true, true, false},
	[]bool{true, false, true, false, false, false, true, true, false, false, false},
	[]bool{true, false, false, false, true, false, true, true, false, false, false},
	[]bool{true, false, false, false, true, false, false, false, true, true, false},
	[]bool{true, false, true, true, false, false, false, true, false, false, false},
	[]bool{true, false, false, false, true, true, false, true, false, false, false},
	[]bool{true, false, false, false, true, true, false, false, false, true, false},
	[]bool{true, true, false, true, false, false, false, true, false, false, false},
	[]bool{true, true, false, false, false, true, false, true, false, false, false},
	[]bool{true, true, false, false, false, true, false, false, false, true, false},
	[]bool{true, false, true, true, false, true, true, true, false, false, false},
	[]bool{true, false, true, true, false, false, false, true, true, true, false},
	[]bool{true, false, false, false, true, true, false, true, true, true, false},
	[]bool{true, false, true, true, true, false, true, true, false, false, false},
	[]bool{true, false, true, true, true, false, false, false, true, true, false},
	[]bool{true, false, false, false, true, true, true, false, true, true, false},
	[]bool{true, true, true, false, true, true, true, false, true, true, false},
	[]bool{true, true, false, true, false, false, false, true, true, true, false},
	[]bool{true, true, false, false, false, true, false, true, true, true, false},
	[]bool{true, true, false, true, true, true, false, true, false, false, false},
	[]bool{true, true, false, true, true, true, false, false, false, true, false},
	[]bool{true, true, false, true, true, true, false, true, true, true, false},
	[]bool{true, true, true, false, true, false, true, true, false, false, false},
	[]bool{true, true, true, false, true, false, false, false, true, true, false},
	[]bool{true, true, true, false, false, false, true, false, true, true, false},
	[]bool{true, true, true, false, true, true, false, true, false, false, false},
	[]bool{true, true, true, false, true, true, false, false, false, true, false},
	[]bool{true, true, true, false, false, false, true, true, false, true, false},
	[]bool{true, true, true, false, true, true, true, true, false, true, false},
	[]bool{true, true, false, false, true, false, false, false, false, true, false},
	[]bool{true, true, true, true, false, false, false, true, false, true, false},
	[]bool{true, false, true, false, false, true, true, false, false, false, false},
	[]bool{true, false, true, false, false, false, false, true, true, false, false},
	[]bool{true, false, false, true, false, true, true, false, false, false, false},
	[]bool{true, false, false, true, false, false, false, false, true, true, false},
	[]bool{true, false, false, false, false, true, false, true, true, false, false},
	[]bool{true, false, false, false, false, true, false, false, true, true, false},
	[]bool{true, false, true, true, false, false, true, false, false, false, false},
	[]bool{true, false, true, true, false, false, false, false, true, false, false},
	[]bool{true, false, false, true, true, false, true, false, false, false, false},
	[]bool{true, false, false, true, true, false, false, false, false, true, false},
	[]bool{true, false, false, false, false, true, true, false, true, false, false},
	[]bool{true, false, false, false, false, true, true, false, false, true, false},
	[]bool{true, true, false, false, false, false, true, false, false, true, false},
	[]bool{true, true, false, false, true, false, true, false, false, false, false},
	[]bool{true, true, true, true, false, true, true, true, false, true, false},
	[]bool{true, true, false, false, false, false, true, false, true, false, false},
	[]bool{true, false, false, false, true, true, true, true, false, true, false},
	[]bool{true, false, true, false, false, true, true, true, true, false, false},
	[]bool{true, false, false, true, false, true, true, true, true, false, false},
	[]bool{true, false, false, true, false, false, true, true, true, true, false},
	[]bool{true, false, true, true, true, true, false, false, true, false, false},
	[]bool{true, false, false, true, true, true, true, false, true, false, false},
	[]bool{true, false, false, true, true, true, true, false, false, true, false},
	[]bool{true, true, true, true, false, true, false, false, true, false, false},
	[]bool{true, true, true, true, false, false, true, false, true, false, false},
	[]bool{true, true, true, true, false, false, true, false, false, true, false},
	[]bool{true, true, false, true, true, false, true, true, true, true, false},
	[]bool{true, true, false, true, true, true, true, false, true, true, false},
	[]bool{true, true, true, true, false, true, true, false, true, true, false},
	[]bool{true, false, true, false, true, true, true, true, false, false, false},
	[]bool{true, false, true, false, false, false, true, true, true, true, false},
	[]bool{true, false, false, false, true, false, true, true, true, true, false},
	[]bool{true, false, true, true, true, true, false, true, false, false, false},
	[]bool{true, false, true, true, true, true, false, false, false, true, false},
	[]bool{true, true, true, true, false, true, false, true, false, false, false},
	[]bool{true, true, true, true, false, true, false, false, false, true, false},
	[]bool{true, false, true, true, true, false, true, true, true, true, false},
	[]bool{true, false, true, true, true, true, false, true, true, true, false},
	[]bool{true, true, true, false, true, false, true, true, true, true, false},
	[]bool{true, true, true, true, false, true, false, true, true, true, false},
	[]bool{true, true, false, true, false, false, false, false, true, false, false},
	[]bool{true, true, false, true, false, false, true, false, false, false, false},
	[]bool{true, true, false, true, false, false, true, true, true, false, false},
	[]bool{true, true, false, false, false, true, true, true, false, true, false, true, true},
}

const startASymbol byte = 103
const startBSymbol byte = 104
const startCSymbol byte = 105

const codeASymbol byte = 101
const codeBSymbol byte = 100
const codeCSymbol byte = 99

const stopSymbol byte = 106

const (
	// FNC1 - Special Function 1
	FNC1 = '\u00f1'
	// FNC2 - Special Function 2
	FNC2 = '\u00f2'
	// FNC3 - Special Function 3
	FNC3 = '\u00f3'
	// FNC4 - Special Function 4
	FNC4 = '\u00f4'
)

const abTable = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"
const bTable = abTable + "`abcdefghijklmnopqrstuvwxyz{|}~\u007F"
const aOnlyTable = "\u0000\u0001\u0002\u0003\u0004" + // NUL, SOH, STX, ETX, EOT
	"\u0005\u0006\u0007\u0008\u0009" + // ENQ, ACK, BEL, BS,  HT
	"\u000A\u000B\u000C\u000D\u000E" + // LF,  VT,  FF,  CR,  SO
	"\u000F\u0010\u0011\u0012\u0013" + // SI,  DLE, DC1, DC2, DC3
	"\u0014\u0015\u0016\u0017\u0018" + // DC4, NAK, SYN, ETB, CAN
	"\u0019\u001A\u001B\u001C\u001D" + // EM,  SUB, ESC, FS,  GS
	"\u001E\u001F" // RS,  US
const aTable = abTable + aOnlyTable
