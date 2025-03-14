package main

import (
	"kafkademo/internal/messaging" // ✅ Updated Import Path
	"log"
)

func main() {
	consumer := messaging.NewAvroConsumer("sms_notification", "delivery-service")
	defer consumer.Close()

	for {
		order, err := consumer.ConsumeAvroMessage()
		if err != nil {
			log.Printf("Failed to decode Avro: %v\n", err)
			continue
		}
		log.Printf("📦 Sending SMS notification: %+v\n", order)
	}
}
