// Package backend defines the Backend data type and the BackendProvider service-discovery interface.
package backend

type Backend struct {
	ID       string
	Metadata map[string]string
}

type BackendProvider interface {
	ListBackends() ([]Backend, error)
}
