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
	"io"

	"github.com/ProtonMail/go-proton-api"
)

type Builder interface {
	NewClient(ctx context.Context, username string, password []byte) (Client, proton.Auth, error)
	Close()
}

type Client interface {
	Auth2FA(ctx context.Context, req proton.Auth2FAReq) error
	AuthDelete(ctx context.Context) error
	GetUser(ctx context.Context) (proton.User, error)
	GetSalts(ctx context.Context) (proton.Salts, error)
	Close()

	GetLabels(ctx context.Context, labelTypes ...proton.LabelType) ([]proton.Label, error)
	GetAddresses(ctx context.Context) ([]proton.Address, error)

	GetGroupedMessageCount(ctx context.Context) ([]proton.MessageGroupCount, error)
	GetMessage(ctx context.Context, messageID string) (proton.Message, error)
	GetMessageMetadataPage(ctx context.Context, page, pageSize int, filter proton.MessageFilter) ([]proton.MessageMetadata, error)
	GetAttachmentInto(ctx context.Context, attachmentID string, reader io.ReaderFrom) error
}
