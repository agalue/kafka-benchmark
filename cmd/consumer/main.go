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
	flag.IntVar(&cfg.PacketsPerSecond, "r", 1000, "Expected nessage generation rate")
	flag.IntVar(&cfg.Workers, "w", 10, "Number of workers part of the consumer group (concurrent go-routines)")
	flag.IntVar(&cfg.StatsPort, "p", 8081, "Port for the prometheus statistics exporter")
	flag.Parse()

	if cfg.PacketsPerSecond <= 0 {
		log.Fatalln("Packet rate cannot be zero.")
	}

	generator := new(service.Consumer)
	if err := generator.Init(cfg); err != nil {
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

	log.Printf("Receiving messages from %s across %d worker(s).", cfg.BootstrapServer, cfg.Workers)
	generator.Start(ctx)
	log.Println("Good bye!")
}
