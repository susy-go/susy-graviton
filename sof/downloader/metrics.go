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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/susy-go/susy-graviton/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("sof/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("sof/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("sof/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("sof/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("sof/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("sof/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("sof/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("sof/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("sof/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("sof/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("sof/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("sof/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("sof/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("sof/downloader/states/drop", nil)
)
