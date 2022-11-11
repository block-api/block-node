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
package params

import "github.com/block-api/block-node/common/types"

// DbDir is a directory inside dataDir (from config file) directory
// where databases data will be kept

const DBDir = "/db/"

// DbNamePrefix prefix for internal databases name
const DBNamePrefix = "bn_"

// DbSys block node internal db
const DBSys = DBNamePrefix + "sys"

// DbSysKnownNodes name of internal database where system keeps
// informations about network topology
const DBSysKnownNodes = DBNamePrefix + "sys_known_nodes"

// DatabasesConfig describes parameters available for databases
type DatabasesConfig map[string]DatabaseConfig

// DatabaseConfig describes single db config options
type DatabaseConfig struct {
	Type     types.DBType `json:"type"`
	Settings any          `json:"settings"`
}

// DatabaseLevelDBConfig settings for LevelDB database
type DatabaseLevelDBConfig struct {
	WriteBufferSize int `json:"write_buffer_size"`
}

// DatabaseSqliteConfig settings for Sqlite database
type DatabaseSqliteConfig struct {
	Path               string `json:"path"`
	Options            string `json:"options"`
	MaxOpenConnections int    `json:"max_open_connections"`
}
