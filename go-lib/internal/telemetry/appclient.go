package telemetry

import "runtime"

type appClient string

const (
	windowsAppClient appClient = "windows-export"
	linuxAppClient   appClient = "linux-export"
	darwinAppClient  appClient = "darwin-export"
	unknownAppClient appClient = "unknown-export"
)

// getAppClient returns the OS-specific app client name.
func getAppClient() appClient {
	switch runtime.GOOS {
	case "darwin":
		return darwinAppClient
	case "linux":
		return linuxAppClient
	case "windows":
		return windowsAppClient
	default:
		return unknownAppClient
	}
}
