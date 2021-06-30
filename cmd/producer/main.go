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
		BootstrapServer:  "localhost:9092",
		Topic:            "Test",
		PacketsPerSecond: 1000,
		Workers:          10,
		StatsPort:        8080,
	}

	flag.StringVar(&cfg.BootstrapServer, "b", cfg.BootstrapServer, "Kafka Boostrap Server")
	flag.IntVar(&cfg.PacketsPerSecond, "r", cfg.PacketsPerSecond, "Number of packets per second to generate")
	flag.IntVar(&cfg.Workers, "w", cfg.Workers, "Number of workers (concurrent go-routines)")
	flag.IntVar(&cfg.StatsPort, "p", cfg.StatsPort, "Port for the prometheus statistics exporter")
	flag.Parse()

	if cfg.PacketsPerSecond <= 0 {
		log.Fatalln("Packet rate cannot be zero.")
	}

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

	log.Printf("Sending messages to %s at target rate of %d packets per seconds across %d worker(s).", cfg.BootstrapServer, cfg.PacketsPerSecond, cfg.Workers)
	producer := &service.Producer{Config: cfg}
	producer.Start(ctx)
	log.Println("Good bye!")
}
