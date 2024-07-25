package metrics

import (
	"gpu-cloudsim/models"
	"gpu-cloudsim/pkg/broker"
	"sync"
)

type MetricsCollector struct {
	metrics []models.Metrics
	mu      sync.Mutex
	broker  *broker.Broker
}

func NewMetricsCollector(broker *broker.Broker) *MetricsCollector {
	return &MetricsCollector{
		metrics: []models.Metrics{},
		broker:  broker,
	}
}

func (m *MetricsCollector) CollectMetrics() models.Metrics {
	return m.broker.GetCurrentMetrics()
}

func (m *MetricsCollector) AddMetrics(metrics models.Metrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = append(m.metrics, metrics)
}

func (m *MetricsCollector) GetLatestMetrics() models.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.metrics) == 0 {
		return models.NewMetrics(0, 0, 0, 0)
	}
	return m.metrics[len(m.metrics)-1]
}
