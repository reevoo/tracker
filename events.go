package tracker

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

// Metadata is used to collect any event-specific context.
type Metadata map[string]interface{}

// An Event is a structure holding information
// about something that has happened in one of our applications.
type Event struct {
	Id       uuid.UUID
	Name     string   `json:"name" binding:"required"`
	Metadata Metadata `json:"metadata"`
}

// Create a new Event.
func NewEvent(name string, metadata Metadata) Event {
	id, _ := uuid.NewV4()

	event := Event{
		Id:       *id,
		Name:     name,
		Metadata: metadata,
	}

	return event
}

// Converts the Event to JSON format.
func (event Event) ToJson() string {
	jsonBytes, _ := json.Marshal(event)
	return string(jsonBytes[:])
}
