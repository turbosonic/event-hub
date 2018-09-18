package dto

import (
	"testing"
	"time"
)

func TestEventValidationFails(t *testing.T) {
	e := Event{}
	valid, errors := e.Validate()

	if valid == true || len(errors) != 3 {
		t.Error("This should have been invalid and had three errors")
	}
}

func TestEventValidationPasses(t *testing.T) {
	e := Event{}
	e.Name = "event:name"
	e.Source = "test-service"
	e.Timestamp = time.Now()

	valid, errors := e.Validate()

	if valid == false || len(errors) != 0 {
		t.Error("This should have been valid and had no errors")
	}
}
