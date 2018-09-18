package dto

import (
	"encoding/json"
	"time"
)

//Event ...an incomming event object
type Event struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Source    string      `json:"source"`
	Timestamp time.Time   `json:"timestamp"`
	Handled   time.Time   `json:"handled_timestamp"`
	Payload   interface{} `json:"payload"`
	RequestID string      `json:"request_id"`
}

func NewEventFromByteArray(bytes *[]byte) (Event, error) {
	e := Event{}
	err := json.Unmarshal(*bytes, &e)
	return e, err
}

// Validate ...returns a true/false value indicating if the event is valid, also contains a list of issues
func (e *Event) Validate() (bool, []string) {
	invalids := []string{}

	if e.Name == "" {
		invalids = append(invalids, "Name is required")
	}

	if e.Source == "" {
		invalids = append(invalids, "Source is required")
	}

	defaultTime := time.Time{}
	if e.Timestamp == defaultTime {
		invalids = append(invalids, "Timestamp is required")
	}

	if len(invalids) > 0 {
		return false, invalids
	}

	return true, []string{}
}

func (e *Event) ToByteArray() ([]byte, error) {
	return json.Marshal(e)
}
