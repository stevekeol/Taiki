package db

import (
	"Taiki/block"
)

// Cursor represents a cursor over key/value pairs and nested buckets of a
// bucket.
//
// Note that open cursors are not tracked on bucket changes and any
// modifications to the bucket, with the exception of Cursor.Delete, invalidates
// the cursor.  After invalidation, the cursor must be repositioned, or the keys
// and values returned may be unpredictable.
type Cursor interface {
	// Bucket returns the bucket the cursor was created for.
	Bucket() Bucket

	// Delete removes the current key/value pair the cursor is at without
	// invalidating the cursor.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrIncompatibleValue if attempted when the cursor points to a
	//     nested bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	Delete() error

	// First positions the cursor at the first key/value pair and returns
	// whether or not the pair exists.
	First() bool

	// Last positions the cursor at the last key/value pair and returns
	// whether or not the pair exists.
	Last() bool

	// Next moves the cursor one key/value pair forward and returns whether
	// or not the pair exists.
	Next() bool

	// Prev moves the cursor one key/value pair backward and returns whether
	// or not the pair exists.
	Prev() bool

	// Seek positions the cursor at the first key/value pair that is greater
	// than or equal to the passed seek key.  Returns whether or not the
	// pair exists.
	Seek(seek []byte) bool

	// Key returns the current key the cursor is pointing to.
	Key() []byte

	// Value returns the current value the cursor is pointing to.  This will
	// be nil for nested buckets.
	Value() []byte
}

// Bucket 表示一个键值对的集合
type Bucket interface {
	// Bucket retrieves a nested bucket with the given key.  Returns nil if
	// the bucket does not exist.
	Bucket(key []byte) Bucket

	// CreateBucket creates and returns a new nested bucket with the given
	// key.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketExists if the bucket already exists
	//   - ErrBucketNameRequired if the key is empty
	//   - ErrIncompatibleValue if the key is otherwise invalid for the
	//     particular implementation
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	CreateBucket(key []byte) (Bucket, error)

	// CreateBucketIfNotExists creates and returns a new nested bucket with
	// the given key if it does not already exist.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketNameRequired if the key is empty
	//   - ErrIncompatibleValue if the key is otherwise invalid for the
	//     particular implementation
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	CreateBucketIfNotExists(key []byte) (Bucket, error)

	// DeleteBucket removes a nested bucket with the given key.  This also
	// includes removing all nested buckets and keys under the bucket being
	// deleted.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrBucketNotFound if the specified bucket does not exist
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	DeleteBucket(key []byte) error

	// ForEach invokes the passed function with every key/value pair in the
	// bucket.  This does not include nested buckets or the key/value pairs
	// within those nested buckets.
	//
	// WARNING: It is not safe to mutate data while iterating with this
	// method.  Doing so may cause the underlying cursor to be invalidated
	// and return unexpected keys and/or values.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The slices returned by this function are only valid during a
	// transaction.  Attempting to access them after a transaction has ended
	// results in undefined behavior.  Additionally, the slices must NOT
	// be modified by the caller.  These constraints prevent additional data
	// copies and allows support for memory-mapped database implementations.
	ForEach(func(k, v []byte) error) error

	// ForEachBucket invokes the passed function with the key of every
	// nested bucket in the current bucket.  This does not include any
	// nested buckets within those nested buckets.
	//
	// WARNING: It is not safe to mutate data while iterating with this
	// method.  Doing so may cause the underlying cursor to be invalidated
	// and return unexpected keys and/or values.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The keys returned by this function are only valid during a
	// transaction.  Attempting to access them after a transaction has ended
	// results in undefined behavior.  This constraint prevents additional
	// data copies and allows support for memory-mapped database
	// implementations.
	ForEachBucket(func(k []byte) error) error

	// Cursor returns a new cursor, allowing for iteration over the bucket's
	// key/value pairs and nested buckets in forward or backward order.
	//
	// You must seek to a position using the First, Last, or Seek functions
	// before calling the Next, Prev, Key, or Value functions.  Failure to
	// do so will result in the same return values as an exhausted cursor,
	// which is false for the Prev and Next functions and nil for Key and
	// Value functions.
	Cursor() Cursor

	// Writable returns whether or not the bucket is writable.
	Writable() bool

	// Put saves the specified key/value pair to the bucket.  Keys that do
	// not already exist are added and keys that already exist are
	// overwritten.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrKeyRequired if the key is empty
	//   - ErrIncompatibleValue if the key is the same as an existing bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	//
	// NOTE: The slices passed to this function must NOT be modified by the
	// caller.  This constraint prevents the requirement for additional data
	// copies and allows support for memory-mapped database implementations.
	Put(key, value []byte) error

	// Get returns the value for the given key.  Returns nil if the key does
	// not exist in this bucket.  An empty slice is returned for keys that
	// exist but have no value assigned.
	//
	// NOTE: The value returned by this function is only valid during a
	// transaction.  Attempting to access it after a transaction has ended
	// results in undefined behavior.  Additionally, the value must NOT
	// be modified by the caller.  These constraints prevent additional data
	// copies and allows support for memory-mapped database implementations.
	Get(key []byte) []byte

	// Delete removes the specified key from the bucket.  Deleting a key
	// that does not exist does not return an error.
	//
	// The interface contract guarantees at least the following errors will
	// be returned (other implementation-specific errors are possible):
	//   - ErrKeyRequired if the key is empty
	//   - ErrIncompatibleValue if the key is the same as an existing bucket
	//   - ErrTxNotWritable if attempted against a read-only transaction
	//   - ErrTxClosed if the transaction has already been closed
	Delete(key []byte) error
}

// 指定由hash,偏移,长度指定的区块的特定区域
type BlockRegion struct {
	Hash   []byte
	Offset uint32
	Len    uint32
}

// 表示数据库事务(交易)，只读/读写，该事务提供一个元数据桶进行读和写操作
// 事务将只提供创建数据库时的数据库视图。
// 不应该长时间运行（因为事务原子操作）
type Tx interface {
	// 返回所有元数据存储的最上面的桶
	Metadata() Bucket
	// 将给定block存进数据库
	// ErrBlockExists-区块已存在；ErrTxNotWritable-尝试只读事务；ErrTxClosed -事务已被关闭
	StoreBlock(block block.Block) error
	// 判断是否有给定hash的区块
	// ErrTxClosed -事务已被关闭
	HasBlock(hash []byte) (bool, error)
	HasBlocks(hashes [][]byte) ([]bool, error)
	// 根据给定hash取出区块头的原始序列化的字节
	FetchBlockHeader(hash []byte) ([]byte, error)
	FetchBlockHeaders(hashes [][]byte) ([][]byte, error)
	// 根据给定hash取出区块的原始序列化的字节
	FetchBlock(hash []byte) ([]byte, error)
	FetchBlocks(hashes [][]byte) ([][]byte, error)

	// 获取给定区块范围的区块的原始序列化的字节
	//
	// For example, it is possible to directly extract Bitcoin transactions
	// and/or scripts from a block with this function.  Depending on the
	// backend implementation, this can provide significant savings by
	// avoiding the need to load entire blocks.

	//   - ErrBlockNotFound if the requested block hash does not exist
	//   - ErrBlockRegionInvalid if the region exceeds the bounds of the
	//     associated block
	//   - ErrTxClosed if the transaction has already been closed
	//   - ErrCorruption if the database has somehow become corrupted
	FetchBlockRegion(region *BlockRegion) ([]byte, error)
	FetchBlockRegions(regions []BlockRegion) ([][]byte, error)

	// 提交 元数据和区块存储的所有变动提交
	Commit() error

	// 回滚 元数据和区块存储的所有变动
	Rollback() error
}

// DB用来区块和相关元数据的存储
type DB interface {
	// 返回数据库驱动类型
	Type() string
	// 开启一个事务（只读/读写，可多次读写）(当已经开启后再次调用该方法会被阻塞)（当调用Rollback或Commit时，该事务必须被关闭）
	Begin(writable bool) (Tx, error)
	// 在一个受控只读事务中，调用传入的函数（此时不能调用Rollback或Commit）
	View(fn func(tx Tx) error) error
	// 在一个受控只读事务中，调用传入的函数（该函数发生错误会导致该事务回滚并退出）（此时不能调用Rollback或Commit）
	Update(fn func(tx Tx) error) error
	// 干净的关闭数据库并同步所有数据（不管是回滚还是提交都会在完成后才关闭）
	Close() error
}
