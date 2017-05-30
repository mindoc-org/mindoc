// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfnt

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// Platform IDs and Platform Specific IDs as per
// https://www.microsoft.com/typography/otspec/name.htm
const (
	pidUnicode   = 0
	pidMacintosh = 1
	pidWindows   = 3

	psidUnicode2BMPOnly        = 3
	psidUnicode2FullRepertoire = 4

	psidMacintoshRoman = 0

	psidWindowsUCS2 = 1
	psidWindowsUCS4 = 10
)

// platformEncodingWidth returns the number of bytes per character assumed by
// the given Platform ID and Platform Specific ID.
//
// Very old fonts, from before Unicode was widely adopted, assume only 1 byte
// per character: a character map.
//
// Old fonts, from when Unicode meant the Basic Multilingual Plane (BMP),
// assume that 2 bytes per character is sufficient.
//
// Recent fonts naturally support the full range of Unicode code points, which
// can take up to 4 bytes per character. Such fonts might still choose one of
// the legacy encodings if e.g. their repertoire is limited to the BMP, for
// greater compatibility with older software, or because the resultant file
// size can be smaller.
func platformEncodingWidth(pid, psid uint16) int {
	switch pid {
	case pidUnicode:
		switch psid {
		case psidUnicode2BMPOnly:
			return 2
		case psidUnicode2FullRepertoire:
			return 4
		}

	case pidMacintosh:
		switch psid {
		case psidMacintoshRoman:
			return 1
		}

	case pidWindows:
		switch psid {
		case psidWindowsUCS2:
			return 2
		case psidWindowsUCS4:
			return 4
		}
	}
	return 0
}

// The various cmap formats are described at
// https://www.microsoft.com/typography/otspec/cmap.htm

var supportedCmapFormat = func(format, pid, psid uint16) bool {
	switch format {
	case 0:
		return pid == pidMacintosh && psid == psidMacintoshRoman
	case 4:
		return true
	case 12:
		// TODO: implement.
	}
	return false
}

func (f *Font) makeCachedGlyphIndex(buf []byte, offset, length uint32, format uint16) ([]byte, error) {
	switch format {
	case 0:
		return f.makeCachedGlyphIndexFormat0(buf, offset, length)
	case 4:
		return f.makeCachedGlyphIndexFormat4(buf, offset, length)
	case 12:
		// TODO: implement, including a cmapEntry32 type (32, not 16).
	}
	panic("unreachable")
}

func (f *Font) makeCachedGlyphIndexFormat0(buf []byte, offset, length uint32) ([]byte, error) {
	if length != 6+256 || offset+length > f.cmap.length {
		return nil, errInvalidCmapTable
	}
	var err error
	buf, err = f.src.view(buf, int(f.cmap.offset+offset), int(length))
	if err != nil {
		return nil, err
	}
	var table [256]byte
	copy(table[:], buf[6:])
	f.cached.glyphIndex = func(f *Font, b *Buffer, r rune) (GlyphIndex, error) {
		// TODO: for this closure to be goroutine-safe, the
		// golang.org/x/text/encoding/charmap API needs to allocate a new
		// Encoder and new []byte buffers, for every call to this closure, even
		// though all we want to do is to encode one rune as one byte. We could
		// possibly add some fields in the Buffer struct to re-use these
		// allocations, but a better solution is to improve the charmap API.
		var dst, src [utf8.UTFMax]byte
		n := utf8.EncodeRune(src[:], r)
		_, _, err = charmap.Macintosh.NewEncoder().Transform(dst[:], src[:n], true)
		if err != nil {
			// The source rune r is not representable in the Macintosh-Roman encoding.
			return 0, nil
		}
		return GlyphIndex(table[dst[0]]), nil
	}
	return buf, nil
}

func (f *Font) makeCachedGlyphIndexFormat4(buf []byte, offset, length uint32) ([]byte, error) {
	const headerSize = 14
	if offset+headerSize > f.cmap.length {
		return nil, errInvalidCmapTable
	}
	var err error
	buf, err = f.src.view(buf, int(f.cmap.offset+offset), headerSize)
	if err != nil {
		return nil, err
	}
	offset += headerSize

	segCount := u16(buf[6:])
	if segCount&1 != 0 {
		return nil, errInvalidCmapTable
	}
	segCount /= 2
	if segCount > maxCmapSegments {
		return nil, errUnsupportedNumberOfCmapSegments
	}

	eLength := 8*uint32(segCount) + 2
	if offset+eLength > f.cmap.length {
		return nil, errInvalidCmapTable
	}
	buf, err = f.src.view(buf, int(f.cmap.offset+offset), int(eLength))
	if err != nil {
		return nil, err
	}
	offset += eLength

	entries := make([]cmapEntry16, segCount)
	for i := range entries {
		entries[i] = cmapEntry16{
			end:    u16(buf[0*len(entries)+0+2*i:]),
			start:  u16(buf[2*len(entries)+2+2*i:]),
			delta:  u16(buf[4*len(entries)+2+2*i:]),
			offset: u16(buf[6*len(entries)+2+2*i:]),
		}
	}

	f.cached.glyphIndex = func(f *Font, b *Buffer, r rune) (GlyphIndex, error) {
		if uint32(r) > 0xffff {
			return 0, nil
		}

		c := uint16(r)
		for i, j := 0, len(entries); i < j; {
			h := i + (j-i)/2
			entry := &entries[h]
			if c < entry.start {
				j = h
			} else if entry.end < c {
				i = h + 1
			} else if entry.offset == 0 {
				return GlyphIndex(c + entry.delta), nil
			} else {
				// TODO: support the glyphIdArray as per
				// https://www.microsoft.com/typography/OTSPEC/cmap.htm
				//
				// This will probably use the *Font and *Buffer arguments.
				return 0, errUnsupportedCmapFormat
			}
		}
		return 0, nil
	}
	return buf, nil
}

type cmapEntry16 struct {
	end, start, delta, offset uint16
}
