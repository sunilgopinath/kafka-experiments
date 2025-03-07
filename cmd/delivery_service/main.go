package main

import (
	"encoding/json"
	"kafkademo/internal/models"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// Create Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "delivery-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"order_shipped"}, nil)

	// Create Kafka producer for next stage
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	topic := "order_delivered" // Store topic name in a variable

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			order := models.OrderEvent{}
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Failed to parse order event: %v\n", err)
				continue
			}

			// Process delivery
			log.Printf("📦 Delivering Order %s\n", order.OrderID)

			// Update status and send to next topic
			order.Status = "order_delivered"
			data, _ := json.Marshal(order)

			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // ✅ Fixed
				Value:          data,
			}, nil)

			if err != nil {
				log.Printf("⚠️ Failed to send order_delivered event: %v\n", err)
			} else {
				log.Printf("✅ Order delivered: %+v\n", order)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
