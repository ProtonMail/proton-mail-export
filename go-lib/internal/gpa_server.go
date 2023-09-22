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
	"fmt"
	"runtime"

	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/go-proton-api/server"
	"github.com/bradenaw/juniper/stream"
	"github.com/bradenaw/juniper/xslices"
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

func (g *GPAServer) CreateTestMessages(userID, addrID, email, password string,
	count int,
) ([]string, error) {
	var messageIDS []string

	const DummyMessage = `To: recipient@pm.me
From: sender@pm.me
Subject: Test
Content-Type: text/plain; charset=utf-8

Test

`
	labelID, err := g.server.CreateLabel(userID, "folder", "", proton.LabelTypeFolder)
	if err != nil {
		return nil, err
	}

	err = withClient(g.ctx, g.server, email, []byte(password), func(ctx context.Context, client *proton.Client) error {
		m, err := createMessagesWithFlags(ctx, client, addrID, labelID, []byte(password), 0, xslices.Repeat([]byte(DummyMessage), count)...)

		messageIDS = m

		return err
	})
	if err != nil {
		return nil, err
	}

	return messageIDS, nil
}

func withClient(
	ctx context.Context,
	s *server.Server,
	username string,
	password []byte,
	fn func(context.Context, *proton.Client) error,
) error { //nolint:unparam
	m := proton.New(
		proton.WithHostURL(s.GetHostURL()),
		proton.WithTransport(proton.InsecureTransport()),
	)

	c, _, err := m.NewClientWithLogin(ctx, username, password)
	if err != nil {
		return err
	}

	defer c.Close()

	return fn(ctx, c)
}

func createMessagesWithFlags(
	ctx context.Context,
	c *proton.Client,
	addrID, labelID string,
	password []byte,
	flags proton.MessageFlag,
	messages ...[]byte,
) ([]string, error) {
	user, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	addr, err := c.GetAddresses(ctx)
	if err != nil {
		return nil, err
	}

	salt, err := c.GetSalts(ctx)
	if err != nil {
		return nil, err
	}

	keyPass, err := salt.SaltForKey(password, user.Keys.Primary().ID)
	if err != nil {
		return nil, err
	}

	_, addrKRs, err := proton.Unlock(user, addr, keyPass, async.NoopPanicHandler{})
	if err != nil {
		return nil, err
	}

	_, ok := addrKRs[addrID]
	if !ok {
		return nil, fmt.Errorf("could not find keyring for address")
	}

	var msgFlags proton.MessageFlag
	if flags == 0 {
		msgFlags = proton.MessageFlagReceived
	} else {
		msgFlags = flags
	}

	str, err := c.ImportMessages(
		ctx,
		addrKRs[addrID],
		runtime.NumCPU(),
		runtime.NumCPU(),
		xslices.Map(messages, func(message []byte) proton.ImportReq {
			return proton.ImportReq{
				Metadata: proton.ImportMetadata{
					AddressID: addrID,
					LabelIDs:  []string{labelID},
					Flags:     msgFlags,
				},
				Message: message,
			}
		})...,
	)
	if err != nil {
		return nil, err
	}

	res, err := stream.Collect(ctx, str)
	if err != nil {
		return nil, err
	}

	return xslices.Map(res, func(res proton.ImportRes) string {
		return res.MessageID
	}), nil
}
