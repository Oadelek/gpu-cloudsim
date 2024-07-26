package models

type Container struct {
	ID            string
	CPURequest    int // in millicores
	MemoryRequest int // in MB
	GPURequest    *GPU
	Priority      int // Priority for scheduling
}

func NewContainer(id string, cpuRequest, memoryRequest int, gpuRequest *GPU, priority int) *Container {
	return &Container{
		ID:            id,
		CPURequest:    cpuRequest,
		MemoryRequest: memoryRequest,
		GPURequest:    gpuRequest,
		Priority:      priority,
	}
}

func (c *Container) Clone() *Container {
	return &Container{
		ID:            c.ID,
		CPURequest:    c.CPURequest,
		MemoryRequest: c.MemoryRequest,
		GPURequest:    c.GPURequest.Clone(),
		Priority:      c.Priority,
	}
}
