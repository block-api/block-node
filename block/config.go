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
package block

import (
	"os"

	"github.com/block-api/block-node/common"
	"github.com/block-api/block-node/params"
)

func loadConfigFile(configFilePath string) (*params.NodeConfig, error) {
	if configFilePath == "" {
		return nil, ErrConfigFileNotFound
	}

	var config = new(params.NodeConfig)

	cfgFile, err := common.OpenFile(configFilePath, common.YML)
	if err != nil {
		return nil, err
	}

	if cfgFile != nil {
		err = cfgFile.Parse(config)
		if err != nil {
			return nil, err
		}
	}

	// check if DATA_DIR is set in env, if yes then override value from config file
	dataDir := os.Getenv("DATA_DIR")
	if dataDir != "" {
		config.DataDir = dataDir
	}

	return config, nil
}
