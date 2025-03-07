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
		"group.id":          "payment-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"order_created"}, nil)

	// Create Kafka producer for next stage
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

			// Process payment (simulate delay)
			log.Printf("💰 Processing payment for Order %s\n", order.OrderID)

			// Update status and send to next topic
			order.Status = "order_paid"
			data, _ := json.Marshal(order)
			topic := "order_paid"

			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          data,
			}, nil)

			if err != nil {
				log.Printf("Failed to send payment event: %v\n", err)
			} else {
				log.Printf("✅ Payment Completed: %+v\n", order)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
