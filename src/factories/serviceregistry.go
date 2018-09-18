package factories

import (
	"log"
	"os"

	"contino.visualstudio.com/event-hub/src/serviceregistry"
	"contino.visualstudio.com/event-hub/src/serviceregistry/clients/http"
	"contino.visualstudio.com/event-hub/src/serviceregistry/clients/static"
)

// ServiceRegistryClient ...generates a concrete ServiceRegistryClient from environment variables
func ServiceRegistryClient() serviceregistry.Client {
	registryclient := os.Getenv("SERVICE_REGISTRY_CLIENT")

	switch registryclient {
	case "static":
		log.Println("[x] Using a static Service Registry")
		return static.New()
	default:
		log.Println("[x] Using the default http Service Registry")
		return http.New()
	}
}
