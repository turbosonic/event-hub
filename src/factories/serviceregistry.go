package factories

import (
	"log"
	"os"

	"github.com/turbosonic/event-hub/src/serviceregistry"
	"github.com/turbosonic/event-hub/src/serviceregistry/clients/http"
	"github.com/turbosonic/event-hub/src/serviceregistry/clients/kubernetes"
	"github.com/turbosonic/event-hub/src/serviceregistry/clients/static"
)

// ServiceRegistryClient ...generates a concrete ServiceRegistryClient from environment variables
func ServiceRegistryClient() serviceregistry.Client {
	registryclient := os.Getenv("SERVICE_REGISTRY_CLIENT")

	switch registryclient {
	case "static":
		log.Println("[x] Using a static Service Registry")
		return static.New()
	case "kubernetes":
		log.Println("[x] Using a Kubernetes annotations as a Service Registry")
		return kubernetes.New()
	default:
		log.Println("[x] Using the default http Service Registry")
		return http.New()
	}
}
