package db

type Snapshot interface {
	// Has retrieves if a key is present in the snapshot backing by a key-value
	// data store.
	Has(key []byte) (bool, error)

	// Get retrieves the given key if it's present in the snapshot backing by
	// key-value data store.
	Get(key []byte) ([]byte, error)

	// Release releases associated resources. Release should always succeed and can
	// be called multiple times without causing error.
	Release()
}

// Snapshotter wraps the Snapshot method of a backing data store.
type Snapshotter interface {
	// NewSnapshot creates a database snapshot based on the current state.
	// The created snapshot will not be affected by all following mutations
	// happened on the database.
	// Note don't forget to release the snapshot once it's used up, otherwise
	// the stale data will never be cleaned up by the underlying compactor.
	NewSnapshot() (Snapshot, error)
}
