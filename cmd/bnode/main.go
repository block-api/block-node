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
package main

import (
	"log"
	"os"

	"github.com/block-api/block-node/cmd/utils"
	"github.com/urfave/cli/v2"
)

var (
	nodeFlags = utils.FlagsMerge(
		utils.ConfigFileFlag,
		utils.DataDirFlag,
		// network flags
		utils.NetworkTransportFlag,
		utils.NetworkBindHostFlag,
		utils.NetworkBindPortFlag,
	)

	app = &cli.App{
		Name:        "block-node",
		Version:     "0.2.0",
		Flags:       nodeFlags,
		Usage:       "Block Node Functional Framework in Go",
		Description: "You can use flags to override options from config file or environment variables",
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
