package tracker

import (
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/masahide/fluent-logger-golang/fluent"
	"io"
	"os"
)

// An EventStore is used to permanently store events
type EventLogger interface {
	Log(event Event) error
}


func NewEventLogger(writer io.Writer) (EventLogger, error) {
  socketPath := os.Getenv("FLUENT_SOCKET")
  if socketPath == "" {
    return newIoEventLogger(writer)
  } else {
    return newFluentEventLogger(socketPath)
  }
}

// An EventLogger outputs events as JSON to an io.Writer
type IoEventLogger struct {
	Writer io.Writer
}

func (l *IoEventLogger) Log(event Event) error {
	_, error := io.WriteString(l.Writer, event.ToJson()+"\n")
	return error
}

func newIoEventLogger(writer io.Writer) (*IoEventLogger, error) {
    if writer == nil {
      writer = os.Stdout
    }
    return &IoEventLogger{Writer: writer}, nil
}


type FluentEventLogger struct {
	logger *fluent.Fluent
}

func newFluentEventLogger(socketPath string) (*FluentEventLogger, error) {
	logger, err := fluent.New(fluent.Config{
		FluentSocketPath: socketPath,
		FluentNetwork:    "unix",
	})

	if err != nil {
		return &FluentEventLogger{}, err
	}
	return &FluentEventLogger{logger: logger}, err
}

func (l *FluentEventLogger) Log(event Event) error {
	error := l.logger.Post("tracker.event", event)
	return error
}
