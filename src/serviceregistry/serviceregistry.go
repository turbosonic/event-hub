package serviceregistry

import (
	"log"
	"time"
)

// Microservice ...details of the events a single microservice subscribes to
type Microservice struct {
	Name        string   `json:"name" yaml:"name"`
	QueueEvents []string `json:"queue_events" yaml:"queue_events"`
	TopicEvents []string `json:"topic_events" yaml:"topic_events"`
}

// Client ...an interface which provides a list of microservices and their event subscriptions
type Client interface {
	GetMicroservices() ([]Microservice, error)
}

// ServiceRegistry ...the struct which contains all microservices in the ecosystem, in memory
type ServiceRegistry struct {
	client        Client
	Microservices []Microservice
}

// New ...creates a new Service Registry from an injected client, does an initial get and stores microservice in memory,
func New(sr Client) *ServiceRegistry {
	s := ServiceRegistry{}
	s.client = sr
	s.getServices()
	startRetrievalInterval(&s)
	return &s
}

func (sr ServiceRegistry) GetMicroservicesByQueueEvents(eventName string) []Microservice {
	var m []Microservice
	for _, s := range sr.Microservices {
		for _, e := range s.QueueEvents {
			if e == eventName {
				m = append(m, s)
			}
		}
	}
	return m
}

func (sr ServiceRegistry) GetMicroservicesByTopicEvents(eventName string) []Microservice {
	var m []Microservice
	for _, s := range sr.Microservices {
		for _, e := range s.TopicEvents {
			if e == eventName {
				m = append(m, s)
			}
		}
	}
	return m
}

func startRetrievalInterval(sr *ServiceRegistry) {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			sr.getServices()
		}
	}()
}

func (sr *ServiceRegistry) getServices() {
	ms, err := sr.client.GetMicroservices()
	if err != nil {
		log.Println("Could not retrieve services:", err, "- trying again in 5 seconds")
	} else if len(ms) != len(sr.Microservices) {
		log.Println(len(ms), "service(s) retrieved, list updated")
		log.Print(ms)
	}
	sr.Microservices = ms
}
