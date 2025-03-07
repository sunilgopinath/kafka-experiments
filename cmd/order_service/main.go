package main

import (
	"encoding/json"
	"kafkademo/internal/models"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	// Sample orders
	orders := []models.OrderEvent{
		{"order1", "user1", "order_created", time.Now().Unix()},
		{"order2", "user2", "order_created", time.Now().Unix()},
	}

	for _, order := range orders {
		data, _ := json.Marshal(order)
		topic := "order_created"

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          data,
		}, nil)

		if err != nil {
			log.Printf("Failed to send order event: %v\n", err)
		} else {
			log.Printf("🛒 Order Created: %+v\n", order)
		}
	}

	p.Flush(5000)
	log.Println("All orders sent.")
}
