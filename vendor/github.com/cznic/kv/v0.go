// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"os"

	"github.com/cznic/lldb"
)

func open00(name string, in *DB) (db *DB, err error) {
	db = in
	if db.alloc, err = lldb.NewAllocator(lldb.NewInnerFiler(db.filer, 16), &lldb.Options{}); err != nil {
		return nil, &os.PathError{Op: "kv.open00", Path: name, Err: err}
	}

	db.alloc.Compress = true
	return
}
