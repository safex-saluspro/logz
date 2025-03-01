package logger

import (
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
}
type ConfigImpl struct {
	VlPort            string
	VlBindAddress     string
	VlAddress         string
	VlPidFile         string
	VlReadTimeout     time.Duration
	VlWriteTimeout    time.Duration
	VlIdleTimeout     time.Duration
	VlDefaultLogPath  string
	VlNotifierManager NotifierManager
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

type ConfigManager interface {
	GetConfig() Config
	GetPidPath() string
	GetConfigPath() string
	LoadConfig() (Config, error)
}
type ConfigManagerImpl struct{ config Config }

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
	home = filepath.Join(home, ".kubex", "logz", "config.json")
	if mkdirErr := os.MkdirAll(filepath.Dir(home), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return home
}
func (cm *ConfigManagerImpl) LoadConfig() (Config, error) {
	viperObj := viper.GetViper()
	if viperObj == nil {
		viperObj = viper.New()
	}
	cfgDir := cm.GetConfigPath()
	viperObj.SetConfigFile(cfgDir)
	if readErr := viperObj.ReadInConfig(); readErr != nil {
		return nil, fmt.Errorf("failed to read config: %w", readErr)
	}
	notifierManager := NewNotifierManager(nil)
	if notifierManager == nil {
		return nil, fmt.Errorf("failed to create notifier manager")
	}
	ntfMgr := *notifierManager

	config := ConfigImpl{
		VlPort:            viperObj.GetString("port"),
		VlBindAddress:     viperObj.GetString("bindAddress"),
		VlAddress:         viperObj.GetString("address"),
		VlPidFile:         viperObj.GetString("pidFile"),
		VlReadTimeout:     viperObj.GetDuration("readTimeout"),
		VlWriteTimeout:    viperObj.GetDuration("writeTimeout"),
		VlIdleTimeout:     viperObj.GetDuration("idleTimeout"),
		VlDefaultLogPath:  viperObj.GetString("defaultLogPath"),
		VlNotifierManager: ntfMgr,
	}
	cm.config = &config

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuração alterada: %s", e.Name)

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
