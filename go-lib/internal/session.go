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

package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"net/http/cookiejar"
)

type SessionLoginState int

var ErrInvalidLoginState = errors.New("invalid login state")

const (
	SessionLoginStateLoggedOut SessionLoginState = iota
	SessionLoginStateAwaitingTOTP
	SessionLoginStatAwaitingMailboxPassword
	SessionLoginStateAwaitingHV
	SessionLoginStateLoggedIn
)

type Session struct {
	group           *async.Group
	m               *proton.Manager
	client          *proton.Client
	loginState      SessionLoginState
	passwordMode    proton.PasswordMode
	mailboxPassword string
}

func NewSession(ctx context.Context, apiURL string) *Session {
	panicHandler := &async.NoopPanicHandler{}
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(fmt.Errorf("unexpected error:%w", err))
	}

	return &Session{
		group: async.NewGroup(ctx, panicHandler),
		m: proton.New(
			proton.WithHostURL(apiURL),
			proton.WithCookieJar(jar),
			proton.WithAppVersion("export"),
			proton.WithLogger(logrus.StandardLogger()),
			proton.WithPanicHandler(panicHandler),
		),
		client: nil,
	}
}

func (s *Session) Close(ctx context.Context) {
	s.group.CancelAndWait()
	if err := s.Logout(ctx); err != nil {
		logrus.WithError(err).Error("Failed to logout")
	}
	if s.client != nil {
		s.client.Close()
	}
	s.m.Close()
	s.mailboxPassword = ""
}

func (s *Session) Login(ctx context.Context, email, password string) error {
	if s.loginState != SessionLoginStateLoggedOut {
		return ErrInvalidLoginState
	}

	client, auth, err := s.m.NewClientWithLogin(ctx, email, []byte(password))
	if err != nil {
		return err
	}

	s.client = client
	s.mailboxPassword = password

	if auth.TwoFA.Enabled&proton.HasTOTP != 0 {
		s.loginState = SessionLoginStateAwaitingTOTP
		return nil
	}

	s.passwordMode = auth.PasswordMode

	if auth.PasswordMode == proton.TwoPasswordMode {
		s.loginState = SessionLoginStatAwaitingMailboxPassword
		return nil
	}

	s.loginState = SessionLoginStateLoggedIn

	return nil
}

func (s *Session) Logout(ctx context.Context) error {
	if s.loginState == SessionLoginStateLoggedOut {
		return ErrInvalidLoginState
	}

	return s.client.AuthDelete(ctx)
}

func (s *Session) SubmitTOTP(ctx context.Context, totp string) error {
	if s.loginState != SessionLoginStateAwaitingTOTP {
		return ErrInvalidLoginState
	}

	if err := s.client.Auth2FA(ctx, proton.Auth2FAReq{TwoFactorCode: totp}); err != nil {
		return err
	}

	if s.passwordMode == proton.TwoPasswordMode {
		s.loginState = SessionLoginStatAwaitingMailboxPassword
	} else {
		s.loginState = SessionLoginStateLoggedIn
	}

	return nil
}

func (s *Session) SubmitMailboxPassword(password string) error {
	if s.loginState != SessionLoginStatAwaitingMailboxPassword {
		return ErrInvalidLoginState
	}

	s.mailboxPassword = password
	s.loginState = SessionLoginStateLoggedIn
	return nil
}

func (s *Session) LoginState() SessionLoginState {
	return s.loginState
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
