package telemetry

func mapBoolStr(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

func (s *Service) SetUsingDefaultExportPath(usingDefaultExportPath bool) {
	s.data.useDefaultExportPath = usingDefaultExportPath
}
