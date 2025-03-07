package main

import (
	"kafkademo/producer"
	"log"
)

func main() {
	log.Println("Starting Notification Producer...")
	producer.Start()
}
