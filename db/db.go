package db

import (
	"database/sql"
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"sync"
)

var lock = new(sync.Mutex)
var database *Database

type Database struct {
	config  *config.Database
	leveldb map[string]*leveldb.DB
	sqlite  map[string]*SQLite
}

// NewDatabase creates new database struct

func NewDatabase(config *config.Database) *Database {
	if database == nil {
		lock.Lock()
		defer lock.Unlock()

		database = &Database{
			config:  config,
			leveldb: make(map[string]*leveldb.DB),
			sqlite:  make(map[string]*SQLite),
		}

		var err error

		if len(config.LevelDB) > 0 {
			for dbName, dbConfig := range config.LevelDB {
				o := &opt.Options{
					WriteBuffer: dbConfig.WriteBufferSize * opt.MiB,
					Filter:      filter.NewBloomFilter(10),
				}

				database.leveldb[dbName], err = leveldb.OpenFile(dbConfig.DbPath, o)

				if err != nil {
					panic(err)
				}

				log.Debug("database " + dbName + ": " + dbConfig.DbPath)
			}
		}

		if len(config.SQLite) > 0 {
			for dbName, dbConfig := range config.SQLite {
				if _, err := os.Stat(dbConfig.DbPath); err != nil {
					file, err := os.Create(dbConfig.DbPath)
					if err != nil {
						log.Panic(err.Error())
					}
					_ = file.Close()
				}

				s3db, err := sql.Open("sqlite3", "file:"+dbConfig.DbPath+"?"+dbConfig.Options)
				if err != nil {
					log.Panic(err.Error())
				}

				err = s3db.Ping()
				if err != nil {
					log.Panic(err.Error())
				}

				database.sqlite[dbName] = &SQLite{
					Db: s3db,
				}

				if dbConfig.MaxOpenConnections > 0 {
					database.sqlite[dbName].Db.SetMaxOpenConns(dbConfig.MaxOpenConnections)
				}

				if _, err := database.sqlite[dbName].Db.Exec(CreateMigrationTable); err != nil {
					log.Panic(err.Error())
				}

				log.Debug("database " + dbName + ": " + dbConfig.DbPath)
			}
		}
	}

	return database
}

func GetDatabase() *Database {
	return database
}

func (db *Database) RunMigrations() error {
	for _, dbSQLite := range db.sqlite {
		err := dbSQLite.RunMigrations()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) RevertMigrations() error {
	for _, dbSQLite := range db.sqlite {
		err := dbSQLite.RevertMigrations()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) GetLevelDB(name string) *leveldb.DB {
	return db.leveldb[name]
}

func (db *Database) GetSQLite(name string) *SQLite {
	return db.sqlite[name]
}
