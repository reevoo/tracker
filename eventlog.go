package tracker

import ()

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// An EventLog outputs events as JSON to STDOUT.
type EventLog struct{}

func (store EventLog) Store(event Event) error {
	return nil
}
