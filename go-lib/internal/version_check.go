package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
)

type VersionInfo struct {
	Major  uint64 `json:"version_major"`
	Minor  uint64 `json:"version_minor"`
	Patch  uint64 `json:"version_patch"`
	String string `json:"version"`
}

func HasNewVersion() bool {
	const versionURL = "https://proton.me/download/export-tool/version.json"

	client := http.Client{
		Timeout: 1 * time.Minute,
	}

	resp, err := client.Get(versionURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to get version info")
		return false
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return false
	}

	var version VersionInfo

	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		logrus.WithError(err).Error("Failed to parser version info")
		return false
	}

	logrus.Debugf("Remote version: %v", version.String)

	versionRemote := semver.New(version.Major, version.Minor, version.Patch, "", "")
	versionLocal := semver.New(ETVersionMajor, ETVersionMinor, ETVersionPatch, "", "")

	return versionLocal.Compare(versionRemote) < 0
}
