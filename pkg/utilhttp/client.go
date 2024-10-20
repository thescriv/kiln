package utilhttp

import (
	"net/http"
	"time"
)

// Client is an interface that exposes a single Do method to perform an HTTP request.
type Client interface {
	// Do performs an HTTP request and request a Response or an error.
	Do(*http.Request) (*http.Response, error)
}

// NewClient returns a Client configured with a timeout.
func NewClient(timeout time.Duration) Client {
	var client = &http.Client{
		Timeout: timeout,
	}

	return client
}
