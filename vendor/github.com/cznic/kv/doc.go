// Copyright 2014 The kv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Package kv implements a simple and easy to use persistent key/value (KV) store.

Changelog

2016-07-11: KV now uses the stable version of lldb. (github.com/cznic/lldb).

The stored KV pairs are sorted in the key collation order defined by an user
supplied 'compare' function (passed as a field in Options).

Keys and Values Limits

Keys, as well as the values associated with them, are opaque []bytes. Maximum
size of a "native" key or value is 65787 bytes. Larger keys or values have to
be composed of the "native" ones in client code.

Database limits

The maximum DB size kv can handle is 2^60 bytes (1 exabyte). See also [4]:
"Block handles".

ACID and transactional properties

Transactions are resource limited. All changes made by a transaction are held
in memory until the top level transaction is committed. ACID[1] implementation
notes/details follows.

Atomicity

A successfully committed transaction appears (by its effects on the database)
to be indivisible ("atomic") iff the transaction is performed in isolation. An
aborted (via RollBack) transaction appears like it never happened under the
same limitation.

Atomic updates to the DB, via functions like Set, Inc, etc., are performed in
their own automatic transaction. If the partial progress of any such function
fails at any point, the automatic transaction is canceled via Rollback before
returning from the function. A non nil error is returned in that case.

Consistency

All reads, including those made from any other concurrent non isolated
transaction(s), performed during a not yet committed transaction, are dirty
reads, i.e.  the data returned are consistent with the in-progress state of the
open transaction, or all of the open transactions. Obviously, conflicts, data
races and inconsistent states can happen, but iff non isolated transactions are
performed.

Performing a Rollback at a nested transaction level properly returns the
transaction state (and data read from the DB) to what it was before the
respective BeginTransaction.

Isolation

Transactions of the atomic updating functions (Set, Put, Delete ...) are always
isolated. Transactions controlled by BeginTransaction/Commit/RollBack, are
isolated iff their execution is serialized.

Durability

Transactions are committed using the two phase commit protocol(2PC)[2] and a
write ahead log(WAL)[3]. DB recovery after a crash is performed automatically
using data from the WAL. Last transaction data, either of an in progress
transaction or a transaction being committed at the moment of the crash, can get
lost.

No protection from non readable files, files corrupted by other processes or by
memory faults or other HW problems, is provided. Always properly backup your DB
data file(s).

Links

Referenced from above:

  [1]: http://en.wikipedia.org/wiki/ACID
  [2]: http://en.wikipedia.org/wiki/2PC
  [3]: http://en.wikipedia.org/wiki/Write_ahead_logging
  [4]: http://godoc.org/github.com/cznic/lldb#Allocator

*/
package kv
