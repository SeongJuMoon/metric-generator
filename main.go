package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	randCoverage = 65536
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
		rand.Seed(time.Now().UnixNano())
		return rand.Float64() * randCoverage
	}

	gaugeRaceOnRuntime := func(count int) {
		if count%10 == 0 {
			val := randomGenerator()
			randGauge.Set(val)
		}
		time.Sleep(time.Second)
	}

	go func() {
		count := 0
		for {
			if count > 100 {
				count = 0
			}
			gaugeRaceOnRuntime(count)
			count = count + 1
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{EnableOpenMetrics: true},
	))
	http.ListenAndServe(":8080", nil)
}
