package tracker

import (
	"io"
	"os"
)

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// An EventLog outputs events as JSON to STDOUT.
type EventLog struct {
	Writer io.Writer
}

func NewEventLog(writer io.Writer) EventLog {
	if writer == nil {
		writer = os.Stdout
	}

	return EventLog{Writer: writer}
}

// Store an event
func (store EventLog) Store(event Event) error {
	_, error := store.Writer.Write([]byte(event.ToJson()))

	return error
}
