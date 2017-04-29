// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"bytes"
	"fmt"
)

type header struct {
	magic    []byte
	ver      byte
	reserved []byte
}

func (h *header) rd(b []byte) error {
	if len(b) != 16 {
		panic("internal error")
	}

	if h.magic = b[:4]; !bytes.Equal(h.magic, []byte(magic)) {
		return fmt.Errorf("Unknown file format")
	}

	b = b[4:]
	h.ver = b[0]
	h.reserved = b[1:]
	return nil
}
