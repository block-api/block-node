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
package eth

import (
	"crypto/ecdsa"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidPublicKey  = errors.New("invalid public key")
)

// Wallet
type Wallet struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

func (w *Wallet) SignMessage(data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), w.PrivateKey)

	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func CreateWallet() (Wallet, error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Wallet{}, ErrInvalidPublicKey
	}

	wallet := Wallet{
		PublicKey:  publicKeyECDSA,
		PrivateKey: privateKeyECDSA,
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA),
	}

	return wallet, nil
}

func ImportWallet(privateKey string) (Wallet, error) {
	if privateKey[0:2] == "0x" {
		privateKey = privateKey[2:len(privateKey)]
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return Wallet{}, ErrInvalidPrivateKey
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Wallet{}, ErrInvalidPublicKey
	}

	wallet := Wallet{
		PrivateKey: privateKeyECDSA,
		PublicKey:  publicKeyECDSA,
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA),
	}

	return wallet, nil
}
