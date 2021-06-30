package service

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Config *Config
}

func (p *Producer) Start(ctx context.Context) {
	p.Config.StartWorkers(ctx, p.startWorker)
}

func (p *Producer) startWorker(ctx context.Context, stats *Stats) {
	kafkaCfg := &kafka.ConfigMap{"bootstrap.servers": p.Config.BootstrapServer}
	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		log.Panicf("Cannot create producer: %v", err)
		return
	}
	ticker := time.NewTicker(p.Config.TickDuration())
	delivery := make(chan kafka.Event)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			close(delivery)
			return
		case <-ticker.C:
			producer.Produce(p.generateMessage(), delivery)
		case e := <-delivery:
			m := e.(*kafka.Message)
			if m.TopicPartition.Error == nil {
				stats.Packets.Inc()
			} else {
				stats.Errors.Inc()
			}
		}
	}
}

func (p *Producer) generateMessage() *kafka.Message {
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Config.Topic, Partition: kafka.PartitionAny},
		Value:          []byte("This is a test"),
	}
}
