package tracker

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
)

// An unrecoverable error in Tracker.
type TrackerError struct {
	name    string
	context map[string]string
}

// Convert a Go error into a TrackerError.
func NewTrackerErrorFromError(err error, context map[string]string) TrackerError {
	return TrackerError{
		name: err.Error(),
	}
}

// Convert a TrackerError to a Map.
// Useful for JSON-based error logging due to its hierarchy.
func (err TrackerError) ToMap() map[string]string {
	all := err.context
	all["name"] = err.name
	return all
}

// Convert a TrackerError to a Go error.
func (err TrackerError) ToError() error {
	return errors.New(err.name)
}

// EventLoggers keep a log of errors.
type ErrorLogger interface {
	LogError(err TrackerError)
}

// Sentry is used to log errors.
type SentryErrorLogger struct{}

// Logs an error to Sentry.
func (SentryErrorLogger) LogError(err TrackerError) {
	packet := raven.NewPacket(err.name, raven.NewException(err.ToError(), raven.NewStacktrace(2, 3, nil)))
	raven.Capture(packet, err.ToMap())
}

// Logs errors to the console.
type ConsoleLogger struct{}

// Logs an error to the console.
func (ConsoleLogger) LogError(err TrackerError) {
	fmt.Printf("[ERROR] %s", err.name)
}
