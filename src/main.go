package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Budry/prometheus-simple-docker-exporter/src/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
)

var (
	namespace   = "docker"
	labels      = []string{"container", "name", "project"}
	memoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "memory_usage_bytes",
		Help:      "Current memory usage in bytes.",
	}, labels,
	)
	memoryLimit = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "memory_limit_bytes",
		Help:      "Memory limit for container in bytes.",
	}, labels,
	)
	cpuUsagePercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cpu_usage_percent",
		Help:      "Current CPU usage in percent.",
	}, labels,
	)
)

func init() {
	prometheus.MustRegister(memoryUsage)
	prometheus.MustRegister(memoryLimit)
	prometheus.MustRegister(cpuUsagePercent)
}
func main() {

	http.Handle("/metrics", promhttp.Handler())

	go func() {
			cli, err := client.NewEnvClient()
			if err != nil {
				panic(err)
			}

			ctx, cancel := context.WithCancel(context.Background())

			containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				panic(err)
			}

			for _, container := range containers {
				go func(container types.Container) {
					stats, err := cli.ContainerStats(ctx, container.ID, true)
					if err != nil {
						panic(err)
					}
					s := &types.Stats{}

					for {
						select {
						case <-ctx.Done():
							stats.Body.Close()
							fmt.Println("Stop logging")
							return
						default:
							if err := json.NewDecoder(stats.Body).Decode(&s); err == io.EOF {
								return
							} else if err != nil {
								cancel()
							}
							memoryUsage.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(float64(s.MemoryStats.Usage))
							memoryLimit.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(float64(s.MemoryStats.Limit))
							cpuUsagePercent.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(utils.CalculateCPUPercentUnix(s.PreCPUStats, s.CPUStats))
							fmt.Println(container.ID, s.CPUStats.CPUUsage.TotalUsage)
						}
					}
				}(container)

				//json.NewDecoder(stats.Body).Decode(s)

				//memoryUsage.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(float64(s.MemoryStats.Usage))
				//memoryLimit.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(float64(s.MemoryStats.Limit))
				//cpuUsagePercent.WithLabelValues(container.ID, container.Names[0], container.Labels["com.docker.compose.project"]).Set(utils.CalculateCPUPercentUnix(s.PreCPUStats, s.CPUStats))
			}
	}()

	err := http.ListenAndServe(":9101", nil)
	if err != nil {
		panic(err)
	}
}
