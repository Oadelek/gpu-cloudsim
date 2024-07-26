package qos

import (
	"fmt"
	"gpu-cloudsim/models"
)

type QoS struct {
	// QoS parameters
	cpuUsageThreshold    float64
	memoryUsageThreshold float64
	gpuUsageThreshold    float64
	ioUsageThreshold     float64
}

func NewQoS(cpuUsageThreshold, memoryUsageThreshold, gpuUsageThreshold, ioUsageThreshold float64) *QoS {
	return &QoS{
		cpuUsageThreshold:    cpuUsageThreshold,
		memoryUsageThreshold: memoryUsageThreshold,
		gpuUsageThreshold:    gpuUsageThreshold,
		ioUsageThreshold:     ioUsageThreshold,
	}
}

func (q *QoS) Monitor(metrics models.Metrics) bool {
	// Monitor QoS requirements and print the violated requirement
	if metrics.CPUUsage > q.cpuUsageThreshold {
		fmt.Printf("CPU usage violated: %f > %f\n", metrics.CPUUsage, q.cpuUsageThreshold)
		return false // QoS requirements not met
	}
	if metrics.MemoryUsage > q.memoryUsageThreshold {
		fmt.Printf("Memory usage violated: %f > %f\n", metrics.MemoryUsage, q.memoryUsageThreshold)
		return false // QoS requirements not met
	}
	if metrics.GPUUsage > q.gpuUsageThreshold {
		fmt.Printf("GPU usage violated: %f > %f\n", metrics.GPUUsage, q.gpuUsageThreshold)
		return false // QoS requirements not met
	}
	if metrics.IOUsage > q.ioUsageThreshold {
		fmt.Printf("IO usage violated: %f > %f\n", metrics.IOUsage, q.ioUsageThreshold)
		return false // QoS requirements not met
	}
	return true // QoS requirements met
}
