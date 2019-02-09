package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Budry/prometheus-simple-docker-exporter/src/uitls"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "docker"
	labels = []string{"container", "name", "project"}
	memoryUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,"", "memory_usage_bytes"),
		"Current memory usage",
		labels, nil,
	)
	memoryLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,"", "memory_limit_bytes"),
		"Current memory usage",
		labels, nil,
	)
	cpuUsagePercent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,"", "cpu_usage_percent"),
		"Current CPU usage in percent",
		labels, nil,
	)
)

type DockerStatsCollector struct{}

func (exporter *DockerStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- memoryUsage
	ch <- memoryLimit
	ch <- cpuUsagePercent
}

func (exporter *DockerStatsCollector) Collect(ch chan<- prometheus.Metric) {

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.Labels["com.docker.compose.project"])
		stats, err := cli.ContainerStats(ctx, container.ID, false)
		if err != nil {
			panic(err)
		}
		s := &types.Stats{}
		json.NewDecoder(stats.Body).Decode(s)

		ch <- prometheus.MustNewConstMetric(memoryUsage, prometheus.GaugeValue, float64(s.MemoryStats.Usage), container.ID, container.Names[0], container.Labels["com.docker.compose.project"])
		ch <- prometheus.MustNewConstMetric(memoryLimit, prometheus.GaugeValue, float64(s.MemoryStats.Limit), container.ID, container.Names[0], container.Labels["com.docker.compose.project"])
		ch <- prometheus.MustNewConstMetric(cpuUsagePercent, prometheus.GaugeValue, uitls.CalculateCPUPercentUnix(s.PreCPUStats, s.CPUStats), container.ID, container.Names[0], container.Labels["com.docker.compose.project"])
	}
}

func NewDockerStatsCollector() *DockerStatsCollector {
	return &DockerStatsCollector{}
}
