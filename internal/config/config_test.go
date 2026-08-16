package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	content := `
listenServer:
  address: ":8080"
routingMethod: "round-robin"
requestTimeout: 5
healthCheckInterval: 10
`
	dir := t.TempDir()
	path := dir + "/config.yaml"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ListenServer.Address != ":8080" {
		t.Errorf("Expected listen address ':8080', got '%s'", cfg.ListenServer.Address)
	}

	if cfg.RoutingMethod != "round-robin" {
		t.Errorf("Expected routing method 'round-robin', got '%s'", cfg.RoutingMethod)
	}

	if cfg.RequestTimeout != int((5 * time.Second).Seconds()) {
		t.Errorf("Expected request timeout 5s, got %v", cfg.RequestTimeout)
	}

	if cfg.HealthCheckInterval != int((10 * time.Second).Seconds()) {
		t.Errorf("Expected health check interval 10s, got %v", cfg.HealthCheckInterval)
	}
}
