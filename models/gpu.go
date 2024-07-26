package models

type GPU struct {
	ID               string
	CUDACores        int
	TensorCores      int
	VRAM             int // in MB
	MemoryBandwidth  int // in GB/s
	TFLOPS           float64
	PowerConsumption int // in watts
}

func NewGPU(id string, cudaCores, tensorCores, vram, memoryBandwidth int, tflops float64, powerConsumption int) *GPU {
	return &GPU{
		ID:               id,
		CUDACores:        cudaCores,
		TensorCores:      tensorCores,
		VRAM:             vram,
		MemoryBandwidth:  memoryBandwidth,
		TFLOPS:           tflops,
		PowerConsumption: powerConsumption,
	}
}

func (g *GPU) Clone() *GPU {
	return &GPU{
		ID:               g.ID,
		CUDACores:        g.CUDACores,
		TensorCores:      g.TensorCores,
		VRAM:             g.VRAM,
		MemoryBandwidth:  g.MemoryBandwidth,
		TFLOPS:           g.TFLOPS,
		PowerConsumption: g.PowerConsumption,
	}
}
