// Package db defines the interface for a Taiki data store
package db

import (
	// "github.com/syndtr/goleveldb/leveldb"
	"io"
)

type KeyValueReader interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
}

type KeyValueWriter interface {
	// use Put instead of Set (leveldb)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
}

type KeyValueStater interface {
	// Stat return a particular internal stat of the database
	Stat(property string) (string, error)
}

type Compacter interface {
	// TODO
	// Compact flattens the underlying datads store for the given key range
	Compact(start []byte, limit []byte) error
}

type KeyValueStore interface {
	KeyValueReader
	KeyValueWriter
	KeyValueStater
	Batcher
	Iteratee
	Compacter
	Snapshotter
	io.Closer
}

// AncientReaderOp contains the methods required to read from immutable ancient data.
type AncientReaderOp interface {
	// HasAncient returns an indicator whether the specified data exists in the
	// ancient store.
	HasAncient(kind string, number uint64) (bool, error)

	// Ancient retrieves an ancient binary blob from the append-only immutable files.
	Ancient(kind string, number uint64) ([]byte, error)

	// AncientRange retrieves multiple items in sequence, starting from the index 'start'.
	// It will return
	//  - at most 'count' items,
	//  - at least 1 item (even if exceeding the maxBytes), but will otherwise
	//   return as many items as fit into maxBytes.
	AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error)

	// Ancients returns the ancient item numbers in the ancient store.
	Ancients() (uint64, error)

	// Tail returns the number of first stored item in the freezer.
	// This number can also be interpreted as the total deleted item numbers.
	Tail() (uint64, error)

	// AncientSize returns the ancient size of the specified category.
	AncientSize(kind string) (uint64, error)
}

// AncientReader is the extended ancient reader interface including 'batched' or 'atomic' reading.
type AncientReader interface {
	AncientReaderOp

	// ReadAncients runs the given read operation while ensuring that no writes take place
	// on the underlying freezer.
	ReadAncients(fn func(AncientReaderOp) error) (err error)
}

// AncientWriter contains the methods required to write to immutable ancient data.
type AncientWriter interface {
	// ModifyAncients runs a write operation on the ancient store.
	// If the function returns an error, any changes to the underlying store are reverted.
	// The integer return value is the total size of the written data.
	ModifyAncients(func(AncientWriteOp) error) (int64, error)

	// TruncateHead discards all but the first n ancient data from the ancient store.
	// After the truncation, the latest item can be accessed it item_n-1(start from 0).
	TruncateHead(n uint64) error

	// TruncateTail discards the first n ancient data from the ancient store. The already
	// deleted items are ignored. After the truncation, the earliest item can be accessed
	// is item_n(start from 0). The deleted items may not be removed from the ancient store
	// immediately, but only when the accumulated deleted data reach the threshold then
	// will be removed all together.
	TruncateTail(n uint64) error

	// Sync flushes all in-memory ancient store data to disk.
	Sync() error

	// MigrateTable processes and migrates entries of a given table to a new format.
	// The second argument is a function that takes a raw entry and returns it
	// in the newest format.
	MigrateTable(string, func([]byte) ([]byte, error)) error
}

// AncientWriteOp is given to the function argument of ModifyAncients.
type AncientWriteOp interface {
	// Append adds an RLP-encoded item.
	Append(kind string, number uint64, item interface{}) error

	// AppendRaw adds an item without RLP-encoding it.
	AppendRaw(kind string, number uint64, item []byte) error
}

// AncientStater wraps the Stat method of a backing data store.
type AncientStater interface {
	// AncientDatadir returns the root directory path of the ancient store.
	AncientDatadir() (string, error)
}

// Reader contains the methods required to read data from both key-value as well as
// immutable ancient data.
type Reader interface {
	KeyValueReader
	AncientReader
}

// Writer contains the methods required to write data to both key-value as well as
// immutable ancient data.
type Writer interface {
	KeyValueWriter
	AncientWriter
}

// Stater contains the methods required to retrieve states from both key-value as well as
// immutable ancient data.
type Stater interface {
	KeyValueStater
	AncientStater
}

// AncientStore contains all the methods required to allow handling different
// ancient data stores backing immutable chain data store.
type AncientStore interface {
	AncientReader
	AncientWriter
	AncientStater
	io.Closer
}

// Database contains all the methods required by the high level database to not
// only access the key-value data store but also the chain freezer.
type Database interface {
	Reader
	Writer
	Batcher
	Iteratee
	Stater
	Compacter
	Snapshotter
	io.Closer
}