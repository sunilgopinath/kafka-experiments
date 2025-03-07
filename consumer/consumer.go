package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
)


const (
	postgresConn = "postgresql://postgres@localhost:5432/user_tracking"
)

type UserEvent struct {
	UserID    string `json:"userId"`
	EventType string `json:"eventType"`
	Timestamp int64  `json:"timestamp"`
}

func main() {

    // Connect to PostgreSQL
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), postgresConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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

    go func() {
        <-sigchan
        fmt.Println("\nShutting down consumer...")
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

			// Insert into PostgreSQL
			_, err = conn.Exec(context.Background(),
				"INSERT INTO user_activity (user_id, event_type, timestamp) VALUES ($1, $2, $3)",
				event.UserID, event.EventType, event.Timestamp)
			if err != nil {
				log.Printf("Failed to insert event into database: %v\n", err)
			} else {
				log.Printf("Inserted event into database: %+v\n", event)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
