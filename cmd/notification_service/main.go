package main

import (
	"kafkademo/internal/messaging"
	"log"
)

func main() {
	consumer := messaging.NewAvroConsumer("notifications", "notification-service")
	defer consumer.Close()

	for {
		message, err := consumer.ConsumeAvroMessage()
		if err != nil {
			log.Printf("Failed to decode Avro: %v\n", err)
			continue
		}
		log.Printf("🔔 New Notification: %+v\n", message)
	}
}
