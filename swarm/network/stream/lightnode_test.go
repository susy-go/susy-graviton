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
package stream

import (
	"testing"

	p2ptest "github.com/susy-go/susy-graviton/p2p/testing"
)

// This test checks the default behavior of the server, that is
// when syncing is enabled.
func TestLigthnodeRequestSubscriptionWithSync(t *testing.T) {
	registryOptions := &RegistryOptions{
		Syncing: SyncingRegisterOnly,
	}
	tester, _, _, teardown, err := newStreamerTester(registryOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	node := tester.Nodes[0]

	syncStream := NewStream("SYNC", FormatSyncBinKey(1), false)

	err = tester.TestExchanges(
		p2ptest.Exchange{
			Label: "RequestSubscription",
			Triggers: []p2ptest.Trigger{
				{
					Code: 8,
					Msg: &RequestSubscriptionMsg{
						Stream: syncStream,
					},
					Peer: node.ID(),
				},
			},
			Expects: []p2ptest.Expect{
				{
					Code: 4,
					Msg: &SubscribeMsg{
						Stream: syncStream,
					},
					Peer: node.ID(),
				},
			},
		})

	if err != nil {
		t.Fatalf("Got %v", err)
	}
}

// This test checks the Lightnode behavior of the server, that is
// when syncing is disabled.
func TestLigthnodeRequestSubscriptionWithoutSync(t *testing.T) {
	registryOptions := &RegistryOptions{
		Syncing: SyncingDisabled,
	}
	tester, _, _, teardown, err := newStreamerTester(registryOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	node := tester.Nodes[0]

	syncStream := NewStream("SYNC", FormatSyncBinKey(1), false)

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RequestSubscription",
		Triggers: []p2ptest.Trigger{
			{
				Code: 8,
				Msg: &RequestSubscriptionMsg{
					Stream: syncStream,
				},
				Peer: node.ID(),
			},
		},
		Expects: []p2ptest.Expect{
			{
				Code: 7,
				Msg: &SubscribeErrorMsg{
					Error: "stream SYNC not registered",
				},
				Peer: node.ID(),
			},
		},
	}, p2ptest.Exchange{
		Label: "RequestSubscription",
		Triggers: []p2ptest.Trigger{
			{
				Code: 4,
				Msg: &SubscribeMsg{
					Stream: syncStream,
				},
				Peer: node.ID(),
			},
		},
		Expects: []p2ptest.Expect{
			{
				Code: 7,
				Msg: &SubscribeErrorMsg{
					Error: "stream SYNC not registered",
				},
				Peer: node.ID(),
			},
		},
	})

	if err != nil {
		t.Fatalf("Got %v", err)
	}
}
