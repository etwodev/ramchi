package middleware

import (
	"net/http"
)

// Middleware defines the interface that wraps an HTTP middleware function
// and provides metadata about the middleware such as its name, status,
// and whether it is experimental.
//
// This interface allows middleware to be managed, enabled/disabled,
// and identified dynamically within the server or application.
type Middleware interface {
	// Method returns the core middleware function of type
	// func(http.Handler) http.Handler that will be applied to HTTP handlers.
	Method() func(http.Handler) http.Handler

	// Status returns whether the middleware is enabled (true) or disabled (false).
	Status() bool

	// Experimental returns true if the middleware is experimental,
	// indicating it might be unstable or under active development.
	Experimental() bool

	// Name returns the identifying name of the middleware.
	Name() string
}
