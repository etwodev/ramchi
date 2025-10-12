package router

import (
	"net/http"
)

// --- Internal structs ---

// route implements the Route interface.
type route struct {
	method       string
	path         string
	status       bool
	experimental bool
	handler      http.HandlerFunc
	middleware   []func(http.Handler) http.Handler
}

// router implements the Router interface.
type router struct {
	status     bool
	prefix     string
	routes     []Route
	middleware []func(http.Handler) http.Handler
}

// --- Route implementation ---

// Handler returns the HTTP handler function for this route.
func (r route) Handler() http.HandlerFunc {
	return r.handler
}

// Method returns the HTTP method (GET, POST, etc.) for this route.
func (r route) Method() string {
	return r.method
}

// Path returns the relative path for this route.
func (r route) Path() string {
	return r.path
}

// Status returns whether the route is enabled.
func (r route) Status() bool {
	return r.status
}

// Experimental returns whether the route is experimental.
func (r route) Experimental() bool {
	return r.experimental
}

// Middleware returns the middleware chain for this route.
func (r route) Middleware() []func(http.Handler) http.Handler {
	return r.middleware
}

// --- Router implementation ---

// Routes returns all routes registered on this router.
func (r router) Routes() []Route {
	return r.routes
}

// Status returns whether the router is enabled.
func (r router) Status() bool {
	return r.status
}

// Prefix returns the base path prefix for all routes in this router.
func (r router) Prefix() string {
	return r.prefix
}

// Middleware returns router-level middleware applied to all routes.
func (r router) Middleware() []func(http.Handler) http.Handler {
	return r.middleware
}

// --- Wrappers for extensibility ---

// RouterWrapper defines a function signature to wrap or modify a Router.
type RouterWrapper func(r Router) Router

// RouteWrapper defines a function signature to wrap or modify a Route.
type RouteWrapper func(r Route) Route

// --- Constructors ---

// NewRouter creates a new Router with the specified prefix, routes, status, and optional middleware.
// Additional options can be applied using RouterWrapper functions.
//
// Example:
//
//	r := NewRouter("/api", routes, true, []func(http.Handler) http.Handler{loggingMiddleware})
func NewRouter(prefix string, routes []Route, status bool, middleware []func(http.Handler) http.Handler, opts ...RouterWrapper) Router {
	var r Router = router{
		status:     status,
		prefix:     prefix,
		routes:     routes,
		middleware: middleware,
	}
	for _, o := range opts {
		r = o(r)
	}
	return r
}

// NewRoute creates a new Route with the given HTTP method, path, status, experimental flag,
// handler function, and optional middleware.
// Additional options can be applied using RouteWrapper functions.
//
// Example:
//
//	r := NewRoute(http.MethodGet, "/ping", true, false, pingHandler, nil)
func NewRoute(method, path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	var r Route = route{
		method:       method,
		path:         path,
		status:       status,
		experimental: experimental,
		handler:      handler,
		middleware:   middleware,
	}
	for _, o := range opts {
		r = o(r)
	}
	return r
}

// --- Convenience functions for HTTP methods ---

// NewGetRoute creates a new GET Route.
//
// Example:
//
//	r := NewGetRoute("/users", true, false, usersHandler, nil)
func NewGetRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodGet, path, status, experimental, handler, middleware, opts...)
}

// NewPostRoute creates a new POST Route.
//
// Example:
//
//	r := NewPostRoute("/users", true, false, createUserHandler, nil)
func NewPostRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPost, path, status, experimental, handler, middleware, opts...)
}

// NewPutRoute creates a new PUT Route.
//
// Example:
//
//	r := NewPutRoute("/users/{id}", true, false, updateUserHandler, nil)
func NewPutRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPut, path, status, experimental, handler, middleware, opts...)
}

// NewDeleteRoute creates a new DELETE Route.
//
// Example:
//
//	r := NewDeleteRoute("/users/{id}", true, false, deleteUserHandler, nil)
func NewDeleteRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodDelete, path, status, experimental, handler, middleware, opts...)
}

// NewPatchRoute creates a new PATCH Route.
//
// Example:
//
//	r := NewPatchRoute("/users/{id}", true, false, patchUserHandler, nil)
func NewPatchRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPatch, path, status, experimental, handler, middleware, opts...)
}

// NewOptionsRoute creates a new OPTIONS Route.
//
// Example:
//
//	r := NewOptionsRoute("/users", true, false, optionsHandler, nil)
func NewOptionsRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodOptions, path, status, experimental, handler, middleware, opts...)
}

// NewHeadRoute creates a new HEAD Route.
//
// Example:
//
//	r := NewHeadRoute("/users", true, false, headHandler, nil)
func NewHeadRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodHead, path, status, experimental, handler, middleware, opts...)
}

// NewConnectRoute creates a new CONNECT Route.
//
// Example:
//
//	r := NewConnectRoute("/users/{id}", true, false, connectHandler, nil)
func NewConnectRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodConnect, path, status, experimental, handler, middleware, opts...)
}

// NewTraceRoute creates a new TRACE Route.
//
// Example:
//
//	r := NewTraceRoute("/users/{id}", true, false, traceHandler, nil)
func NewTraceRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodTrace, path, status, experimental, handler, middleware, opts...)
}
