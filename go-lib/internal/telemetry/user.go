package telemetry

import (
	"errors"
	"strings"

	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

const (
	planUnknown    = "unknown"
	planOther      = "other"
	planBusiness   = "business"
	planIndividual = "individual"
	planGroup      = "group"
	planFree       = "free"
)

func mapUserPlan(planName string) string {
	if planName == "" {
		return planUnknown
	}
	switch strings.TrimSpace(strings.ToLower(planName)) {
	case "free":
		return planFree
	case "mail2022":
		return planIndividual
	case "bundle2022":
		return planIndividual
	case "family2022":
		return planGroup
	case "visionary2022":
		return planGroup
	case "mailpro2022":
		return planBusiness
	case "planbiz2024":
		return planBusiness
	case "bundlepro2022":
		return planBusiness
	case "bundlepro2024":
		return planBusiness
	case "duo2024":
		return planGroup
	default:
		return planOther
	}
}

func (s *Service) setUserPlan() {
	userData, err := s.client.GetOrganizationData(s.ctx)
	if err != nil {
		var protonErr *proton.APIError
		// 2501 - corresponds to field does not exist in DB -> free user
		if errors.As(err, &protonErr) && protonErr.Code == 2501 {
			s.data.userPlan = planFree
			return
		}

		logrus.WithError(err).Info("Failed to get user organization data")
		return
	}

	s.data.userPlan = mapUserPlan(userData.Organization.PlanName)
}
