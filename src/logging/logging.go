package logging

import (
	"time"

	"contino.visualstudio.com/event-hub/src/dto"
)

type EventLog struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Source          string            `json:"source"`
	Timestamp       time.Time         `json:"timestamp"`
	Handled         time.Time         `json:"handled_timestamp"`
	Payload         interface{}       `json:"payload"`
	RequestID       string            `json:"request_id"`
	Valid           bool              `json:"valid"`
	InvalidMessages []string          `json:"invalid_message,omitempty"`
	DeliveredQueues []string          `json:"delivered_queues"`
	DeliveredTopics []string          `json:"delivered_topics"`
	FailedQueues    map[string]string `json:"failed_queues"`
	FailedTopics    map[string]string `json:"failed_topics"`
}

func NewEventLog(e *dto.Event) EventLog {
	return EventLog{
		e.ID,
		e.Name,
		e.Source,
		e.Timestamp,
		e.Handled,
		e.Payload,
		e.RequestID,
		true,
		[]string{},
		[]string{},
		[]string{},
		make(map[string]string),
		make(map[string]string),
	}
}

type LoggingClient interface {
	LogEvent(*EventLog)
}

type Logger struct {
	LoggingClient LoggingClient
}

func NewLogger(lc LoggingClient) *Logger {
	return &Logger{
		lc,
	}
}
