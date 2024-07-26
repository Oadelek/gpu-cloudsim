package scheduler

import (
	"fmt"
	"gpu-cloudsim/models"
	"sort"
)

type BinPackingStrategy struct{}

func (b *BinPackingStrategy) Allocate(containers []models.Container, gpus []models.GPU) error {
	for i := range containers {
		for j := range gpus {
			if gpus[j].CUDACores >= containers[i].GPURequest.CUDACores &&
				gpus[j].VRAM >= containers[i].GPURequest.VRAM &&
				gpus[j].MemoryBandwidth >= containers[i].GPURequest.MemoryBandwidth &&
				gpus[j].PowerConsumption >= containers[i].GPURequest.PowerConsumption {

				containers[i].GPURequest = &gpus[j]
				break
			}
		}
	}
	return nil
}

func (b *BinPackingStrategy) Schedule(containers []*models.Container, hosts []*models.Host) error {
	// Sort containers by resource requirements (descending order)
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].CPURequest > containers[j].CPURequest
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
