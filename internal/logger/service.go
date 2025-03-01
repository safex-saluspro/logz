package logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
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
	pidFile = "logz_srv.pid"
)

var (
	lSrv         *http.Server
	lClient      *http.Client
	lSocket      *zmq4.Socket
	lDBus        *dbus.Conn
	globalLogger *Logger // Logger global para o serviço
	startTime    = time.Now()
)

func Run() error {
	// Verifica se já há um serviço rodando para evitar múltiplas instâncias
	if IsRunning() {
		if stopErr := shutdown(); stopErr != nil {
			return stopErr
		}
	}

	// Inicializa o ConfigManager e carrega a configuração
	configManager := NewConfigManager()
	if configManager == nil {
		return errors.New("failed to initialize config manager")
	}
	cfgMgr := *configManager

	config, err := cfgMgr.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Inicializa o Logger global com a configuração
	initializeGlobalLogger(config)

	// Configura o servidor HTTP
	mux := http.NewServeMux()
	if err := registerHandlers(mux); err != nil {
		return err
	}

	lSrv = &http.Server{
		Addr:         config.Address(),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		IdleTimeout:  config.IdleTimeout(),
	}

	// Inicia o servidor HTTP
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		globalLogger.Info(fmt.Sprintf("Service running on %s", config.Address()), nil)
		if err := lSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			globalLogger.Error(fmt.Sprintf("Service encountered an error: %v", err), nil)
		}
	}()

	<-stop
	return shutdown()
}
func Start(port string) error {
	if IsRunning() {
		return errors.New("service already running (pid file exists: " + getPidPath() + ")")
	}

	// Usa o Viper para carregar configuração de execução
	vpr := viper.GetViper()
	if vpr == nil {
		return errors.New("viper not initialized")
	}

	cmd := exec.Command(os.Args[0], "service", "spawn", "-c", vpr.ConfigFileUsed())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	file, err := os.OpenFile(getPidPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open PID file: %w", err)
	}
	defer file.Close()

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return errors.New("another process is writing to the PID file")
	}

	pid := cmd.Process.Pid
	pidData := fmt.Sprintf("%d\n%s", pid, port)
	if _, writeErr := file.Write([]byte(pidData)); writeErr != nil {
		return fmt.Errorf("failed to write PID data: %w", writeErr)
	}

	globalLogger.Info(fmt.Sprintf("Service started with pid %d", pid), nil)
	return nil
}
func Stop() error {
	pid, port, pidPath, err := GetServiceInfo()
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	time.Sleep(1 * time.Second)
	if err := os.Remove(pidPath); err != nil {
		return err
	}

	globalLogger.Info(fmt.Sprintf("Service with pid %d and port %s stopped", pid, port), nil)
	return nil
}

func Server() *http.Server {
	return lSrv
}
func Client() *http.Client {
	if lClient == nil {
		lClient = &http.Client{}
	}
	return lClient
}
func Socket() *zmq4.Socket {
	if lSocket == nil {
		lSocket, _ = zmq4.NewSocket(zmq4.PUB)
	}
	return lSocket
}
func DBus() *dbus.Conn {
	if lDBus == nil {
		lDBus, _ = dbus.SystemBus()
	}
	return lDBus
}

func getPidPath() string {
	if envPath := os.Getenv("LOGZ_PID_PATH"); envPath != "" {
		return envPath
	}
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil {
		cacheDir = "/tmp"
	}
	cacheDir = filepath.Join(cacheDir, "logz", pidFile)
	if mkdirErr := os.MkdirAll(filepath.Dir(cacheDir), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return cacheDir
}

func IsRunning() bool {
	_, err := os.Stat(getPidPath())
	return err == nil
}
func GetServiceInfo() (int, string, string, error) {
	pidPath := getPidPath()

	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, "", "", os.ErrNotExist
	}

	lines := strings.Split(string(data), "\n")
	pid, pidErr := strconv.Atoi(lines[0])
	if pidErr != nil {
		return 0, "", "", os.ErrInvalid
	}

	port := "unknown"
	if len(lines) > 1 {
		port = lines[1]
	}

	return pid, port, pidPath, nil
}
func registerHandlers(mux *http.ServeMux) error {
	integrations := viper.GetStringMap("integrations")
	if integrations == nil {
		return errors.New("no integrations configured")
	}

	for path := range integrations {
		if !viper.GetBool("integrations." + path + ".enabled") {
			continue
		}

		healthPath, _ := url.JoinPath("/", path, "/health")
		metricsPath, _ := url.JoinPath("/", path, "/metrics")
		callbackPath, _ := url.JoinPath("/", path, "/receive")

		mux.HandleFunc(healthPath, healthHandler)
		mux.HandleFunc(metricsPath, metricsHandler)
		mux.HandleFunc(callbackPath, callbackHandler)
	}

	return nil
}
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limita o tamanho do payload para evitar abusos
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if _, ok := payload["message"]; !ok {
		http.Error(w, "Missing 'message' in payload", http.StatusBadRequest)
		return
	}

	globalLogger.Info(fmt.Sprintf("Callback received: %v", payload), nil)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success","message":"Callback processed"}`))
}
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(startTime).String()
	response := fmt.Sprintf("OK\nUptime: %s\n", uptime)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}
func metricsHandler(w http.ResponseWriter, _ *http.Request) {
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
func shutdown() error {
	globalLogger.Info("Shutting down service gracefully...", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := lSrv.Shutdown(ctx); err != nil {
		globalLogger.Error(fmt.Sprintf("Service shutdown failed: %v", err), nil)
		return fmt.Errorf("shutdown process failed: %w", err)
	}

	globalLogger.Info("Service stopped gracefully.", nil)
	return nil
}

func initializeGlobalLogger(config Config) {
	if globalLogger == nil {
		globalLogger = NewLogger(config)
	}
}
