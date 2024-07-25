package qos

import (
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
	// Monitor QoS requirements
	if metrics.CPUUsage > q.cpuUsageThreshold ||
		metrics.MemoryUsage > q.memoryUsageThreshold ||
		metrics.GPUUsage > q.gpuUsageThreshold ||
		metrics.IOUsage > q.ioUsageThreshold {
		return false // QoS requirements not met
	}
	return true // QoS requirements met
}
