package producer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Notification struct (imported from consumer package)
type Notification struct {
	UserID    string `json:"userId"`
	Type      string `json:"type"`     // email, sms, push
	Priority  string `json:"priority"` // high, low
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Start sends notifications to Kafka
func Start() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	notifications := []Notification{
		{"user1", "email", "high", "Your OTP is 123456", time.Now().Unix()},
		{"user2", "sms", "low", "Get 10% off on your next purchase!", time.Now().Unix()},
		{"user3", "push", "high", "Suspicious login detected!", time.Now().Unix()},
		{"user4", "email", "low", "Welcome to our service!", time.Now().Unix()},
	}

	for _, notification := range notifications {
		data, _ := json.Marshal(notification)
		topic := "low_priority_notifications"
		if notification.Priority == "high" {
			topic = "high_priority_notifications"
		}

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          data,
		}, nil)

		if err != nil {
			log.Printf("Failed to send notification: %v\n", err)
		} else {
			log.Printf("Sent notification: %+v\n", notification)
		}
	}

	p.Flush(5000)
	log.Println("All notifications sent.")
}
