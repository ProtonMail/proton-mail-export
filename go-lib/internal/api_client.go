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

package internal

import (
	"context"
	"github.com/ProtonMail/go-proton-api"
)

type APIClientBuilder interface {
	NewClient(ctx context.Context, username string, password []byte) (APIClient, proton.Auth, error)
	Close()
}

type APIClient interface {
	Auth2FA(ctx context.Context, req proton.Auth2FAReq) error
	AuthDelete(ctx context.Context) error
	Close()
}
