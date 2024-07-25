package scheduler

import (
	"fmt"
	"gpu-cloudsim/models"
	"sort"
)

type PrioritySchedulingStrategy struct{}

func (p *PrioritySchedulingStrategy) Schedule(containers []*models.Container, hosts []*models.Host) error {
	// Sort containers by priority (higher priority first)
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Priority > containers[j].Priority
	})

	for _, container := range containers {
		allocated := false
		for _, host := range hosts {
			if canAllocate(container, host) {
				host.AddContainer(container)
				allocated = true
				break
			}
		}
		if !allocated {
			return fmt.Errorf("unable to allocate resources for container %s", container.ID)
		}
	}
	return nil
}

func canAllocate(container *models.Container, host *models.Host) bool {
	// Check if the host has enough resources for the container
	return host.CPUCores >= container.CPURequest &&
		host.Memory >= container.MemoryRequest &&
		len(host.GPUs) > 0 &&
		host.GPUs[0].CUDACores >= container.GPURequest.CUDACores &&
		host.GPUs[0].VRAM >= container.GPURequest.VRAM
}
