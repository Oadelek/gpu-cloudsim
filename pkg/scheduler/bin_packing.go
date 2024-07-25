package scheduler

import (
	"gpu-cloudsim/models"
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
