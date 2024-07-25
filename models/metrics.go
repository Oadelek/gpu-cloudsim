package models

type Metrics struct {
	CPUUsage    float64
	MemoryUsage float64
	GPUUsage    float64
	IOUsage     float64
}

func NewMetrics(cpuUsage, memoryUsage, gpuUsage, ioUsage float64) Metrics {
	return Metrics{
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		GPUUsage:    gpuUsage,
		IOUsage:     ioUsage,
	}
}
