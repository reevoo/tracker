package logger

import (
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/masahide/fluent-logger-golang/fluent"
)

type FluentLogger struct {
	logger *fluent.Fluent
}

func (l *FluentLogger) Log(message interface{}) error {
	error := l.logger.Post("tracker.event", message)
	return error
}

func newFluentLogger(socketPath string) (*FluentLogger, error) {
	logger, err := fluent.New(fluent.Config{
		FluentSocketPath: socketPath,
		FluentNetwork:    "unix",
	})

	if err != nil {
		return &FluentLogger{}, err
	}
	return &FluentLogger{logger: logger}, err
}
