package event

import (
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
)

// An Event is a structure holding information
// about something that has happened in one of our applications.
type Event map[string][]string

// Create a new Event.
func New(params map[string][]string) Event {
	id, _ := uuid.NewV4()
	params["id"] = []string{id.String()}

	return params
}

// Returns true if no items are in the Event.
func (event Event) Empty() bool {
	// An Event has an ID by default,
	// so an empty event has 1 entry.
	return len(event) == 1
}
