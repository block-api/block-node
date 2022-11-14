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

// Package db
package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/db/sqlite"
	"github.com/block-api/block-node/params"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	ErrDbManagerInstantiated = errors.New("database manager is alredy instantiated")
	ErrDbNameExist           = errors.New("database with provided name already exist")
)

var lock = new(sync.Mutex)
var manager *Manager

type Manager struct {
	config    *params.NodeConfig
	databases map[string]bool
	leveldb   map[string]*LevelDB
	sqlite    map[string]*SQLite
}

// NewManager creates new Manager struct
func NewManager(config *params.NodeConfig) (*Manager, error) {
	if manager == nil {
		lock.Lock()
		defer lock.Unlock()

		manager = &Manager{
			config:    config,
			databases: make(map[string]bool),
			leveldb:   make(map[string]*LevelDB),
			sqlite:    make(map[string]*SQLite),
		}

		err := manager.startSysDatabases()
		if err != nil {
			return nil, err
		}

		for dbName, dbConfig := range config.Database {
			if manager.databases[dbName] {
				return nil, ErrDbNameExist
			}

			settingsBytes, err := json.Marshal(dbConfig.Settings)
			if err != nil {
				return nil, err
			}

			if dbConfig.Type == types.DbLevelDB {
				var dbSettings params.DatabaseLevelDBConfig

				err := json.Unmarshal(settingsBytes, &dbSettings)
				if err != nil {
					return nil, err
				}

				options := &opt.Options{
					WriteBuffer: dbSettings.WriteBufferSize * opt.MiB,
					Filter:      filter.NewBloomFilter(10),
				}

				dbInstance, err := leveldb.OpenFile(config.DataDir+params.DBDir+dbName, options)
				if err != nil {
					return nil, err
				}

				manager.leveldb[dbName] = &LevelDB{
					DB: dbInstance,
				}
				manager.databases[dbName] = true

				continue
			}

			if dbConfig.Type == types.DbSqlite {
				var dbSettings params.DatabaseSqliteConfig

				err := json.Unmarshal(settingsBytes, &dbSettings)
				if err != nil {
					return nil, err
				}

				dbPath := config.DataDir + params.DBDir + dbName + ".sqlite"
				if _, err := os.Stat(dbPath); err != nil {
					file, err := os.Create(dbPath)
					if err != nil {
						return nil, err
					}

					err = file.Close()
					if err != nil {
						return nil, err
					}
				}

				dbInstance, err := sql.Open("sqlite3", "file:"+dbPath+"?"+dbSettings.Options)
				if err != nil {
					return nil, err
				}

				err = dbInstance.Ping()
				if err != nil {
					return nil, err
				}

				manager.sqlite[dbName] = &SQLite{
					DB:         dbInstance,
					migrations: make([]sqlite.SQLMigration, 0),
				}
				manager.databases[dbName] = true

				if dbSettings.MaxOpenConnections > 0 {
					manager.sqlite[dbName].DB.SetMaxOpenConns(dbSettings.MaxOpenConnections)
				}

				if _, err := manager.sqlite[dbName].DB.Exec(sqlite.CreateMigrationTable); err != nil {
					return nil, err
				}

				manager.databases[dbName] = true

				continue
			}
		}

		return manager, nil
	}
	return nil, ErrDbManagerInstantiated
}

func GetManager() *Manager {
	return manager
}

func (db *Manager) startSysDatabases() error {
	if db.databases[params.DBSysKnownNodes] || db.databases[params.DBSys] {
		return ErrDbNameExist
	}

	dbSys, err := leveldb.OpenFile(db.config.DataDir+params.DBDir+params.DBSys, nil)
	if err != nil {
		return err
	}
	db.databases[params.DBSys] = true
	db.leveldb[params.DBSys] = &LevelDB{
		DB: dbSys,
	}

	dbKnownNodes, err := leveldb.OpenFile(db.config.DataDir+params.DBDir+params.DBSysKnownNodes, nil)
	if err != nil {
		return err
	}

	db.databases[params.DBSysKnownNodes] = true
	db.leveldb[params.DBSysKnownNodes] = &LevelDB{
		DB: dbKnownNodes,
	}

	return nil
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

func (db *Manager) GetLevelDB(name string) *LevelDB {
	return db.leveldb[name]
}

func (db *Manager) GetSQLite(name string) *SQLite {
	return db.sqlite[name]
}
