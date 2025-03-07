package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "security.protocol": "PLAINTEXT",
    })
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer p.Close()

    topic := "user_activity"

    events := []map[string]interface{}{
        {"userId": "u1", "eventType": "page_view", "timestamp": time.Now().Unix()},
        {"userId": "u2", "eventType": "click", "timestamp": time.Now().Unix()},
        {"userId": "u3", "eventType": "purchase", "timestamp": time.Now().Unix()},
        {"userId": "u1", "eventType": "page_view", "timestamp": time.Now().Unix()},
        {"userId": "u2", "eventType": "click", "timestamp": time.Now().Unix()},
        {"userId": "u3", "eventType": "purchase", "timestamp": time.Now().Unix()},
        {"userId": "u1", "eventType": "page_view", "timestamp": time.Now().Unix()},
        {"userId": "u2", "eventType": "click", "timestamp": time.Now().Unix()},
        {"userId": "u3", "eventType": "purchase", "timestamp": time.Now().Unix()},
        {"userId": "u1", "eventType": "page_view", "timestamp": time.Now().Unix()},
        {"userId": "u2", "eventType": "click", "timestamp": time.Now().Unix()},
        {"userId": "u3", "eventType": "purchase", "timestamp": time.Now().Unix()},
        {"userId": "u1", "eventType": "page_view", "timestamp": time.Now().Unix()},
        {"userId": "u2", "eventType": "click", "timestamp": time.Now().Unix()},
        {"userId": "u3", "eventType": "purchase", "timestamp": time.Now().Unix()},
    }

    for _, event := range events {
        data, _ := json.Marshal(event)
        key := event["userId"].(string) // Use userId as the key

        err = p.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
            Key:           []byte(key), // Key ensures partitioning by user
            Value:         data,
        }, nil)

        if err != nil {
            log.Fatalf("Failed to send message: %v", err)
        }
    }

    p.Flush(5000)
    log.Println("Sent user activity events")
}
