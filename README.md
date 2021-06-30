Kafka Benchmark
====

The intention is to create a fairly simple Kafka Producer and a Kafka Consumer to analyze how fast a given Kafka cluster can be by comparing the receiving rate against the production rate.

For testing purposes, I created a docker environment to demonstrate how everything works. Still, the idea is to use the Go applications against a real cluster and collect metrics from other fields.

Currently, only one topic is modified, and the message content is a fixed text. The idea would be to introduce some complexity by creating dynamic content sent in a Protobuf message containing a timestamp when the message was created and comparing that date within the consumer to see how much it took since it was generated until it was generated was consumed.

The ultimate goal is understanding how the Kafka Brokers performance is impacted by:

* JVM Settings (specifically heap and GC)
* Linux Distribution
* Kernel Tuning (specifically network settings)
* Underlaying Hardware (CPU, Memory, Disk, Filesystem Types)