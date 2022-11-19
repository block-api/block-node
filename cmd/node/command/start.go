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

// Package command
package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/block-api/block-node/block"
	"github.com/block-api/block-node/log"
	"github.com/urfave/cli/v2"
)

func StartNode(cCtx *cli.Context) error {
	node, err := block.NewNode()
	if err != nil {
		return err
	}

	log.Default("starting block node v" + cCtx.App.Version)
	log.Default("node id: " + node.ID())

	node.Start()

	var osSignal chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-osSignal:
			log.Warning("shutting down, please wait")
			node := block.GetNode()

			if node != nil {
				block.GetNode().Stop()
			}

			os.Exit(0)
		}
	}
}
