# **Event-Driven Microservices with Kafka**

## **🚀 Overview**
This project implements an **event-driven microservices architecture** using **Apache Kafka**. We built a **scalable E-commerce Order Processing System** with **Event Sourcing**, where each microservice publishes and consumes **domain events** asynchronously.

## **🛠️ Technologies Used**
- **Go (Golang)** - Microservice implementation
- **Apache Kafka** - Event streaming platform
- **Docker & Docker Compose** - Service orchestration
- **PostgreSQL (Optional)** - Database for persistence
- **Kafka Consumer Groups** - Parallel processing
- **Exponential Backoff & Retry Mechanisms** - Fault tolerance
- **Dead Letter Queue (DLQ)** - Handling failures gracefully
  
## **📌 Architecture**

```
USER → Order Service → order_created → Kafka
              ↳ Payment Service → order_paid → Kafka
              ↳ Inventory Service → inventory_reserved/inventory_failed → Kafka
              ↳ Shipping Service → order_shipped → Kafka
              ↳ Delivery Service → order_delivered → Kafka
              ↳ Notification Service → Sends user updates
```

- **Each service listens for relevant Kafka topics.**
- **Microservices communicate via domain events (not direct API calls).**
- **Kafka retains a history of events (Event Sourcing).**
- **Event Replay allows restoring state from Kafka topics.**

## **📂 Project Structure**
```
.
├── cmd
│   ├── order_service        # Handles order creation
│   │   └── main.go
│   ├── payment_service      # Handles payment processing
│   │   └── main.go
│   ├── shipping_service     # Handles shipping
│   │   └── main.go
│   ├── delivery_service     # Handles order delivery
│   │   └── main.go
│   ├── inventory_service    # Manages stock availability
│   │   └── main.go
│   ├── notification_service # Sends user notifications
│   │   └── main.go
│
├── internal
│   ├── models               # Shared domain models
│   │   └── order.go
│   ├── kafka                # Kafka utility functions
│   │   ├── producer.go
│   │   ├── consumer.go
│
├── docker-compose.yml        # Kafka + Zookeeper setup
├── go.mod                    # Go dependencies
├── go.sum                    # Go dependency hashes
└── README.md                 # Project documentation
```

## **📝 Order Event Model** (Shared Across Microservices)
```go
package models

type OrderEvent struct {
    OrderID   string `json:"orderId"`
    UserID    string `json:"userId"`
    Status    string `json:"status"`
    Timestamp int64  `json:"timestamp"`
}
```

---

# **🚀 Getting Started**

```bash
$> git clone git@github.com:sunilgopinath/kafka-experiments.git
```

## **1️⃣ Setup Kafka with Docker**
Run the following command to start Kafka & Zookeeper:
```sh
docker-compose up -d
```

✅ **Verify Kafka is running:**
```sh
docker ps | grep kafka
CONTAINER ID   IMAGE                 COMMAND                  CREATED          STATUS          PORTS                    NAMES
5a992aa5f7bf   bitnami/kafka:3.6.1   "/opt/bitnami/script…"   15 minutes ago   Up 15 minutes   0.0.0.0:9092->9092/tcp   kafka-experiments-kafka-1
```


## Create topic

```bash
$> docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic payments --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
```

### Test topic creation

```bash
❯ docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --list --bootstrap-server localhost:9092
payments
```

## Usage

terminal 1

```bash
❯ go run producer/producer.go

Message delivered to payments[1]@0
```


terminal 2

```bash
❯ go run consumer/consumer.go
Received: {"amount":100,"paymentId":"p1","userId":"u1"}
```

## Write to database upon consumer

Please see [database documentation](./docs/db.md) for setup

## Writing to Dead Letter Queue

```bash
docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic dead_letter_queue --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
```

### Test DLQ

```bash
docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --list --bootstrap-server localhost:9092
```

## Notification System

Create two queues
1. high priority
2. low priority

```bash
docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic high_priority_notifications --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092

docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic low_priority_notifications --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
```

## Order Service

### **2️⃣ Create Kafka Topics**
```sh
docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic order_created --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092

docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic order_paid --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092

docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic order_shipped --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092

docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --create --topic order_delivered --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
```
✅ **Verify topics:**
```sh
docker exec -it $(docker ps -q -f name=kafka) kafka-topics.sh --list --bootstrap-server localhost:9092
```

---

# **🛠️ Running the Microservices**

## **Start All Services in Separate Terminals**
```sh
go run cmd/order_service/main.go
```
```sh
go run cmd/payment_service/main.go
```
```sh
go run cmd/shipping_service/main.go
```
```sh
go run cmd/delivery_service/main.go
```
```sh
go run cmd/inventory_service/main.go
```
```sh
go run cmd/notification_service/main.go
```

✅ **Expected Output:**
```
🛒 Order Created: order1
💰 Payment Processed: order1
✅ Inventory Reserved for Order order1
📦 Shipping Order: order1
🚚 Order Delivered: order1
📩 Notification sent: Order shipped
```

---

# **📌 Event Replay (Event Sourcing)**
If a service **crashes**, we can replay events to rebuild its state.
```sh
docker exec -it $(docker ps -q -f name=kafka) kafka-console-consumer.sh --topic order_created --from-beginning --bootstrap-server localhost:9092
```
✅ **Expected Output:**
```
{"orderId":"order1","userId":"user1","status":"order_created","timestamp":1711345827}
```

---

# **📊 Monitoring & Debugging**
✅ **View consumer group assignments:**
```sh
docker exec -it $(docker ps -q -f name=kafka) kafka-consumer-groups.sh --describe --group order-service --bootstrap-server localhost:9092
```

✅ **View Dead Letter Queue (DLQ) events:**
```sh
docker exec -it $(docker ps -q -f name=kafka) kafka-console-consumer.sh --topic dead_letter_queue --from-beginning --bootstrap-server localhost:9092
```

---

# **🛠️ Next Steps**
🔹 **Add Retries & Exponential Backoff for Failures**
🔹 **Implement DLQ Handling for Fault-Tolerant Processing**
🔹 **Optimize Performance with Kafka Streams & Schema Registry**

🚀 **Happy Coding!** 🎯