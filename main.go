package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	weight = 65536
)

var (
	randGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		ConstLabels: prometheus.Labels{"version": "v1.0"},
		Help:        "Random returns gauge metrics",
		Name:        "random_gauge_value",
	})
)

func init() {
	prometheus.MustRegister(randGauge)
}

func main() {

	randomGenerator := func() float64 {
		rand.Seed(time.Hour.Milliseconds())
		return rand.Float64() * weight
	}

	go func() {
		count := 0
		for {
			val := randomGenerator()
			count = count + 1
			if count%10 == 0 {
				randGauge.Add(val * 5)
			} else {
				randGauge.Sub(val / 2)
			}
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{EnableOpenMetrics: true},
	))
	http.ListenAndServe(":8080", nil)
}
