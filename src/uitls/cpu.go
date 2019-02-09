package uitls

import "github.com/docker/docker/api/types"

func CalculateCPUPercentUnix(previousCPUStats types.CPUStats, actualCPUStates types.CPUStats) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(actualCPUStates.CPUUsage.TotalUsage) - float64(previousCPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(actualCPUStates.SystemUsage) - float64(previousCPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(actualCPUStates.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}