package db

import (
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/log"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	DbDiscovery = "discovery"
)

type Database struct {
	config  *config.ConfigDatabase
	leveldb map[string]*leveldb.DB
}

// NewDatabase creates new database struct
func NewDatabase(config *config.ConfigDatabase) Database {
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
