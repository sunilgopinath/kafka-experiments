package main

import (
	"kafkademo/internal/messaging" // ✅ Updated Import Path
	"log"
)

func main() {
	consumer := messaging.NewAvroConsumer("push_notifications", "push-service")
	defer consumer.Close()

	for {
		message, err := consumer.ConsumeAvroMessage()
		if err != nil {
			log.Printf("Failed to decode Avro: %v\n", err)
			continue
		}
		log.Printf("📲 Sending Push Notification: %+v\n", message)
	}
}
