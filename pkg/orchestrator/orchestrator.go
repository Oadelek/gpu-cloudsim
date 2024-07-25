package orchestrator

import (
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/broker"
	"gpu-cloudsim/pkg/metrics"
	"gpu-cloudsim/pkg/qos"
	"log"
	"sort"
	"sync"
	"time"
)

type Orchestrator struct {
	Broker           *broker.Broker
	MetricsCollector *metrics.MetricsCollector
	QoSMonitor       *qos.QoS
}

func NewOrchestrator(broker *broker.Broker, metricsCollector *metrics.MetricsCollector, qosMonitor *qos.QoS) *Orchestrator {
	return &Orchestrator{
		Broker:           broker,
		MetricsCollector: metricsCollector,
		QoSMonitor:       qosMonitor,
	}
}

func (o *Orchestrator) Run(containers []*models.Container, duration time.Duration) error {
	err := o.Broker.AllocateResources(containers)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start metrics collection
	go func() {
		defer wg.Done()
		o.collectMetrics(duration)
	}()

	// Start QoS monitoring
	go func() {
		defer wg.Done()
		o.monitorQoS(duration)
	}()

	wg.Wait()
	return nil
}

func (o *Orchestrator) collectMetrics(duration time.Duration) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	end := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			metrics := o.MetricsCollector.CollectMetrics()
			o.MetricsCollector.AddMetrics(metrics)
		default:
			if time.Now().After(end) {
				return
			}
		}
	}
}

func (o *Orchestrator) monitorQoS(duration time.Duration) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	end := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			metrics := o.MetricsCollector.GetLatestMetrics()
			if !o.QoSMonitor.Monitor(metrics) {
				// QoS violated, trigger reallocation
				o.triggerReallocation()
			}
		default:
			if time.Now().After(end) {
				return
			}
		}
	}
}

func (o *Orchestrator) triggerReallocation() {
	log.Println("Triggering reallocation due to QoS violation")

	// Get all hosts and their current load
	hosts := o.Broker.Hosts
	hostLoads := make(map[*models.Host]float64)

	for _, host := range hosts {
		load := o.calculateHostLoad(host)
		hostLoads[host] = load
	}

	// Sort hosts by load (descending)
	sortedHosts := make([]*models.Host, len(hosts))
	copy(sortedHosts, hosts)
	sort.Slice(sortedHosts, func(i, j int) bool {
		return hostLoads[sortedHosts[i]] > hostLoads[sortedHosts[j]]
	})

	// Try to migrate containers from the most loaded hosts
	for _, sourceHost := range sortedHosts {
		if hostLoads[sourceHost] < 0.8 { // Don't migrate if load is below 80%
			break
		}

		for _, container := range sourceHost.Containers {
			if destHost := o.findSuitableHost(container, hostLoads); destHost != nil {
				o.migrateContainer(container, sourceHost, destHost)

				// Update load values
				hostLoads[sourceHost] = o.calculateHostLoad(sourceHost)
				hostLoads[destHost] = o.calculateHostLoad(destHost)

				log.Printf("Migrated container %s from host %s to host %s\n",
					container.ID, sourceHost.ID, destHost.ID)

				// Break if source host is no longer overloaded
				if hostLoads[sourceHost] < 0.8 {
					break
				}
			}
		}
	}
}

func (o *Orchestrator) calculateHostLoad(host *models.Host) float64 {
	totalCPU := float64(host.CPUCores)
	totalMemory := float64(host.Memory)
	usedCPU := 0.0
	usedMemory := 0.0

	for _, container := range host.Containers {
		usedCPU += float64(container.CPURequest)
		usedMemory += float64(container.MemoryRequest)
	}

	cpuLoad := usedCPU / totalCPU
	memoryLoad := usedMemory / totalMemory

	// Return the higher of CPU or memory load
	if cpuLoad > memoryLoad {
		return cpuLoad
	}
	return memoryLoad
}

func (o *Orchestrator) findSuitableHost(container *models.Container, hostLoads map[*models.Host]float64) *models.Host {
	for _, host := range o.Broker.Hosts {
		if hostLoads[host] >= 0.8 {
			continue // Skip overloaded hosts
		}

		if canAllocate(container, host) {
			return host
		}
	}
	return nil
}

func (o *Orchestrator) migrateContainer(container *models.Container, sourceHost, destHost *models.Host) {
	// Remove container from source host
	sourceHost.RemoveContainer(container.ID)

	// Add container to destination host
	destHost.AddContainer(container)

	// Update container's GPU if necessary
	if len(destHost.GPUs) > 0 {
		container.GPURequest = destHost.GPUs[0]
	}
}

func canAllocate(container *models.Container, host *models.Host) bool {
	// Check if the host has enough resources for the container
	return host.CPUCores >= container.CPURequest &&
		host.Memory >= container.MemoryRequest &&
		len(host.GPUs) > 0 &&
		host.GPUs[0].CUDACores >= container.GPURequest.CUDACores &&
		host.GPUs[0].VRAM >= container.GPURequest.VRAM
}
