package router

import (
	"net/http"
)

// --- Internal structs ---

type route struct {
	method       string
	path         string
	status       bool
	experimental bool
	handler      http.HandlerFunc
	middleware   []func(http.Handler) http.Handler
}

type router struct {
	status     bool
	prefix     string
	routes     []Route
	middleware []func(http.Handler) http.Handler
}

// --- Route implementation ---

func (r route) Handler() http.HandlerFunc {
	return r.handler
}

func (r route) Method() string {
	return r.method
}

func (r route) Path() string {
	return r.path
}

func (r route) Status() bool {
	return r.status
}

func (r route) Experimental() bool {
	return r.experimental
}

func (r route) Middleware() []func(http.Handler) http.Handler {
	return r.middleware
}

// --- Router implementation ---

func (r router) Routes() []Route {
	return r.routes
}

func (r router) Status() bool {
	return r.status
}

func (r router) Prefix() string {
	return r.prefix
}

func (r router) Middleware() []func(http.Handler) http.Handler {
	return r.middleware
}

// --- Wrappers for extensibility ---

type RouterWrapper func(r Router) Router
type RouteWrapper func(r Route) Route

// --- Constructors ---

// NewRouter creates a new Router with a prefix, status flag, routes, and optional middleware.
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

// NewRoute creates a new Route with method, path, status, middleware, and handler.
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

// --- Convenience functions for each HTTP verb ---

func NewGetRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodGet, path, status, experimental, handler, middleware, opts...)
}

func NewPostRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPost, path, status, experimental, handler, middleware, opts...)
}

func NewPutRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPut, path, status, experimental, handler, middleware, opts...)
}

func NewDeleteRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodDelete, path, status, experimental, handler, middleware, opts...)
}

func NewOptionsRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodOptions, path, status, experimental, handler, middleware, opts...)
}

func NewHeadRoute(path string, status, experimental bool, handler http.HandlerFunc, middleware []func(http.Handler) http.Handler, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodHead, path, status, experimental, handler, middleware, opts...)
}
