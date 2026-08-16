// Package config loads load balancer configuration from YAML files and environment variables.
package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenServer struct {
		Address string `yaml:"address"`
	} `yaml:"listenServer"`
	RoutingMethod       string `yaml:"routingMethod"`
	RequestTimeout      int    `yaml:"requestTimeout"`
	HealthCheckInterval int    `yaml:"healthCheckInterval"`
}

type ConfigLoader interface {
	LoadConfig(path string) (*Config, error)
}

func LoadConfig(path string) (*Config, error) {
	// Load configuration from YAML file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	// Load configuration from environment variables
	// cfg.ListenServer.Address = os.Getenv("LISTEN_ADDRESS")
	// cfg.RoutingMethod = os.Getenv("ROUTING_METHOD")
	// cfg.RequestTimeout = os.Getenv("REQUEST_TIMEOUT")
	// cfg.HealthCheckInterval = os.Getenv("HEALTH_CHECK_INTERVAL")

	return &cfg, nil
}
