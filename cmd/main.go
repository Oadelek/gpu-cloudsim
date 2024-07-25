package main

import (
	"fmt"
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/broker"
	"gpu-cloudsim/pkg/metrics"
	"gpu-cloudsim/pkg/orchestrator"
	"gpu-cloudsim/pkg/qos"
	"gpu-cloudsim/pkg/scheduler"
	"log"
	"os"
	"time"
)

func main() {
	// Initialize logging
	logFile := "simulation.log"
	logFileHandle, err := openLogFile(logFile)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFileHandle.Close()
	log.SetOutput(logFileHandle)

	// Create scheduler
	schedulingStrategy := &scheduler.PrioritySchedulingStrategy{}

	// Create broker
	b := broker.NewBroker(schedulingStrategy)

	// Create hosts
	host1 := models.NewHost("host-1", 16, 32768) // 16 cores, 32 GB RAM
	host2 := models.NewHost("host-2", 32, 65536) // 32 cores, 64 GB RAM

	// Add GPUs to hosts
	gpu1 := models.NewGPU("gpu-1", 3584, 224, 8192, 900, 13.4, 250)
	gpu2 := models.NewGPU("gpu-2", 4352, 272, 16384, 1200, 18.6, 300)
	host1.AddGPU(gpu1)
	host2.AddGPU(gpu2)

	b.AddHost(host1)
	b.AddHost(host2)

	// Create containers with different priorities
	containers := []*models.Container{
		models.NewContainer("container-1", 2000, 2048, gpu1, 1),
		models.NewContainer("container-2", 1500, 1024, gpu1, 2),
		models.NewContainer("container-3", 1000, 512, gpu2, 3),
	}

	// Create metrics collector
	metricsCollector := metrics.NewMetricsCollector()

	// Create QoS monitor
	qosMonitor := qos.NewQoS(70.0, 1500.0) // Example thresholds

	// Create orchestrator
	orch := orchestrator.NewOrchestrator(b, metricsCollector, qosMonitor)

	// Run orchestrator
	simulationDuration := 5 * time.Minute
	err = orch.Run(containers, simulationDuration)
	if err != nil {
		log.Fatalf("Error running orchestrator: %v", err)
	}

	// Print final metrics and QoS status
	finalMetrics := metricsCollector.GetLatestMetrics()
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
