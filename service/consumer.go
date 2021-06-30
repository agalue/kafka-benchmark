package service

import (
	"context"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Config *Config
}

func (gen *Consumer) Start(ctx context.Context) {
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

func (gen *Consumer) startWorker(ctx context.Context, stats *Stats) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":       gen.Config.BootstrapServer,
		"group.id":                "Benchmark-Task",
		"auto.commit.interval.ms": 1000,
	}
	consumer, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		log.Panicf("Cannot create consumer: %v", err)
		return
	}
	if err := consumer.Subscribe(gen.Config.Topic, nil); err != nil {
		log.Panicf("Cannot subscribe to topic %s: %v", gen.Config.Topic, err)
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
