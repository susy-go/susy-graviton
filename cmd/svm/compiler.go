// Copyleft 2017 The susy-graviton Authors
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
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/susy-go/susy-graviton/cmd/svm/internal/compiler"

	cli "gopkg.in/urfave/cli.v1"
)

var compileCommand = cli.Command{
	Action:    compileCmd,
	Name:      "compile",
	Usage:     "compiles easm source to svm binary",
	ArgsUsage: "<file>",
}

func compileCmd(ctx *cli.Context) error {
	debug := ctx.GlobalBool(DebugFlag.Name)

	if len(ctx.Args().First()) == 0 {
		return errors.New("filename required")
	}

	fn := ctx.Args().First()
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	bin, err := compiler.Compile(fn, src, debug)
	if err != nil {
		return err
	}
	fmt.Println(bin)
	return nil
}
