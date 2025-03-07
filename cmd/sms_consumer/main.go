package main

import (
	"kafkademo/consumer"
	"log"
)

func main() {
	log.Println("Starting SMS Consumer...")
	consumer.StartConsumer("sms")
}
