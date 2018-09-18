package serviceregistry

import (
	"testing"
)

type client struct {
}

func (c client) GetMicroservices() ([]Microservice, error) {
	m := make([]Microservice, 2)

	m[0] = Microservice{
		"microservice-one",
		[]string{"event-one", "event-two"},
		[]string{"event-one", "event-two"},
	}

	m[1] = Microservice{
		"microservice-two",
		[]string{"event-one", "event-three"},
		[]string{"event-one", "event-three"},
	}
	return m, nil
}

func TestGetServiceRegistry(t *testing.T) {
	c := client{}
	sr := New(c)

	if len(sr.Microservices) != 2 {
		t.Errorf("Did not load all microservices")
	}

}
