package router

import (
	"net/http"
)

type Router interface {
	// Routes returns all registered routes
	Routes() []Route
	// Status returns whether the router is active
	Status() bool
	// Prefix returns the base path
	Prefix() string
	// Middleware returns router-level middleware
	Middleware() []func(http.Handler) http.Handler
}

type Route interface {
	// Handler is the HTTP handler function
	Handler() http.HandlerFunc
	// Method is the HTTP verb (GET, POST, etc.)
	Method() string
	// Path is the relative route path
	Path() string
	// Status returns whether the route is active
	Status() bool
	// Experimental returns whether the route is experimental
	Experimental() bool
	// Middleware returns route-level middleware
	Middleware() []func(http.Handler) http.Handler
}
