package main

import (
	"kafkademo/internal/messaging" // ✅ Updated Import Path
	"log"
)

func main() {
	consumer := messaging.NewAvroConsumer("order_shipped", "delivery-service")
	defer consumer.Close()

	for {
		order, err := consumer.ConsumeAvroMessage()
		if err != nil {
			log.Printf("Failed to decode Avro: %v\n", err)
			continue
		}
		log.Printf("🚚 Order Shipped: %+v\n", order)
	}
}
