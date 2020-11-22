package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	randCoverage   = 65536
	multiplyWeight = 3
	divideWeight   = 2
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
		return rand.Float64() * randCoverage
	}

	gaugeRaceOnRuntime := func(count int) {
		val := randomGenerator()
		if count%10 == 0 {
			randGauge.Add(val * multiplyWeight)
		} else {
			randGauge.Sub(val / divideWeight)
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
