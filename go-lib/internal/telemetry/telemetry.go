package telemetry

import (
	"context"

	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

type userClient interface {
	GetUserSettings(ctx context.Context) (proton.UserSettings, error)
	SendDataEvent(ctx context.Context, req proton.SendStatsReq) error
	GetOrganizationData(ctx context.Context) (proton.OrganizationResponse, error)
}

type helperData struct {
	useDefaultExportPath bool
	userPlan             string
}

type Service struct {
	ctx    context.Context
	client userClient

	clientInitialized bool
	telemetryChecked  bool
	telemetryDisabled bool
	data              helperData
}

func NewService(disableTelemetry bool) *Service {
	return &Service{
		data: helperData{
			userPlan:             planUnknown,
			useDefaultExportPath: false,
		},
		telemetryDisabled: disableTelemetry,
	}
}

func (s *Service) checkTelemetrySettings() {
	if s.telemetryDisabled || s.telemetryChecked {
		return
	}

	userSettings, err := s.client.GetUserSettings(s.ctx)
	if err != nil {
		logrus.WithError(err).Info("Failed to get user telemetry settings")
		return
	}

	s.telemetryDisabled = userSettings.Telemetry == proton.SettingDisabled
	s.telemetryChecked = true
}

func (s *Service) Initialize(ctx context.Context, userClient userClient) {
	s.ctx = ctx
	s.client = userClient
	s.clientInitialized = true
	s.checkTelemetrySettings()
	s.setUserPlan()
}

func (s *Service) withTelemetry(fn func()) {
	if !s.clientInitialized {
		return
	}

	s.checkTelemetrySettings()
	if s.telemetryDisabled {
		return
	}

	fn()
}
