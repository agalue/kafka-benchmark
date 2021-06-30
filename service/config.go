package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	BootstrapServer  string
	StatsPort        int
	Workers          int
	PacketsPerSecond int
}

func (cfg *Config) TickDuration() time.Duration {
	return time.Duration((cfg.Workers * 1000000000) / cfg.PacketsPerSecond)
}

type Stats struct {
	Packets prometheus.Counter
	Errors  prometheus.Counter
}

func (s *Stats) Init(statsPort int) {
	s.Packets = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_messages",
		Help: "The total number of successfully messages sent/received",
	})
	s.Errors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_errors",
		Help: "The total number of errors received while sending messages",
	})
	prometheus.MustRegister(s.Packets, s.Errors)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "ok", "msg": "use /metrics"}`))
		})
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", statsPort), nil)
	}()
}
