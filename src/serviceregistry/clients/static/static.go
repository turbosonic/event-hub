package static

import (
	"io/ioutil"
	"log"
	"path/filepath"

	sr "github.com/turbosonic/event-hub/src/serviceregistry"
	yaml "gopkg.in/yaml.v2"
)

// ServiceRegistryClient ...an http ServiceRegistryClient
type ServiceRegistryClient struct {
}

// GetMicroservices ...gets the current microservices from the Service Registry's API
func (c ServiceRegistryClient) GetMicroservices() ([]sr.Microservice, error) {
	var m []sr.Microservice

	filename, _ := filepath.Abs("./serviceregistry/clients/static/services.yaml")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Could not read services.yaml file")
		panic("Could not read services.yaml file")
	}

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Println("Could not unmarshal yaml")
		panic("Could not unmarshal yaml")
	}

	return m, nil
}

// New ...creates a new http Service Registry Client
func New() ServiceRegistryClient {
	return ServiceRegistryClient{}
}
