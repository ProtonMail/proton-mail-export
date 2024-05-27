package mail

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/sirupsen/logrus"
)

var mailFolderRegExp = regexp.MustCompile(`^mail_\d{8}_\d{6}$`)

type RestoreTask struct {
	ctx       context.Context
	ctxCancel func()
	backupDir string
	session   *session.Session
	log       *logrus.Entry
}

func NewRestoreTask(ctx context.Context, backupDir string, session *session.Session) *RestoreTask {
	ctx, cancel := context.WithCancel(ctx)

	return &RestoreTask{
		ctx:       ctx,
		ctxCancel: cancel,
		backupDir: backupDir,
		session:   session,
		log:       logrus.WithField("backup", "mail").WithField("userID", session.GetUser().ID),
	}
}

func (r *RestoreTask) Run(_ Reporter) error {
	defer r.log.Info("Finished")
	r.log.WithField("backupDir", r.backupDir).Info("Starting")

	if err := r.validateBackupDir(); err != nil {
		return err
	}

	return nil
}

func (r *RestoreTask) validateBackupDir() error {
	r.log.Info("Verifying backup folder")

	dirEntry, err := os.ReadDir(r.backupDir)
	if err != nil {
		return err
	}

	var importableCount int
	var dirs []string = nil
	for _, entry := range dirEntry {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		name := entry.Name()
		if entry.IsDir() {
			if mailFolderRegExp.MatchString(name) {
				r.log.WithField("name", name).Info("Found a potential backup sub-folder")
				dirs = append(dirs, name)
			}
			continue
		}

		if !strings.HasSuffix(name, ".eml") {
			if !strings.HasSuffix(name, ".metadata.json") {
				r.log.WithField("fileName", name).Warn("Ignoring unknown file")
			}
			continue
		}

		jsonFile := strings.TrimSuffix(name, ".eml") + ".metadata.json"
		stats, err := os.Stat(filepath.Join(r.backupDir, jsonFile))
		if err != nil {
			r.log.WithError(err).WithField("jsonFile", jsonFile).Warn("EML file has no associated JSON file")
			continue
		}
		if stats.IsDir() {
			r.log.WithField("jsonFile", jsonFile).Warn("JSON file is directory")
			continue

		}
		importableCount++
	}

	if importableCount > 0 {
		r.log.WithField("mailCount", importableCount).Info("Importable emails found")
		return nil
	}

	if len(dirs) == 0 {
		return errors.New("no importable mail found")
	}

	if len(dirs) > 1 {
		return errors.New("the specified folder contains more than one backup sub-folder")
	}

	r.log.WithField("folderName", dirs[0]).Info("A potential backup sub-folder has been found and will be inspected")
	r.backupDir = filepath.Join(r.backupDir, dirs[0])
	return r.validateBackupDir()
}
