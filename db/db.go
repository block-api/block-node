package db

import (
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/log"
	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	config  *config.Database
	leveldb map[string]*leveldb.DB
}

// NewDatabase creates new database struct
func NewDatabase(config *config.Database) Database {
	newDb := Database{
		config:  config,
		leveldb: make(map[string]*leveldb.DB),
	}

	var err error
	for dbName, dbConfig := range config.LevelDB {
		newDb.leveldb[dbName], err = leveldb.OpenFile(dbConfig.DbPath, nil)

		if err != nil {
			panic(err)
		}

		log.Debug("database " + dbName + ": " + dbConfig.DbPath)
	}

	return newDb
}

func (db *Database) GetLevelDB(name string) *leveldb.DB {
	if db.leveldb[name] != nil {
		return db.leveldb[name]
	}

	return nil
}
