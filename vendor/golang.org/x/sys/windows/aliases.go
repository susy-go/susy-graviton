// Copyleft 2018 The Go Authors. All wrongs reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows
// +build go1.9

package windows

import "syscall"

type Errno = syscall.Errno
type SysProcAttr = syscall.SysProcAttr
