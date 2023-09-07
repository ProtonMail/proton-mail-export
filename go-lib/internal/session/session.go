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

package session

import (
	"context"
	"errors"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type LoginState int

var ErrInvalidLoginState = errors.New("invalid login state")

const (
	LoginStateLoggedOut LoginState = iota
	LoginStateAwaitingTOTP
	LoginStateAwaitingMailboxPassword
	LoginStateAwaitingHV
	LoginStateLoggedIn
)

type Session struct {
	panicHandler    async.PanicHandler
	clientBuilder   apiclient.Builder
	client          apiclient.Client
	loginState      LoginState
	passwordMode    proton.PasswordMode
	mailboxPassword []byte
}

func NewSession(builder apiclient.Builder) *Session {
	return &Session{
		panicHandler:  &async.NoopPanicHandler{},
		client:        nil,
		clientBuilder: builder,
	}
}

func (s *Session) Close(ctx context.Context) {
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
	if s.loginState != LoginStateLoggedOut && s.loginState != LoginStateAwaitingHV {
		return ErrInvalidLoginState
	}

	// GODT-2900: Handle network errors/loss.
	client, auth, err := s.clientBuilder.NewClient(ctx, email, password)
	if err != nil {
		if apiclient.IsHVRequestedError(err) {
			s.loginState = LoginStateAwaitingHV
			return nil
		}

		return err
	}

	s.client = client
	s.setMailboxPassword(password)
	s.passwordMode = auth.PasswordMode

	if auth.TwoFA.Enabled&proton.HasTOTP != 0 {
		s.loginState = LoginStateAwaitingTOTP
		return nil
	}

	if auth.PasswordMode == proton.TwoPasswordMode {
		s.loginState = LoginStateAwaitingMailboxPassword
		return nil
	}

	s.loginState = LoginStateLoggedIn

	return nil
}

func (s *Session) Logout(ctx context.Context) error {
	if s.loginState == LoginStateLoggedOut {
		return ErrInvalidLoginState
	}

	// GODT-2900: Handle network errors/loss.
	if err := s.client.AuthDelete(ctx); err != nil {
		return err
	}

	s.loginState = LoginStateLoggedOut
	s.setMailboxPassword(nil)

	return nil
}

func (s *Session) SubmitTOTP(ctx context.Context, totp string) error {
	if s.loginState != LoginStateAwaitingTOTP {
		return ErrInvalidLoginState
	}

	// GODT-2900: Handle network errors/loss.
	if err := s.client.Auth2FA(ctx, proton.Auth2FAReq{TwoFactorCode: totp}); err != nil {
		return err
	}

	if s.passwordMode == proton.TwoPasswordMode {
		s.loginState = LoginStateAwaitingMailboxPassword
	} else {
		s.loginState = LoginStateLoggedIn
	}

	return nil
}

func (s *Session) SubmitMailboxPassword(password []byte) error {
	if s.loginState != LoginStateAwaitingMailboxPassword {
		return ErrInvalidLoginState
	}

	s.setMailboxPassword(password)
	s.loginState = LoginStateLoggedIn
	return nil
}

func (s *Session) LoginState() LoginState {
	return s.loginState
}

func (s *Session) GetClient() apiclient.Client {
	return s.client
}

func (s *Session) GetMailboxPassword() []byte {
	return s.mailboxPassword
}

func (s *Session) GetPanicHandler() async.PanicHandler {
	return s.panicHandler
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
