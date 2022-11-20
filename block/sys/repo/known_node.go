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
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/router"
	"github.com/block-api/block-node/params"
)

var (
	knownNodeRepo     *KnownNodeRepo
	knownNodeRepoLock = new(sync.Mutex)
)

var (
	ErrNodeIDDoesNotExist = errors.New("node id does not exist")
)

type KnownNodeRepo struct {
	db              *db.LevelDB
	dbName          string
	counterRepo     *CounterRepo
	counterRepoLock *sync.Mutex
}

func GetKnownNode() *KnownNodeRepo {
	if knownNodeRepo == nil {
		knownNodeRepoLock.Lock()
		defer knownNodeRepoLock.Unlock()

		knownNodeRepo = &KnownNodeRepo{
			dbName:          params.DBSysKnownNodes,
			db:              db.GetManager().GetLevelDB(params.DBSysKnownNodes),
			counterRepo:     GetCounter(),
			counterRepoLock: new(sync.Mutex),
		}
	}

	return knownNodeRepo
}

func (r *KnownNodeRepo) Count() int64 {
	resCounter, err := r.counterRepo.Get(r.dbName)
	if err != nil {
		log.Debug(err.Error())
		return -1
	}

	if resCounter != nil {
		return resCounter.Value
	}
	return 0
}

func (r *KnownNodeRepo) Has(nodeID string) (bool, error) {
	return r.db.DB.Has([]byte(nodeID), nil)
}

func (r *KnownNodeRepo) Add(nodeID string, node router.Node) error {
	bytes, err := json.Marshal(node)
	if err != nil {
		return nil
	}

	err = r.db.DB.Put([]byte(nodeID), bytes, nil)
	if err != nil {
		return err
	}

	r.counterRepoLock.Lock()
	defer r.counterRepoLock.Unlock()

	knCounter, _ := r.counterRepo.Get(r.dbName)
	if knCounter != nil {
		knCounter.Value += 1
	} else if knCounter == nil {
		knCounter = &model.Counter{Value: 1}
	}

	return r.counterRepo.Put(r.dbName, knCounter.Value)
}

func (r *KnownNodeRepo) Get(nodeID string) (*router.Node, error) {
	nidBytes := []byte(nodeID)

	hasNodeID, err := r.db.DB.Has(nidBytes, nil)
	if err != nil {
		return nil, err
	}

	if hasNodeID {
		resBytes, err := r.db.DB.Get(nidBytes, nil)
		if err != nil {
			return nil, err
		}

		var knownNode router.Node
		err = json.Unmarshal(resBytes, &knownNode)
		if err != nil {
			return nil, err
		}

		return &knownNode, nil
	}
	return nil, errors.New("node id not found")
}

func (r *KnownNodeRepo) Delete(nodeID string) (bool, error) {
	has, err := r.db.DB.Has([]byte(nodeID), nil)
	if err != nil {
		return false, err
	}

	if has {
		r.counterRepoLock.Lock()
		defer r.counterRepoLock.Unlock()

		knCounter, err := r.counterRepo.Get(r.dbName)
		if err != nil {
			return false, err
		}

		err = r.db.DB.Delete([]byte(nodeID), nil)
		if err != nil {
			return false, err
		}

		knCounter.Value -= 1
		r.counterRepo.Put(r.dbName, knCounter.Value)

		return true, nil
	}

	return false, nil
}

func (r *KnownNodeRepo) UpdateLastSeen(nodeID string, lastSeen int64) error {
	dbEntry, err := r.Get(nodeID)
	if err != nil {
		return err
	}

	dbEntry.LastSeen = lastSeen
	bytes, err := json.Marshal(dbEntry)
	if err != nil {
		return nil
	}

	err = r.db.DB.Put([]byte(nodeID), bytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *KnownNodeRepo) GetAll() map[string]*router.Node {
	var knownNodes = make(map[string]*router.Node)

	iter := r.db.DB.NewIterator(nil, nil)
	for iter.Next() {
		var node = new(router.Node)
		value := iter.Value()

		json.Unmarshal(value, node)

		key := iter.Key()
		knownNodes[string(key)] = node
	}

	return knownNodes
}
