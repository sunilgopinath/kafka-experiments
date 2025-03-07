package main

import (
	"kafkademo/consumer"
	"log"
)

func main() {
	log.Println("Starting Email Consumer...")
	consumer.StartConsumer("email")
}
