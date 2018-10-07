package azureservicebus

import (
	"context"
	"log"
	"os"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type MQClient struct {
	ns     *servicebus.Namespace
	queues map[string]*servicebus.Queue
	topics map[string]*servicebus.Topic
}

func (mqc MQClient) Listen(c chan []byte) {
	queueName := os.Getenv("INGRESS_QUEUE_NAME")

	if queueName == "" {
		queueName = "event-hub"
	}

	qm := mqc.ns.NewQueueManager()

	existingQueues, err := qm.List(context.Background())
	if err != nil {
		panic(err)
	}

	found := false
	for _, q := range existingQueues {
		if q.Name == queueName {
			found = true
		}
	}

	if !found {
		_, err := qm.Put(context.Background(), queueName)
		if err != nil {
			panic(err)
		}
	}

	q, err := mqc.ns.NewQueue(context.Background(), queueName)
	if err != nil {
		panic(err)
	}

	for {

		q.Receive(context.Background(),
			func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
				c <- msg.Data
				return msg.Complete()
			})
	}
}

func (mqc MQClient) SendToQueue(queueName string, event *[]byte) error {
	queueName = "queue/" + queueName
	_, found := mqc.queues[queueName]
	if found == false {
		qm := mqc.ns.NewQueueManager()

		existingQueues, err := qm.List(context.Background())
		if err != nil {
			panic(err)
		}

		found := false
		for _, q := range existingQueues {
			if q.Name == queueName {
				found = true
			}
		}

		if !found {
			_, err := qm.Put(context.Background(), queueName)
			if err != nil {
				panic(err)
			}
		}

		q, err := mqc.ns.NewQueue(context.Background(), queueName)
		if err != nil {
			log.Println("Create queue error:", err)
			return err
		}
		mqc.queues[queueName] = q
	}

	q := mqc.queues[queueName]

	return q.Send(context.Background(), servicebus.NewMessage(*event))
}

func (mqc MQClient) SendToTopic(topicName string, event *[]byte) error {
	topicName = "topic/" + topicName

	_, found := mqc.topics[topicName]
	if found == false {
		tm := mqc.ns.NewTopicManager()

		existingTopics, err := tm.List(context.Background())
		if err != nil {
			panic(err)
		}

		found := false
		for _, t := range existingTopics {
			if t.Name == topicName {
				found = true
			}
		}

		if !found {
			_, err := tm.Put(context.Background(), topicName)
			if err != nil {
				panic(err)
			}
		}

		t, err := mqc.ns.NewTopic(context.Background(), topicName)
		if err != nil {
			log.Println("Create Topic error:", err)
			return err
		}
		mqc.topics[topicName] = t
	}

	t := mqc.topics[topicName]

	return t.Send(context.Background(), servicebus.NewMessage(*event))
}

func New() MQClient {
	mqc := MQClient{}
	mqc.queues = make(map[string]*servicebus.Queue)
	mqc.topics = make(map[string]*servicebus.Topic)
	connect(&mqc)
	return mqc
}

func connect(mqc *MQClient) {
	connStr := os.Getenv("AZURE_SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		panic(err)
	}
	mqc.ns = ns
}
