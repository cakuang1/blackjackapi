package handlers

import (
	"blackjackapi/models"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"net/http"
	"strings"
)

// documentation from upstash kafka.

func (h *Handler) KafkaSSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	vars := mux.Vars(r)
	tableID := vars["tableID"]

	_, err := models.GetSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Failed to retrieve table from Redis.Ensure tableID has been created and is correct", http.StatusInternalServerError)
		return
	}
	connected := fmt.Sprintf("Connected to table %s\n\n", tableID)
	if _, err := w.Write([]byte(connected)); err != nil {
		log.Printf("Error writing SSE event to response: %v", err)
		http.Error(w, "Failed to write connection message", http.StatusInternalServerError)
		return
	}

	mechanism, _ := scram.Mechanism(scram.SHA512, "bG9naWNhbC1iYXNzLTEwMzczJEjezQWu6R_uRLWL8ASgp3SIKfxNkd4qfoF7yRs", "YTcxZGUxNWYtMjYxZS00YjFmLWE5MDktODBkOGYwMDJjNWE5")
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"logical-bass-10373-us1-kafka.upstash.io:9092"},
		GroupID: tableID,
		Topic:   "broadcast",
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
	})
	defer kafkaReader.Close()
	ctx := context.Background()
	for {
		message, err := kafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}
		// Check if the message key matches the desired key (tableID)
		replacedString := strings.ReplaceAll(string(message.Value), "+", " ")
		_, err = w.Write([]byte(replacedString))
		if err != nil {
			log.Printf("Error writing SSE event to response: %v", err)
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

}
