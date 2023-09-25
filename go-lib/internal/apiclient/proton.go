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
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type ProtonCallbacks interface {
	OnNetworkRestored()
	OnNetworkLost()
}

type ProtonAPIClientBuilder struct {
	manager  *proton.Manager
	callback ProtonCallbacks
}

func NewProtonAPIClientBuilder(apiURL string, panicHandler async.PanicHandler, callbacks ProtonCallbacks) (*ProtonAPIClientBuilder, error) {
	cookieJar, err := newCookieJar(apiURL)
	if err != nil {
		return nil, err
	}

	b := &ProtonAPIClientBuilder{
		manager: proton.New(
			proton.WithHostURL(apiURL),
			proton.WithAppVersion("Other"),
			proton.WithLogger(logrus.StandardLogger()),
			proton.WithPanicHandler(panicHandler),
			proton.WithCookieJar(cookieJar),
		),
		callback: callbacks,
	}

	b.manager.AddStatusObserver(func(status proton.Status) {
		if status == proton.StatusDown {
			logrus.Info("Connection to proton servers lost")
			if callbacks != nil {
				callbacks.OnNetworkLost()
			}
		} else {
			logrus.Info("Connection to proton servers restored")
			if callbacks != nil {
				callbacks.OnNetworkRestored()
			}
		}
	})

	return b, nil
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

func newCookieJar(apiURL string) (*cookiejar.Jar, error) {
	url, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	for name, value := range map[string]string{
		"hhn":  getProtectedHostname(),
		"tz":   getTimeZone(),
		"lng":  getSystemLang(),
		"arch": getHostArch(),
	} {
		jar.SetCookies(url, []*http.Cookie{{Name: name, Value: value, Secure: true}})
	}

	return jar, nil
}
