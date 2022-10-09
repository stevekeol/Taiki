package main

import (
	"Taiki/core/tdb"
	"Taiki/db"
)

func loadDatabase(cfg *config) (db.KeyValueStore, error) {
	var kvdb db.KeyValueStore
	var err error

	switch cfg.DbType {
	case "memorydb":
		{
			// log.Debug("Creating block database in memory.")
			kvdb, err = tdb.NewMemoryDatabase()
			log.Info("[2/3]database opened(memorydb).")
		}
	case "leveldb":
		{
			// log.Debug("Creating block database in LevelDB.")
			kvdb, err = tdb.NewLevelDBDatabase(cfg.DataDir, 0, 0, false) // TODO 此处的参数写死了
			log.Info("[2/3]database opened(leveldb).")
		}
	}
	if err != nil {
		return nil, err
	}

	return kvdb, nil
}
