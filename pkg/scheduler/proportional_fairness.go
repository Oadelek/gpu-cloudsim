package scheduler

import (
	"fmt"
	"gpu-cloudsim/models"
)

type ProportionalFairnessStrategy struct{}

func (p *ProportionalFairnessStrategy) Allocate(containers []models.Container, gpus []models.GPU) error {
	// Implement proportional fairness algorithm
	// Example: Allocate GPUs based on the proportional needs of containers
	totalNeed := 0
	for _, container := range containers {
		totalNeed += container.GPURequest.CUDACores // Adjust based on needs
	}

	totalGPUCores := 0
	for _, gpu := range gpus {
		totalGPUCores += gpu.CUDACores
	}

	for i := range containers {
		allocated := false
		requiredCores := int(float64(containers[i].GPURequest.CUDACores) / float64(totalNeed) * float64(totalGPUCores))

		for j := range gpus {
			if gpus[j].CUDACores >= requiredCores &&
				gpus[j].VRAM >= containers[i].GPURequest.VRAM &&
				gpus[j].MemoryBandwidth >= containers[i].GPURequest.MemoryBandwidth &&
				gpus[j].PowerConsumption >= containers[i].GPURequest.PowerConsumption {

				// Allocate GPU to container
				containers[i].GPURequest = &gpus[j]

				// Update remaining GPU resources
				gpus[j].CUDACores -= requiredCores
				gpus[j].VRAM -= containers[i].GPURequest.VRAM
				gpus[j].MemoryBandwidth -= containers[i].GPURequest.MemoryBandwidth
				gpus[j].PowerConsumption -= containers[i].GPURequest.PowerConsumption

				allocated = true
				break
			}
		}
		if !allocated {
			return fmt.Errorf("unable to allocate resources for container %s", containers[i].ID)
		}
	}
	return nil
}
