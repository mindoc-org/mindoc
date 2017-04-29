// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cznic/lldb"
)

const (
	// BeginUpdate/EndUpdate/Rollback will be no-ops. All operations
	// updating a DB will be written immediately including partial updates
	// during operation's progress. If any update fails, the DB can become
	// unusable. The same applies to DB crashes and/or any other non clean
	// DB shutdown.
	_ACIDNone = iota

	// Enable transactions. BeginUpdate/EndUpdate/Rollback will be
	// effective. All operations on the DB will be automatically performed
	// within a transaction. Operations will thus either succeed completely
	// or have no effect at all - they will be rollbacked in case of any
	// error. If any update fails the DB will not be corrupted. DB crashes
	// and/or any other non clean DB shutdown may still render the DB
	// unusable.
	_ACIDTransactions

	// Enable durability. Same as ACIDTransactions plus enables 2PC and
	// WAL.  Updates to the DB will be first made permanent in a WAL and
	// only after that reflected in the DB. A DB will automatically recover
	// from crashes and/or any other non clean DB shutdown. Only last
	// uncommitted transaction (transaction in progress ATM of a crash) can
	// get lost.
	//
	// NOTE: Options.GracePeriod may extend the span of a single
	// transaction to a batch of multiple transactions.
	//
	// NOTE2: Non zero GracePeriod requires GOMAXPROCS > 1 to work. Dbm
	// checks GOMAXPROCS in such case and if the value is 1 it
	// automatically sets GOMAXPROCS = 2.
	_ACIDFull
)

// Options are passed to the DB create/open functions to amend the behavior of
// those functions. The compatibility promise is the same as of struct types in
// the Go standard library - introducing changes can be made only by adding new
// exported fields, which is backward compatible as long as client code uses
// field names to assign values of imported struct types literals.
type Options struct {
	// Compare compares x and y. Compare may be nil, then bytes.Compare is
	// used instead.
	//
	// Compare returns:
	//
	//	-1 if x <  y
	//	 0 if x == y
	//	+1 if x >  y
	Compare func(x, y []byte) int

	// Locker specifies a function to lock a named file.
	// On success it returns an io.Closer to release the lock.
	// If nil, a default implementation is used.
	Locker func(name string) (io.Closer, error)

	// See the ACID* constants documentation.
	_ACID int

	// The write ahead log pathname. Applicable iff ACID == ACIDFull. May
	// be left empty in which case an unspecified pathname will be chosen,
	// which is computed from the DB name and which will be in the same
	// directory as the DB. Moving or renaming the DB while it is shut down
	// will break it's connection to the automatically computed name.
	// Moving both the files (the DB and the WAL) into another directory
	// with no renaming is safe.
	//
	// On creating a new DB the WAL file must not exist or it must be
	// empty. It's not safe to write to a non empty WAL file as it may
	// contain unprocessed DB recovery data.
	WAL string

	// Time to collect transactions before committing them into the WAL.
	// Applicable iff ACID == ACIDFull. All updates are held in memory
	// during the grace period so it should not be more than few seconds at
	// most.
	//
	// Recommended value for GracePeriod is 1 second.
	//
	// NOTE: Using small GracePeriod values will make DB updates very slow.
	// Zero GracePeriod will make every single update a separate 2PC/WAL
	// transaction.  Values smaller than about 100-200 milliseconds
	// (particularly for mechanical, rotational HDs) are not recommended
	// and they may not be always honored.
	_GracePeriod time.Duration
	wal          *os.File
	lock         io.Closer

	noClone bool // test hook

	// VerifyDbBeforeOpen turns on structural verification of the DB before
	// it is opened. This verification may legitimately fail if the DB
	// crashed and a yet-to-be-processed non empty WAL file exists.
	VerifyDbBeforeOpen bool

	// VerifyDbAfterOpen turns on structural verification of the DB after
	// it is opened and possibly recovered from WAL.
	VerifyDbAfterOpen bool

	// VerifyDbBeforeClose turns on structural verification of the DB
	// before it is closed.
	VerifyDbBeforeClose bool

	// VerifyDbAfterClose turns on structural verification of the DB after
	// it is closed.
	VerifyDbAfterClose bool

	// Turns on verification of every single mutation of the DB. Before any
	// such mutation a snapshot of the DB is created and the specific
	// mutation operation and parameters are recorded. After the mutation
	// the whole DB is verified. If the verification fails the last known
	// good state (the snapshot discussed above) and the corrupted state
	// are "core" dumped to a well known location (TBD).
	//
	//MAYBE ParanoidUpdates bool
}

func (o *Options) locker(dbname string) (io.Closer, error) {
	if o == nil || o.Locker == nil {
		return defaultLocker(dbname)
	}
	return o.Locker(dbname)
}

func (o *Options) clone() *Options {
	if o.noClone {
		return o
	}

	r := &Options{}
	*r = *o
	return r
}

func (o *Options) check(dbname string, new, lock bool) (err error) {
	if lock {
		if o.lock, err = o.locker(dbname); err != nil {
			return
		}
	}

	if o.VerifyDbBeforeOpen && !new {
		if err = verifyDbFile(dbname); err != nil {
			return
		}
	}

	switch o._ACID {
	default:
		panic("internal error")
	case _ACIDTransactions:
	case _ACIDFull:
		o._GracePeriod = time.Second
		if o.WAL == "" {
			o.WAL = o.walName(dbname, o.WAL)
		}

		switch new {
		case true:
			if o.wal, err = os.OpenFile(o.WAL, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err != nil {
				if os.IsExist(err) {
					fi, e := os.Stat(o.WAL)
					if e != nil {
						return e
					}

					if sz := fi.Size(); sz != 0 {
						return fmt.Errorf("cannot create DB %q: non empty WAL file %q (size %d) exists", dbname, o.WAL, sz)
					}

					o.wal, err = os.OpenFile(o.WAL, os.O_RDWR, 0666)
				}
				return
			}
		case false:
			if o.wal, err = os.OpenFile(o.WAL, os.O_RDWR, 0666); err != nil {
				if os.IsNotExist(err) {
					if o.wal, err = os.OpenFile(o.WAL, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err != nil {
						return fmt.Errorf("cannot open DB %q: failed to create  WAL file %q: %v", dbname, o.WAL, err)
					}

					err = nil
				}
				return err
			}
		}
	}

	return err
}

func (o *Options) walName(dbname, wal string) (r string) {
	if wal != "" {
		return filepath.Clean(wal)
	}

	base := filepath.Base(filepath.Clean(dbname))
	h := sha1.New()
	io.WriteString(h, base)
	return filepath.Join(filepath.Dir(dbname), fmt.Sprintf(".%x", h.Sum(nil)))
}

func (o *Options) acidFiler(db *DB, f lldb.Filer) (r lldb.Filer, err error) {
	switch o._ACID {
	default:
		panic("internal error")
	case _ACIDTransactions:
		if r, err = lldb.NewRollbackFiler(
			f,
			func(sz int64) error {
				return f.Truncate(sz)
			},
			f,
		); err != nil {
			return nil, err
		}

		return r, nil
	case _ACIDFull:
		if r, err = lldb.NewACIDFiler(f, o.wal); err != nil {
			return nil, err
		}

		db.acidState = stIdle
		db.gracePeriod = o._GracePeriod
		if o._GracePeriod == 0 {
			panic("internal error")
		}
		return r, nil
	}
}
