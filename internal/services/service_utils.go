package services

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func newConfig(port string) (*Config, error) {
	if port == "" {
		port = os.Getenv("LOGZ_PORT")
		if port == "" {
			port = defaultPort
		}
	} else {
		if err := validatePort(port); err != nil {
			return nil, err
		}
	}
	bindAddress := os.Getenv("LOGZ_BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = defaultBindAddress
	}
	newCfg := Config{
		Port:           port,
		BindAddress:    bindAddress,
		Address:        net.JoinHostPort(bindAddress, port),
		PidFile:        getPidPath(),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		DefaultLogPath: "./logs",
		Integrations:   map[string]struct{ enabled bool }{"prometheus": {enabled: true}},
	}
	return &newCfg, nil
}

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

func generateConfig(cfgDir, cfgFile, port string) error {
	cfgPath := filepath.Join(cfgDir, cfgFile)
	cfgDefault, cfgDefaultErr := newConfig(port)
	if cfgDefaultErr != nil {
		return fmt.Errorf("failed to create default config: %w", cfgDefaultErr)
	}
	cfgNewFile, cfgNewFileErr := os.Create(cfgPath)
	if cfgNewFileErr != nil {
		return fmt.Errorf("failed to create config file: %w", cfgNewFileErr)
	}
	defer func(cfgNewFile *os.File) {
		_ = cfgNewFile.Close()
	}(cfgNewFile)
	if cfgFile == "" {
		cfgFile = "service_config.json"
	}
	jsonData, jsonErr := json.Marshal(cfgDefault)
	if jsonErr != nil {
		return fmt.Errorf("failed to marshal default config: %w", jsonErr)
	}
	if _, writeErr := cfgNewFile.Write(jsonData); writeErr != nil {
		return fmt.Errorf("failed to write default config: %w", writeErr)
	}
	return nil
}

func getConfigType(configFile string) string {
	ext := filepath.Ext(configFile)
	switch ext {
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	default:
		return ""
	}
}

func GetConfig(port string) (string, string, error) {
	var cfgDir, cfgFile, cfgType string
	if configPath := os.Getenv("LOGZ_CONFIG_PATH"); configPath != "" {
		cfgDir = filepath.Dir(configPath)
		cfgFile = filepath.Base(configPath)
	} else {
		usrConfigDir, usrConfigDirErr := os.UserConfigDir()
		if usrConfigDirErr == nil {
			cfgDir = filepath.Join(usrConfigDir, "logz")
		} else {
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				cacheDir, cacheDirErr := os.UserCacheDir()
				if cacheDirErr != nil {
					return "", "", os.ErrNotExist
				} else {
					cfgDir = filepath.Join(cacheDir, "logz")
				}
			} else {
				cfgDir = filepath.Join(home, "logz")
			}
		}
		if _, statErr := os.Stat(cfgDir); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll(cfgDir, 0755); mkdirErr != nil {
				return "", "", fmt.Errorf("failed to create config directory: %w", mkdirErr)
			}
		}
		cfgFileSearch, cfgFileSearchErr := filepath.Glob(filepath.Join(cfgDir, "service_config.*"))
		if cfgFileSearchErr == nil && len(cfgFileSearch) > 0 {
			cfgDir = filepath.Dir(cfgFileSearch[0])
			cfgFile = filepath.Base(cfgFileSearch[0])
		}
	}
	if cfgFile == "" {
		if generateErr := generateConfig(cfgDir, cfgFile, port); generateErr != nil {
			return "", "", fmt.Errorf("failed to generate config file: %w", generateErr)
		}
		cfgFile = "service_config.json"
		cfgType = "json"
	} else {
		cfgType = getConfigType(cfgFile)
		if cfgType == "" {
			return "", "", os.ErrInvalid
		}
	}
	return filepath.Join(cfgDir, cfgFile), cfgType, nil
}

func loadConfig(port string) error {
	cfgPath, cfgType, cfgErr := GetConfig(port)
	if cfgErr != nil {
		return cfgErr
	}
	viper.SetConfigFile(cfgPath)
	viper.SetConfigType(cfgType)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("erro ao ler config: %w", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuração alterada: %s", e.Name)
		_ = Run()
	})
	return nil
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
