package elasticsearch

import (
	"context"
	"os"
	"time"

	"github.com/turbosonic/event-hub/src/logging"
	elastic "gopkg.in/olivere/elastic.v5"
)

type elasticsearchLogger struct {
	client *elastic.Client
}

func New() *elasticsearchLogger {
	elasticURL := os.Getenv("LOGGING_ELASTIC_URL")
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticURL))
	if err != nil {
		panic(err)
	}

	esl := elasticsearchLogger{}
	esl.client = client

	return &esl
}

func (esl *elasticsearchLogger) LogEvent(l *logging.EventLog) {
	index := "event-hub-" + time.Now().Format("2006-01-02")
	ctx := context.Background()
	esl.client.Index().Index(index).Type("event").Id(l.ID).BodyJson(&l).Do(ctx)
}
