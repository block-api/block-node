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

// Package sqlite
package sqlite

const CreateMigrationTable string = `
	CREATE TABLE IF NOT EXISTS migration (
		name TEXT NOT NULL PRIMARY KEY,
		created_at INTEGER NOT NULL
	);
`

const NewMigrationEntry string = `
	INSERT INTO migration (name, created_at) VALUES (?, ?)
`

const FindMigrationEntry string = `
	SELECT name FROM migration WHERE name = ?
`

const DeleteMigrationEntry string = `
	DELETE FROM migration WHERE name = ?
`

type SQLMigration struct {
	Name      string
	UpQuery   string
	DownQuery string
}

func NewSQLMigration(name string, upQuery string, downQuery string) SQLMigration {
	return SQLMigration{
		Name:      name,
		UpQuery:   upQuery,
		DownQuery: downQuery,
	}
}
