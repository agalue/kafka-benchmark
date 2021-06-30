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

	cfg := new(service.Config)
	flag.StringVar(&cfg.BootstrapServer, "b", "localhost:9092", "Kafka Boostrap Server")
	flag.IntVar(&cfg.PacketsPerSecond, "r", 1000, "Number of packets per second to generate")
	flag.IntVar(&cfg.Workers, "w", 10, "Number of workers (concurrent go-routines)")
	flag.IntVar(&cfg.StatsPort, "p", 8080, "Port for the prometheus statistics exporter")
	flag.Parse()

	if cfg.PacketsPerSecond <= 0 {
		log.Fatalln("Packet rate cannot be zero.")
	}

	producer := new(service.Producer)
	if err := producer.Init(cfg); err != nil {
		log.Fatalln(err)
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
	producer.Start(ctx)
	log.Println("Good bye!")
}
