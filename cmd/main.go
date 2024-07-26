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
	hosts := []*models.Host{
		models.NewHost("host-1", 32, 65536),   // 32 cores, 64 GB RAM
		models.NewHost("host-2", 64, 131072),  // 64 cores, 128 GB RAM
		models.NewHost("host-3", 48, 98304),   // 48 cores, 96 GB RAM
		models.NewHost("host-4", 96, 262144),  // 96 cores, 256 GB RAM
		models.NewHost("host-5", 128, 524288), // 128 cores, 512 GB RAM
		// Add more hosts as needed
	}

	// Create GPUs
	gpus := []*models.GPU{
		models.NewGPU("gpu-1", 3584, 224, 8192, 900, 13.4, 250),
		models.NewGPU("gpu-2", 4352, 272, 16384, 1200, 18.6, 300),
		models.NewGPU("gpu-3", 5120, 320, 24576, 1500, 21.2, 350),
		models.NewGPU("gpu-4", 6144, 384, 32768, 1800, 24.6, 400),
		models.NewGPU("gpu-5", 7168, 448, 40960, 2100, 28.3, 450),
		models.NewGPU("gpu-6", 8192, 512, 49152, 2400, 32.1, 500),
		models.NewGPU("gpu-7", 3072, 192, 6144, 800, 11.3, 200),
		models.NewGPU("gpu-8", 3840, 240, 12288, 1000, 15.7, 275),
		models.NewGPU("gpu-9", 4608, 288, 20480, 1300, 19.9, 325),
		models.NewGPU("gpu-10", 5376, 336, 28672, 1600, 23.5, 375),
		// Add more GPUs as needed
	}

	// Distribute GPUs among hosts
	for i, host := range hosts {
		host.AddGPU(gpus[i*2%len(gpus)])
		host.AddGPU(gpus[(i*2+1)%len(gpus)])
		b.AddHost(host)
	}

	fmt.Println("Brokers, hosts and GPUs have been created and assigned")

	// Create containers with different priorities and realistic resource requests
	// We'll use only a few GPUs for these containers
	containers := []*models.Container{
		models.NewContainer("container-1", 8000, 16384, gpus[0], 1),  // 8 cores, 16 GB RAM, using gpu-1
		models.NewContainer("container-2", 4000, 8192, gpus[0], 2),   // 4 cores, 8 GB RAM, using gpu-1
		models.NewContainer("container-3", 2000, 4096, gpus[1], 3),   // 2 cores, 4 GB RAM, using gpu-2
		models.NewContainer("container-4", 6000, 12288, gpus[1], 1),  // 6 cores, 12 GB RAM, using gpu-2
		models.NewContainer("container-5", 3000, 6144, gpus[2], 2),   // 3 cores, 6 GB RAM, using gpu-3
		models.NewContainer("container-6", 5000, 10240, gpus[2], 3),  // 5 cores, 10 GB RAM, using gpu-3
		models.NewContainer("container-7", 7000, 14336, gpus[0], 1),  // 7 cores, 14 GB RAM, using gpu-1
		models.NewContainer("container-8", 4500, 9216, gpus[1], 2),   // 4.5 cores, 9 GB RAM, using gpu-2
		models.NewContainer("container-9", 3500, 7168, gpus[2], 3),   // 3.5 cores, 7 GB RAM, using gpu-3
		models.NewContainer("container-10", 5500, 11264, gpus[0], 1), // 5.5 cores, 11 GB RAM, using gpu-1
	}

	// Create QoS monitor with realistic thresholds
	cpuThreshold := 80.0    // CPU usage threshold in percentage
	memoryThreshold := 85.0 // Memory usage threshold in percentage
	gpuThreshold := 95.0    // GPU usage threshold in percentage
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
