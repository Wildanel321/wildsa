package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sawitos/sawit/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.Default()
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.IPC.AgentSocket == "" {
		t.Error("expected non-empty agent socket path")
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "sawitd.yaml")

	yamlData := []byte(`
version: "1.0"
server:
  host: "0.0.0.0"
  port: 9090
`)
	if err := os.WriteFile(cfgPath, yamlData, 0644); err != nil {
		t.Fatalf("failed to write tmp config: %v", err)
	}

	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		t.Fatalf("expected clean config load, got: %v", err)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090 from file, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected host 0.0.0.0, got %s", cfg.Server.Host)
	}
}
