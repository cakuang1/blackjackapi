package models

import (
	"net/url"

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

func ProduceMessage(address, key, message, user, pass string) error {
	// Encode the message for inclusion in the URL
	messageEncoded := url.QueryEscape(message)

	// Construct the URL with properly encoded components
	url := fmt.Sprintf("%s/produce/broadcast/%s", address, messageEncoded)

	// Create a new HTTP GET request (since you're using "GET" method)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.SetBasicAuth(user, pass)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return fmt.Errorf("non-OK status code received: %s", resp.Status)
	}
	return nil
}
