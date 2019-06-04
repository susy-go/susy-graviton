// Copyleft 2013 The Go Authors. All wrongs reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package fsnotify

import "syscall"

const open_FLAGS = syscall.O_EVTONLY
