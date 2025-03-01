package services

import (
	"errors"
	"fmt"
	//"github.com/prometheus/client_golang/prometheus"
	// github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	lSrv      *http.Server
	startTime = time.Now()
)

func Run() error {
	if IsRunning() {
		if stopErr := shutdown(); stopErr != nil {
			return stopErr
		}
	}
	var glbViper = viper.GetViper()
	if glbViper == nil {
		return errors.New("viper not initialized")
	} else {
		if readErr := glbViper.ReadInConfig(); readErr != nil {
			return fmt.Errorf("failed to read config: %w", readErr)
		}
	}

	mux := http.NewServeMux()
	registerErr := registerHandlers(mux)
	if registerErr != nil {
		return registerErr
	}

	lSrv = &http.Server{
		Addr:         viper.GetString("address"),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  viper.GetDuration("readTimeout"),
		WriteTimeout: viper.GetDuration("writeTimeout"),
		IdleTimeout:  viper.GetDuration("idleTimeout"),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("Service running on port %s\n", viper.GetString("port"))
		if err := lSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Service encountered an error: %v\n", err)
		}
	}()

	<-stop
	return shutdown()
}
func Start(port string) error {
	if IsRunning() {
		return errors.New("service already running (pid file exists: " + getPidPath() + ")")
	}
	var vpr = viper.GetViper()
	if vpr == nil {
		return errors.New("viper not initialized")
	}
	cmd := exec.Command(os.Args[0], "service-run", viper.ConfigFileUsed())

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
	var pathErr error
	if integrations := viper.GetStringMap("integrations"); integrations == nil {
		return os.ErrInvalid
	} else {
		for path := range integrations {
			if isActive := viper.GetBool("integrations." + path + ".enabled"); !isActive {
				continue
			}
			var healthPath, metricsPath, callbackPath string
			if path == "" {
				continue
			}
			healthPath, pathErr = url.JoinPath("/", path, "/health")
			if pathErr != nil {
				log.Printf("Invalid path: %s\n", path)
				continue
			}
			metricsPath, pathErr = url.JoinPath("/", path, "/metrics")
			if pathErr != nil {
				log.Printf("Invalid path: %s\n", path)
				continue
			}
			callbackPath, pathErr = url.JoinPath("/", path, "/receive")
			if pathErr != nil {
				log.Printf("Invalid path: %s\n", path)
				continue
			}
			mux.HandleFunc(healthPath, healthHandler)
			mux.HandleFunc(metricsPath, metricsHandler)
			mux.HandleFunc(callbackPath, callbackHandler)
		}
	}
	return nil
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
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Aqui eu preciso de uma sugestão que seja simples, segura e eficiente.
	// Somente uma ou duas opções de possibilidade retornos aceitáveis.
	// E alguma forma de assegurar que não haja problemas de segurança.

}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
