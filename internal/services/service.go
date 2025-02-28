package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	pidFile     = "service.pid"
	defaultPort = "9999"
)

// Run inicia o servidor HTTP e bloqueia até receber um sinal de término.
func Run() error {
	mux := http.NewServeMux()

	// Endpoints para integrações (stubs)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/discord", discordHandler)
	mux.HandleFunc("/slack", slackHandler)
	mux.HandleFunc("/telegram", telegramHandler)
	mux.HandleFunc("/metrics", prometheusMetricsHandler)
	mux.HandleFunc("/grafana", grafanaHandler)

	port := defaultPort
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Canal para sinais de desligamento
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine
	go func() {
		fmt.Printf("Service running on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Service encountered an error: %v\n", err)
		}
	}()

	// Bloqueia até receber um sinal de parada
	<-stop
	fmt.Println("Shutting down service gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("service shutdown failed: %w", err)
	}

	log.Println("Service stopped")
	return nil
}

// Start lança o serviço em segundo plano (daemoniza) se não estiver rodando.
func Start() error {
	if _, err := os.Stat(pidFile); err == nil {
		return errors.New("service already running (pid file exists)")
	}

	// Executa o mesmo binário com o argumento "service-run"
	cmd := exec.Command(os.Args[0], "service-run")
	// Redireciona a saída (você pode ajustar se preferir redirecionar para um arquivo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	// Escreve o PID em um arquivo para permitir o "stop"
	pid := cmd.Process.Pid
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("failed to write pid file: %w", err)
	}
	fmt.Printf("Service started with pid %d", pid)
	return nil
}

// Stop encerra o serviço em execução lendo o PID e enviando um SIGTERM.
func Stop() error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return errors.New("service not running (pid file not found)")
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("invalid pid file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	// Aguarda um pouco e remove o arquivo de PID
	time.Sleep(1 * time.Second)
	removeErr := os.Remove(pidFile)
	if removeErr != nil {
		return removeErr
	}
	log.Printf("Service with pid %d stopped", pid)
	return nil
}

// healthHandler retorna status OK.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func discordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Discord integration endpoint"))
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Slack integration endpoint"))
}

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telegram integration endpoint"))
}

func grafanaHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Grafana integration endpoint"))
}

// prometheusMetricsHandler expõe as métricas no formato Prometheus.
func prometheusMetricsHandler(w http.ResponseWriter, r *http.Request) {
	pm := GetPrometheusManager()
	if !pm.IsEnabled() {
		http.Error(w, "Prometheus integration is not enabled", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for name, value := range pm.GetMetrics() {
		// Cada métrica é exposta como um gauge.
		fmt.Fprintf(w, "# HELP %s Custom metric from Logz\n", name)
		fmt.Fprintf(w, "# TYPE %s gauge\n", name)
		fmt.Fprintf(w, "%s %f\n", name, value)
	}
}
