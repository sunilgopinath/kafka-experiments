# Kafka-Based Microservices with Avro & Schema Registry

## Overview
This project implements a **Kafka-based microservices architecture** using Avro serialization and the Confluent Schema Registry. The system follows an **event-driven approach** with multiple services communicating via Kafka topics.

## 🚀 Features
- **Kafka-based messaging** with producers and consumers.
- **Schema Registry Integration** for Avro-based message serialization.
- **Event-driven microservices** with order processing, payments, inventory, shipping, and notifications.
- **Distributed Transactions (Saga Pattern)** for order management.

---
## 🏗 Architecture
### **Microservices & Topics**
| Service              | Consumes From         | Produces To             |
|----------------------|----------------------|--------------------------|
| `order_service`      | -                    | `order_created`         |
| `payment_service`    | `order_created`      | `order_paid`            |
| `inventory_service`  | `order_created`      | `inventory_reserved`    |
| `shipping_service`   | `order_paid`         | `order_shipped`         |
| `delivery_service`   | `order_shipped`      | `order_delivered`       |
| `notification_service` | `order_*` events   | -                        |

### **Schema Registry Integration**
All events are serialized in **Avro** before being published to Kafka.
- **Schema Registry URL**: `http://localhost:8081`
- **Registered Subjects**:
  - `order_events-value`
  - `order_created-value`
  - `order_paid-value`
  - `order_shipped-value`

---
## 🛠 Setup
### **1️⃣ Start Kafka & Schema Registry**
```sh
docker-compose up -d
```
Ensure both services are running:
```sh
docker ps
```

### **2️⃣ Verify Schema Registry**
```sh
curl -X GET http://localhost:8081/subjects
```
Expected output:
```json
["order_events-value"]
```

### **3️⃣ Create Kafka Topics**
```sh
docker exec -it <KAFKA_CONTAINER_ID> kafka-topics \
    --create --topic order_created --partitions 3 --replication-factor 1 \
    --bootstrap-server localhost:9092
```
Repeat for other topics: `order_paid`, `order_shipped`, `order_delivered`, `inventory_reserved`.

---
## 📦 Services
### **Order Service** (Producer)
```sh
go run cmd/order_service/main.go
```
Publishes Avro-encoded `order_created` events.

### **Payment Service** (Consumer)
```sh
go run cmd/payment_service/main.go
```
Consumes `order_created` events and produces `order_paid`.

### **Shipping Service** (Consumer)
```sh
go run cmd/shipping_service/main.go
```
Consumes `order_paid` events and produces `order_shipped`.

---
## 📜 Avro Schema Example
### **schemas/order_schema.avsc**
```json
{
  "type": "record",
  "name": "OrderEvent",
  "fields": [
    { "name": "orderId", "type": "string" },
    { "name": "userId", "type": "string" },
    { "name": "status", "type": "string" },
    { "name": "timestamp", "type": "long" }
  ]
}
```

---
## 🏗 Internal Components
### **Kafka Consumer (Avro Decoding)**
File: `internal/messaging/consumer.go`
```go
package messaging

import (
    "log"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "github.com/linkedin/goavro/v2"
    "github.com/riferrei/srclient"
)

const schemaRegistryURL = "http://localhost:8081"

type KafkaConsumer struct {
    Consumer  *kafka.Consumer
    AvroCodec *goavro.Codec
    Schema    *srclient.Schema
    Topic     string
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
    
    schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
    schema, err := schemaRegistryClient.GetLatestSchema("order_events-value")
    if err != nil {
        log.Fatalf("Failed to fetch schema: %v", err)
    }

    avroCodec, err := goavro.NewCodec(schema.Schema())
    if err != nil {
        log.Fatalf("Failed to create Avro codec: %v", err)
    }

    return &KafkaConsumer{Consumer: c, AvroCodec: avroCodec, Schema: schema, Topic: topic}
}

func (kc *KafkaConsumer) ConsumeAvroMessage() (map[string]interface{}, error) {
    msg, err := kc.Consumer.ReadMessage(-1)
    if err != nil {
        return nil, err
    }
    
    native, _, err := kc.AvroCodec.NativeFromBinary(msg.Value[5:]) // Skip schema ID
    if err != nil {
        return nil, err
    }
    return native.(map[string]interface{}), nil
}

func (kc *KafkaConsumer) Close() {
    kc.Consumer.Close()
}
```

---
## ✅ Next Steps
- Implement **dead letter queues (DLQ)** for failed messages.
- Optimize **consumer scaling** with Kafka consumer groups.
- Implement **transactional guarantees** using Kafka transactions.

---
## 🛠 Troubleshooting
### **Check if Kafka is Running**
```sh
docker ps | grep kafka
```

### **Check Topics**
```sh
docker exec -it <KAFKA_CONTAINER_ID> kafka-topics --list --bootstrap-server localhost:9092
```

### **Check Messages in a Topic**
```sh
docker exec -it <KAFKA_CONTAINER_ID> kafka-console-consumer --topic order_created --from-beginning --bootstrap-server localhost:9092
```

---
## 🎯 Conclusion
This system enables **event-driven microservices** with **Avro serialization** and **Kafka Schema Registry**, ensuring **message consistency** and **compatibility** between services. 🚀

