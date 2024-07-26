package models

type Host struct {
	ID         string
	Containers []*Container
	GPUs       []*GPU
	CPUCores   int
	Memory     int // in MB
}

func NewHost(id string, cpuCores, memory int) *Host {
	return &Host{
		ID:         id,
		Containers: []*Container{},
		GPUs:       []*GPU{},
		CPUCores:   cpuCores,
		Memory:     memory,
	}
}

func (h *Host) Clone() *Host {
	clonedHost := &Host{
		ID:         h.ID,
		CPUCores:   h.CPUCores,
		Memory:     h.Memory,
		GPUs:       make([]*GPU, len(h.GPUs)),
		Containers: make([]*Container, len(h.Containers)),
	}

	// Deep copy GPUs
	for i, gpu := range h.GPUs {
		clonedHost.GPUs[i] = gpu.Clone()
	}

	// Deep copy Containers
	for i, container := range h.Containers {
		clonedHost.Containers[i] = container.Clone()
	}

	return clonedHost
}

func (h *Host) AddContainer(c *Container) {
	h.Containers = append(h.Containers, c)
}

func (h *Host) AddGPU(g *GPU) {
	h.GPUs = append(h.GPUs, g)
}

func (h *Host) RemoveContainer(containerID string) {
	for i, c := range h.Containers {
		if c.ID == containerID {
			h.Containers = append(h.Containers[:i], h.Containers[i+1:]...)
			break
		}
	}
}

func (h *Host) GetCPUUsage() float64 {
	var totalUsage float64
	for _, container := range h.Containers {
		totalUsage += float64(container.CPURequest) / 1000 // Convert millicores to cores
	}
	return (totalUsage / float64(h.CPUCores)) * 100 // Return as percentage
}

func (h *Host) GetMemoryUsage() float64 {
	var totalUsage int
	for _, container := range h.Containers {
		totalUsage += container.MemoryRequest
	}
	return (float64(totalUsage) / float64(h.Memory)) * 100 // Return as percentage
}

func (h *Host) GetGPUUsage() float64 {
	if len(h.GPUs) == 0 {
		return 0
	}

	totalGPUCores := 0
	for _, gpu := range h.GPUs {
		totalGPUCores += gpu.CUDACores
	}

	usedGPUCores := 0
	for _, container := range h.Containers {
		usedGPUCores += container.GPURequest.CUDACores
	}

	return (float64(usedGPUCores) / float64(totalGPUCores)) * 100
}

func (h *Host) GetIOUsage() float64 {
	// Implementing I/O usage might require additional tracking mechanisms
	// This is a placeholder implementation
	return 50.0 // Return a constant value for now
}
