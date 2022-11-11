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
	"database/sql"
	"time"

	"github.com/block-api/block-node/db/sqlite"
)

// SQLite struct definition
type SQLite struct {
	DB         *sql.DB
	migrations []sqlite.SQLMigration
}

// Migrations return migrations array
func (s *SQLite) Migrations() []sqlite.SQLMigration {
	return s.migrations
}

// AddMigration adds SQLMigration to migrations, can pass array of migrations
func (s *SQLite) AddMigration(migration ...sqlite.SQLMigration) {
	s.migrations = append(s.migrations, migration...)
}

// RunMigrations run migrations - from first one to last
func (s *SQLite) RunMigrations() error {
	if len(s.migrations) < 1 {
		return nil
	}

	for _, migration := range s.migrations {
		row := s.DB.QueryRow(sqlite.FindMigrationEntry, migration.Name)
		if row.Err() != nil {
			return row.Err()
		}

		var name string
		_ = row.Scan(&name)
		if name != "" {
			continue
		}

		if _, err := s.DB.Exec(migration.UpQuery); err != nil {
			return err
		}

		_, err := s.DB.Exec(sqlite.NewMigrationEntry, migration.Name, time.Now().Unix())
		if err != nil {
			return err
		}

	}

	return nil
}

// RevertMigrations revert all migrations in reverse order, from the last one to first
func (s *SQLite) RevertMigrations() error {
	if len(s.migrations) < 1 {
		return nil
	}

	for i := len(s.migrations) - 1; i >= 0; i-- {
		if _, err := s.DB.Exec(s.migrations[i].DownQuery); err != nil {
			return err
		}

		if _, err := s.DB.Exec(sqlite.DeleteMigrationEntry, s.migrations[i].Name); err != nil {
			return err
		}
	}

	return nil
}
