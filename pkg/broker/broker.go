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
