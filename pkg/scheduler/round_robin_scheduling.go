package scheduler

import (
	"fmt"
	"gpu-cloudsim/models"
)

type RoundRobinStrategy struct {
	currentHostIndex int
}

func (r *RoundRobinStrategy) Schedule(containers []*models.Container, hosts []*models.Host) error {
	for _, container := range containers {
		allocated := false
		for i := 0; i < len(hosts); i++ {
			hostIndex := (r.currentHostIndex + i) % len(hosts)
			host := hosts[hostIndex]
			if canAllocate(container, host) {
				host.AddContainer(container)
				allocated = true
				r.currentHostIndex = (hostIndex + 1) % len(hosts)
				break
			}
		}
		if !allocated {
			return fmt.Errorf("unable to allocate resources for container %s", container.ID)
		}
	}
	return nil
}
