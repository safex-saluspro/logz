package logger

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPort        = "9999"
	defaultBindAddress = "0.0.0.0"
	defaultLogPath     = "stdout"
	defaultMode        = ModeStandalone
)

type Config interface {
	Port() string
	BindAddress() string
	Address() string
	PidFile() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	DefaultLogPath() string
	NotifierManager() NotifierManager
	Mode() LogMode
	Level() string
	Format() string
}
type ConfigImpl struct {
	VlLevel           LogLevel
	VlFormat          LogFormat
	VlPort            string
	VlBindAddress     string
	VlAddress         string
	VlPidFile         string
	VlReadTimeout     time.Duration
	VlWriteTimeout    time.Duration
	VlIdleTimeout     time.Duration
	VlDefaultLogPath  string
	VlNotifierManager NotifierManager
	VlMode            LogMode
}

func (c *ConfigImpl) Port() string                     { return c.VlPort }
func (c *ConfigImpl) BindAddress() string              { return c.VlBindAddress }
func (c *ConfigImpl) Address() string                  { return c.VlAddress }
func (c *ConfigImpl) PidFile() string                  { return c.VlPidFile }
func (c *ConfigImpl) ReadTimeout() time.Duration       { return c.VlReadTimeout }
func (c *ConfigImpl) WriteTimeout() time.Duration      { return c.VlWriteTimeout }
func (c *ConfigImpl) IdleTimeout() time.Duration       { return c.VlIdleTimeout }
func (c *ConfigImpl) DefaultLogPath() string           { return c.VlDefaultLogPath }
func (c *ConfigImpl) NotifierManager() NotifierManager { return c.VlNotifierManager }
func (c *ConfigImpl) Mode() LogMode                    { return c.VlMode }
func (c *ConfigImpl) Level() string                    { return string(c.VlLevel) }
func (c *ConfigImpl) Format() string                   { return string(c.VlFormat) }

type ConfigManager interface {
	GetConfig() Config
	GetPidPath() string
	GetConfigPath() string
	LoadConfig() (Config, error)
}
type ConfigManagerImpl struct {
	config Config
}

func (cm *ConfigManagerImpl) GetConfig() Config { return cm.config }
func (cm *ConfigManagerImpl) GetPidPath() string {
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil {
		cacheDir = "/tmp"
	}
	cacheDir = filepath.Join(cacheDir, "logz", cm.config.PidFile())
	if mkdirErr := os.MkdirAll(filepath.Dir(cacheDir), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return cacheDir
}
func (cm *ConfigManagerImpl) GetConfigPath() string {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	configPath := filepath.Join(home, ".kubex", "logz", "config.json")
	if mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return configPath
}
func (cm *ConfigManagerImpl) LoadConfig() (Config, error) {
	configPath := cm.GetConfigPath()
	if err := ensureConfigExists(configPath); err != nil {
		return nil, fmt.Errorf("failed to ensure config exists: %w", err)
	}

	viperObj := viper.New()
	viperObj.SetConfigFile(configPath)
	if readErr := viperObj.ReadInConfig(); readErr != nil {
		return nil, fmt.Errorf("failed to read config: %w", readErr)
	}

	notifierManager := NewNotifierManager(nil)
	if notifierManager == nil {
		return nil, fmt.Errorf("failed to create notifier manager")
	}

	// Construir o Config com valores do arquivo ou padrões
	mode := LogMode(viperObj.GetString("mode"))
	if mode != ModeService && mode != ModeStandalone {
		mode = defaultMode
	}

	config := ConfigImpl{
		VlPort:            getOrDefault(viperObj.GetString("port"), defaultPort),
		VlBindAddress:     getOrDefault(viperObj.GetString("bindAddress"), defaultBindAddress),
		VlAddress:         fmt.Sprintf("%s:%s", defaultBindAddress, defaultPort),
		VlPidFile:         viperObj.GetString("pidFile"),
		VlReadTimeout:     viperObj.GetDuration("readTimeout"),
		VlWriteTimeout:    viperObj.GetDuration("writeTimeout"),
		VlIdleTimeout:     viperObj.GetDuration("idleTimeout"),
		VlDefaultLogPath:  getOrDefault(viperObj.GetString("defaultLogPath"), defaultLogPath),
		VlNotifierManager: *notifierManager,
		VlMode:            mode,
	}

	cm.config = &config

	viperObj.WatchConfig()
	viperObj.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuração alterada: %s", e.Name)
		// Atualizar Config dinamicamente, se necessário
	})

	return cm.config, nil
}

func NewConfigManager() *ConfigManager {
	cfgMgr := &ConfigManagerImpl{}

	if cfg, err := cfgMgr.LoadConfig(); err != nil || cfg == nil {
		log.Printf("Erro ao carregar configuração: %v\n", err)
		return nil
	}

	var cfgM ConfigManager
	cfgM = cfgMgr

	return &cfgM
}

func ensureConfigExists(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Gera uma configuração padrão
		defaultConfig := ConfigImpl{
			VlPort:            defaultPort,
			VlBindAddress:     defaultBindAddress,
			VlAddress:         fmt.Sprintf("%s:%s", defaultBindAddress, defaultPort),
			VlPidFile:         "logz_srv.pid",
			VlReadTimeout:     15 * time.Second,
			VlWriteTimeout:    15 * time.Second,
			VlIdleTimeout:     60 * time.Second,
			VlDefaultLogPath:  defaultLogPath,
			VlNotifierManager: *NewNotifierManager(nil),
			VlMode:            defaultMode,
		}
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		if writeErr := os.WriteFile(configPath, data, 0644); writeErr != nil {
			return fmt.Errorf("failed to create default config: %w", writeErr)
		}
	}
	return nil
}
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
