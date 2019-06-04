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

// +build race

package testutil

// RaceEnabled is true when -race flag is provided to the go tool.
// See norace.go for more.
const RaceEnabled = true
