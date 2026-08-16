// Package health implements active and passive health checking of backends.
package health

import (
	"masit/internal/backend"
	"masit/internal/connector"
)

type HealthStatus int

const (
	StateUnknown HealthStatus = iota
	StateUp
	StateDown
)

type HealthChecker interface {
	CheckHealth(backendConnector connector.BackendConnector, backendSnapshot backend.Backend) (HealthStatus)
}