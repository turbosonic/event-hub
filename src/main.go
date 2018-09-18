package main

import (
	"fmt"
	"log"
	"time"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"

	"contino.visualstudio.com/event-hub/src/dto"
	"contino.visualstudio.com/event-hub/src/factories"
	"contino.visualstudio.com/event-hub/src/logging"
	"contino.visualstudio.com/event-hub/src/mq"
	"contino.visualstudio.com/event-hub/src/serviceregistry"
)

func main() {

	// Load the environmentals
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Set up a connection to the logging client
	// Providing a single interface for all types of logging
	lgc := factories.NewLoggingClient()
	logger := logging.NewLogger(lgc)

	// Set up the connection to the Service Registry
	// This loads all Microservice subscriptions in to memory
	// and refreshes them every 5 seconds
	src := factories.ServiceRegistryClient()
	sr := serviceregistry.New(src)

	// Set up the connection to the MQ
	// This provides an interface to the message broker
	mqc := factories.MQClient()
	msq := mq.New(mqc)

	c := make(chan []byte)
	// from here we just need to start listening and pass through the function which emits the event
	go mqc.Listen(c)

	for {
		go handleEvent(<-c, sr, msq, logger)
	}
}

func handleEvent(eventByteArray []byte, sr *serviceregistry.ServiceRegistry, msq *mq.MQ, logger *logging.Logger) {
	e, err := dto.NewEventFromByteArray(&eventByteArray)
	if err != nil {
		log.Println("Could not create event from message: ", fmt.Sprint(&eventByteArray))
		return
	}

	// give the event an ID
	u, _ := uuid.NewV4()
	e.ID = fmt.Sprint(u)

	// add a handled time
	e.Handled = time.Now()

	// create a log event
	logEvent := logging.NewEventLog(&e)

	// validate it
	valid, errors := e.Validate()
	if valid == false {
		// if it fails validation we should log it
		logEvent.Valid = false
		logEvent.InvalidMessages = errors
		logger.LoggingClient.LogEvent(&logEvent)
		return
	}

	// turn the event in to a byte array ready to send
	eba, err := e.ToByteArray()
	if err != nil {
		logEvent.Valid = false
		logEvent.InvalidMessages = append(logEvent.InvalidMessages, err.Error())
		logger.LoggingClient.LogEvent(&logEvent)
		return
	}

	// get the services which need to be sent to queues
	qmsvs := sr.GetMicroservicesByQueueEvents(e.Name)

	// distribute it to the queues
	for _, m := range qmsvs {
		err := msq.Client.SendToQueue(m.Name, &eba)
		if err != nil {
			logEvent.FailedQueues[m.Name] = err.Error()
		} else {
			logEvent.DeliveredQueues = append(logEvent.DeliveredQueues, m.Name)
		}
	}

	// get the services which the message need to be sent to topics
	tmsvs := sr.GetMicroservicesByTopicEvents(e.Name)

	// distribute it to topics
	for _, m := range tmsvs {
		err := msq.Client.SendToTopic(m.Name, &eba)
		if err != nil {
			logEvent.FailedTopics[m.Name] = err.Error()
		} else {
			logEvent.DeliveredTopics = append(logEvent.DeliveredTopics, m.Name)
		}
	}

	// log it
	logger.LoggingClient.LogEvent(&logEvent)
}
