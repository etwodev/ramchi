package router

import (
	"net/http"
)

type preRouter struct {
	status 			bool
	routes  		[]Route
}

type preRoute struct {
	method  string
	path    string
	status  bool
	experimental bool
	handler http.HandlerFunc
}

// RouterWrapper wraps a router with extra functionality .
// It is passed in when creating a new router.
type RouterWrapper func(r Router) Router

// RouteWrapper wraps a route with extra functionality.
// It is passed in when creating a new route.
type RouteWrapper func(r Route) Route

// Routes returns an array of routes
func (p preRouter) Routes() []Route {
	return p.routes
}

// Status returns whether the router should be enabled.
func (p preRouter) Status() bool {
	return p.status
}

// Function returns the function route applies.
func (p preRoute) Handler() http.HandlerFunc {
	return p.handler
}

// Method returns the http method that the route responds to.
func (p preRoute) Method() string {
	return p.method
}

// Path returns the subpath where the route responds to.
func (p preRoute) Path() string {
	return p.path
}

// Status returns whether the route should be enabled.
func (p preRoute) Status() bool {
	return p.status
}

// Experimental returns whether the route is enabled.
func (p preRoute) Experimental() bool {
	return p.experimental
}

// NewRouter initializes a new local router for the system.
func NewRouter(routes []Route, status bool, opts ...RouterWrapper) Router {
	var r Router = preRouter{status, routes}
	for _, o := range opts {
		r = o(r)
	}
	return r
}

// NewRoute initializes a new local route for the router.
func NewRoute(method string, path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	var r Route = preRoute{method, path, status, experimental, handler}
	for _, o := range opts {
		r = o(r)
	}
	return r
}

// NewGetRoute initializes a new route with the http method GET.
func NewGetRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodGet, path, status, experimental, handler, opts...)
}

// NewPostRoute initializes a new route with the http method POST.
func NewPostRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPost, path, status, experimental, handler, opts...)
}

// NewPutRoute initializes a new route with the http method PUT.
func NewPutRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPut, path, status, experimental, handler, opts...)
}

// NewDeleteRoute initializes a new route with the http method DELETE.
func NewDeleteRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodDelete, path, status, experimental, handler, opts...)
}

// NewOptionsRoute initializes a new route with the http method OPTIONS.
func NewOptionsRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodOptions, path, status, experimental, handler, opts...)
}

// NewHeadRoute initializes a new route with the http method HEAD.
func NewHeadRoute(path string, status bool, experimental bool, handler http.HandlerFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodHead, path, status, experimental, handler, opts...)
}