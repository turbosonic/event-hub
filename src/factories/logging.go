package factories

import (
	"fmt"
	"log"
	"os"

	"github.com/turbosonic/event-hub/src/logging"
	"github.com/turbosonic/event-hub/src/logging/clients/applicationinsights"
	"github.com/turbosonic/event-hub/src/logging/clients/elasticsearch"
	"github.com/turbosonic/event-hub/src/logging/clients/influxdb"
	"github.com/turbosonic/event-hub/src/logging/clients/stdout"
)

func NewLoggingClient() logging.LoggingClient {
	lgc := os.Getenv("LOGGING_CLIENT")

	switch lgc {
	case "applicationinsights":
		log.Println("[x] Using Application insights for logging")
		return applicationinsights.New()
	case "elasticsearch":
		log.Println("[x] Using ElasticSearch for logging")
		return elasticsearch.New()
	case "influxdb":
		fmt.Println("[x] Logging to InfluxDB")
		return influxdb.New()
	default:
		log.Println("[x] Using stdout for logging")
		return stdout.New()
	}
}
