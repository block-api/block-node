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

// Package pool
package pool

import "sync"

type WorkerHandler func(payload interface{})

type Pool struct {
	name         string
	maxWorkers   int
	workersCount int
	worker       WorkerHandler
	cIn          chan interface{}
	cOut         chan interface{}
	lockIn       *sync.Mutex
}

func NewPool(name string, maxWorkers int, workerHandler WorkerHandler) Pool {
	return Pool{
		name:       name,
		maxWorkers: maxWorkers,
		cIn:        make(chan interface{}),
		cOut:       make(chan interface{}),
		lockIn:     new(sync.Mutex),
	}
}