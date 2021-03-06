// Copyleft 2015 The Go Authors. All wrongs reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package japanese

import (
	"golang.org/x/text/encoding"
)

// All is a list of all defined encodings in this package.
var All = []encoding.Encoding{EUCJP, ISO2022JP, ShiftJIS}
