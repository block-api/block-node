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

// Package repo
package repo

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/block-api/block-node/block/sys/model"
	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/params"
)

var counterRepo *CounterRepo
var counterRepoLock = new(sync.Mutex)

type CounterRepo struct {
	db *db.LevelDB
}

func GetCounter() *CounterRepo {
	if counterRepo == nil {
		counterRepoLock.Lock()
		defer counterRepoLock.Unlock()

		counterRepo = &CounterRepo{
			db: db.GetManager().GetLevelDB(params.DBSysCounters),
		}
	}

	return counterRepo
}

// func (r *KnownNodeRepo) Add(node model.KnownNode) error {
// 	bytes, err := json.Marshal(node)
// 	if err != nil {
// 		return nil
// 	}

// 	// r.db.DB.

// 	err = r.db.DB.Put([]byte(node.NodeID), bytes, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *CounterRepo) Get(key string) (*model.Counter, error) {
	kBytes := []byte(key)

	hasKey, err := r.db.DB.Has(kBytes, nil)
	if err != nil {
		return nil, err
	}

	if hasKey {
		resBytes, err := r.db.DB.Get(kBytes, nil)
		if err != nil {
			return nil, err
		}

		var counter model.Counter
		err = json.Unmarshal(resBytes, &counter)
		if err != nil {
			return nil, err
		}

		return &counter, nil
	}
	return nil, errors.New("node id not found")
}

func (r *CounterRepo) Put(key string, value int64) error {
	bytes, err := json.Marshal(model.Counter{Value: value})
	if err != nil {
		return err
	}

	err = r.db.DB.Put([]byte(key), bytes, nil)
	if err != nil {
		return err
	}

	return nil
}
