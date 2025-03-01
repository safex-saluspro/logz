package services

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPort        = "9999"
	defaultBindAddress = "0.0.0.0"
	pidFile            = "logz_srv.pid"
)

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

func shutdown() error {
	fmt.Println("Shutting down service gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := lSrv.Shutdown(ctx); err != nil {
		log.Printf("Service shutdown failed: %v\n", err)
		return fmt.Errorf("shutdown process failed: %w", err)
	}
	log.Println("Service stopped gracefully.")
	return nil
}
