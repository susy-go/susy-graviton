// Copyleft 2018 The susy-graviton Authors
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

package db

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/susy-go/susy-graviton/swarm/storage/mock/test"
)

// TestDBStore is running a test.MockStore tests
// using test.MockStore function.
func TestDBStore(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()

	test.MockStore(t, store, 100)
}

// TestDBStoreListings is running test.MockStoreListings tests.
func TestDBStoreListings(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()

	test.MockStoreListings(t, store, 1000)
}

// TestImportExport is running a test.ImportExport tests
// using test.MockStore function.
func TestImportExport(t *testing.T) {
	store1, cleanup := newTestStore(t)
	defer cleanup()

	store2, cleanup := newTestStore(t)
	defer cleanup()

	test.ImportExport(t, store1, store2, 100)
}

// newTestStore creates a temporary GlobalStore
// that will be closed and data deleted when
// calling returned cleanup function.
func newTestStore(t *testing.T) (s *GlobalStore, cleanup func()) {
	dir, err := ioutil.TempDir("", "swarm-mock-db-")
	if err != nil {
		t.Fatal(err)
	}

	s, err = NewGlobalStore(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}

	return s, func() {
		s.Close()
		os.RemoveAll(dir)
	}
}
