package kubernetes

import (
	"encoding/json"
	"net/http"
	"strings"

	sr "github.com/turbosonic/event-hub/src/serviceregistry"
)

type deployments struct {
	Items []deployment
}

type deployment struct {
	Metadata metadata
}

type metadata struct {
	Name        string
	Annotations map[string]string
}

// ServiceRegistryClient ...an http ServiceRegistryClient
type ServiceRegistryClient struct {
}

// GetMicroservices ...gets the current microservices from the Service Registry's API
func (c ServiceRegistryClient) GetMicroservices() ([]sr.Microservice, error) {
	var d deployments

	err := getJSON("http://localhost:8001/apis/apps/v1/deployments", &d)

	s := make([]sr.Microservice, len(d.Items))

	for i, deploy := range d.Items {
		s[i] = sr.Microservice{
			Name:        deploy.Metadata.Name,
			QueueEvents: strings.Split(deploy.Metadata.Annotations["queue_events"], ","),
			TopicEvents: strings.Split(deploy.Metadata.Annotations["topic_events"], ","),
		}
	}

	return s, err
}

// New ...creates a new http Service Registry Client
func New() ServiceRegistryClient {
	return ServiceRegistryClient{}
}

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
