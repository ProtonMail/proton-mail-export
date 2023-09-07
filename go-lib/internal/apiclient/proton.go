// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
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

package apiclient

import (
	"context"
	"errors"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type ProtonAPIClientBuilder struct {
	manager *proton.Manager
}

func NewProtonAPIClientBuilder(apiURL string, panicHandler async.PanicHandler) *ProtonAPIClientBuilder {
	return &ProtonAPIClientBuilder{
		manager: proton.New(
			proton.WithHostURL(apiURL),
			proton.WithAppVersion("export"),
			proton.WithLogger(logrus.StandardLogger()),
			proton.WithPanicHandler(panicHandler),
		),
	}
}

func (p *ProtonAPIClientBuilder) NewClient(ctx context.Context, username string, password []byte) (Client, proton.Auth, error) {
	return p.manager.NewClientWithLogin(ctx, username, password)
}

func (p *ProtonAPIClientBuilder) Close() {
	p.manager.Close()
}

func IsHVRequestedError(err error) bool {
	if err == nil {
		return false
	}

	var protonErr *proton.APIError

	if !errors.As(err, &protonErr) {
		return false
	}

	return protonErr.Code == 9001
}
