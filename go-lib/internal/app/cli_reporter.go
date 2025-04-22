package app

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
			0,
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetWidth(100),
		),
	}
}

func (m *cliReporter) SetMessageTotal(total uint64) {
	m.progressbar.Reset()
	m.progressbar.ChangeMax64(int64(total)) //nolint:gosec // no potential to overflow.
	m.totalMessageCount.Store(total)
}

func (m *cliReporter) SetMessageProcessed(total uint64) {
	_ = m.progressbar.Set64(int64(total)) //nolint:gosec // no potential to overflow.
	m.currentMessageCount.Store(total)
}

func (m *cliReporter) OnProgress(delta int) {
	_ = m.currentMessageCount.Add(uint64(delta)) //nolint:gosec // yet again, we shouldn't overflow.
	_ = m.progressbar.Add(delta)
}
