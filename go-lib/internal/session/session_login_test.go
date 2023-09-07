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

package session

import (
	"context"
	"testing"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const TestUserEmail = "foo@bar.com"

var TestUserPassword = []byte("12345")

func TestSessionLogin_SinglePasswordMode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword)).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())
	require.Equal(t, TestUserPassword, session.mailboxPassword)
}

func TestSessionLogin_LoginAfterLoginIsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())

	err := session.Login(ctx, TestUserEmail, TestUserPassword)
	require.ErrorIs(t, err, ErrInvalidLoginState)
}

func TestSessionLogin_TwoPasswordMode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{
		PasswordMode: proton.TwoPasswordMode,
	}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingMailboxPassword, session.LoginState())

	mailboxPassword := []byte("some password")

	require.NoError(t, session.SubmitMailboxPassword(mailboxPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())

	require.Equal(t, mailboxPassword, session.mailboxPassword)
}

func TestSessionLogin_SinglePasswordModeWithTOTP(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{
		TwoFA: proton.TwoFAInfo{
			Enabled: proton.HasTOTP,
		},
	}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()

	const totpCode = "01245"

	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().Close()
	client.EXPECT().Auth2FA(gomock.Any(), gomock.Eq(proton.Auth2FAReq{
		TwoFactorCode: totpCode,
	})).Return(nil)

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingTOTP, session.LoginState())

	require.NoError(t, session.SubmitTOTP(ctx, totpCode))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())
}

func TestSessionLogin_TwoPasswordModeWithTOTP(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{
		TwoFA: proton.TwoFAInfo{
			Enabled: proton.HasTOTP,
		},
		PasswordMode: proton.TwoPasswordMode,
	}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()

	const totpCode = "01245"

	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().Close()
	client.EXPECT().Auth2FA(gomock.Any(), gomock.Eq(proton.Auth2FAReq{
		TwoFactorCode: totpCode,
	})).Return(nil)

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingTOTP, session.LoginState())

	require.NoError(t, session.SubmitTOTP(ctx, totpCode))
	require.Equal(t, LoginStateAwaitingMailboxPassword, session.LoginState())

	mailboxPassword := []byte("some password")

	require.NoError(t, session.SubmitMailboxPassword(mailboxPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())

	require.Equal(t, mailboxPassword, session.mailboxPassword)
}

func TestSessionLogin_Logout(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil).Times(1)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())

	require.NoError(t, session.Logout(ctx))
	require.Equal(t, LoginStateLoggedOut, session.LoginState())
}

func TestSessionLogin_CatchHVError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq([]byte(TestUserPassword))).Return(
		nil,
		clientAuth,
		&proton.APIError{
			Status: 421,
			Code:   9001,
		},
	)
	clientBuilder.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingHV, session.LoginState())
}
