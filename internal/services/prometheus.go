package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// RegEx válida para nomes de métricas conforme as regras do Prometheus.
var (
	metricNameRegex = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
)

// validateMetricName assegura que o nome da métrica esteja correto.
func validateMetricName(name string) error {
	if !metricNameRegex.MatchString(name) {
		return fmt.Errorf("invalid metric name '%s': must match [a-zA-Z_:][a-zA-Z0-9_:]*", name)
	}
	return nil
}

// Metric representa uma métrica com seu valor e metadados associados.
type Metric struct {
	Value    float64           `json:"value"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// PrometheusManager gerencia todas as métricas expostas e a sua persistência.
type PrometheusManager struct {
	enabled         bool
	metrics         map[string]Metric
	mutex           sync.RWMutex
	metricsFile     string          // caminho do arquivo de persistência
	exportWhitelist map[string]bool // Se não vazio, apenas estas métricas serão exportadas para Prometheus
}

// Instância singleton do PrometheusManager.
var prometheusManagerInstance *PrometheusManager

// getMetricsFilePath define o caminho para o arquivo de métricas.
func getMetricsFilePath() string {
	if envPath := os.Getenv("LOGZ_METRICS_FILE"); envPath != "" {
		return envPath
	}
	home, err := os.UserCacheDir()
	if err != nil {
		home = "/tmp"
	}
	dir := filepath.Join(home, "kubex", "logz")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "metrics.json")
}

// loadMetrics carrega as métricas persistidas, se existirem.
func (pm *PrometheusManager) loadMetrics() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	data, err := os.ReadFile(pm.metricsFile)
	if err != nil {
		if os.IsNotExist(err) {
			pm.metrics = make(map[string]Metric)
			return nil
		}
		return err
	}
	var loaded map[string]Metric
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	pm.metrics = loaded
	return nil
}

// saveMetrics persiste as métricas atuais no arquivo definido.
func (pm *PrometheusManager) saveMetrics() error {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	data, err := json.MarshalIndent(pm.metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pm.metricsFile, data, 0644)
}

// GetPrometheusManager retorna a instância singleton do PrometheusManager.
func GetPrometheusManager() *PrometheusManager {
	if prometheusManagerInstance == nil {
		prometheusManagerInstance = &PrometheusManager{
			enabled:         false,
			metrics:         make(map[string]Metric),
			metricsFile:     getMetricsFilePath(),
			exportWhitelist: make(map[string]bool),
		}
		if err := prometheusManagerInstance.loadMetrics(); err != nil {
			fmt.Printf("Warning: could not load metrics: %v\n", err)
		}
		if os.Getenv("LOGZ_PROMETHEUS_ENABLED") == "true" {
			prometheusManagerInstance.enabled = true
		}
	}
	return prometheusManagerInstance
}

// SetExportWhitelist permite que o cliente defina quais métricas devem ser exportadas para Prometheus.
func (pm *PrometheusManager) SetExportWhitelist(metrics []string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.exportWhitelist = make(map[string]bool)
	for _, m := range metrics {
		pm.exportWhitelist[m] = true
	}
	fmt.Println("Export whitelist updated for Prometheus metrics.")
}

// Enable ativa a funcionalidade de métricas.
func (pm *PrometheusManager) Enable() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if pm.enabled {
		fmt.Println("Prometheus metrics are already enabled.")
		return
	}
	pm.enabled = true
	fmt.Println("Prometheus metrics enabled.")
}

// Disable desativa a funcionalidade Prometheus.
func (pm *PrometheusManager) Disable() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if !pm.enabled {
		fmt.Println("Prometheus metrics are already disabled.")
		return
	}
	pm.enabled = false
	fmt.Println("Prometheus metrics disabled.")
}

// IsEnabled retorna se as métricas estão habilitadas.
func (pm *PrometheusManager) IsEnabled() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.enabled
}

// AddMetric adiciona ou atualiza uma métrica com o valor e metadados fornecidos.
func (pm *PrometheusManager) AddMetric(name string, value float64, metadata map[string]string) {
	if err := validateMetricName(name); err != nil {
		fmt.Printf("Error adding metric: %v\n", err)
		return
	}
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics[name] = Metric{
		Value:    value,
		Metadata: metadata,
	}
	fmt.Printf("Metric '%s' added/updated with value: %f\n", name, value)
	if err := pm.saveMetrics(); err != nil {
		fmt.Printf("Error saving metrics: %v\n", err)
	}
}

// RemoveMetric remove uma métrica pelo nome.
func (pm *PrometheusManager) RemoveMetric(name string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.metrics, name)
	fmt.Printf("Metric '%s' removed.\n", name)
	if err := pm.saveMetrics(); err != nil {
		fmt.Printf("Error saving metrics: %v\n", err)
	}
}

// IncrementMetric incrementa a métrica especificada em delta; cria a métrica se necessário.
func (pm *PrometheusManager) IncrementMetric(name string, delta float64) {
	if err := validateMetricName(name); err != nil {
		fmt.Printf("Error incrementing metric: %v\n", err)
		return
	}
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	metric, exists := pm.metrics[name]
	if !exists {
		metric = Metric{Value: 0, Metadata: nil}
	}
	metric.Value += delta
	pm.metrics[name] = metric
	fmt.Printf("Metric '%s' incremented by %f, new value: %f\n", name, delta, metric.Value)
	if err := pm.saveMetrics(); err != nil {
		fmt.Printf("Error saving metrics: %v\n", err)
	}
}

// ListMetrics exibe todas as métricas registradas, incluindo metadados se houver.
func (pm *PrometheusManager) ListMetrics() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	if len(pm.metrics) == 0 {
		fmt.Println("No metrics registered.")
		return
	}
	fmt.Println("Registered metrics:")
	for name, metric := range pm.metrics {
		fmt.Printf("- %s: %f", name, metric.Value)
		if len(metric.Metadata) > 0 {
			metadataJSON, _ := json.Marshal(metric.Metadata)
			fmt.Printf(" (metadata: %s)", string(metadataJSON))
		}
		fmt.Println()
	}
}

// GetMetrics retorna uma cópia das métricas no formato {metric_name: value}, aplicando o filtro de exportWhitelist se definido.
func (pm *PrometheusManager) GetMetrics() map[string]float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	copiedMetrics := make(map[string]float64, len(pm.metrics))
	for k, metric := range pm.metrics {
		// Se o whitelist estiver definido (não vazio), exportar apenas os itens nele contidos.
		if len(pm.exportWhitelist) > 0 {
			if !pm.exportWhitelist[k] {
				continue
			}
		}
		copiedMetrics[k] = metric.Value
	}
	return copiedMetrics
}
