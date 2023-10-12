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
	"fmt"
	"strings"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/reporter"
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
	prevLoginState  LoginState
	passwordMode    proton.PasswordMode
	mailboxPassword []byte
	callbacks       Callbacks
	reporter        reporter.Reporter
	hvDetails       *proton.APIHVDetails
	user            proton.User
}

func NewSession(
	builder apiclient.Builder,
	callbacks Callbacks,
	panicHandler async.PanicHandler,
	reporter reporter.Reporter,
) *Session {
	return &Session{
		panicHandler:   panicHandler,
		client:         nil,
		clientBuilder:  builder,
		callbacks:      callbacks,
		reporter:       reporter,
		loginState:     LoginStateLoggedOut,
		prevLoginState: LoginStateLoggedOut,
	}
}

func (s *Session) Close(ctx context.Context) {
	defer async.HandlePanic(s.panicHandler)

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
	if email == "crash@bandicoot" {
		panic("Crash Time")
	}

	if s.loginState != LoginStateLoggedOut && s.loginState != LoginStateAwaitingHV {
		return ErrInvalidLoginState
	}

	logrus.Debugf("Performing login for user %v", email)

	client, auth, err := s.clientBuilder.NewClient(ctx, email, password, s.hvDetails)
	if err != nil {
		if s.checkHVRequest(err) {
			return nil
		}

		logrus.WithError(err).Error("Failed to login")
		return err
	}

	client = apiclient.NewAutoRetryClient(client, &apiclient.SleepRetryStrategyBuilder{})
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

	// We can now get the user.
	if err := s.loadUser(ctx); err != nil {
		if s.checkHVRequest(err) {
			return nil
		}

		logrus.WithError(err).Error("Failed to get user")
		return fmt.Errorf("failed to load user: %w", err)
	}

	return nil
}

func (s *Session) Logout(ctx context.Context) error {
	if s.loginState == LoginStateLoggedOut {
		return ErrInvalidLoginState
	}

	logrus.Debugf("Logging out")

	if err := s.client.AuthDelete(ctx); err != nil {
		logrus.WithError(err).Error("Failed to logout")
		return err
	}

	s.loginState = LoginStateLoggedOut
	s.prevLoginState = LoginStateLoggedOut
	s.setMailboxPassword(nil)

	return nil
}

func (s *Session) SubmitTOTP(ctx context.Context, totp string) error {
	if s.loginState != LoginStateAwaitingTOTP {
		return ErrInvalidLoginState
	}

	logrus.Debugf("Submitting TOTP code")

	if err := s.client.Auth2FA(ctx, proton.Auth2FAReq{TwoFactorCode: totp}); err != nil {
		logrus.WithError(err).Error("Failed to Submit totp")
		return err
	}

	if s.passwordMode == proton.TwoPasswordMode {
		s.loginState = LoginStateAwaitingMailboxPassword
	} else {
		s.loginState = LoginStateLoggedIn
	}

	// We can now get the user.
	if err := s.loadUser(ctx); err != nil {
		if s.checkHVRequest(err) {
			return nil
		}

		logrus.WithError(err).Error("Failed to get user")
		return fmt.Errorf("failed to load user: %w", err)
	}

	return nil
}

func (s *Session) SubmitMailboxPassword(password []byte) error {
	if s.loginState != LoginStateAwaitingMailboxPassword {
		return ErrInvalidLoginState
	}

	logrus.Debugf("Submitting Mailbox Password")

	s.setMailboxPassword(password)
	s.loginState = LoginStateLoggedIn
	return nil
}

func (s *Session) GetHVSolveURL() (string, error) {
	if s.loginState != LoginStateAwaitingHV || s.hvDetails == nil {
		return "", ErrInvalidLoginState
	}

	return fmt.Sprintf("https://verify.proton.me/?methods=%v&token=%v",
		strings.Join(s.hvDetails.Methods, ","),
		s.hvDetails.Token), nil
}

func (s *Session) MarkHVSolved(ctx context.Context) error {
	if s.loginState != LoginStateAwaitingHV || s.hvDetails == nil {
		return ErrInvalidLoginState
	}

	s.loginState = s.prevLoginState
	s.prevLoginState = LoginStateLoggedOut

	if s.loginState == LoginStateLoggedIn {
		return s.loadUser(ctx)
	}

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

func (s *Session) GetReporter() reporter.Reporter {
	return s.reporter
}

func (s *Session) GetUser() *proton.User {
	return &s.user
}

func (s *Session) checkHVRequest(err error) bool {
	if details := apiclient.GetHVData(err); details != nil {
		s.prevLoginState = s.loginState
		s.loginState = LoginStateAwaitingHV
		s.hvDetails = details
		return true
	}

	return false
}

func (s *Session) setMailboxPassword(p []byte) {
	if s.mailboxPassword != nil {
		zeroSlice(s.mailboxPassword)
	}

	s.mailboxPassword = p
}

func (s *Session) loadUser(ctx context.Context) error {
	logrus.Debug("Getting user info")
	u, err := s.client.GetUserWithHV(ctx, s.hvDetails)
	if err != nil {
		if s.checkHVRequest(err) {
			return nil
		}

		return err
	}

	s.user = u
	return nil
}

func zeroSlice(s []byte) {
	for i := 0; i < len(s); i++ {
		s[i] = 0
	}
}
