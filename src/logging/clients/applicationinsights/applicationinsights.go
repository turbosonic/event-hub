package applicationinsights

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"contino.visualstudio.com/event-hub/src/logging"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

type applicationInsightsLogger struct {
	client appinsights.TelemetryClient
}

func New() *applicationInsightsLogger {
	ail := applicationInsightsLogger{}
	ail.client = appinsights.NewTelemetryClient(os.Getenv("APPLICATIONINSIGHTS_INTRUMENTATION_KEY"))

	return &ail
}

func (ail applicationInsightsLogger) LogEvent(l *logging.EventLog) {
	event := appinsights.NewEventTelemetry(l.Name)
	event.Timestamp = l.Handled
	event.Properties["ID"] = l.ID
	event.Properties["Source"] = l.Source
	event.Properties["Submitted"] = l.Timestamp.String()
	event.Properties["Handled"] = l.Handled.String()
	event.Properties["RequestID"] = l.RequestID
	event.Properties["Valid"] = strconv.FormatBool(l.Valid)
	event.Properties["InvalidMessages"] = strings.Join(l.InvalidMessages, ", ")
	event.Properties["DeliveredQueues"] = strings.Join(l.DeliveredQueues, ", ")
	event.Properties["DeliveredTopics"] = strings.Join(l.DeliveredTopics, ", ")

	payload, err := json.Marshal(l.Payload)
	if err != nil {
		event.Properties["Payload"] = "Could not convert payload to json"
	} else {
		event.Properties["Payload"] = string(payload)
	}

	ail.client.Track(event)
}
