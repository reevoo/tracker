package tracker

import (
	"fmt"
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/getsentry/raven-go"
	"io"
	"os"
)

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.SetRelease(os.Getenv("SENTRY_RELEASE"))
}

// An unrecoverable error in Tracker.
type TrackerError struct {
	Name    string
	Context map[string]string
}

// Returns the error description in full.
func (err TrackerError) Error() string {
	return err.Name
}

// Convert a Go error into a TrackerError.
func NewTrackerErrorFromError(err error, context map[string]string) TrackerError {
	return TrackerError{
		Name:    err.Error(),
		Context: context,
	}
}

// Convert a TrackerError to a Map.
// Useful for JSON-based error logging due to its hierarchy.
func (err TrackerError) ToMap() map[string]string {
	if err.Context != nil {
		err.Context["name"] = err.Name
	}
	return err.Context
}

// EventLoggers keep a log of errors.
type ErrorLogger interface {
	LogError(err TrackerError)
}

// Sentry is used to log errors.
type SentryErrorLogger struct {
	client interface {
		CaptureError(err error, tags map[string]string) string
	}
}

// Logs an error to Sentry.
func (s SentryErrorLogger) LogError(err TrackerError) {
	s.client.CaptureError(err, err.ToMap())
}

// Logs errors to the console.
type ConsoleLogger struct {
	writer io.Writer
}

// Logs an error to the console.
func (c ConsoleLogger) LogError(err TrackerError) {
	c.writer.Write([]byte(fmt.Sprintf("[ERROR] %s", err.Error())))
}
