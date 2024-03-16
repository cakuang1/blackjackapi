package models

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"net/http"
	"os"
	"time"
)

// documentation from upstash kafka.

func KafkaConsumer(groupId string) {
	mechanism, _ := scram.Mechanism(scram.SHA512, os.Getenv("KAFKAPASSWORD"), os.Getenv("KAFKASALT"))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"logical-bass-10373-us1-kafka.upstash.io:9092"},
		GroupID: groupId,
		Topic:   "broadcast",
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120) // Increase the timeout
	defer cancel()
	for {
		message, _ := r.ReadMessage(ctx)
		fmt.Println(message.Partition, message.Offset, string(message.Value))
	}
}

func produceMessage(key, message string) error {
	url := fmt.Sprintf("https://logical-bass-10373-us1-rest-kafka.upstash.io/produce/%s/%s", "broadcast", key)
	// Create a new HTTP POST request with the message as the payload
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(message))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	// Set the request headers
	req.SetBasicAuth(os.Getenv("KAFKAUSERNAME"), os.Getenv("KAFKAPASSWORD"))
	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK status code received: %s", resp.Status)
	}

	fmt.Println("Message successfully produced to Kafka topic with key.")
	return nil
}
