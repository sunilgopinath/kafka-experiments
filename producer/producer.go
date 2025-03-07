package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "security.protocol": "PLAINTEXT",
    })
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer p.Close()

    topic := "payments"
    payment := map[string]interface{}{"paymentId": "p1", "userId": "u1", "amount": 100}
    data, _ := json.Marshal(payment)

    deliveryChan := make(chan kafka.Event)

    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          data,
    }, deliveryChan)

    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }

    e := <-deliveryChan
    m := e.(*kafka.Message)

    if m.TopicPartition.Error != nil {
        log.Fatalf("Delivery failed: %v", m.TopicPartition.Error)
    } else {
        fmt.Printf("Message delivered to %v\n", m.TopicPartition)
    }

    p.Flush(5000)
}
