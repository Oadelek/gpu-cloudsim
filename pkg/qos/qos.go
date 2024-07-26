package qos

import (
	"gpu-cloudsim/models"
	"log"
	"time"
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

func (q *QoS) Monitor(metrics models.Metrics, logger *log.Logger) bool {
	violations := 0

	// Monitor QoS requirements and print the violated requirement
	if metrics.CPUUsage > q.cpuUsageThreshold {
		violations++
		logger.Printf("CPU usage violated: %f > %f\n", metrics.CPUUsage, q.cpuUsageThreshold)
	}
	if metrics.MemoryUsage > q.memoryUsageThreshold {
		violations++
		logger.Printf("Memory usage violated: %f > %f\n", metrics.MemoryUsage, q.memoryUsageThreshold)
	}
	if metrics.GPUUsage > q.gpuUsageThreshold {
		violations++
		logger.Printf("GPU usage violated: %f > %f\n", metrics.GPUUsage, q.gpuUsageThreshold)
	}
	if metrics.IOUsage > q.ioUsageThreshold {
		violations++
		logger.Printf("IO usage violated: %f > %f\n", metrics.IOUsage, q.ioUsageThreshold)
	}

	if violations > 0 {
		logger.Printf("Time: %s, QoS Violations: %d\n", time.Now().Format("15:04:05"), violations)
		return false // QoS requirements not met
	}

	return true // QoS requirements met
}
