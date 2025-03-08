package main

import (
	"kafkademo/internal/messaging"
	"log"
)

func main() {
    consumer := messaging.NewAvroConsumer("order_created", "payment-service")
    defer consumer.Close()

    for {
        order, err := consumer.ConsumeAvroMessage()
        if err != nil {
            log.Printf("❌ Failed to consume message: %v", err)
            continue
        }
        log.Printf("💰 Processing Payment for Order: %+v", order)
    }
}