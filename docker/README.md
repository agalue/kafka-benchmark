Docker Test Environment
====

The Docker environment exists to verify the usability of the beanchmark environment for Kafka, meaning this should server as a test environment for learning purposes.

You should build a real Kafka Cluster, and use this as a reference to configure all the necessary components to stress test Kafka and analyze how it behaves.

## Environment Variables

* `KAFKA_PARTITIONS`, the default number of partitions for new topics (defaults to `10`)
* `KAFKA_HEAP`, the JVM Heap Memory for Kafka (defaults to `4g`)
* `PRODUCER_RATE`, the number of messages per second to send to Kafka (defaults to `1000`)
* `PRODUCER_WORKERS`, the number of workers that will be producing messages (defaults to `10`)
* `CONSUMER_WORKERS`, the number of workers that will be receiving messages as part of the same consumer group (defaults to `10`). It is advised to make this value less or equal than `KAFKA_PARTITIONS`.

## Exposed Applications

* CMAK / Kafka Manager (port 9000)
* Prometheus (port 9090)
* Producer Prometheus Statistics (port 8180)
* Consumer Prometheus Statistics (port 8181)
* Grafana (port 3000)