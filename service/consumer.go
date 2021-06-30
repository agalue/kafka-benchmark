package service

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	config   *Config
	consumer *kafka.Consumer
	topic    string
}

func (gen *Consumer) Init(cfg *Config) error {
	var err error
	gen.config = cfg
	gen.topic = "Test"

	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":       gen.config.BootstrapServer,
		"group.id":                "Benchmark-Task",
		"auto.commit.interval.ms": 1000,
	}
	if gen.consumer, err = kafka.NewConsumer(kafkaCfg); err != nil {
		return fmt.Errorf("could not create consumer: %v", err)
	}
	if err := gen.consumer.Subscribe("Test", nil); err != nil {
		return fmt.Errorf("cannot subscribe to topic Test: %v", err)
	}
	return nil
}

func (gen *Consumer) Start(ctx context.Context) {
	stats := new(Stats)
	stats.Init(gen.config.StatsPort)
	for {
		select {
		case <-ctx.Done():
			gen.consumer.Close()
			return
		default:
			event := gen.consumer.Poll(100)
			switch event.(type) {
			case *kafka.Message:
				stats.Packets.Inc()
			case kafka.Error:
				stats.Errors.Inc()
			}
		}
	}
}
