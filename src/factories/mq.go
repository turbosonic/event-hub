package factories

import (
	"log"
	"os"

	"github.com/turbosonic/event-hub/src/mq"
	"github.com/turbosonic/event-hub/src/mq/clients/activemq"
	"github.com/turbosonic/event-hub/src/mq/clients/azureservicebus"
)

// MQClient ...generates a concrete MQClient from environment variables
func MQClient() mq.Client {
	mqc := os.Getenv("MQ_CLIENT")

	switch mqc {
	case "azureservicebus":
		log.Println("[x] Using Azure Service Bus as a message broker")
		return azureservicebus.New()
	default:
		log.Println("[x] Using ActiveMQ as a message broker")
		return activemq.New()
	}
}
