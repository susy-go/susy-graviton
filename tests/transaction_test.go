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
	"math/big"
	"testing"

	"github.com/susy-go/susy-graviton/params"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	txt := new(testMatcher)
	txt.config(`^Homestead/`, params.ChainConfig{
		HomesteadBlock: big.NewInt(0),
	})
	txt.config(`^SIP155/`, params.ChainConfig{
		HomesteadBlock: big.NewInt(0),
		SIP150Block:    big.NewInt(0),
		SIP155Block:    big.NewInt(0),
		SIP158Block:    big.NewInt(0),
		ChainID:        big.NewInt(1),
	})
	txt.config(`^Byzantium/`, params.ChainConfig{
		HomesteadBlock: big.NewInt(0),
		SIP150Block:    big.NewInt(0),
		SIP155Block:    big.NewInt(0),
		SIP158Block:    big.NewInt(0),
		ByzantiumBlock: big.NewInt(0),
	})

	txt.walk(t, transactionTestDir, func(t *testing.T, name string, test *TransactionTest) {
		cfg := txt.findConfig(name)
		if err := txt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
