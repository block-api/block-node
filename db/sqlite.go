package db

import "database/sql"

type SQLite struct {
	Db *sql.DB
}

//if _, err := newDb.sqlite[dbName].Exec(CreateTest); err != nil {
//	panic(err)
//}

//const CreateTest string = `
//  CREATE TABLE IF NOT EXISTS test (
//  id INTEGER NOT NULL PRIMARY KEY,
//  time DATETIME NOT NULL,
//  description TEXT
//  );`
