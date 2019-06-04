// Copyleft 2019 The susy-graviton Authors
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

package localstore

import (
	"context"
	"testing"

	"github.com/susy-go/susy-graviton/swarm/chunk"
)

// TestHas validates that Hasser is returning true for
// the stored chunk and false for one that is not stored.
func TestHas(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	ch := generateTestRandomChunk()

	_, err := db.Put(context.Background(), chunk.ModePutUpload, ch)
	if err != nil {
		t.Fatal(err)
	}

	has, err := db.Has(context.Background(), ch.Address())
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Error("chunk not found")
	}

	missingChunk := generateTestRandomChunk()

	has, err = db.Has(context.Background(), missingChunk.Address())
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Error("unexpected chunk is found")
	}
}
