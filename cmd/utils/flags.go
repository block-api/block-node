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
package utils

import (
	"github.com/block-api/block-node/params"
	"github.com/urfave/cli/v2"
)

var (
	GeneralCategoryFlag = "General"
	NetworkCategoryFlag = "Network"
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:        "config-file",
		Category:    GeneralCategoryFlag,
		Required:    false,
		DefaultText: params.DefaultDataDir + "/config.yml",
		Usage:       "Absolute path to config.yml file",
	}

	DataDirFlag = &cli.StringFlag{
		Name:        "data-dir",
		Category:    GeneralCategoryFlag,
		Required:    false,
		DefaultText: params.DefaultDataDir,
		Usage:       "Absolute path to directory where data are stored",
	}

	NetworkTransportFlag = &cli.StringFlag{
		Name:        "transport",
		Category:    NetworkCategoryFlag,
		DefaultText: string(params.DefaultNetworkTransport),
		Usage:       "Network transport method eg tcp / redis",
	}

	NetworkBindHostFlag = &cli.StringFlag{
		Name:        "bind-host",
		Category:    NetworkCategoryFlag,
		DefaultText: params.DefaultNetworkBindHost,
		Usage:       "Set bind host",
	}

	NetworkBindPortFlag = &cli.StringFlag{
		Name:        "bind-port",
		Category:    NetworkCategoryFlag,
		DefaultText: params.DefaultNetworkBindPort,
		Usage:       "Set bind port",
	}
)
