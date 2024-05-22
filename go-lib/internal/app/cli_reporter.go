package app

import "C"
import (
	"sync/atomic"

	"github.com/schollz/progressbar/v3"
)

type cliReporter struct {
	totalMessageCount   atomic.Uint64
	currentMessageCount atomic.Uint64
	progressbar         *progressbar.ProgressBar
}

func newCliReporter() *cliReporter {
	return &cliReporter{
		progressbar: progressbar.NewOptions64(
			1000,
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionFullWidth(),
		),
	}
}

func (m *cliReporter) SetMessageTotal(total uint64) {
	m.progressbar.ChangeMax64(int64(total))
	m.totalMessageCount.Store(total)
}

func (m *cliReporter) SetMessageDownloaded(total uint64) {
	_ = m.progressbar.Set64(int64(total))
	m.currentMessageCount.Store(total)
}

func (m *cliReporter) OnProgress(delta int) {
	_ = m.currentMessageCount.Add(uint64(delta))
	_ = m.progressbar.Add(delta)
}
