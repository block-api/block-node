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
	"errors"
	"os"
	"sync"

	"github.com/block-api/block-node/wallet/eth"
)

var (
	ErrNodeAccountAlreadyInstantiatied = errors.New("node account is already instantiated")
	ErrNodePrivateKeyMissing           = errors.New("node account private key (NODE_PRIVATE_KEY) is missing")
	ErrNodeAccountInvalidPrivateKey    = errors.New("node account private key is invalid")
)

var (
	nodeAccount     *NodeAccount
	nodeAccountLock = new(sync.Mutex)
)

// NodeAccount is an account assigned to current node. Accounts are ETH wallets with
// public and private keys used to sign and verify messages sent between nodes.
// Upon startup it will check for provided environment variable NODE_PRIVATE_KEY,
// and in DATA_DIR directory (for json file with account credentials) and validate them.
//
// If none of these is present it will not allow to start node instance.
// To generate new node account you can use `new-account` command for cli eg block-node new-account.
//
// There is allowed only one instance of NodeAccount in the system.
// If you wish to use ETH wallets in your project you can use wallet/eth package instead.
type NodeAccount struct {
	wallet *eth.Wallet
}

// NewNodeAccount is responsible for loading node account information and validating them
// it wll try to restore node account from NODE_PRIVATE_KEY environment variable
func NewNodeAccount() (*NodeAccount, error) {
	if nodeAccount == nil {
		nodeAccountLock.Lock()
		defer nodeAccountLock.Unlock()

		nodePrivateKey := os.Getenv("NODE_PRIVATE_KEY")
		if nodePrivateKey == "" {
			return nil, ErrNodePrivateKeyMissing
		}

		wallet, err := eth.ImportWallet(nodePrivateKey)
		if err != nil {
			if err == eth.ErrInvalidPrivateKey {
				return nil, ErrNodeAccountInvalidPrivateKey
			}
			return nil, err
		}

		nodeAccount = &NodeAccount{
			wallet: &wallet,
		}

		return nodeAccount, nil
	}
	return nil, ErrNodeAccountAlreadyInstantiatied
}
