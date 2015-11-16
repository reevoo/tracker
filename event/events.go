package event

import (
	"bytes"
	"encoding/json"
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
	"net/url"
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

// Returns the ID of an event.
func (event Event) Id() string {
	return event["id"][0]
}

// Returns true if no items are in the Event.
func (event Event) Empty() bool {
	// An Event has an ID by default,
	// so an empty event has 1 entry.
	return len(event) == 1
}

// Converts the Event to JSON format.
func (event Event) ToJson() string {
	jsonBytes, _ := json.Marshal(event)
	return string(jsonBytes[:])
}

// Converts the Event to query parameters.
// Only used for testing; hence why we ignore errors...
func (event Event) ToParams() string {
	var buf bytes.Buffer
	for key, values := range event {
		for _, value := range values {
			buf.WriteString(url.QueryEscape(key))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(value))
			buf.WriteByte('&')
		}
	}
	s := buf.String()
	return s[0 : len(s)-1]
}
