// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cznic/lldb"
)

func verifyAllocator(a *lldb.Allocator) error {
	bits, err := ioutil.TempFile("", "kv-verify-")
	if err != nil {
		return err
	}

	sf := lldb.NewSimpleFileFiler(bits)

	defer func() {
		nm := bits.Name()
		sf.Close()
		os.Remove(nm)
	}()

	var lerr error
	if err = a.Verify(
		sf,
		func(err error) bool {
			lerr = err
			return false
		},
		nil,
	); err != nil {
		return err
	}

	if lerr != nil {
		return lerr
	}

	t, err := lldb.OpenBTree(a, nil, 1)
	if err != nil {
		return err
	}

	e, err := t.SeekFirst()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return err
	}

	for {
		_, _, err := e.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}

func verifyDbFile(fn string) error {
	f, err := os.OpenFile(fn, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	sf := lldb.NewSimpleFileFiler(f)
	if f == nil {
		return fmt.Errorf("cannot create %s", fn)
	}

	defer sf.Close()

	a, err := lldb.NewAllocator(lldb.NewInnerFiler(sf, 16), &lldb.Options{})
	if err != nil {
		return err
	}

	return verifyAllocator(a)
}
