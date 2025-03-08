package messaging

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

const schemaRegistryURL = "http://localhost:8081"

type KafkaConsumer struct {
    Consumer *kafka.Consumer
    Topic    string
}

func NewAvroConsumer(topic, groupID string) *KafkaConsumer {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9094",
        "group.id":          groupID,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        log.Fatalf("Failed to create consumer: %v", err)
    }

    c.SubscribeTopics([]string{topic}, nil)
    return &KafkaConsumer{Consumer: c, Topic: topic}
}

func (kc *KafkaConsumer) ConsumeAvroMessage() (map[string]interface{}, error) {
    msg, err := kc.Consumer.ReadMessage(-1)
    if err != nil {
        return nil, fmt.Errorf("consumer error: %w", err)
    }

    if len(msg.Value) < 5 || msg.Value[0] != 0x0 {
        return nil, fmt.Errorf("invalid Avro message format")
    }

    // Extract schema ID
    schemaID := int(binary.BigEndian.Uint32(msg.Value[1:5]))
    log.Printf("Received schema ID: %d", schemaID) // Debug log

    // Fetch schema
    srClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
    schema, err := srClient.GetSchema(schemaID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch schema ID %d: %w", schemaID, err)
    }

    // Create codec
    codec, err := goavro.NewCodec(schema.Schema())
    if err != nil {
        return nil, fmt.Errorf("failed to create Avro codec: %w", err)
    }

    // Decode Avro
    native, _, err := codec.NativeFromBinary(msg.Value[5:])
    if err != nil {
        return nil, fmt.Errorf("failed to decode Avro: %w", err)
    }

    return native.(map[string]interface{}), nil
}

func (kc *KafkaConsumer) Close() {
    kc.Consumer.Close()
}