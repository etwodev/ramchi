package router

import (
	"net/http"
)

type Router interface {
	// Returns the list of all routes
	Routes() []Route
	// Is the router enabled
	Status() bool
	// The router prefix
	Prefix() string
}

type Route interface {
	// Handler returns the function the route applies
	Handler() http.HandlerFunc
	// Method returns the http method the route corresponds to
	Method() string
	// Path returns the subpath where the route responds
	Path() string
	// Status returns whether the route is enabled
	Status() bool
	// Experimental returns whether the route is experimental
	Experimental() bool
}
