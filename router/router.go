package router

import (
	"net/http"
	"regexp"
)

// RequestHandler is a direct function call that handles a http request
type RequestHandler func(*Request)

// RequestRouter defines rounter handles should be able to run
type RequestRouter interface {
	handle(r *Request)
}

// this structs holds route the pattern to match
// and the corresponding handler
type route struct {
	pattern         *regexp.Regexp
	subrouter       RequestRouter
	middlewareChain *MiddlewareChain
}

// Router type holds all the data to handle and correctly route requests to direct
// handler or other subrouters
type Router struct {
	subRouters      []route
	routes          []route
	fallbackHandler RequestHandler
	middlewareChain *MiddlewareChain
	middlewares     []Middleware
}

//*********************************************************************************************************************

// MakeRouter creates a new router and initialize its values
func MakeRouter() *Router {
	r := &Router{subRouters: make([]route, 0), routes: make([]route, 0), fallbackHandler: nil}
	r.middlewareChain = MakeMiddlewareChain(nil, r.handle) // setup chain to be able to execute handle operations
	return r
}

// AddSubRouter adds a new subrouter to the current roter with the desired prefix.
// For example we want to add all api handlers under the same prefix without writing it every time
// we just need to create a new router that handles all realtive paths
// and the pass is as a subrouter
//
// Example: we create a router that handles /test and /getme
// and set it as a subrouter of /api, with this we will obtain two handlers /api/test /api/getme
func (router *Router) AddSubRouter(pattern string, subRouter RequestRouter) {

	// match strings like /test /test/sadad
	r := route{pattern: regexp.MustCompile("^(" + pattern + ")(/.*)?$"), subrouter: subRouter}
	router.subRouters = append(router.subRouters, r)

}

// AddHandler create a standard relative handler to the Request relative path
// this function is used to add normal handlers that are not router and actually make a reply
func (router *Router) AddHandler(pattern string, handler RequestHandler) {

	router.AddHandlerChain(pattern, nil, handler)
}

// AddHandlerChain creates a middleware chain before executing the main handler
func (router *Router) AddHandlerChain(pattern string, middlewares []Middleware, handler RequestHandler) {

	r := route{pattern: regexp.MustCompile("^(" + pattern + ")(/.*)?$"), subrouter: nil,
		middlewareChain: MakeMiddlewareChain(middlewares, handler)}
	router.routes = append(router.routes, r)
}

// SetFallbackHandler sets the router to use a custom handler if nothings is found
// by default the fallback route is nill and nothing will happen if a query fails
func (router *Router) SetFallbackHandler(handler RequestHandler) {

	router.fallbackHandler = handler
}

// UseMiddleware adds a new global router middleware that will be applied to every route and subrouter
// this gives great flexibilty and power
func (router *Router) UseMiddleware(middleware Middleware) {
	if middleware == nil {
		return
	}

	router.middlewares = append(router.middlewares, middleware)
	router.middlewareChain = MakeMiddlewareChain(router.middlewares, router.handle)

}

// handle a request
// first check subrouters and then proced to fixed routes
// after all the checks return true if the route was handled and false in not
// this is used to understand then the process should stop
func (router *Router) handle(r *Request) {

	if r.relativePath == "" {
		r.relativePath = r.reader.URL.Path
	}

	// search subrouters
	for _, sr := range router.subRouters {
		if matches := sr.pattern.FindStringSubmatch(r.relativePath); len(matches) > 1 {
			// one group means have a path like /test
			if len(matches) == 1 {
				r.relativePath = "/"
			} else { // we have a subpath that is set in group 2
				r.relativePath = matches[2]
			}
			sr.subrouter.handle(r)
			// check if handled
			if r.isHandled {
				return
			}
		}
	}

	//search handlers
	//TODO: support paramters w/ caputure groups
	for _, sr := range router.routes {
		if matches := sr.pattern.FindStringSubmatch(r.relativePath); len(matches) > 0 {
			sr.middlewareChain.Next(r)
			return
		}
	}

	// nothing found fallback to default route if not nill
	if router.fallbackHandler != nil {
		router.fallbackHandler(r)
	}

}

// ServeHTTP is used to implement http.Handler interface
// to enable the router to handle direct server calls
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//router.handle(MakeRequest(w, r))
	router.middlewareChain.Next(MakeRequest(w, r))

}
