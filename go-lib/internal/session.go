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
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type SessionLoginState int

var ErrInvalidLoginState = errors.New("invalid login state")

const (
	SessionLoginStateLoggedOut SessionLoginState = iota
	SessionLoginStateAwaitingTOTP
	SessionLoginStateAwaitingMailboxPassword
	SessionLoginStateAwaitingHV
	SessionLoginStateLoggedIn
)

type Session struct {
	group           *async.Group
	clientBuilder   APIClientBuilder
	client          APIClient
	loginState      SessionLoginState
	passwordMode    proton.PasswordMode
	mailboxPassword []byte
}

func NewSession(ctx context.Context, builder APIClientBuilder) *Session {
	panicHandler := &async.NoopPanicHandler{}
	return &Session{
		group:         async.NewGroup(ctx, panicHandler),
		client:        nil,
		clientBuilder: builder,
	}
}

func (s *Session) Close(ctx context.Context) {
	s.group.CancelAndWait()
	if s.client != nil {
		if err := s.Logout(ctx); err != nil {
			logrus.WithError(err).Error("Failed to logout")
		}

		s.client.Close()
	}
	s.clientBuilder.Close()
	s.setMailboxPassword(nil)
}

func (s *Session) Login(ctx context.Context, email string, password []byte) error {
	if s.loginState != SessionLoginStateLoggedOut && s.loginState != SessionLoginStateAwaitingHV {
		return ErrInvalidLoginState
	}

	client, auth, err := s.clientBuilder.NewClient(ctx, email, password)
	if err != nil {
		if isHVRequestedError(err) {
			s.loginState = SessionLoginStateAwaitingHV
			return nil
		}

		return err
	}

	s.client = client
	s.setMailboxPassword(password)
	s.passwordMode = auth.PasswordMode

	if auth.TwoFA.Enabled&proton.HasTOTP != 0 {
		s.loginState = SessionLoginStateAwaitingTOTP
		return nil
	}

	if auth.PasswordMode == proton.TwoPasswordMode {
		s.loginState = SessionLoginStateAwaitingMailboxPassword
		return nil
	}

	s.loginState = SessionLoginStateLoggedIn

	return nil
}

func (s *Session) Logout(ctx context.Context) error {
	if s.loginState == SessionLoginStateLoggedOut {
		return ErrInvalidLoginState
	}

	if err := s.client.AuthDelete(ctx); err != nil {
		return err
	}

	s.loginState = SessionLoginStateLoggedOut
	s.setMailboxPassword(nil)

	return nil
}

func (s *Session) SubmitTOTP(ctx context.Context, totp string) error {
	if s.loginState != SessionLoginStateAwaitingTOTP {
		return ErrInvalidLoginState
	}

	if err := s.client.Auth2FA(ctx, proton.Auth2FAReq{TwoFactorCode: totp}); err != nil {
		return err
	}

	if s.passwordMode == proton.TwoPasswordMode {
		s.loginState = SessionLoginStateAwaitingMailboxPassword
	} else {
		s.loginState = SessionLoginStateLoggedIn
	}

	return nil
}

func (s *Session) SubmitMailboxPassword(password []byte) error {
	if s.loginState != SessionLoginStateAwaitingMailboxPassword {
		return ErrInvalidLoginState
	}

	s.setMailboxPassword(password)
	s.loginState = SessionLoginStateLoggedIn
	return nil
}

func (s *Session) LoginState() SessionLoginState {
	return s.loginState
}

func (s *Session) setMailboxPassword(p []byte) {
	if s.mailboxPassword != nil {
		zeroSlice(s.mailboxPassword)
	}

	s.mailboxPassword = p
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func zeroSlice(s []byte) {
	for i := 0; i < len(s); i++ {
		s[i] = 0
	}
}
