type ProportionalFairnessStrategy struct{}

func (p *ProportionalFairnessStrategy) Allocate(containers []models.Container, gpus []models.GPU) error {
	// Implement proportional fairness algorithm
	// Example: Allocate GPUs based on the proportional needs of containers
	totalNeed := 0.0
	for _, container := range containers {
		totalNeed += container.GPURequest.CUDACores // Adjust based on needs
	}

	for i := range containers {
		allocated := false
		for j := range gpus {
			if gpus[j].CUDACores >= containers[i].GPURequest.CUDACores &&
				gpus[j].VRAM >= containers[i].GPURequest.VRAM &&
				gpus[j].MemoryBandwidth >= containers[i].GPURequest.MemoryBandwidth &&
				gpus[j].PowerConsumption >= containers[i].GPURequest.PowerConsumption {

				containers[i].GPURequest = &gpus[j]
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
