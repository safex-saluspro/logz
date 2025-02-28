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
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	defaultPort = "9999"
	pidFile     = "logz_srv.pid"
)

type Config struct {
	Port           string
	PidFile        string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	DefaultLogPath string
}

func getConfig(port string) *Config {
	if port == "" {
		port = defaultPort
	}
	return &Config{
		Port:           defaultPort,
		PidFile:        getPidPath(),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		DefaultLogPath: "./logs",
	}
}

var startTime = time.Now()

func getPidPath() string {
	if envPath := os.Getenv("LOGZ_PID_PATH"); envPath != "" {
		return envPath
	}
	home, homeErr := os.UserCacheDir()
	if homeErr != nil {
		home = "/tmp"
	}
	return filepath.Join(home, pidFile)
}

func validatePort(port string) error {
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("invalid port: %s (must be between 1 and 65535)", port)
	}
	return nil
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/discord", discordHandler)
	mux.HandleFunc("/slack", slackHandler)
	mux.HandleFunc("/telegram", telegramHandler)
	mux.HandleFunc("/metrics", prometheusMetricsHandler)
	mux.HandleFunc("/grafana", grafanaHandler)
}

func Run(port string) error {
	if IsRunning() {
		return errors.New("service already running (pid file exists)")
	}

	if err := validatePort(port); err != nil {
		return err
	}

	mux := http.NewServeMux()
	registerHandlers(mux)
	config := getConfig(port)

	srv := &http.Server{
		Addr:         "127.0.0.1:" + config.Port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Service running on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Service encountered an error: %v\n", err)
		}
	}()

	<-stop
	fmt.Println("Shutting down service gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Service shutdown failed: %v\n", err)
		return fmt.Errorf("shutdown process failed: %w", err)
	}
	log.Println("Service stopped gracefully.")

	return nil
}

func Start(port string) error {
	if IsRunning() {
		return errors.New("service already running (pid file exists)")
	}

	if port == "" {
		port = defaultPort
	}
	cmd := exec.Command(os.Args[0], "service-run", "--port", port)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	file, err := os.OpenFile(getPidPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open PID file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return errors.New("another process is writing to the PID file")
	}

	pid := cmd.Process.Pid
	pidData := fmt.Sprintf("%d\n%s", pid, port)
	if _, err := file.Write([]byte(pidData)); err != nil {
		return fmt.Errorf("failed to write PID data: %w", err)
	}
	fmt.Printf("Service started with pid %d\n", pid)
	return nil
}

func Stop() error {
	pid, port, pidPath, err := GetServiceInfo()
	if err != nil {
		return err
	}

	process, processErr := os.FindProcess(pid)
	if processErr != nil {
		return fmt.Errorf("failed to find process: %w", processErr)
	}

	if signalErr := process.Signal(syscall.SIGTERM); signalErr != nil {
		return fmt.Errorf("failed to stop process: %w", signalErr)
	}

	time.Sleep(1 * time.Second)
	if removeErr := os.Remove(pidPath); removeErr != nil {
		return removeErr
	}
	log.Printf("Service with pid %d and port %s stopped", pid, port)
	return nil
}

func IsRunning() bool {
	_, err := os.Stat(getPidPath())
	return err == nil
}

func GetServiceInfo() (int, string, string, error) {
	pidPath := getPidPath()

	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, "", "", errors.New("service not running (PID file not found)")
	}

	lines := strings.Split(string(data), "\n")
	pid, pidErr := strconv.Atoi(lines[0])
	if pidErr != nil {
		return 0, "", "", fmt.Errorf("invalid PID in PID file: %w", pidErr)
	}

	port := "unknown"
	if len(lines) > 1 {
		port = lines[1]
	}

	return pid, port, pidPath, nil
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(startTime).String()
	response := fmt.Sprintf("OK\nUptime: %s\n", uptime)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func discordHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Discord integration endpoint"))
}

func slackHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Slack integration endpoint"))
}

func telegramHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Telegram integration endpoint"))
}

func grafanaHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Grafana integration endpoint"))
}

func prometheusMetricsHandler(w http.ResponseWriter, _ *http.Request) {
	pm := GetPrometheusManager()
	if !pm.IsEnabled() {
		http.Error(w, "Prometheus integration is not enabled", http.StatusForbidden)
		return
	}

	metrics := pm.GetMetrics()
	if len(metrics) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for name, value := range metrics {
		if _, err := fmt.Fprintf(w, "# HELP %s Custom metric from Logz\n# TYPE %s gauge\n%s %f\n", name, name, name, value); err != nil {
			fmt.Println(fmt.Sprintf("Error writing metric '%s': %v", name, err))
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
