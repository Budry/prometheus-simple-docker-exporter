package main

import (
	"github.com/Budry/prometheus-simple-docker-exporter/src/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {

	exporter := collectors.NewDockerStatsCollector()

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(":3000", nil))
}
