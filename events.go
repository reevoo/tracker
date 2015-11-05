package tracker

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

// An Event is a structure holding information
// about something that has happened in one of our applications.
type Event struct {
	Id       uuid.UUID
	Name     string                 `json:"name" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Converts the Event to JSON format.
func (event Event) ToJson() string {
	jsonBytes, _ := json.Marshal(event)
	return string(jsonBytes[:])
}
