package tracker

import (
	"fmt"
	"io"
	"os"
)

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// An EventLog outputs events as JSON to STDOUT.
type EventLog struct {
	writer io.Writer
}

func NewEventLog(writer io.Writer) {
	if writer == nil {
		writer = os.Stdout
	}
}

// Store an event
func (store EventLog) Store(event Event) error {
	fmt.Println(event.ToJson())

	return nil
}
