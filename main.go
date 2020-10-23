package main

import (
	"math/rand"
	"net/http"

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
	randGauge.Set(rand.Float64() * weight)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
