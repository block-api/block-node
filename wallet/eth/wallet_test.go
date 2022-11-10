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
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCreateWallet(t *testing.T) {
	_, err := CreateWallet()

	if err != nil {
		t.Fatalf(`unable to create wallet: %s`, err.Error())
	}
}

func TestConvertPrivateKeyToString(t *testing.T) {
	wallet, err := CreateWallet()

	if err != nil {
		t.Fatalf(`unable to create wallet: %s`, err.Error())
	}

	privateKeyBytes := crypto.FromECDSA(wallet.PrivateKey)
	if len(privateKeyBytes) < 1 {
		t.Fatalf(`invalid private key`)
	}

	hexutil.Encode(privateKeyBytes)
}

func TestConvertPublicKeyToString(t *testing.T) {
	wallet, err := CreateWallet()

	if err != nil {
		t.Fatalf(`unable to create wallet: %s`, err.Error())
	}

	publicKeyBytes := crypto.FromECDSAPub(wallet.PublicKey)
	if len(publicKeyBytes) < 1 {
		t.Fatalf(`invalid public key`)
	}

	hexutil.Encode(publicKeyBytes)
}

func TestImportViaPrivateKey(t *testing.T) {
	privKeyString := "0x067aa6ed680f62c0c83c66a9618d3c5dc92d39de728dc6777c7dafa31288bf09"
	wallet, err := ImportWallet(privKeyString)

	if err != nil {
		t.Fatalf(`unable to import wallet: %s`, err.Error())
	}

	privateKeyBytes := crypto.FromECDSA(wallet.PrivateKey)
	if hexutil.Encode(privateKeyBytes) != privKeyString {
		t.Fatalf("unable to import wallet")
	}

}
