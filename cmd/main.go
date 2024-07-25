package main

import (
	"fmt"
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/broker"
	"gpu-cloudsim/pkg/orchestrator"
	"gpu-cloudsim/pkg/qos"
	"gpu-cloudsim/pkg/scheduler"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting GPU Cloudsim...")

	// Initialize logging
	logFile := "simulation.log"
	logFileHandle, err := openLogFile(logFile)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFileHandle.Close()
	log.SetOutput(logFileHandle)
	fmt.Println("Log file opened...")

	// Create scheduler
	schedulingStrategy := &scheduler.PrioritySchedulingStrategy{}

	// Create broker
	b := broker.NewBroker(schedulingStrategy)

	// Create hosts
	host1 := models.NewHost("host-1", 32, 65536)  // 32 cores, 64 GB RAM
	host2 := models.NewHost("host-2", 64, 131072) // 64 cores, 128 GB RAM

	// Add GPUs to hosts
	gpu1 := models.NewGPU("gpu-1", 3584, 224, 8192, 900, 13.4, 250)
	gpu2 := models.NewGPU("gpu-2", 4352, 272, 16384, 1200, 18.6, 300)
	host1.AddGPU(gpu1)
	host2.AddGPU(gpu2)

	b.AddHost(host1)
	b.AddHost(host2)

	fmt.Println("Brokers, hosts and GPUs have been created and assigned")

	// Create containers with different priorities and realistic resource requests
	containers := []*models.Container{
		models.NewContainer("container-1", 8000, 16384, gpu1, 1), // 8 cores, 16 GB RAM
		models.NewContainer("container-2", 4000, 8192, gpu1, 2),  // 4 cores, 8 GB RAM
		models.NewContainer("container-3", 2000, 4096, gpu2, 3),  // 2 cores, 4 GB RAM
	}

	// Create QoS monitor with realistic thresholds
	cpuThreshold := 80.0    // CPU usage threshold in percentage
	memoryThreshold := 85.0 // Memory usage threshold in percentage
	gpuThreshold := 90.0    // GPU usage threshold in percentage
	ioThreshold := 75.0     // IO usage threshold in percentage
	qosMonitor := qos.NewQoS(cpuThreshold, memoryThreshold, gpuThreshold, ioThreshold)

	// Create orchestrator
	orch := orchestrator.NewOrchestrator(b, qosMonitor)

	fmt.Println("Containers, QoS monitor and orchestrator have all been created")

	// Run orchestrator
	simulationDuration := 5 * time.Minute
	err = orch.Run(containers, simulationDuration)
	if err != nil {
		log.Fatalf("Error running orchestrator: %v", err)
	}

	fmt.Println("Orchestration finished, metrics now being collected...")

	// Simulate workload changes
	go simulateWorkloadChanges(b, orch, simulationDuration)

	// Wait for simulation to complete
	time.Sleep(simulationDuration)

	// Print final metrics and QoS status
	finalMetrics := b.GetCurrentMetrics()
	log.Printf("Final metrics: %+v\n", finalMetrics)

	isQoSMet := qosMonitor.Monitor(finalMetrics)
	if isQoSMet {
		log.Println("QoS requirements met.")
	} else {
		log.Println("QoS requirements not met.")
	}

	fmt.Println("Simulation complete. Check the log file for details.")
}

func openLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func simulateWorkloadChanges(b *broker.Broker, orch *orchestrator.Orchestrator, duration time.Duration) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	end := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			// Simulate random workload changes
			for _, host := range b.Hosts {
				for _, container := range host.Containers {
					container.CPURequest = int(float64(container.CPURequest) * (0.8 + rand.Float64()*0.4))       // +/- 20%
					container.MemoryRequest = int(float64(container.MemoryRequest) * (0.8 + rand.Float64()*0.4)) // +/- 20%
				}
			}
			log.Println("Workload changed. Triggering reallocation...")
			orch.TriggerReallocation() // Call the orchestrator's method instead
		default:
			if time.Now().After(end) {
				return
			}
			time.Sleep(time.Second)
		}
	}
}
