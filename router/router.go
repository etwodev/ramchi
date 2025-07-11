package router

import (
	"net/http"
)

// Router represents an HTTP router that manages multiple routes,
// their associated middleware, and a base URL prefix.
//
// It provides methods to retrieve all registered routes,
// check if the router is active, get its base path prefix,
// and access router-level middleware.
//
// Implementations of Router are responsible for handling
// request dispatching to the appropriate routes.
type Router interface {
	// Routes returns a slice of all registered Route instances
	// managed by this router.
	Routes() []Route

	// Status returns true if the router is active and should
	// handle requests; false otherwise.
	Status() bool

	// Prefix returns the base URL path prefix for all routes
	// in this router, e.g., "/api/v1".
	Prefix() string

	// Middleware returns the slice of middleware functions
	// applied at the router level. Middleware functions have
	// the signature func(http.Handler) http.Handler and wrap
	// the handler chain for all routes under this router.
	Middleware() []func(http.Handler) http.Handler
}

// Route represents a single HTTP endpoint with a method, path,
// handler function, status, experimental flag, and optional middleware.
//
// It encapsulates all data necessary for routing a request to its handler
// with possible route-specific middleware.
type Route interface {
	// Handler returns the HTTP handler function (http.HandlerFunc)
	// that will be invoked when this route matches an incoming request.
	Handler() http.HandlerFunc

	// Method returns the HTTP method (e.g., "GET", "POST", "PUT", "DELETE")
	// that this route responds to.
	Method() string

	// Path returns the relative route path (e.g., "/users/{id}").
	// This is appended to the Router's Prefix when building full route paths.
	Path() string

	// Status returns true if the route is enabled and should
	// accept incoming requests; false otherwise.
	Status() bool

	// Experimental returns true if the route is experimental,
	// allowing for conditional enabling or feature flagging.
	Experimental() bool

	// Middleware returns route-level middleware functions specific
	// to this route, which will be executed after router-level middleware.
	Middleware() []func(http.Handler) http.Handler
}
