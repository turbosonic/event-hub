package mq

type MQ struct {
	Client Client
}

// Client ...an common interface for all message brokers
type Client interface {
	Listen(chan []byte)
	SendToQueue(string, *[]byte) error
	SendToTopic(string, *[]byte) error
}

func New(client Client) *MQ {
	m := MQ{}
	m.Client = client
	return &m
}
