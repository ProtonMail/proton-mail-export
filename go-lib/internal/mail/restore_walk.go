package mail

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func (r *RestoreTask) walkBackupDir(fn func(emlPath string)) error {
	return filepath.Walk(r.backupDir, func(path string, info fs.FileInfo, err error) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		if err != nil {
			logrus.WithError(err).WithField("path", path).Warn("Cannot inspect path. Skipping.")
			return nil
		}

		if info.IsDir() && (path != r.backupDir) { // we skip any dir that is not the root dir.
			return filepath.SkipDir
		}

		emlPath := filepath.Join(r.backupDir, info.Name())
		if !strings.HasSuffix(emlPath, emlExtension) {
			return nil
		}

		if _, err := os.Stat(emlToMetadataFilename(emlPath)); errors.Is(err, os.ErrNotExist) {
			logrus.WithField("path", emlPath).Warn("Skipping EML file with no associated metadata file.")
			return nil
		}

		fn(emlPath)

		return nil
	})
}

func (r *RestoreTask) getTimestampedBackupDirs() ([]string, error) {
	var result []string
	err := filepath.Walk(r.backupDir, func(path string, info fs.FileInfo, err error) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		if err != nil {
			return nil //nolint:nilerr // we proceed in case of errors
		}

		name := info.Name()
		if (err != nil) || !info.IsDir() || (path == r.backupDir) {
			return nil //nolint:nilerr // ignore errors, files, and the walk's root folder
		}

		if mailFolderRegExp.MatchString(name) {
			result = append(result, filepath.Join(r.backupDir, name))
		}

		return fs.SkipDir // we do not recurse into dirs
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
