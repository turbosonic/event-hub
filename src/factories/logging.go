package factories

import (
	"log"
	"os"

	"contino.visualstudio.com/event-hub/src/logging"
	"contino.visualstudio.com/event-hub/src/logging/clients/applicationinsights"
	"contino.visualstudio.com/event-hub/src/logging/clients/elasticsearch"
	"contino.visualstudio.com/event-hub/src/logging/clients/stdout"
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
	default:
		log.Println("[x] Using stdout for logging")
		return stdout.New()
	}
}
