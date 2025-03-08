package main

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

const schemaRegistryURL = "http://localhost:8081"

func main() {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9094",
    })
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer p.Close()

    // Schema Registry client
    srClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)

    // Fetch existing schema
    schema, err := srClient.GetSchema(1) // Use ID 1 from your registry
    if err != nil {
        log.Fatalf("Failed to fetch schema ID 1: %v", err)
    }

    // Create Avro codec
    codec, err := goavro.NewCodec(schema.Schema())
    if err != nil {
        log.Fatalf("Failed to create Avro codec: %v", err)
    }

    // Sample orders
    orders := []map[string]interface{}{
        {"orderId": "order1", "userId": "user1", "status": "order_created", "timestamp": time.Now().Unix()},
        {"orderId": "order2", "userId": "user2", "status": "order_created", "timestamp": time.Now().Unix()},
    }

    topic := "order_created"
    for _, order := range orders {
        // Encode to Avro
        avroData, err := codec.BinaryFromNative(nil, order)
        if err != nil {
            log.Printf("Failed to encode Avro: %v", err)
            continue
        }

        // Prefix with magic byte and schema ID 1
        msg := make([]byte, 5+len(avroData))
        msg[0] = 0x0 // Magic byte
        binary.BigEndian.PutUint32(msg[1:5], uint32(1)) // Schema ID 1
        copy(msg[5:], avroData)

        // Produce message
        err = p.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
            Value:          msg,
        }, nil)
        if err != nil {
            log.Printf("Failed to send order event: %v", err)
        } else {
            log.Printf("🛒 Order Created (Avro): %+v", order)
        }
    }

    p.Flush(5000)
    log.Println("All orders sent.")
}