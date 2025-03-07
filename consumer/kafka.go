package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// StartConsumer listens for notifications and processes only the given notification type (email, sms, push)
func StartConsumer(notificationType string) {
	groupID := notificationType + "-service" // Group ID for scaling

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          groupID,  // Multiple instances will be part of this group
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"high_priority_notifications", "low_priority_notifications"}, nil)

	// Handle shutdown gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigchan
		fmt.Println("\nShutting down", notificationType, "consumer...")
		c.Close()
		os.Exit(0)
	}()

	// Process messages
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
