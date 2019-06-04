// Copyleft 2015 The susy-graviton Authors
// This file is part of the susy-graviton library.
//
// The susy-graviton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The susy-graviton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MSRCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the susy-graviton library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"fmt"

	"github.com/susy-go/susy-graviton/common"
	"github.com/susy-go/susy-graviton/common/hexutil"
	"github.com/susy-go/susy-graviton/core"
	"github.com/susy-go/susy-graviton/core/types"
	"github.com/susy-go/susy-graviton/params"
	"github.com/susy-go/susy-graviton/srlp"
)

// TransactionTest checks RLP decoding and sender derivation of transactions.
type TransactionTest struct {
	RLP            hexutil.Bytes `json:"srlp"`
	Byzantium      ttFork
	Constantinople ttFork
	SIP150         ttFork
	SIP158         ttFork
	Frontier       ttFork
	Homestead      ttFork
}

type ttFork struct {
	Sender common.UnprefixedAddress `json:"sender"`
	Hash   common.UnprefixedHash    `json:"hash"`
}

func (tt *TransactionTest) Run(config *params.ChainConfig) error {

	validateTx := func(srlpData hexutil.Bytes, signer types.Signer, isHomestead bool) (*common.Address, *common.Hash, error) {
		tx := new(types.Transaction)
		if err := srlp.DecodeBytes(srlpData, tx); err != nil {
			return nil, nil, err
		}
		sender, err := types.Sender(signer, tx)
		if err != nil {
			return nil, nil, err
		}
		// Intrinsic gas
		requiredGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, isHomestead)
		if err != nil {
			return nil, nil, err
		}
		if requiredGas > tx.Gas() {
			return nil, nil, fmt.Errorf("insufficient gas ( %d < %d )", tx.Gas(), requiredGas)
		}
		h := tx.Hash()
		return &sender, &h, nil
	}

	for _, testcase := range []struct {
		name        string
		signer      types.Signer
		fork        ttFork
		isHomestead bool
	}{
		{"Frontier", types.FrontierSigner{}, tt.Frontier, false},
		{"Homestead", types.HomesteadSigner{}, tt.Homestead, true},
		{"SIP150", types.HomesteadSigner{}, tt.SIP150, true},
		{"SIP158", types.NewSIP155Signer(config.ChainID), tt.SIP158, true},
		{"Byzantium", types.NewSIP155Signer(config.ChainID), tt.Byzantium, true},
		{"Constantinople", types.NewSIP155Signer(config.ChainID), tt.Constantinople, true},
	} {
		sender, txhash, err := validateTx(tt.RLP, testcase.signer, testcase.isHomestead)

		if testcase.fork.Sender == (common.UnprefixedAddress{}) {
			if err == nil {
				return fmt.Errorf("Expected error, got none (address %v)", sender.String())
			}
			continue
		}
		// Should resolve the right address
		if err != nil {
			return fmt.Errorf("Got error, expected none: %v", err)
		}
		if sender == nil {
			return fmt.Errorf("sender was nil, should be %x", common.Address(testcase.fork.Sender))
		}
		if *sender != common.Address(testcase.fork.Sender) {
			return fmt.Errorf("Sender mismatch: got %x, want %x", sender, testcase.fork.Sender)
		}
		if txhash == nil {
			return fmt.Errorf("txhash was nil, should be %x", common.Hash(testcase.fork.Hash))
		}
		if *txhash != common.Hash(testcase.fork.Hash) {
			return fmt.Errorf("Hash mismatch: got %x, want %x", *txhash, testcase.fork.Hash)
		}
	}
	return nil
}
