package services

import (
	"fmt"
	"sync"
)

// PrometheusManager gerencia todas as métricas expostas para o Prometheus.
type PrometheusManager struct {
	enabled bool
	metrics map[string]float64
	mutex   sync.RWMutex
}

// Instância singleton do PrometheusManager.
var prometheusManagerInstance *PrometheusManager

// GetPrometheusManager retorna a instância singleton do gerenciador de métricas.
func GetPrometheusManager() *PrometheusManager {
	if prometheusManagerInstance == nil {
		prometheusManagerInstance = &PrometheusManager{
			enabled: false,
			metrics: make(map[string]float64),
		}
	}
	return prometheusManagerInstance
}

// Enable ativa a funcionalidade Prometheus se configurada via variável de ambiente.
func (pm *PrometheusManager) Enable() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.enabled = true
	fmt.Println("Prometheus metrics enabled.")
}

// Disable desativa a funcionalidade Prometheus.
func (pm *PrometheusManager) Disable() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.enabled = false
	fmt.Println("Prometheus metrics disabled.")
}

// IsEnabled verifica se o Prometheus está habilitado.
func (pm *PrometheusManager) IsEnabled() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.enabled
}

// AddMetric adiciona ou atualiza uma métrica.
func (pm *PrometheusManager) AddMetric(name string, value float64) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics[name] = value
	fmt.Printf("Metric '%s' added/updated with value: %f\n", name, value)
}

// RemoveMetric remove uma métrica pelo nome.
func (pm *PrometheusManager) RemoveMetric(name string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.metrics, name)
	fmt.Printf("Metric '%s' removed.\n", name)
}

// ListMetrics exibe todas as métricas registradas.
func (pm *PrometheusManager) ListMetrics() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	if len(pm.metrics) == 0 {
		fmt.Println("No metrics registered.")
		return
	}
	fmt.Println("Registered metrics:")
	for name, value := range pm.metrics {
		fmt.Printf("- %s: %f\n", name, value)
	}
}

// GetMetrics retorna todas as métricas no formato esperado pelo Prometheus.
func (pm *PrometheusManager) GetMetrics() map[string]float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	copiedMetrics := make(map[string]float64, len(pm.metrics))
	for k, v := range pm.metrics {
		copiedMetrics[k] = v
	}
	return copiedMetrics
}
