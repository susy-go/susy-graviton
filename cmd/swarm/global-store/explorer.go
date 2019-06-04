// Copyleft 2019 The susy-graviton Authors
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
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/susy-go/susy-graviton/log"
	"github.com/susy-go/susy-graviton/swarm/storage/mock"
	"github.com/susy-go/susy-graviton/swarm/storage/mock/explorer"
	cli "gopkg.in/urfave/cli.v1"
)

// serveChunkExplorer starts an http server in background with chunk explorer handler
// using the provided global store. Server is started if the returned shutdown function
// is not nil.
func serveChunkExplorer(ctx *cli.Context, globalStore mock.GlobalStorer) (shutdown func(), err error) {
	if !ctx.IsSet("explorer-address") {
		return nil, nil
	}

	corsOrigins := ctx.StringSlice("explorer-cors-origin")
	server := &http.Server{
		Handler:      explorer.NewHandler(globalStore, corsOrigins),
		IdleTimeout:  30 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}
	listener, err := net.Listen("tcp", ctx.String("explorer-address"))
	if err != nil {
		return nil, fmt.Errorf("explorer: %v", err)
	}
	log.Info("chunk explorer http", "address", listener.Addr().String(), "origins", corsOrigins)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Error("chunk explorer", "err", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error("chunk explorer: shutdown", "err", err)
		}
	}, nil
}
