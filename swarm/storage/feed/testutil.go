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

package feed

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/susy-go/susy-graviton/p2p/enode"
	"github.com/susy-go/susy-graviton/swarm/chunk"
	"github.com/susy-go/susy-graviton/swarm/storage"
	"github.com/susy-go/susy-graviton/swarm/storage/localstore"
)

const (
	testDbDirName = "feeds"
)

type TestHandler struct {
	*Handler
}

func (t *TestHandler) Close() {
	t.chunkStore.Close()
}

type mockNetFetcher struct{}

func (m *mockNetFetcher) Request(hopCount uint8) {
}
func (m *mockNetFetcher) Offer(source *enode.ID) {
}

func newFakeNetFetcher(context.Context, storage.Address, *sync.Map) storage.NetFetcher {
	return &mockNetFetcher{}
}

// NewTestHandler creates Handler object to be used for testing purposes.
func NewTestHandler(datadir string, params *HandlerParams) (*TestHandler, error) {
	path := filepath.Join(datadir, testDbDirName)
	fh := NewHandler(params)

	db, err := localstore.New(filepath.Join(path, "chunks"), make([]byte, 32), nil)
	if err != nil {
		return nil, err
	}

	localStore := chunk.NewValidatorStore(db, storage.NewContentAddressValidator(storage.MakeHashFunc(feedsHashAlgorithm)), fh)

	netStore, err := storage.NewNetStore(localStore, nil)
	if err != nil {
		return nil, err
	}
	netStore.NewNetFetcherFunc = newFakeNetFetcher
	fh.SetStore(netStore)
	return &TestHandler{fh}, nil
}
