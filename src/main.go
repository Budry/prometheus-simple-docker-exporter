package main

import (
	"github.com/Budry/prometheus-simple-docker-exporter/src/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {

	prometheus.MustRegister(collectors.NewDockerStatsCollector())
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":9101", nil)
	if err != nil {
		panic(err)
	}
}
