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
	// Check CPU cores (convert millicores to cores)
	if float64(host.CPUCores) < float64(container.CPURequest)/1000 {
		return false
	}

	// Check Memory
	if host.Memory < container.MemoryRequest {
		return false
	}

	// Check if any GPU on the host meets the requirements
	for _, gpu := range host.GPUs {
		if gpu.CUDACores >= container.GPURequest.CUDACores &&
			gpu.VRAM >= container.GPURequest.VRAM &&
			gpu.MemoryBandwidth >= container.GPURequest.MemoryBandwidth {
			return true // Found a suitable GPU
		}
	}

	return false
}
