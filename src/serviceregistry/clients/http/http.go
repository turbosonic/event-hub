package http

import (
	"encoding/json"
	"net/http"
	"os"

	sr "github.com/turbosonic/event-hub/src/serviceregistry"
)

// ServiceRegistryClient ...an http ServiceRegistryClient
type ServiceRegistryClient struct {
}

// GetMicroservices ...gets the current microservices from the Service Registry's API
func (c ServiceRegistryClient) GetMicroservices() ([]sr.Microservice, error) {
	var m []sr.Microservice

	err := getJSON(os.Getenv("SERVICE_REGISTRY_URL"), &m)

	return m, err
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
