// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is Free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Mail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

//go:build gpa_server

package internal

import (
	"context"
	"github.com/ProtonMail/go-proton-api/server"
)

type GPAServer struct {
	ctx    context.Context
	cancel func()
	server *server.Server
}

func NewGPAServer(ctx context.Context) *GPAServer {
	srv := server.New(server.WithTLS(false))

	ctx, cancel := context.WithCancel(ctx)

	return &GPAServer{
		ctx:    ctx,
		cancel: cancel,
		server: srv,
	}
}

func (g *GPAServer) Close() {
	g.cancel()
	g.server.Close()
}

func (g *GPAServer) CreateUser(email, password string) (string, string, error) {
	return g.server.CreateUser(email, []byte(password))
}

func (g *GPAServer) GetURL() string {
	return g.server.GetHostURL()
}
