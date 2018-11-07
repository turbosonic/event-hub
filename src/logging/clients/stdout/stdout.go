package stdout

import (
	"fmt"

	"github.com/turbosonic/event-hub/src/logging"
)

type stdOutLogger struct {
}

func New() *stdOutLogger {
	sol := stdOutLogger{}
	return &sol
}

func (std stdOutLogger) LogEvent(l *logging.EventLog) {
	fmt.Printf("%+v\n", *l)
}
