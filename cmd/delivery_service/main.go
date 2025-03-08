package main

import (
	"kafkademo/internal/messaging"
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
		log.Printf("🚚 Order Delivered: %+v\n", order)
	}
}
