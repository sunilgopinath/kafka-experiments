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
		"group.id":          "shipping-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"order_paid"}, nil)

	// Create Kafka producer for next stage
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	topic := "order_shipped" // Store topic name in a variable

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			order := models.OrderEvent{}
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Failed to parse order event: %v\n", err)
				continue
			}

			// Process shipping
			log.Printf("📦 Shipping Order %s\n", order.OrderID)

			// Update status and send to next topic
			order.Status = "order_shipped"
			data, _ := json.Marshal(order)

			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // Fixed!
				Value:          data,
			}, nil)

			if err != nil {
				log.Printf("Failed to send shipping event: %v\n", err)
			} else {
				log.Printf("🚚 Order Shipped: %+v\n", order)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
