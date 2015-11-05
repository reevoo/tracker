package tracker

import (
	"fmt"
	"github.com/getsentry/raven-go"
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
func NewTrackerErrorFromError(err error, Context map[string]string) TrackerError {
	return TrackerError{
		Name: err.Error(),
	}
}

// Convert a TrackerError to a Map.
// Useful for JSON-based error logging due to its hierarchy.
func (err TrackerError) ToMap() map[string]string {
	all := err.Context
	all["name"] = err.Name
	return all
}

// EventLoggers keep a log of errors.
type ErrorLogger interface {
	LogError(err TrackerError)
}

// Sentry is used to log errors.
type SentryErrorLogger struct{}

// Logs an error to Sentry.
func (SentryErrorLogger) LogError(err TrackerError) {
	packet := raven.NewPacket(err.Error(), raven.NewException(err, raven.NewStacktrace(2, 3, nil)))
	raven.Capture(packet, err.ToMap())
}

// Logs errors to the console.
type ConsoleLogger struct{}

// Logs an error to the console.
func (ConsoleLogger) LogError(err TrackerError) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("[ERROR] %s", err.Error()))

}
