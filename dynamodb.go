package tracker

import ()

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// TODO: Write real client!
type DynamoDBEventStore struct{}

func (store DynamoDBEventStore) Store(event Event) error {
	return nil
}
