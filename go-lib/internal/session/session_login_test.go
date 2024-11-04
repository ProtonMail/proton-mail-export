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
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const TestUserEmail = "foo@bar.com"

var TestUserPassword = []byte("12345")

func TestSessionLogin_SinglePasswordMode(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)
	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
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

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)
	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
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

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)
	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingMailboxPassword, session.LoginState())

	mailboxPassword := []byte("some password")

	require.NoError(t, session.SubmitMailboxPassword(&AlwaysValidMailboxPasswordValidator{}, mailboxPassword))
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

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)

	const totpCode = "01245"

	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)

	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()
	client.EXPECT().Auth2FA(gomock.Any(), gomock.Eq(proton.Auth2FAReq{
		TwoFactorCode: totpCode,
	})).Return(nil)

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
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

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	const totpCode = "01245"

	client.EXPECT().AuthDelete(gomock.Any()).Return(nil)
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)
	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()
	client.EXPECT().Auth2FA(gomock.Any(), gomock.Eq(proton.Auth2FAReq{
		TwoFactorCode: totpCode,
	})).Return(nil)

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingTOTP, session.LoginState())

	require.NoError(t, session.SubmitTOTP(ctx, totpCode))
	require.Equal(t, LoginStateAwaitingMailboxPassword, session.LoginState())

	mailboxPassword := []byte("some password")

	require.NoError(t, session.SubmitMailboxPassword(&AlwaysValidMailboxPasswordValidator{}, mailboxPassword))
	require.Equal(t, LoginStateLoggedIn, session.LoginState())

	require.Equal(t, mailboxPassword, session.mailboxPassword)
}

func TestSessionLogin_Logout(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	client := apiclient.NewMockClient(mockCtrl)
	clientBuilder := apiclient.NewMockBuilder(mockCtrl)
	clientAuth := proton.Auth{}

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		client,
		clientAuth,
		nil,
	)
	clientBuilder.EXPECT().Close()
	client.EXPECT().AuthDelete(gomock.Any()).Return(nil).Times(1)
	client.EXPECT().GetUserWithHV(gomock.Any(), gomock.Any()).Return(proton.User{}, nil)
	client.EXPECT().GetSalts(gomock.Any()).Return(proton.Salts{}, nil)
	client.EXPECT().GetUserSettings(gomock.Any()).Return(proton.UserSettings{}, nil)
	client.EXPECT().GetOrganizationData(gomock.Any()).Return(proton.OrganizationResponse{}, nil)
	client.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
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

	clientBuilder.EXPECT().NewClient(gomock.Any(), gomock.Eq(TestUserEmail), gomock.Eq(TestUserPassword), gomock.Any()).Return(
		nil,
		clientAuth,
		&proton.APIError{
			Status:  422,
			Code:    9001,
			Details: []byte(`{"Methods": ["captcha"],"Token":"token"}`),
		},
	)
	clientBuilder.EXPECT().Close()

	ctx := context.Background()
	session := NewSession(clientBuilder, nil, &async.NoopPanicHandler{}, &reporter.NullReporter{}, false)
	defer session.Close(ctx)

	require.NoError(t, session.Login(ctx, TestUserEmail, TestUserPassword))
	require.Equal(t, LoginStateAwaitingHV, session.LoginState())
}

type AlwaysValidMailboxPasswordValidator struct{}

func (a AlwaysValidMailboxPasswordValidator) IsValid(_ []byte) bool {
	return true
}
