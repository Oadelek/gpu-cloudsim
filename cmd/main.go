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

const (
	simulationDuration     = 5 * time.Minute
	workloadChangeInterval = 30 * time.Second
	numHosts               = 100 // Number of hosts
	numGPUs                = 200 // Number of GPUs
	numContainers          = 500 // Number of containers
)

func main() {
	fmt.Println("Starting GPU Cloudsim...")

	hosts := createHosts()
	gpus := createGPUs()
	containers := createContainers(gpus)
	qosMonitor := createQoSMonitor()

	// Define scheduling strategies
	strategies := map[string]scheduler.Scheduler{
		"Priority":   &scheduler.PrioritySchedulingStrategy{},
		"BinPacking": &scheduler.BinPackingStrategy{},
		"RoundRobin": &scheduler.RoundRobinStrategy{},
	}

	// Run simulation for each strategy
	for name, strategy := range strategies {
		fmt.Printf("Running simulation with %s strategy...\n", name)
		runSimulation(name, strategy, hosts, containers, qosMonitor, gpus)
	}

	fmt.Println("All simulations complete. Check the log files for details.")
}

func runSimulation(name string, strategy scheduler.Scheduler, hosts []*models.Host, containers []*models.Container, qosMonitor *qos.QoS, gpus []*models.GPU) {
	logFileName := fmt.Sprintf("%s_simulation.log", name)
	logger := setupLogger(logFileName)

	b := broker.NewBroker(strategy)

	// Distribute GPUs among hosts
	for i, host := range hosts {
		host.AddGPU(gpus[i*2%len(gpus)])
		host.AddGPU(gpus[(i*2+1)%len(gpus)])
		b.AddHost(host)
	}

	orch := orchestrator.NewOrchestrator(b, qosMonitor, logger)

	logger.Printf("Starting %s simulation\n", name)

	err := orch.Run(containers, simulationDuration)
	if err != nil {
		logger.Printf("Error running orchestrator: %v", err)
		return
	}

	// Simulate workload changes
	go simulateWorkloadChanges(orch, simulationDuration, logger)

	// Wait for simulation to complete
	time.Sleep(simulationDuration)

	// Print final metrics and QoS status
	finalMetrics := b.GetCurrentMetrics()
	logger.Printf("Final metrics: %+v\n", finalMetrics)

	isQoSMet := qosMonitor.Monitor(finalMetrics, logger)
	if isQoSMet {
		logger.Println("QoS requirements met.")
	} else {
		logger.Println("QoS requirements not met.")
	}

	logger.Printf("%s simulation complete.\n", name)
}

func simulateWorkloadChanges(orch *orchestrator.Orchestrator, duration time.Duration, logger *log.Logger) {
	ticker := time.NewTicker(workloadChangeInterval)
	defer ticker.Stop()

	end := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			// Simulate random workload changes
			for _, host := range orch.Broker.Hosts {
				for _, container := range host.Containers {
					container.CPURequest = int(float64(container.CPURequest) * (0.8 + rand.Float64()*0.4))       // +/- 20%
					container.MemoryRequest = int(float64(container.MemoryRequest) * (0.8 + rand.Float64()*0.4)) // +/- 20%
				}
			}
			logger.Println("Workload changed. Triggering reallocation...")
			orch.TriggerReallocation()
		default:
			if time.Now().After(end) {
				return
			}
			time.Sleep(time.Second)
		}
	}
}

func createHosts() []*models.Host {
	hosts := make([]*models.Host, numHosts)
	for i := 0; i < numHosts; i++ {
		hosts[i] = models.NewHost(fmt.Sprintf("host-%d", i+1), 32+rand.Intn(128), 65536+rand.Intn(524288))
	}
	return hosts
}

func createGPUs() []*models.GPU {
	gpus := make([]*models.GPU, numGPUs)
	for i := 0; i < numGPUs; i++ {
		gpus[i] = models.NewGPU(fmt.Sprintf("gpu-%d", i+1), 3584+rand.Intn(8192), 224+rand.Intn(512), 8192+rand.Intn(49152), 900+rand.Intn(2400), 13.4+rand.Float64()*18.7, 250+rand.Intn(250))
	}
	return gpus
}

func createContainers(gpus []*models.GPU) []*models.Container {
	containers := make([]*models.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = models.NewContainer(fmt.Sprintf("container-%d", i+1), 1000+rand.Intn(16000), 2048+rand.Intn(32768), gpus[rand.Intn(len(gpus))], 1+rand.Intn(3))
	}
	return containers
}

func createQoSMonitor() *qos.QoS {
	cpuThreshold := 80.0    // CPU usage threshold in percentage
	memoryThreshold := 85.0 // Memory usage threshold in percentage
	gpuThreshold := 95.0    // GPU usage threshold in percentage
	ioThreshold := 75.0     // IO usage threshold in percentage
	return qos.NewQoS(cpuThreshold, memoryThreshold, gpuThreshold, ioThreshold)
}

func setupLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file %s: %v", filename, err)
	}
	return log.New(file, "", log.LstdFlags)
}
