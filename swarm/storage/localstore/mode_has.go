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
	"time"

	"github.com/susy-go/susy-graviton/metrics"
	"github.com/susy-go/susy-graviton/swarm/chunk"
)

// Has returns true if the chunk is stored in database.
func (db *DB) Has(ctx context.Context, addr chunk.Address) (bool, error) {
	metricName := "localstore.Has"

	metrics.GetOrRegisterCounter(metricName, nil).Inc(1)
	defer totalTimeMetric(metricName, time.Now())

	has, err := db.retrievalDataIndex.Has(addressToItem(addr))
	if err != nil {
		metrics.GetOrRegisterCounter(metricName+".error", nil).Inc(1)
	}
	return has, err
}
