package service

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Config *Config
}

func (c *Consumer) Start(ctx context.Context) {
	c.Config.StartWorkers(ctx, c.startWorker)
}

func (c *Consumer) startWorker(ctx context.Context, stats *Stats) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":       c.Config.BootstrapServer,
		"group.id":                "Benchmark-Task",
		"auto.commit.interval.ms": 1000,
	}
	consumer, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		log.Panicf("Cannot create consumer: %v", err)
		return
	}
	if err := consumer.Subscribe(c.Config.Topic, nil); err != nil {
		log.Panicf("Cannot subscribe to topic %s: %v", c.Config.Topic, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			consumer.Close()
			return
		default:
			event := consumer.Poll(100)
			switch event.(type) {
			case *kafka.Message:
				stats.Packets.Inc()
			case kafka.Error:
				stats.Errors.Inc()
			}
		}
	}
}
