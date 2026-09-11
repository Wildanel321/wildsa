package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string       `yaml:"version" json:"version"`
	Server  ServerConfig `yaml:"server" json:"server"`
	IPC     IPCConfig    `yaml:"ipc" json:"ipc"`
	Logging LogConfig    `yaml:"logging" json:"logging"`
	Health  HealthConfig `yaml:"health" json:"health"`
}

type ServerConfig struct {
	Host       string `yaml:"host" json:"host"`
	Port       int    `yaml:"port" json:"port"`
	TLSEnabled bool   `yaml:"tls_enabled" json:"tls_enabled"`
	TLSCert    string `yaml:"tls_cert" json:"tls_cert"`
	TLSKey     string `yaml:"tls_key" json:"tls_key"`
}

type IPCConfig struct {
	AgentSocket    string `yaml:"agent_socket" json:"agent_socket"`
	TimeoutSeconds int    `yaml:"timeout_seconds" json:"timeout_seconds"`
}

type LogConfig struct {
	Level string `yaml:"level" json:"level"`
	File  string `yaml:"file" json:"file"`
}

type HealthConfig struct {
	IntervalSeconds      int     `yaml:"interval_seconds" json:"interval_seconds"`
	CPUWarningPercent    float64 `yaml:"cpu_warning_percent" json:"cpu_warning_percent"`
	CPUCriticalPercent   float64 `yaml:"cpu_critical_percent" json:"cpu_critical_percent"`
	MemoryWarningPercent float64 `yaml:"memory_warning_percent" json:"memory_warning_percent"`
	MemoryCriticalPercent float64 `yaml:"memory_critical_percent" json:"memory_critical_percent"`
	DiskWarningPercent   float64 `yaml:"disk_warning_percent" json:"disk_warning_percent"`
	DiskCriticalPercent  float64 `yaml:"disk_critical_percent" json:"disk_critical_percent"`
}

// Default returns sensible defaults for SawitOS daemon
func Default() *Config {
	return &Config{
		Version: "1.0",
		Server: ServerConfig{
			Host:       "127.0.0.1",
			Port:       8080,
			TLSEnabled: false,
		},
		IPC: IPCConfig{
			AgentSocket:    "/run/sawit/sawit-agent.sock",
			TimeoutSeconds: 10,
		},
		Logging: LogConfig{
			Level: "info",
			File:  "/var/log/sawit/sawitd.log",
		},
		Health: HealthConfig{
			IntervalSeconds:       30,
			CPUWarningPercent:     85.0,
			CPUCriticalPercent:    95.0,
			MemoryWarningPercent:  85.0,
			MemoryCriticalPercent: 95.0,
			DiskWarningPercent:    85.0,
			DiskCriticalPercent:   95.0,
		},
	}
}

// LoadFromFile attempts to parse configuration file at path, falling back to defaults if not found
func LoadFromFile(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
