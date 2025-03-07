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
		"group.id":          "notification-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"order_paid", "order_shipped"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			order := models.OrderEvent{}
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Failed to parse order event: %v\n", err)
				continue
			}

			// Send notification
			if order.Status == "order_paid" {
				log.Printf("📩 Sending Payment Confirmation to user %s for Order %s\n", order.UserID, order.OrderID)
			} else if order.Status == "order_shipped" {
				log.Printf("🚚 Sending Shipping Update to user %s for Order %s\n", order.UserID, order.OrderID)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
