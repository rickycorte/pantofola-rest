/*
   Copyright 2020 rickycorte

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package router

import (
	"log"
	"net/http"
	"regexp"
)

const (
	httpGET          = 0
	httpPOST         = 1
	httpPUT          = 2
	httpPATCH        = 3
	httpDELETE       = 4
	httpTotalMethods = 5
)

// Router type holds all the data to handle and correctly route requests to direct
// handler or other subrouters
type pathNode struct {
	staticRoutes          map[string]*pathNode
	parameterHandler      *pathNode
	directMiddlewareChain *MiddlewareChain
	parameterName         string
}

// Router is the main block of the api and hold all registered paths
// this shoul be used instead of the default server mux
type Router struct {
	pathTrees       [httpTotalMethods]*pathNode
	middlewares     []Middleware
	middlewareChain *MiddlewareChain
	index           RequestHandler
	fallback        RequestHandler
}

//*********************************************************************************************************************

// generate a tree from a path and a method
func (r *Router) setPath(method int, path string, middlewares []Middleware, handler RequestHandler) {

	// split path in /*** data
	regex := regexp.MustCompile(`((/:?[a-zA-Z0-9]+))`)
	matches := regex.FindAllString(path, -1)

	//log.Println("Parsing " + path)

	currentNode := r.pathTrees[method]

	// create first node if not exist
	if currentNode == nil {
		currentNode = &pathNode{}
		r.pathTrees[method] = currentNode
	}

	for i, v := range matches {

		// check if this is a parameter
		if len(v) > 2 && v[0:2] == "/:" {

			if currentNode.parameterHandler != nil && currentNode.parameterHandler.parameterName != v[2:] {
				panic("Paramter name mismatch for route " + path)
			}
			//check if there is no handler
			if currentNode.parameterHandler == nil {
				currentNode.parameterHandler = &pathNode{parameterName: v[2:]}
				//log.Println("Created parameter node: " + v[2:])
			}
			currentNode = currentNode.parameterHandler
			//log.Println("Added paramter path node: " + v[2:])

		} else {
			// map need to be generated
			if currentNode.staticRoutes == nil {
				currentNode.staticRoutes = make(map[string]*pathNode)
			}
			// static part of path
			existNode := currentNode.staticRoutes[v]
			// check if exist
			if existNode == nil {
				currentNode.staticRoutes[v] = &pathNode{}
				//log.Println("Created node: " + v)
			}

			currentNode = currentNode.staticRoutes[v]
			//log.Println("Added path node: " + v)
		}

		// time to assign the handler only to last node
		if i == len(matches)-1 {
			currentNode.directMiddlewareChain = MakeMiddlewareChain(middlewares, handler)
		}
	}

}

// convert method from string into a mapped int
func methodToInt(method string) int {
	switch method {
	case "GET":
		return httpGET
	case "POST":
		return httpPOST
	case "PUT":
		return httpPUT
	case "PATCH":
		return httpPATCH
	case "delete":
		return httpDELETE
	default:
		return -1
	}
}

// parse a request url and call the right handler
func (r *Router) executeHandler(req *Request) {

	url := req.reader.URL.Path
	method := methodToInt(req.reader.Method)
	if method == -1 {
		//TODO: custimuze error function
		req.Reply(500, nil, "Unsupported method")
		log.Fatal("Unsupported method: " + req.reader.Method)
		return
	}

	//log.Println("Executing: " + req.reader.Method + " " + url)

	// return index page
	if len(url) == 0 || url == "/" {
		r.index(req)
		return
	}

	currentNode := r.pathTrees[method]
	lastSlash := 0
	// start from one to skip the first /
	for i := 1; i < len(url); i++ {
		// do something only when a / is found (or end of url is reached)
		if url[i] == '/' || i == len(url)-1 {
			var sch string
			// grab the paramter from the url with a slice
			if i != len(url)-1 {
				sch = url[lastSlash:i]
			} else {
				sch = url[lastSlash:]
			}
			//log.Print("Searching node: " + sch)

			// first search static nodes
			var staticNode *pathNode
			if currentNode.staticRoutes != nil {
				staticNode = currentNode.staticRoutes[sch]
			} else {
				staticNode = nil
			}

			if staticNode != nil {
				currentNode = staticNode
				//log.Println("Moved to static node: " + sch)
			} else if currentNode.parameterHandler != nil { // then check if the value could be a paramter
				currentNode = currentNode.parameterHandler
				//log.Println("Added new parameter: " + sch[1:])
				req.SetParameter(currentNode.parameterName, sch[1:])
			} else { // in nothing is found then we can print a not found message
				r.fallback(req) // not found any possible match
				return
				//log.Println("Not found: " + sch[1:])
			}

			lastSlash = i // update last slash pos after operations
		}
	}

	// when we are here we are in the last node of the url so we can execute the action
	if currentNode.directMiddlewareChain != nil {
		currentNode.directMiddlewareChain.Next(req)
	} else {
		log.Fatal("Missing direct handler for " + url)
	}
}

// default error page
func defaultFallback(req *Request) {
	req.Reply(404, nil, "Not found")
}

// default index page
func defaultIndex(req *Request) {
	req.Reply(200, nil, "Welcome to Pantofola-Rest!")
}

//*********************************************************************************************************************

// MakeRouter creates a new router with no middleware and with default index and error pages
func MakeRouter() *Router {
	r := &Router{fallback: defaultFallback, index: defaultIndex}
	r.middlewareChain = MakeMiddlewareChain(nil, r.executeHandler)
	return r
}

// Use adds a global middleware to the router
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
	r.middlewareChain = MakeMiddlewareChain(r.middlewares, r.executeHandler)
}

// SetFallback sets the default error page used when a element is not found
func (r *Router) SetFallback(handler RequestHandler) {
	r.fallback = handler
}

// SetIndex sets the default error page
func (r *Router) SetIndex(handler RequestHandler) {
	r.index = handler
}

func (r *Router) GET(path string, middlewares []Middleware, handler RequestHandler) {
	r.setPath(httpGET, path, middlewares, handler)
}

func (r *Router) POST(path string, middlewares []Middleware, handler RequestHandler) {
	r.setPath(httpPOST, path, middlewares, handler)
}

func (r *Router) PUT(path string, middlewares []Middleware, handler RequestHandler) {
	r.setPath(httpPUT, path, middlewares, handler)
}

func (r *Router) PATCH(path string, middlewares []Middleware, handler RequestHandler) {
	r.setPath(httpPATCH, path, middlewares, handler)
}

func (r *Router) DELETE(path string, middlewares []Middleware, handler RequestHandler) {
	r.setPath(httpDELETE, path, middlewares, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router.middlewareChain.Next(MakeRequest(w, r))
}
