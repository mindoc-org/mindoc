// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/cznic/fileutil"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/lldb"
)

const (
	magic = "\x60\xdbKV"
)

const (
	stDisabled = iota // stDisabled must be zero
	stIdle
	stCollecting
	stIdleArmed
	stCollectingArmed
	stCollectingTriggered
	stEndUpdateFailed
)

func init() {
	if stDisabled != 0 {
		panic("stDisabled != 0")
	}
}

// DB represents the database (the KV store).
type DB struct {
	acidNest      int             // Grace period nesting level
	acidState     int             // Grace period FSM state.
	acidTimer     *time.Timer     // Grace period timer
	alloc         *lldb.Allocator // The machinery. Wraps filer
	bkl           sync.Mutex      // Big Kernel Lock
	closeMu       sync.Mutex      // Close() coordination
	closed        bool            // it was
	filer         lldb.Filer      // Wraps f
	gracePeriod   time.Duration   // WAL grace period
	isMem         bool            // No signal capture
	lastCommitErr error           // from failed EndUpdate
	lock          io.Closer       // The DB file lock
	opts          *Options
	root          *lldb.BTree // The KV layer
	wal           *os.File    // WAL if any
}

// CreateFromFiler is like Create but accepts an arbitrary backing storage
// provided by filer.
//
// For the meaning of opts please see documentation of Options.
func CreateFromFiler(filer lldb.Filer, opts *Options) (db *DB, err error) {
	opts = opts.clone()
	opts._ACID = _ACIDFull
	return create(filer, opts, false)
}

// Create creates the named DB file mode 0666 (before umask). The file must not
// already exist. If successful, methods on the returned DB can be used for
// I/O; the associated file descriptor has mode os.O_RDWR. If there is an
// error, it will be of type *os.PathError.
//
// For the meaning of opts please see documentation of Options.
func Create(name string, opts *Options) (db *DB, err error) {
	opts = opts.clone()
	opts._ACID = _ACIDFull
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}

	return CreateFromFiler(lldb.NewSimpleFileFiler(f), opts)
}

func create(filer lldb.Filer, opts *Options, isMem bool) (db *DB, err error) {
	defer func() {
		if db != nil {
			db.opts = opts
		}
	}()
	defer func() {
		lock := opts.lock
		if err != nil && lock != nil {
			lock.Close()
			db = nil
		}
	}()

	if err = opts.check(filer.Name(), true, !isMem); err != nil {
		return
	}

	b := [16]byte{byte(magic[0]), byte(magic[1]), byte(magic[2]), byte(magic[3]), 0x00} // ver 0x00
	if n, err := filer.WriteAt(b[:], 0); n != 16 {
		return nil, &os.PathError{Op: "kv.create.WriteAt", Path: filer.Name(), Err: err}
	}

	db = &DB{lock: opts.lock}

	filer = lldb.NewInnerFiler(filer, 16)
	if filer, err = opts.acidFiler(db, filer); err != nil {
		return nil, err
	}

	db.filer = filer
	if err = filer.BeginUpdate(); err != nil {
		return
	}

	defer func() {
		if e := filer.EndUpdate(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	if db.alloc, err = lldb.NewAllocator(filer, &lldb.Options{}); err != nil {
		return nil, &os.PathError{Op: "kv.create", Path: filer.Name(), Err: err}
	}

	db.alloc.Compress = true
	db.isMem = isMem
	var h int64
	if db.root, h, err = lldb.CreateBTree(db.alloc, opts.Compare); err != nil {
		return
	}

	if h != 1 {
		panic("internal error")
	}

	db.wal = opts.wal
	return
}

// CreateMem creates a new instance of an in-memory DB not backed by a disk
// file. Memory DBs are resource limited as they are completely held in memory
// and are not automatically persisted.
//
// For the meaning of opts please see documentation of Options.
func CreateMem(opts *Options) (db *DB, err error) {
	opts = opts.clone()
	opts._ACID = _ACIDTransactions
	f := lldb.NewMemFiler()
	return create(f, opts, true)
}

// CreateTemp creates a new temporary DB in the directory dir with a basename
// beginning with prefix and name ending in suffix. If dir is the empty string,
// CreateTemp uses the default directory for temporary files (see os.TempDir).
// Multiple programs calling CreateTemp simultaneously will not choose the same
// file name for the DB. The caller can use Name() to find the pathname of the
// DB file. It is the caller's responsibility to remove the file when no longer
// needed.
//
// For the meaning of opts please see documentation of Options.
func CreateTemp(dir, prefix, suffix string, opts *Options) (db *DB, err error) {
	opts = opts.clone()
	opts._ACID = _ACIDFull
	f, err := fileutil.TempFile(dir, prefix, suffix)
	if err != nil {
		return
	}

	return create(lldb.NewSimpleFileFiler(f), opts, false)
}

// Open opens the named DB file for reading/writing. If successful, methods on
// the returned DB can be used for I/O; the associated file descriptor has mode
// os.O_RDWR. If there is an error, it will be of type *os.PathError.
//
// Note: While a DB is opened, it is locked and cannot be simultaneously opened
// again.
//
// For the meaning of opts please see documentation of Options.
func Open(name string, opts *Options) (db *DB, err error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return OpenFromFiler(lldb.NewSimpleFileFiler(f), opts)
}

// OpenFromFiler is like Open but it accepts an arbitrary backing storage
// provided by filer.
func OpenFromFiler(filer lldb.Filer, opts *Options) (db *DB, err error) {
	opts = opts.clone()
	opts._ACID = _ACIDFull
	defer func() {
		if db != nil {
			db.opts = opts
		}
	}()
	defer func() {
		lock := opts.lock
		if err != nil && lock != nil {
			lock.Close()
			db = nil
		}
		if err != nil {
			if db != nil {
				db.Close()
				db = nil
			}
		}
	}()

	name := filer.Name()
	if err = opts.check(name, false, true); err != nil {
		return
	}

	sz, err := filer.Size()
	if err != nil {
		return
	}

	if sz%16 != 0 {
		return nil, &os.PathError{Op: "kv.Open:", Path: name, Err: fmt.Errorf("file size %d(%#x) is not 0 (mod 16)", sz, sz)}
	}

	var b [16]byte
	if n, err := filer.ReadAt(b[:], 0); n != 16 || err != nil {
		return nil, &os.PathError{Op: "kv.Open.ReadAt", Path: name, Err: err}
	}

	var h header
	if err = h.rd(b[:]); err != nil {
		return nil, &os.PathError{Op: "kv.Open:validate header", Path: name, Err: err}
	}

	db = &DB{lock: opts.lock}
	if filer, err = opts.acidFiler(db, filer); err != nil {
		return nil, err
	}

	db.filer = filer
	switch h.ver {
	default:
		return nil, &os.PathError{Op: "kv.Open", Path: name, Err: fmt.Errorf("unknown/unsupported kv file format version %#x", h.ver)}
	case 0x00:
		if _, err = open00(name, db); err != nil {
			return nil, err
		}
	}

	db.root, err = lldb.OpenBTree(db.alloc, opts.Compare, 1)
	db.wal = opts.wal
	if opts.VerifyDbAfterOpen {
		err = verifyAllocator(db.alloc)
	}
	return
}

// Close closes the DB, rendering it unusable for I/O. It returns an error, if
// any. Failing to call Close before exiting a program can lose the last open
// or being committed transaction.
//
// Successful Close is idempotent.
func (db *DB) Close() (err error) {
	db.closeMu.Lock()
	defer db.closeMu.Unlock()
	if db.closed {
		return
	}

	db.closed = true

	if err = db.enter(); err != nil {
		return
	}

	doLeave := true
	defer func() {
		db.wal = nil
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		if doLeave {
			db.leave(&err)
		}
	}()

	if db.acidTimer != nil {
		db.acidTimer.Stop()
	}

	var e error
	for db.acidNest > 0 {
		db.acidNest--
		if e = db.filer.EndUpdate(); err == nil {
			err = e
		}
	}

	doLeave = false
	if e = db.leave(&err); err == nil {
		err = e
	}
	if db.opts.VerifyDbBeforeClose {
		if e = verifyAllocator(db.alloc); err == nil {
			err = e
		}
	}
	if e = db.close(); err == nil {
		err = e
	}
	if lock := db.lock; lock != nil {
		if e = lock.Close(); err == nil {
			err = e
		}
	}
	if wal := db.wal; wal != nil {
		e = wal.Close()
		db.wal = nil
		if err == nil {
			err = e
		}
	}
	return
}

func (db *DB) close() (err error) {
	// We are safe to close due to locked db.closeMu, but not safe against
	// any other goroutine concurrently calling other exported db methods,
	// causing a race[0] in the db.enter() mechanism. So we must lock
	// db.bkl.
	//
	//  [0]: https://github.com/cznic/kv/issues/17#issuecomment-31960658
	db.bkl.Lock()
	defer db.bkl.Unlock()

	if db.isMem { // lldb.MemFiler
		return
	}

	err = db.filer.Sync()
	if e := db.filer.Close(); err == nil {
		err = e
	}
	if db.opts.VerifyDbAfterClose {
		if e := verifyDbFile(db.Name()); err == nil {
			err = e
		}
	}
	return
}

// Name returns the name of the DB file.
func (db *DB) Name() string {
	return db.filer.Name()
}

// Size returns the size of the DB file.
func (db *DB) Size() (sz int64, err error) {
	db.bkl.Lock()
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		db.bkl.Unlock()
	}()

	return db.filer.Size()
}

func (db *DB) enter() (err error) {
	db.bkl.Lock()
	switch db.acidState {
	default:
		panic("internal error")
	case stDisabled:
		db.acidNest++
		if db.acidNest == 1 {
			if err = db.filer.BeginUpdate(); err != nil {
				return err
			}
		}
	case stIdle:
		if err = db.filer.BeginUpdate(); err != nil {
			return err
		}

		db.acidNest = 1
		db.acidTimer = time.AfterFunc(db.gracePeriod, db.timeout)
		db.acidState = stCollecting
	case stCollecting:
		db.acidNest++
	case stIdleArmed:
		db.acidNest = 1
		db.acidState = stCollectingArmed
	case stCollectingArmed:
		db.acidNest++
	case stCollectingTriggered:
		db.acidNest++
	case stEndUpdateFailed:
		return db.leave(&err)
	}

	return nil
}

func (db *DB) leave(err *error) error {
	switch db.acidState {
	default:
		panic("internal error")
	case stDisabled:
		db.acidNest--
		if db.acidNest == 0 {
			if e := db.filer.EndUpdate(); e != nil && *err == nil {
				*err = e
			}
		}
	case stCollecting:
		db.acidNest--
		if db.acidNest == 0 {
			db.acidState = stIdleArmed
		}
	case stCollectingArmed:
		db.acidNest--
		if db.acidNest == 0 {
			db.acidState = stIdleArmed
		}
	case stCollectingTriggered:
		db.acidNest--
		if db.acidNest == 0 {
			if e := db.filer.EndUpdate(); e != nil && *err == nil {
				*err = e
			}
			db.acidState = stIdle
		}
	case stEndUpdateFailed:
		db.bkl.Unlock()
		return fmt.Errorf("Last transaction commit failed: %v", db.lastCommitErr)
	}

	if *err != nil {
		db.filer.Rollback() // return the original, input error
	}
	db.bkl.Unlock()
	return *err
}

func (db *DB) timeout() {
	db.closeMu.Lock()
	defer db.closeMu.Unlock()
	if db.closed {
		return
	}

	db.bkl.Lock()
	defer db.bkl.Unlock()

	switch db.acidState {
	default:
		panic("internal error")
	case stIdle:
		panic("internal error")
	case stCollecting:
		db.acidState = stCollectingTriggered
	case stIdleArmed:
		if err := db.filer.EndUpdate(); err != nil { // If EndUpdate fails, no WAL was written (automatic Rollback)
			db.acidState = stEndUpdateFailed
			db.lastCommitErr = err
			return
		}

		db.acidState = stIdle
	case stCollectingArmed:
		db.acidState = stCollectingTriggered
	case stCollectingTriggered:
		panic("internal error")
	}
}

// BeginTransaction starts a new transaction. Every call to BeginTransaction
// must be eventually "balanced" by exactly one call to Commit or Rollback (but
// not both). Calls to BeginTransaction may nest.
//
// BeginTransaction is atomic and it is safe for concurrent use by multiple
// goroutines (if/when that makes sense).
func (db *DB) BeginTransaction() (err error) {
	if err = db.enter(); err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		db.leave(&err)
	}()

	db.acidNest++
	return db.filer.BeginUpdate()
}

// Commit commits the current transaction. If the transaction is the top level
// one, then all of the changes made within the transaction are atomically made
// persistent in the DB.  Invocation of an unbalanced Commit is an error.
//
// Commit is atomic and it is safe for concurrent use by multiple goroutines
// (if/when that makes sense).
func (db *DB) Commit() (err error) {
	if err = db.enter(); err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		db.leave(&err)
	}()

	db.acidNest--
	return db.filer.EndUpdate()
}

// Rollback cancels and undoes the innermost transaction level. If the
// transaction is the top level one, then no of the changes made within the
// transactions are persisted. Invocation of an unbalanced Rollback is an
// error.
//
// Rollback is atomic and it is safe for concurrent use by multiple goroutines
// (if/when that makes sense).
func (db *DB) Rollback() (err error) {
	if err = db.enter(); err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		db.leave(&err)
	}()

	db.acidNest--
	return db.filer.Rollback()
}

// Verify attempts to find any structural errors in DB wrt the organization of
// it as defined by lldb.Allocator. Any problems found are reported to 'log'
// except non verify related errors like disk read fails etc. If 'log' returns
// false or the error doesn't allow to (reliably) continue, the verification
// process is stopped and an error is returned from the Verify function.
// Passing a nil log works like providing a log function always returning
// false. Any non-structural errors, like for instance Filer read errors, are
// NOT reported to 'log', but returned as the Verify's return value, because
// Verify cannot proceed in such cases. Verify returns nil only if it fully
// completed verifying DB without detecting any error.
//
// It is recommended to limit the number reported problems by returning false
// from 'log' after reaching some limit. Huge and corrupted DB can produce an
// overwhelming error report dataset.
//
// The verifying process will scan the whole DB at least 3 times (a trade
// between processing space and time consumed). It doesn't read the content of
// free blocks above the head/tail info bytes. If the 3rd phase detects lost
// free space, then a 4th scan (a faster one) is performed to precisely report
// all of them.
//
// Statistics are returned via 'stats' if non nil. The statistics are valid
// only if Verify succeeded, ie. it didn't reported anything to log and it
// returned a nil error.
func (db *DB) Verify(log func(error) bool, stats *lldb.AllocStats) (err error) {
	bitmapf, err := fileutil.TempFile("", "verifier", ".tmp")
	if err != nil {
		return
	}

	defer func() {
		tn := bitmapf.Name()
		bitmapf.Close()
		os.Remove(tn)
	}()

	bitmap := lldb.NewSimpleFileFiler(bitmapf)

	if err = db.enter(); err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
		db.leave(&err)
	}()

	return db.alloc.Verify(bitmap, log, stats)
}

// Delete deletes key and its associated value from the DB.
//
// Delete is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Delete(key []byte) (err error) {
	if err = db.enter(); err != nil {
		return
	}

	err = db.root.Delete(key)
	return db.leave(&err)
}

// Extract is a combination of Get and Delete. If the key exists in the DB, it
// is returned (like Get) and also deleted from the DB in a more efficient way
// which doesn't search for the key twice. The returned slice may be a
// sub-slice of buf if buf was large enough to hold the entire content.
// Otherwise, a newly allocated slice will be returned. It is valid to pass a
// nil buf.
//
// Extract is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Extract(buf, key []byte) (value []byte, err error) {
	if err = db.enter(); err != nil {
		return
	}

	value, err = db.root.Extract(buf, key)
	db.leave(&err)
	return
}

// First returns the first KV pair in the DB, if it exists. Otherwise key ==
// nil and value == nil.
//
// First is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) First() (key, value []byte, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	return db.root.First()
}

// Get returns the value associated with key, or nil if no such value exists.
// The returned slice may be a sub-slice of buf if buf was large enough to hold
// the entire content. Otherwise, a newly allocated slice will be returned. It
// is valid to pass a nil buf.
//
// Get is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Get(buf, key []byte) (value []byte, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	return db.root.Get(buf, key)
}

// Last returns the last KV pair of the DB, if it exists. Otherwise key ==
// nil and value == nil.
//
// Last is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Last() (key, value []byte, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	return db.root.Last()
}

// Put combines Get and Set in a more efficient way where the DB is searched
// for the key only once. The upd(ater) receives the current (key, old-value),
// if that exists or (key, nil) otherwise. It can then return a (new-value,
// true, nil) to create or overwrite the existing value in the KV pair, or
// (whatever, false, nil) if it decides not to create or not to update the
// value of the KV pair.
//
//	db.Set(k, v)
//
// conceptually equals
//
//	db.Put(k, func(k, v []byte){ return v, true }([]byte, bool))
//
// modulo the differing return values.
//
// The returned slice may be a sub-slice of buf if buf was large enough to hold
// the entire content. Otherwise, a newly allocated slice will be returned. It
// is valid to pass a nil buf.
//
// Put is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Put(buf, key []byte, upd func(key, old []byte) (new []byte, write bool, err error)) (old []byte, written bool, err error) {
	if err = db.enter(); err != nil {
		return
	}

	old, written, err = db.root.Put(buf, key, upd)
	db.leave(&err)
	return
}

// Seek returns an enumerator positioned on the first key/value pair whose key
// is 'greater than or equal to' the given key. There may be no such pair, in
// which case the Next,Prev methods of the returned enumerator will always
// return io.EOF.
//
// Seek is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Seek(key []byte) (enum *Enumerator, hit bool, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	enum0, hit, err := db.root.Seek(key)
	if err != nil {
		return
	}

	enum = &Enumerator{
		db:   db,
		enum: enum0,
	}
	return
}

// SeekFirst returns an enumerator positioned on the first KV pair in the DB,
// if any. For an empty DB, err == io.EOF is returned.
//
// SeekFirst is atomic and it is safe for concurrent use by multiple
// goroutines.
func (db *DB) SeekFirst() (enum *Enumerator, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	enum0, err := db.root.SeekFirst()
	if err != nil {
		return
	}

	enum = &Enumerator{
		db:   db,
		enum: enum0,
	}
	return
}

// SeekLast returns an enumerator positioned on the last KV pair in the DB,
// if any. For an empty DB, err == io.EOF is returned.
//
// SeekLast is atomic and it is safe for concurrent use by multiple
// goroutines.
func (db *DB) SeekLast() (enum *Enumerator, err error) {
	db.bkl.Lock()
	defer db.bkl.Unlock()
	enum0, err := db.root.SeekLast()
	if err != nil {
		return
	}

	enum = &Enumerator{
		db:   db,
		enum: enum0,
	}
	return
}

// Set sets the value associated with key. Any previous value, if existed, is
// overwritten by the new one.
//
// Set is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Set(key, value []byte) (err error) {
	if err = db.enter(); err != nil {
		return
	}

	err = db.root.Set(key, value)
	db.leave(&err)
	return
}

// Enumerator captures the state of enumerating a DB. It is returned from the
// Seek* methods. Multiple enumerations may be in progress simultaneously.  The
// enumerator is aware of any mutations made to the tree in the process of
// enumerating it and automatically resumes the enumeration.
//
// Multiple concurrently executing enumerations may be in progress.
type Enumerator struct {
	db   *DB
	enum *lldb.BTreeEnumerator
}

// Next returns the currently enumerated KV pair, if it exists and moves to the
// next KV in the key collation order. If there is no KV pair to return, err ==
// io.EOF is returned.
//
// Next is atomic and it is safe for concurrent use by multiple goroutines.
func (e *Enumerator) Next() (key, value []byte, err error) {
	e.db.bkl.Lock()
	defer e.db.bkl.Unlock()
	return e.enum.Next()
}

// Prev returns the currently enumerated KV pair, if it exists and moves to the
// previous KV in the key collation order. If there is no KV pair to return,
// err == io.EOF is returned.
//
// Prev is atomic and it is safe for concurrent use by multiple goroutines.
func (e *Enumerator) Prev() (key, value []byte, err error) {
	e.db.bkl.Lock()
	defer e.db.bkl.Unlock()
	return e.enum.Prev()
}

// Inc atomically increments the value associated with key by delta and
// returns the new value. If the value doesn't exists before calling Inc or if
// the value is not an [8]byte, the value is considered to be zero before peforming Inc.
//
// Inc is atomic and it is safe for concurrent use by multiple goroutines.
func (db *DB) Inc(key []byte, delta int64) (val int64, err error) {
	if err = db.enter(); err != nil {
		return
	}

	defer db.leave(&err)

	pbuf := buffer.Get(8)
	defer buffer.Put(pbuf)
	_, _, err = db.root.Put(
		*pbuf,
		key,
		func(key []byte, old []byte) (new []byte, write bool, err error) {
			write = true
			if len(old) == 8 {
				val = int64(binary.BigEndian.Uint64(old))
			} else {
				old = make([]byte, 8)
				val = 0
			}
			val += delta
			binary.BigEndian.PutUint64(old, uint64(val))
			new = old
			return
		},
	)

	return
}

// WALName returns the name of the WAL file in use or an empty string for memory
// or closed databases.
func (db *DB) WALName() string {
	if f := db.wal; f != nil {
		return f.Name()
	}

	return ""
}
