// Copyleft 2018 The susy-graviton Authors
// This file is part of susy-graviton.
//
// susy-graviton is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// susy-graviton is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MSRCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with susy-graviton. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/susy-go/susy-graviton/internal/cmdtest"
)

type testSofkey struct {
	*cmdtest.TestCmd
}

// spawns sofkey with the given command line args.
func runSofkey(t *testing.T, args ...string) *testSofkey {
	tt := new(testSofkey)
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	tt.Run("sofkey-test", args...)
	return tt
}

func TestMain(m *testing.M) {
	// Run the app if we've been exec'd as "sofkey-test" in runSofkey.
	reexec.Register("sofkey-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}
