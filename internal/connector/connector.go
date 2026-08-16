// Package connector defines the BackendConnector interface and its implementations (TCP, HTTP, GCP managed instance).
package connector

import (
	"net/http"
	"masit/internal/backend"
)

type BackendConnector interface {
	Forward(req *http.Request, backend backend.Backend) (*http.Response, error)
}