// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccitt

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"
)

func compareImages(t *testing.T, img0 image.Image, img1 image.Image) {
	t.Helper()

	b0 := img0.Bounds()
	b1 := img1.Bounds()
	if b0 != b1 {
		t.Fatalf("bounds differ: %v vs %v", b0, b1)
	}

	for y := b0.Min.Y; y < b0.Max.Y; y++ {
		for x := b0.Min.X; x < b0.Max.X; x++ {
			c0 := img0.At(x, y)
			c1 := img1.At(x, y)
			if c0 != c1 {
				t.Fatalf("pixel at (%d, %d) differs: %v vs %v", x, y, c0, c1)
			}
		}
	}
}

func decodePNG(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func TestMaxCodeLength(t *testing.T) {
	br := bitReader{}
	size := unsafe.Sizeof(br.bits)
	size *= 8 // Convert from bytes to bits.

	// Check that the size of the bitReader.bits field is large enough to hold
	// nextBitMaxNBits bits.
	if size < nextBitMaxNBits {
		t.Fatalf("size: got %d, want >= %d", size, nextBitMaxNBits)
	}

	// Check that bitReader.nextBit will always leave enough spare bits in the
	// bitReader.bits field such that the decode function can unread up to
	// maxCodeLength bits.
	if want := size - nextBitMaxNBits; maxCodeLength > want {
		t.Fatalf("maxCodeLength: got %d, want <= %d", maxCodeLength, want)
	}

	// The decode function also assumes that, when saving bits to possibly
	// unread later, those bits fit inside a uint32.
	if maxCodeLength > 32 {
		t.Fatalf("maxCodeLength: got %d, want <= %d", maxCodeLength, 32)
	}
}

func testDecodeTable(t *testing.T, decodeTable [][2]int16, codes []code, values []uint32) {
	// Build a map from values to codes.
	m := map[uint32]string{}
	for _, code := range codes {
		m[code.val] = code.str
	}

	// Build the encoded form of those values in LSB order.
	enc := []byte(nil)
	bits := uint8(0)
	nBits := uint32(0)
	for _, v := range values {
		code := m[v]
		if code == "" {
			panic("unmapped code")
		}
		for _, c := range code {
			bits |= uint8(c&1) << nBits
			nBits++
			if nBits == 8 {
				enc = append(enc, bits)
				bits = 0
				nBits = 0
			}
		}
	}
	if nBits > 0 {
		enc = append(enc, bits)
	}

	// Decode that encoded form.
	got := []uint32(nil)
	r := &bitReader{
		r: bytes.NewReader(enc),
	}
	finalValue := values[len(values)-1]
	for {
		v, err := decode(r, decodeTable)
		if err != nil {
			t.Fatalf("after got=%d: %v", got, err)
		}
		got = append(got, v)
		if v == finalValue {
			break
		}
	}

	// Check that the round-tripped values were unchanged.
	if !reflect.DeepEqual(got, values) {
		t.Fatalf("\ngot:  %v\nwant: %v", got, values)
	}
}

func TestModeDecodeTable(t *testing.T) {
	testDecodeTable(t, modeDecodeTable[:], modeCodes, []uint32{
		modePass,
		modeV0,
		modeV0,
		modeVL1,
		modeVR3,
		modeVL2,
		modeExt,
		modeVL1,
		modeH,
		modeVL1,
		modeVL1,
		modeEOL,
	})
}

func TestWhiteDecodeTable(t *testing.T) {
	testDecodeTable(t, whiteDecodeTable[:], whiteCodes, []uint32{
		0, 1, 256, 7, 128, 3, 2560,
	})
}

func TestBlackDecodeTable(t *testing.T) {
	testDecodeTable(t, blackDecodeTable[:], blackCodes, []uint32{
		63, 64, 63, 64, 64, 63, 22, 1088, 2048, 7, 6, 5, 4, 3, 2, 1, 0,
	})
}

func TestDecodeInvalidCode(t *testing.T) {
	// The bit stream is:
	// 1 010 000000011011
	// Packing that LSB-first gives:
	// 0b_1101_1000_0000_0101
	src := []byte{0x05, 0xD8}

	decodeTable := modeDecodeTable[:]
	r := &bitReader{
		r: bytes.NewReader(src),
	}

	// "1" decodes to the value 2.
	if v, err := decode(r, decodeTable); v != 2 || err != nil {
		t.Fatalf("decode #0: got (%v, %v), want (2, nil)", v, err)
	}

	// "010" decodes to the value 6.
	if v, err := decode(r, decodeTable); v != 6 || err != nil {
		t.Fatalf("decode #0: got (%v, %v), want (6, nil)", v, err)
	}

	// "00000001" is an invalid code.
	if v, err := decode(r, decodeTable); v != 0 || err != errInvalidCode {
		t.Fatalf("decode #0: got (%v, %v), want (0, %v)", v, err, errInvalidCode)
	}

	// The bitReader should not have advanced after encountering an invalid
	// code. The remaining bits should be "000000011011".
	remaining := []byte(nil)
	for {
		bit, err := r.nextBit()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("nextBit: %v", err)
		}
		remaining = append(remaining, uint8('0'+bit))
	}
	if got, want := string(remaining), "000000011011"; got != want {
		t.Fatalf("remaining bits: got %q, want %q", got, want)
	}
}

func TestReadRegular(t *testing.T) { testRead(t, false) }
func TestReadInvert(t *testing.T)  { testRead(t, true) }

func testRead(t *testing.T, invert bool) {
	t.Helper()

	const width, height = 153, 55
	opts := &Options{
		Invert: invert,
	}

	got := ""
	{
		f, err := os.Open("testdata/bw-gopher.ccitt_group3")
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer f.Close()
		gotBytes, err := ioutil.ReadAll(NewReader(f, MSB, Group3, width, height, opts))
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		got = string(gotBytes)
	}

	want := ""
	{
		img, err := decodePNG("testdata/bw-gopher.png")
		if err != nil {
			t.Fatalf("decodePNG: %v", err)
		}
		gray, ok := img.(*image.Gray)
		if !ok {
			t.Fatalf("decodePNG: got %T, want *image.Gray", img)
		}
		bounds := gray.Bounds()
		if w := bounds.Dx(); w != width {
			t.Fatalf("width: got %d, want %d", w, width)
		}
		if h := bounds.Dy(); h != height {
			t.Fatalf("height: got %d, want %d", h, height)
		}

		// Prepare to extend each row's width to a multiple of 8, to simplify
		// packing from 1 byte per pixel to 1 bit per pixel.
		extended := make([]byte, (width+7)&^7)

		wantBytes := []byte(nil)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rowPix := gray.Pix[(y-bounds.Min.Y)*gray.Stride:]
			rowPix = rowPix[:width]
			copy(extended, rowPix)

			// Pack from 1 byte per pixel to 1 bit per pixel, MSB first.
			byteValue := uint8(0)
			for x, pixel := range extended {
				byteValue |= (pixel & 0x80) >> uint(x&7)
				if (x & 7) == 7 {
					wantBytes = append(wantBytes, byteValue)
					byteValue = 0
				}
			}
		}
		if invert {
			invertBytes(wantBytes)
		}
		want = string(wantBytes)
	}

	// We expect a width of 153 pixels, which is 20 bytes per row (at 1 bit per
	// pixel, plus 7 final bits of padding). Check that want is 20 * height
	// bytes long, and if got != want, format them to split at every 20 bytes.

	if n := len(want); n != 20*height {
		t.Fatalf("len(want): got %d, want %d", n, 20*height)
	}

	format := func(s string) string {
		b := []byte(nil)
		for row := 0; len(s) >= 20; row++ {
			b = append(b, fmt.Sprintf("row%02d: %02X\n", row, s[:20])...)
			s = s[20:]
		}
		if len(s) > 0 {
			b = append(b, fmt.Sprintf("%02X\n", s)...)
		}
		return string(b)
	}

	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", format(got), format(want))
	}
}

func TestDecodeIntoGray(t *testing.T) {
	for _, tt := range []struct {
		fileName string
		sf       SubFormat
		w, h     int
	}{
		{"testdata/bw-gopher.ccitt_group3", Group3, 153, 55},
		{"testdata/bw-gopher.ccitt_group4", Group4, 153, 55},
	} {
		t.Run(tt.fileName, func(t *testing.T) {
			testDecodeIntoGray(t, tt.fileName, MSB, tt.sf, tt.w, tt.h, nil)
		})
	}
}

func testDecodeIntoGray(t *testing.T, fileName string, order Order, sf SubFormat, width int, height int, opts *Options) {
	t.Helper()

	f, err := os.Open(filepath.FromSlash(fileName))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer f.Close()

	got := image.NewGray(image.Rect(0, 0, width, height))
	if err := DecodeIntoGray(got, f, order, sf, opts); err != nil {
		t.Fatalf("DecodeIntoGray: %v", err)
	}

	baseName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	want, err := decodePNG(baseName + ".png")
	if err != nil {
		t.Fatal(err)
	}

	compareImages(t, got, want)
}
