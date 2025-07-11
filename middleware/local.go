package middleware

import "net/http"

// --- Internal structs ---

// middleware implements the Middleware interface and holds
// the core logic of an HTTP middleware along with metadata.
type middleware struct {
	method       func(http.Handler) http.Handler
	name         string
	status       bool
	experimental bool
}

// --- Wrapper for extensibility ---

// MiddlewareWrapper is a function type that accepts and returns
// a Middleware. It is used to wrap or decorate a Middleware with
// additional functionality during initialization.
//
// Example usage:
//
//	func WithLogging(next Middleware) Middleware {
//	    return NewMiddleware(func(h http.Handler) http.Handler {
//	        // wrap handler with logging logic here
//	    }, next.Name(), next.Status(), next.Experimental())
//	}
type MiddlewareWrapper func(r Middleware) Middleware

// --- Middleware implementation ---

// Method returns the actual middleware function of type
// func(http.Handler) http.Handler which will be applied to requests.
func (p middleware) Method() func(http.Handler) http.Handler {
	return p.method
}

// Name returns the identifier string of the middleware.
func (p middleware) Name() string {
	return p.name
}

// Status returns whether the middleware is enabled (true) or disabled (false).
func (p middleware) Status() bool {
	return p.status
}

// Experimental returns true if the middleware is experimental,
// indicating it might be unstable or under active development.
func (p middleware) Experimental() bool {
	return p.experimental
}

// --- Constructors ---

// NewMiddleware constructs a new Middleware instance with the provided
// middleware function, name, enabled status, and experimental flag.
// Additional MiddlewareWrapper options can be passed to decorate or
// modify the middleware before returning.
//
// Example:
//
//	authMiddleware := NewMiddleware(AuthFunc, "Auth", true, false)
//	loggingMiddleware := NewMiddleware(LoggingFunc, "Logging", true, false,
//	  WithLoggingDecorator)
//
// Parameters:
//   - method: the HTTP middleware handler function.
//   - name: a descriptive name for the middleware.
//   - status: whether the middleware should be enabled.
//   - experimental: whether the middleware is experimental.
//   - opts: zero or more MiddlewareWrapper functions to wrap/modify the middleware.
//
// Returns:
//
//	A Middleware interface implementation.
func NewMiddleware(
	method func(http.Handler) http.Handler,
	name string,
	status bool,
	experimental bool,
	opts ...MiddlewareWrapper,
) Middleware {
	var m Middleware = middleware{method, name, status, experimental}
	for _, o := range opts {
		m = o(m)
	}
	return m
}
