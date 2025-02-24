package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/exec"
	"strconv"
)

var (
	infoCount    = prometheus.NewCounter(prometheus.CounterOpts{Name: "log_info_count", Help: "Number of info logs"})
	warnCount    = prometheus.NewCounter(prometheus.CounterOpts{Name: "log_warn_count", Help: "Number of warn logs"})
	errorCount   = prometheus.NewCounter(prometheus.CounterOpts{Name: "log_error_count", Help: "Number of error logs"})
	debugCount   = prometheus.NewCounter(prometheus.CounterOpts{Name: "log_debug_count", Help: "Number of debug logs"})
	successCount = prometheus.NewCounter(prometheus.CounterOpts{Name: "log_success_count", Help: "Number of success logs"})
)

func updateMetrics() {
	infoCount.Inc()
	warnCount.Inc()
	errorCount.Inc()
	debugCount.Inc()
	successCount.Inc()
}

func initPrometheuz() error {
	setPrometheuzConfigErr := setPrometheuzConfig()
	if setPrometheuzConfigErr != nil {
		return setPrometheuzConfigErr
	}
	prometheus.MustRegister(infoCount, warnCount, errorCount, debugCount, successCount)
	return nil
}

func setPrometheuzConfig() error {
	prometheusConfig := `
scrape_configs:
  - job_name: 'logz'
	static_configs:
	  - targets: ['localhost:2112']
`
	// Verifica se o arquivo de configuração do Prometheus existe.
	ptmCheckConfigExistsCmd := exec.Command("type", "-f", "/etc/prometheus/prometheus.yml")
	ptmCheckConfigExistsErr := ptmCheckConfigExistsCmd.Run()
	if ptmCheckConfigExistsErr != nil {
		// Se não existir, cria o arquivo de configuração do Prometheus e insere as configurações necessárias.
		ptmCreateConfigCmd := exec.Command("echo", prometheusConfig, ">", "/etc/prometheus/prometheus.yml")
		ptmCreateConfigErr := ptmCreateConfigCmd.Run()
		if ptmCreateConfigErr != nil {
			return ptmCreateConfigErr
		}
	} else {
		// Se existir, verifica se as configurações necessárias estão presentes.
		ptmCheckConfigCmd := exec.Command("grep", "logz", "/etc/prometheus/prometheus.yml")
		ptmCheckConfigErr := ptmCheckConfigCmd.Run()
		if ptmCheckConfigErr != nil {
			// Se não estiverem, insere as configurações necessárias.
			ptmInsertConfigCmd := exec.Command("echo", prometheusConfig, ">>", "/etc/prometheus/prometheus.yml")
			ptmInsertConfigErr := ptmInsertConfigCmd.Run()
			if ptmInsertConfigErr != nil {
				return ptmInsertConfigErr
			}
		}
	}
	return nil
}

func Prometheuz(route string, port int) error {
	if route == "" {
		route = "metrics"
	}
	if port == 0 {
		port = 2112
	}
	initPrometheuzErr := initPrometheuz()
	if initPrometheuzErr != nil {
		return initPrometheuzErr
	}
	http.Handle("/"+route, promhttp.Handler())
	listenAndServeErr := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if listenAndServeErr != nil {
		return listenAndServeErr
	}
	return nil
}
