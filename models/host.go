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
