package main

import (
	"kafkademo/consumer"
	"log"
)

func main() {
	log.Println("Starting Push Notification Consumer...")
	consumer.StartConsumer("push")
}
