package main

import (
	"kafkademo/internal/messaging"
	"log"
)

func main() {
	consumer := messaging.NewAvroConsumer("high_priority_notifications", "email-service")
	defer consumer.Close()

	for {
		notification, err := consumer.ConsumeAvroMessage()
		if err != nil {
			log.Printf("Failed to decode Avro: %v\n", err)
			continue
		}
		log.Printf("📧 Sending Email Notification: %+v\n", notification)
	}
}
