package orchestrator

import (
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/broker"
	"gpu-cloudsim/pkg/metrics"
	"gpu-cloudsim/pkg/qos"
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
	// Implement reallocation logic here
	// This could involve migrating containers between hosts
}
