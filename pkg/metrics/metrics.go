package metrics

import (
	"gpu-cloudsim/models"
	"math/rand"
	"sync"
)

type MetricsCollector struct {
	metrics []models.Metrics
	mu      sync.Mutex
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: []models.Metrics{},
	}
}

func (m *MetricsCollector) CollectMetrics() models.Metrics {
	// Simulate metrics collection
	return models.NewMetrics(
		rand.Float64()*100,
		rand.Float64()*100,
		rand.Float64()*100,
		rand.Float64()*100,
	)
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
