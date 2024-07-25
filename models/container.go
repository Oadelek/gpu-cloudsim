package models

type Container struct {
	ID            string
	CPURequest    int // in millicores
	MemoryRequest int // in MB
	GPURequest    *GPU
}

func NewContainer(id string, cpuRequest, memoryRequest int, gpuRequest *GPU) *Container {
	return &Container{
		ID:            id,
		CPURequest:    cpuRequest,
		MemoryRequest: memoryRequest,
		GPURequest:    gpuRequest,
	}
}
