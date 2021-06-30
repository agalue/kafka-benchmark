package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	config   *Config
	producer *kafka.Producer
	topic    string
}

func (gen *Producer) Init(cfg *Config) error {
	var err error
	gen.config = cfg
	gen.topic = "Test"

	kafkaCfg := &kafka.ConfigMap{"bootstrap.servers": gen.config.BootstrapServer}
	if gen.producer, err = kafka.NewProducer(kafkaCfg); err != nil {
		return fmt.Errorf("could not create producer: %v", err)
	}
	return nil
}

func (gen *Producer) Start(ctx context.Context) {
	stats := new(Stats)
	stats.Init(gen.config.StatsPort)
	wg := new(sync.WaitGroup)
	wg.Add(gen.config.Workers)
	for i := 0; i < gen.config.Workers; i++ {
		go func() {
			defer wg.Done()
			gen.startWorker(ctx, stats)
		}()
	}
	wg.Wait()
}

func (gen *Producer) startWorker(ctx context.Context, stats *Stats) {
	ticker := time.NewTicker(gen.config.TickDuration())
	delivery := make(chan kafka.Event)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			close(delivery)
			return
		case <-ticker.C:
			gen.producer.Produce(gen.generateMessage(), delivery)
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

func (gen *Producer) generateMessage() *kafka.Message {
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &gen.topic, Partition: kafka.PartitionAny},
		Value:          []byte("This is a test"),
	}
}
