package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/agalue/kafka-benchmark/service"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg := &service.Config{
		BootstrapServer: "localhost:9092",
		Topic:           "Test",
		Workers:         10,
		StatsPort:       8081,
	}

	flag.StringVar(&cfg.BootstrapServer, "b", cfg.BootstrapServer, "Kafka Boostrap Server")
	flag.IntVar(&cfg.Workers, "w", cfg.Workers, "Number of workers part of the consumer group (concurrent go-routines)")
	flag.IntVar(&cfg.StatsPort, "p", cfg.StatsPort, "Port for the prometheus statistics exporter")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
	}()
	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	log.Printf("Receiving messages from %s across %d worker(s).", cfg.BootstrapServer, cfg.Workers)
	consumer := &service.Consumer{Config: cfg}
	consumer.Start(ctx)
	log.Println("Good bye!")
}
