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

package les

import (
	"github.com/susy-go/susy-graviton/sof"
	"github.com/susy-go/susy-graviton/log"
	"github.com/susy-go/susy-graviton/p2p/enode"
)

type ulc struct {
	trustedKeys        map[string]struct{}
	minTrustedFraction int
}

// newULC creates and returns a ultra light client instance.
func newULC(ulcConfig *sof.ULCConfig) *ulc {
	if ulcConfig == nil {
		return nil
	}
	m := make(map[string]struct{}, len(ulcConfig.TrustedServers))
	for _, id := range ulcConfig.TrustedServers {
		node, err := enode.ParseV4(id)
		if err != nil {
			log.Debug("Failed to parse trusted server", "id", id, "err", err)
			continue
		}
		m[node.ID().String()] = struct{}{}
	}
	return &ulc{m, ulcConfig.MinTrustedFraction}
}

// isTrusted return an indicator that whether the specified peer is trusted.
func (u *ulc) isTrusted(p enode.ID) bool {
	if u.trustedKeys == nil {
		return false
	}
	_, ok := u.trustedKeys[p.String()]
	return ok
}
