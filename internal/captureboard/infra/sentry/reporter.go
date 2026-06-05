package sentry

import (
	"time"

	"capture-board-selector/internal/captureboard/domain"
	"github.com/getsentry/sentry-go"
)

type Reporter struct{}

func NewReporter(dsn string) domain.ErrorReporter {
	sentry.Init(sentry.ClientOptions{Dsn: dsn})
	return &Reporter{}
}

func (r *Reporter) Report(err error) {
	sentry.CaptureException(err)
}

func (r *Reporter) Flush() {
	sentry.Flush(2 * time.Second)
}
