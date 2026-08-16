// Package routing defines the Routing interface and its implementations (round-robin, least-connections).
package routing

import (
	"masit/internal/backend"
)

type Routing interface {
	SelectBackend(backends []backend.Backend) (backend.Backend, error)
}