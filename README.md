# kafka-experiments

## Getting started

```bash
$> git clone git@github.com:sunilgopinath/kafka-experiments.git
$> docker-compose up -d
```

### Test installation

```bash
❯ docker ps -a
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