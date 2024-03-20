package handlers

import (
    "blackjackapi/models"
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
    // Set HTTP headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Get tableID from request parameters
    vars := mux.Vars(r)
    tableID := vars["tableID"]

    // Get session information
    _, err := models.GetSession(h.Context, tableID, h.Client)
    if err != nil {
        http.Error(w, "Failed to retrieve table from Redis. Ensure tableID has been created and is correct", http.StatusInternalServerError)
        return
    }

    // Write connection message
    connected := fmt.Sprintf("Connected to table %s\n\n", tableID)
    _, err = w.Write([]byte(connected))
    if err != nil {
        log.Printf("Error writing SSE event to response: %v", err)
        return
    }

    // Create Kafka reader
    mechanism, _ := scram.Mechanism(scram.SHA512, h.KAFKAUSERNAME, h.KAFKAPASSWORD)
    kafkaReader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     []string{h.KAFKAADDRESS},
        GroupID:     tableID,
        Topic:       "broadcast",
        Dialer:      &kafka.Dialer{SASLMechanism: mechanism, TLS: &tls.Config{}},
        StartOffset: kafka.LastOffset,
        // Disable auto-commit to manually commit offsets
        CommitInterval: 0,
    })

    // Close Kafka reader when function returns
    defer kafkaReader.Close()

    // Read messages from Kafka and send them to the client
    ctx := r.Context()
    for {
        select {
        case <-ctx.Done():
            return // Stop processing if the request is canceled
        default:
            message, err := kafkaReader.ReadMessage(ctx)
            if err != nil {
                log.Printf("Error reading message from Kafka: %v", err)
                return
            }
            
            // Check if the message key matches the desired tableID
            if string(message.Key) != tableID {
                continue
            }

            // Write message to the client
            replacedString := strings.ReplaceAll(string(message.Value), "+", " ")
            _, err = w.Write([]byte(replacedString + "\n\n"))
            if err != nil {
                log.Printf("Error writing SSE event to response: %v", err)
                return
            }
            
            // Flush the response writer to ensure data is sent immediately
            if f, ok := w.(http.Flusher); ok {
                f.Flush()
            }

            // Manually commit the offset to mark the message as consumed
            err = kafkaReader.CommitMessages(ctx, message)
            if err != nil {
                log.Printf("Error committing message offset: %v", err)
                return
            }
        }
    }
}
