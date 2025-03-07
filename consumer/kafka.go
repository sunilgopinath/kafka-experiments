package consumer

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// StartConsumer listens for notifications and processes only the given notification type (email, sms, push)
func StartConsumer(notificationType string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          notificationType + "-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"high_priority_notifications", "low_priority_notifications"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var notification Notification
			if err := json.Unmarshal(msg.Value, &notification); err != nil {
				log.Printf("Failed to parse notification: %v\n", err)
				continue
			}

			// Process only messages of the correct type
			if notification.Type == notificationType {
				log.Printf("🔔 Processing %s notification for user %s: %s\n", notification.Type, notification.UserID, notification.Message)
			}
		} else {
			log.Printf("Consumer error: %v\n", err)
			break
		}
	}
}
