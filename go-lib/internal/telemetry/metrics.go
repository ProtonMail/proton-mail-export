package telemetry

import (
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

const (
	tasksMeasurementGroup = "mail.any.export_tool_tasks"
	taskStartedEvent      = "taskStarted"
	taskFinishedEvent     = "taskFinished"

	sessionMeasurementGroup = "mail.any.export_tool_session"
	sessionStartEvent       = "sessionStart"

	taskTypeExport  = "export"
	taskTypeRestore = "restore"
)

func generateTaskStartMetric(taskType string, useDefaultDirectory bool, userPlanType string) proton.SendStatsReq {
	return proton.SendStatsReq{
		MeasurementGroup: tasksMeasurementGroup,
		Event:            taskStartedEvent,
		Values:           map[string]any{},
		Dimensions: map[string]any{
			"type":                taskType,
			"useDefaultDirectory": mapBoolStr(useDefaultDirectory),
			"userPlanType":        userPlanType,
		},
	}
}

func (s *Service) SendExportStart() {
	s.withTelemetry(func() {
		metric := generateTaskStartMetric(taskTypeExport, s.data.useDefaultExportPath, s.data.userPlan)
		if err := s.client.SendDataEvent(s.ctx, metric); err != nil {
			logrus.WithError(err).Info("Failed to send Export start telemetry metric")
		}
	})
}

func (s *Service) SendRestoreStart() {
	s.withTelemetry(func() {
		metric := generateTaskStartMetric(taskTypeRestore, s.data.useDefaultExportPath, s.data.userPlan)
		if err := s.client.SendDataEvent(s.ctx, metric); err != nil {
			logrus.WithError(err).Info("Failed to send Restore start telemetry metric")
		}
	})
}

func generateTaskFinishedMetric(
	taskType string,
	useDefaultDirectory bool,
	userPlanType string,
	taskCancelledByUser bool,
	taskCancelledByError bool,
	durationSeconds int,
	totalMessages int,
	failedMessages int,
	successfullMessages int,
) proton.SendStatsReq {
	return proton.SendStatsReq{
		MeasurementGroup: tasksMeasurementGroup,
		Event:            taskFinishedEvent,
		Values: map[string]any{
			"durationSeconds":    durationSeconds,
			"totalMessages":      totalMessages,
			"failedMessages":     failedMessages,
			"successfulMessages": successfullMessages,
		},
		Dimensions: map[string]any{
			"type":                taskType,
			"useDefaultDirectory": mapBoolStr(useDefaultDirectory),
			"userPlanType":        userPlanType,
			"cancelledByUser":     mapBoolStr(taskCancelledByUser),
			"cancelledByError":    mapBoolStr(taskCancelledByError),
		},
	}
}

func (s *Service) SendExportFinished(taskCancelledByUser, taskCancelledByError bool,
	durationSeconds, totalMessageCount, failedMessageCount, successfulMessageCount int) {
	s.withTelemetry(func() {
		metric := generateTaskFinishedMetric(
			taskTypeExport,
			s.data.useDefaultExportPath,
			s.data.userPlan,
			taskCancelledByUser,
			taskCancelledByError,
			durationSeconds,
			totalMessageCount,
			failedMessageCount,
			successfulMessageCount,
		)
		if err := s.client.SendDataEvent(s.ctx, metric); err != nil {
			logrus.WithError(err).Info("Failed to send Export finished telemetry metric")
		}
	})
}

func (s *Service) SendRestoreFinished(taskCancelledByUser, taskCancelledByError bool,
	durationSeconds, totalMessageCount, failedMessageCount, successfulMessageCount int) {
	s.withTelemetry(func() {
		metric := generateTaskFinishedMetric(
			taskTypeRestore,
			s.data.useDefaultExportPath,
			s.data.userPlan,
			taskCancelledByUser,
			taskCancelledByError,
			durationSeconds,
			totalMessageCount,
			failedMessageCount,
			successfulMessageCount,
		)
		if err := s.client.SendDataEvent(s.ctx, metric); err != nil {
			logrus.WithError(err).Info("Failed to send Restore finished telemetry metric")
		}
	})
}

func GenerateProcessStartMetric(
	etOperation,
	etDir,
	etUserPassword,
	etUserMailboxPassword,
	etTotpCode,
	etUserEmail bool,
) proton.SendStatsReq {
	return proton.SendStatsReq{
		MeasurementGroup: sessionMeasurementGroup,
		Event:            sessionStartEvent,
		Values:           map[string]any{},
		Dimensions: map[string]any{
			"et_operation":             mapBoolStr(etOperation),
			"et_dir":                   mapBoolStr(etDir),
			"et_user_password":         mapBoolStr(etUserPassword),
			"et_user_mailbox_password": mapBoolStr(etUserMailboxPassword),
			"et_totp_code":             mapBoolStr(etTotpCode),
			"et_user_email":            mapBoolStr(etUserEmail),
		},
	}
}
