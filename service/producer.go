package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Config *Config
}

func (gen *Producer) Start(ctx context.Context) {
	stats := new(Stats)
	stats.Init(gen.Config.StatsPort)
	wg := new(sync.WaitGroup)
	wg.Add(gen.Config.Workers)
	for i := 0; i < gen.Config.Workers; i++ {
		go func() {
			defer wg.Done()
			gen.startWorker(ctx, stats)
		}()
	}
	wg.Wait()
}

func (gen *Producer) startWorker(ctx context.Context, stats *Stats) {
	kafkaCfg := &kafka.ConfigMap{"bootstrap.servers": gen.Config.BootstrapServer}
	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		log.Panicf("Cannot create producer: %v", err)
		return
	}
	ticker := time.NewTicker(gen.Config.TickDuration())
	delivery := make(chan kafka.Event)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			close(delivery)
			return
		case <-ticker.C:
			producer.Produce(gen.generateMessage(), delivery)
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
		TopicPartition: kafka.TopicPartition{Topic: &gen.Config.Topic, Partition: kafka.PartitionAny},
		Value:          []byte("This is a test"),
	}
}
