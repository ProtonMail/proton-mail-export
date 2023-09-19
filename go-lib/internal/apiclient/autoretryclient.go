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
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type AutoRetryClient struct {
	client               Client
	retryStrategyBuilder RetryStrategyBuilder
}

func NewAutoRetryClient(client Client, builder RetryStrategyBuilder) *AutoRetryClient {
	return &AutoRetryClient{client: client, retryStrategyBuilder: builder}
}

func (arc *AutoRetryClient) Auth2FA(ctx context.Context, req proton.Auth2FAReq) error {
	return arc.repeatRequest(ctx, func(ctx context.Context, client Client) error {
		return client.Auth2FA(ctx, req)
	})
}

func (arc *AutoRetryClient) AuthDelete(ctx context.Context) error {
	return arc.repeatRequest(ctx, func(ctx context.Context, client Client) error {
		return client.AuthDelete(ctx)
	})
}

func (arc *AutoRetryClient) GetUser(ctx context.Context) (proton.User, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) (proton.User, error) {
		return client.GetUser(ctx)
	})
}

func (arc *AutoRetryClient) GetSalts(ctx context.Context) (proton.Salts, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) (proton.Salts, error) {
		return client.GetSalts(ctx)
	})
}

func (arc *AutoRetryClient) Close() {
	arc.client.Close()
}

func (arc *AutoRetryClient) GetLabels(ctx context.Context, labelTypes ...proton.LabelType) ([]proton.Label, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) ([]proton.Label, error) {
		return client.GetLabels(ctx, labelTypes...)
	})
}

func (arc *AutoRetryClient) GetAddresses(ctx context.Context) ([]proton.Address, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) ([]proton.Address, error) {
		return client.GetAddresses(ctx)
	})
}

func (arc *AutoRetryClient) GetGroupedMessageCount(ctx context.Context) ([]proton.MessageGroupCount, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) ([]proton.MessageGroupCount, error) {
		return client.GetGroupedMessageCount(ctx)
	})
}

func (arc *AutoRetryClient) GetMessage(ctx context.Context, messageID string) (proton.Message, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) (proton.Message, error) {
		return client.GetMessage(ctx, messageID)
	})
}

func (arc *AutoRetryClient) GetMessageMetadataPage(ctx context.Context, page, pageSize int, filter proton.MessageFilter) ([]proton.MessageMetadata, error) {
	return repeatRequestTyped(ctx, arc, func(ctx context.Context, client Client) ([]proton.MessageMetadata, error) {
		return client.GetMessageMetadataPage(ctx, page, pageSize, filter)
	})
}

func (arc *AutoRetryClient) GetAttachmentInto(ctx context.Context, attachmentID string, reader io.ReaderFrom) error {
	return arc.repeatRequest(ctx, func(ctx context.Context, client Client) error {
		return client.GetAttachmentInto(ctx, attachmentID, reader)
	})
}

func (arc *AutoRetryClient) repeatRequest(ctx context.Context, req func(ctx context.Context, client Client) error) error {
	retryStrategy := arc.retryStrategyBuilder.NewRetryStrategy()
	for {
		err := req(ctx, arc.client)
		if err != nil {
			if !isRetrieableError(err) {
				return err
			}

			retryStrategy.HandleRetry(ctx)
			continue
		}

		return nil
	}
}

func repeatRequestTyped[T any](ctx context.Context, arc *AutoRetryClient, req func(ctx context.Context, client Client) (T, error)) (T, error) {
	var result T
	var err error
	err = arc.repeatRequest(ctx, func(ctx context.Context, client Client) error {
		result, err = req(ctx, client)

		return err
	})

	return result, err
}

func isRetrieableError(err error) bool {
	if netErr := new(proton.NetError); errors.As(err, &netErr) {
		// Context cancelled is wrapped in the proton network error. Check here to make sure.
		if errors.Is(netErr.Cause, context.Canceled) {
			return false
		}

		logrus.WithError(err).Debug("Retry due to network error")
		return true
	}

	// Catch all for uncategorized net errors that may slip through.
	if netErr := new(net.OpError); errors.As(err, &netErr) {
		logrus.WithError(err).Debug("Retry due to uncategorized network error")
		return true
	}

	// If the error is an unexpected EOF, return error to retry later.
	if errors.Is(err, io.ErrUnexpectedEOF) {
		logrus.WithError(err).Debug("Retry due to unexpected EOF")
		return true
	}

	// If the error is a server-side issue, return error to retry later.
	if apiErr := new(proton.APIError); errors.As(err, &apiErr) {
		if apiErr.Status == 429 || apiErr.Status >= 500 {
			logrus.WithError(err).Debug("Retry due to unexpected 429/5xx")
			return true
		}
	}

	return false
}

type RetryStrategyBuilder interface {
	// NewRetryStrategy can be called from any go-routine.
	NewRetryStrategy() RetryStrategy
}

// RetryStrategy is meant to be used in the scope of on goroutine for the lifetime of one specific request.
type RetryStrategy interface {
	HandleRetry(ctx context.Context)
}

type SleepRetryStrategyBuilder struct{}

func (r SleepRetryStrategyBuilder) NewRetryStrategy() RetryStrategy {
	return &SleepRetryStrategy{index: 0}
}

type SleepRetryStrategy struct {
	index int
}

func (s *SleepRetryStrategy) HandleRetry(ctx context.Context) {
	sleepCtx(ctx, s.nextWaitTime())
}

func (s *SleepRetryStrategy) nextWaitTime() time.Duration {
	last := len(expWaitTimes) - 1
	if s.index >= last {
		s.index = last
	}

	nextWaitTime := expWaitTimes[s.index] + jitter(10)

	s.index++

	return nextWaitTime
}

//nolint:gochecknoglobals
var expWaitTimes = []time.Duration{
	20 * time.Second,
	40 * time.Second,
	80 * time.Second,
	160 * time.Second,
	300 * time.Second,
	600 * time.Second,
}

func jitter(max int) time.Duration {
	return time.Duration(rand.Intn(max)) * time.Second //nolint:gosec
}

func sleepCtx(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
