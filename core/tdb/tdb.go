// tdb是对底层db的对外暴露
// TODO：目前其实是为了避免lebeldb和db的循环依赖，而临时出此下策；但也不见得就不对，需要再次审计
package tdb

import (
	db "Taiki/db"
	"Taiki/db/leveldb"
	"Taiki/db/memorydb"
)

func NewMemoryDatabase() (db.KeyValueStore, error) {
	return memorydb.New(), nil
}

func NewLevelDBDatabase(file string, cache int, handles int, readonly bool) (db.KeyValueStore, error) {
	db, err := leveldb.New(file, cache, handles, readonly)
	if err != nil {
		return nil, err
	}
	return db, nil
}
