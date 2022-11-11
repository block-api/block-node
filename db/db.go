// Copyright 2022 The block-node Authors
// This file is part of the block-node library.
//
// The block-node library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The block-node library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the block-node library. If not, see <http://www.gnu.org/licenses/>.
package db

import (
	"sync"

	"github.com/block-api/block-node/params"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syndtr/goleveldb/leveldb"
)

var lock = new(sync.Mutex)
var manager *Manager

type Manager struct {
	config  *params.DatabaseConfig
	leveldb map[string]*leveldb.DB
	sqlite  map[string]*SQLite
}

// NewManager creates new Manager struct
func NewManager(config *params.DatabaseConfig) *Manager {
	if manager == nil {
		lock.Lock()
		defer lock.Unlock()

		manager = &Manager{
			config:  config,
			leveldb: make(map[string]*leveldb.DB),
			sqlite:  make(map[string]*SQLite),
		}

		// var err error
		// if len(config.LevelDB) > 0 {
		// 	for dbName, dbConfig := range config.LevelDB {
		// 		o := &opt.Options{
		// 			WriteBuffer: dbConfig.WriteBufferSize * opt.MiB,
		// 			Filter:      filter.NewBloomFilter(10),
		// 		}

		// 		manager.leveldb[dbName], err = leveldb.OpenFile(dbConfig.DbPath, o)

		// 		if err != nil {
		// 			panic(err)
		// 		}

		// 		log.Debug("database " + dbName + ": " + dbConfig.DbPath)
		// 	}
		// }

		// if len(config.SQLite) > 0 {
		// 	for dbName, dbConfig := range config.SQLite {
		// 		if _, err := os.Stat(dbConfig.DbPath); err != nil {
		// 			file, err := os.Create(dbConfig.DbPath)
		// 			if err != nil {
		// 				log.Panic(err.Error())
		// 			}
		// 			_ = file.Close()
		// 		}

		// 		s3db, err := sql.Open("sqlite3", "file:"+dbConfig.DbPath+"?"+dbConfig.Options)
		// 		if err != nil {
		// 			log.Panic(err.Error())
		// 		}

		// 		err = s3db.Ping()
		// 		if err != nil {
		// 			log.Panic(err.Error())
		// 		}

		// 		manager.sqlite[dbName] = &SQLite{
		// 			Db: s3db,
		// 		}

		// 		if dbConfig.MaxOpenConnections > 0 {
		// 			manager.sqlite[dbName].Db.SetMaxOpenConns(dbConfig.MaxOpenConnections)
		// 		}

		// 		if _, err := manager.sqlite[dbName].Db.Exec(CreateMigrationTable); err != nil {
		// 			log.Panic(err.Error())
		// 		}

		// 		log.Debug("database " + dbName + ": " + dbConfig.DbPath)
		// 	}
		// }
	}

	return manager
}

func GetDatabase() *Manager {
	return manager
}

func (db *Manager) RunMigrations() error {
	for _, dbSQLite := range db.sqlite {
		err := dbSQLite.RunMigrations()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Manager) RevertMigrations() error {
	for _, dbSQLite := range db.sqlite {
		err := dbSQLite.RevertMigrations()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Manager) GetLevelDB(name string) *leveldb.DB {
	return db.leveldb[name]
}

func (db *Manager) GetSQLite(name string) *SQLite {
	return db.sqlite[name]
}
