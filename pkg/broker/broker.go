package broker

import (
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/scheduler"
)

type Broker struct {
	Hosts     []*models.Host
	Scheduler scheduler.Scheduler
}

func NewBroker(scheduler scheduler.Scheduler) *Broker {
	return &Broker{
		Hosts:     []*models.Host{},
		Scheduler: scheduler,
	}
}

func (b *Broker) AddHost(host *models.Host) {
	b.Hosts = append(b.Hosts, host)
}

func (b *Broker) AllocateResources(containers []*models.Container) error {
	return b.Scheduler.Schedule(containers, b.Hosts)
}

func (b *Broker) GetCurrentMetrics() models.Metrics {
	var cpuUsage, memoryUsage, gpuUsage, ioUsage float64

	for _, host := range b.Hosts {
		cpuUsage += host.GetCPUUsage()
		memoryUsage += host.GetMemoryUsage()
		gpuUsage += host.GetGPUUsage()
		ioUsage += host.GetIOUsage()
	}

	totalHosts := float64(len(b.Hosts))
	return models.NewMetrics(
		cpuUsage/totalHosts,
		memoryUsage/totalHosts,
		gpuUsage/totalHosts,
		ioUsage/totalHosts,
	)
}
