package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsClient{Conn: nc}, nil
}

func (n *NatsClient) Publish(subject string, msg interface{}) {
	// Convert the message to JSON
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message to JSON: %v", err)
		return
	}

	// Publish the JSON message
	if err := n.Conn.Publish(subject, jsonMsg); err != nil {
		log.Printf("Error publishing message to subject %s: %v", subject, err)
	}
}
