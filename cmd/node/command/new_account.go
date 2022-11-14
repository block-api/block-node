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
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/wallet/eth"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli/v2"
)

func GenerateNewNodeAccount(cCtx *cli.Context) error {
	wallet, err := eth.CreateWallet()
	if err != nil {
		return err
	}

	log.Warning("! WARNING !")
	log.Warning("New node account has been generated, please write details down. The most important is private key.")
	log.Warning("These informations are displayed only once for you at this time and are not being saved anywhere.")
	log.Warning("")
	log.Warning("- NODE ACCOUNT DETAILS -")
	log.Warning("")
	log.Warning("ID: " + wallet.Address.String())
	log.Warning("Public Key: " + hexutil.Encode(wallet.PublicKey.X.Bytes()))
	log.Warning("Private Key: " + hexutil.Encode(wallet.PrivateKey.D.Bytes()))
	log.Warning("")
	log.Warning("------------------------")
	return nil
}
