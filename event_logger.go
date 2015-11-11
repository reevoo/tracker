package tracker

import (
	"io"
	"os"
)

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// An EventLogger outputs events as JSON to STDOUT.
type EventLogger struct {
	Writer io.Writer
}

func NewEventLogger(writer io.Writer) EventLogger {
	if writer == nil {
		writer = os.Stdout
	}

	return EventLogger{Writer: writer}
}

// Store an event
func (store EventLogger) Store(event Event) error {
	_, error := io.WriteString(store.Writer, event.ToJson() + "\n")
	return error
}
