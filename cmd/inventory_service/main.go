package main

import (
	"encoding/json"
	"kafkademo/internal/models"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Inventory DB simulation
var productStock = map[string]int{
	"product1": 10,
	"product2": 0, // Out of stock
}

func main() {
	// Create Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "inventory-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"order_created"}, nil)

	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			order := models.OrderEvent{}
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Failed to parse order event: %v\n", err)
				continue
			}

			// Simulate inventory check
			topic := "inventory_reserved"
			if productStock["product1"] > 0 {
				productStock["product1"]-- // Deduct from stock
				log.Printf("✅ Inventory Reserved for Order %s\n", order.OrderID)
			} else {
				topic = "inventory_failed"
				log.Printf("❌ Inventory Not Available for Order %s\n", order.OrderID)
			}

			// Publish event
			data, _ := json.Marshal(order)
			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          data,
			}, nil)

			if err != nil {
				log.Printf("Failed to send inventory event: %v\n", err)
			} else {
				log.Printf("📦 Sent inventory event: %+v\n", order)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
