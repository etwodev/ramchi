package middleware

import "net/http"

type preMiddleware struct {
	method       func(http.Handler) http.Handler
	name         string
	status       bool
	experimental bool
}

// MiddlewareWrapper wraps a middleware with extra functionality.
// It is passed in when creating a new middleware.
type MiddlewareWrapper func(r Middleware) Middleware

// Method returns the method that the middleware enacts.
func (p preMiddleware) Method() func(http.Handler) http.Handler {
	return p.method
}

// Name returns the identification of the middleware
func (p preMiddleware) Name() string {
	return p.name
}

// Status returns whether the middleware is enabled.
func (p preMiddleware) Status() bool {
	return p.status
}

// Experimental returns whether the middleware is experimental or not.
func (p preMiddleware) Experimental() bool {
	return p.experimental
}

// NewMiddleware initializes a new local middleware for the server.
func NewMiddleware(method func(http.Handler) http.Handler, name string, status bool, experimental bool, opts ...MiddlewareWrapper) Middleware {
	var m Middleware = preMiddleware{method, name, status, experimental}
	for _, o := range opts {
		m = o(m)
	}
	return m
}
