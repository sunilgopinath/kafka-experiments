package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
)

const (
	postgresConn  = "postgresql://postgres@localhost:5432/user_tracking"
	batchSize     = 10 // Number of events before batch insert
)

type UserEvent struct {
	UserID    string `json:"userId"`
	EventType string `json:"eventType"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	// Connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), postgresConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

    // Create Kafka producer for DLQ
    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
    })
    if err != nil {
        log.Fatalf("Failed to create DLQ producer: %v", err)
    }
    defer producer.Close()

	// Set up Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "user-activity-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"user_activity"}, nil)

	// Handle shutdown gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	eventBatch := make([]UserEvent, 0, batchSize)

	// Graceful shutdown handler
	go func() {
		<-sigchan
		fmt.Println("\nShutting down consumer...")

		// Flush remaining batch before exiting
		if len(eventBatch) > 0 {
			insertBatch(conn, eventBatch, producer)
		}

		c.Close()
		os.Exit(0)
	}()

	// Read messages from Kafka
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received: %s\n", string(msg.Value))

			// Parse message
			var event UserEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to parse event: %v\n", err)
				continue
			}

			// Add event to batch
			eventBatch = append(eventBatch, event)

			// If batch reaches batchSize, insert into PostgreSQL
			if len(eventBatch) >= batchSize {
				insertBatch(conn, eventBatch, producer)
				eventBatch = eventBatch[:0] // Reset batch
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}

// insertBatch inserts multiple events into PostgreSQL in a single query
func insertBatch(conn *pgx.Conn, batch []UserEvent, producer *kafka.Producer) {
	maxRetries := 5
	baseDelay := 100 * time.Millisecond // Start with 100ms
	jitterFactor := 0.3                 // Add up to 30% random jitter

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempt %d: Inserting batch of %d events...\n", attempt, len(batch))

		// Prepare SQL statement
		sql := "INSERT INTO user_activity (user_id, event_type, timestamp) VALUES "
		args := make([]interface{}, 0, len(batch)*3)
		placeholders := ""

		for i, event := range batch {
			placeholders += fmt.Sprintf("($%d, $%d, $%d),", i*3+1, i*3+2, i*3+3)
			args = append(args, event.UserID, event.EventType, event.Timestamp)
		}
		sql = sql + placeholders[:len(placeholders)-1] // Remove last comma

		// Execute batch insert
		_, err := conn.Exec(context.Background(), sql, args...)
		if err == nil {
			log.Printf("Successfully inserted batch of %d events\n", len(batch))
			return
		}

		log.Printf("Failed to insert batch (attempt %d/%d): %v\n", attempt, maxRetries, err)

		// Calculate exponential backoff with jitter
		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt))) // Exponential increase
		jitter := time.Duration(rand.Float64() * jitterFactor * float64(backoffTime)) // Add random jitter
		sleepTime := backoffTime + jitter

		log.Printf("Retrying in %v...\n", sleepTime)
		time.Sleep(sleepTime)
	}

	log.Printf("Max retries reached. Sending batch to Dead Letter Queue (DLQ).")

	// Send failed batch to Dead Letter Queue
	sendToDLQ(batch, producer)
}

func sendToDLQ(batch []UserEvent, producer *kafka.Producer) {
	topic := "dead_letter_queue"

	for _, event := range batch {
		data, _ := json.Marshal(event)

		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          data,
		}, nil)

		if err != nil {
			log.Printf("Failed to send event to DLQ: %v\n", err)
		} else {
			log.Printf("Sent failed event to DLQ: %s\n", string(data))
		}
	}

	producer.Flush(5000) // Ensure all messages are sent
}
