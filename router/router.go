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
	"strings"
)

const (
	httpGET          = 0
	httpPOST         = 1
	httpPUT          = 2
	httpPATCH        = 3
	httpDELETE       = 4
	httpTotalMethods = 5
)

const poolSize = 500
const maxPoolSize = 10000

// RequestHandler is a direct function call that handles a http request
type RequestHandler func(http.ResponseWriter, *http.Request, *ParameterList)

// contrainer of a order list of path nodes all with the same size
type pathContainer []*pathNode

// Router type holds all the data to handle and correctly route requests to direct
// handler or other subrouters
type pathNode struct {
	staticRoutes     []pathContainer
	parameterHandler *pathNode
	handler          RequestHandler
	name             string
}

// Router is the main block of the api and hold all registered paths
// this shoul be used instead of the default server mux
type Router struct {
	pathTrees        [httpTotalMethods]*pathNode
	index            RequestHandler
	fallback         RequestHandler
	notAllowedMethod RequestHandler
	maxParamters     int
	paramPool        ParametersPool
}

//*********************************************************************************************************************
// pathContainer

// get an element with subpath es: "/api" from a container with binary search
func (pc *pathContainer) get(subpath string) *pathNode {

	// binary search
	low := 0
	high := len(*pc) - 1
	var mid, cmp int

	for low <= high {
		mid = (low + high) / 2
		cmp = strings.Compare(subpath, (*pc)[mid].name)
		if cmp == 0 {
			return (*pc)[mid]
		} else if cmp < 0 {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return nil
}

// create or set an item with specified subpath
func setStaticNode(pc pathContainer, subpath string, node *pathNode) pathContainer {

	if node == nil {
		panic("Node must not be nil")
	}

	node.name = subpath

	// inserion sort
	pc = append(pc, node)
	i := 0
	for i = len(pc) - 1; i > 0 && strings.Compare(pc[i].name, pc[i-1].name) < 0; i-- {
		temp := pc[i]
		pc[i] = pc[i-1]
		pc[i-1] = temp
	}

	return pc
}

//*********************************************************************************************************************
// router

// generate a tree from a path and a method
func (r *Router) setPath(method int, path string, handler RequestHandler) {

	currentNode := r.pathTrees[method]

	// create first node if not exist
	if currentNode == nil {
		currentNode = &pathNode{}
		r.pathTrees[method] = currentNode
	}

	if method == -1 {
		panic("Unsupported method for path " + path)
	}

	if path == "/" {
		r.index = handler
		return
	}

	lastSlash := 0
	paramCount := 0
	pathSize := len(path)

	for i := 1; i < pathSize; i++ {
		// do something only when a / is found (or end of url is reached)
		if path[i] == '/' || i == pathSize-1 {
			var sch string
			// grab the paramter from the url with a slice
			if i != pathSize-1 {
				sch = path[lastSlash:i]
			} else {
				sch = path[lastSlash:]
			}

			// check if this is a parameter
			if len(sch) > 2 && sch[0:2] == "/:" {

				if currentNode.parameterHandler != nil && currentNode.parameterHandler.name != sch[2:] {
					panic("Paramter name mismatch for route " + path)
				}
				//check if there is no handler
				if currentNode.parameterHandler == nil {
					currentNode.parameterHandler = &pathNode{name: sch[2:]}
				}
				currentNode = currentNode.parameterHandler
				paramCount++
			} else {

				subSize := len(sch)
				// allocate or reallocate array
				if currentNode.staticRoutes == nil || subSize >= len(currentNode.staticRoutes) {
					temp := make([]pathContainer, subSize-len(currentNode.staticRoutes)+5)
					currentNode.staticRoutes = append(currentNode.staticRoutes, temp...)
				}
				// static part of path
				container := currentNode.staticRoutes[subSize]
				// check if exist and add a new element if there is node
				if container == nil || container.get(sch) == nil {
					currentNode.staticRoutes[subSize] = setStaticNode(container, sch, &pathNode{})
				}

				currentNode = currentNode.staticRoutes[subSize].get(sch)
			}

			// time to assign the handler only to last node
			if i == pathSize-1 {
				currentNode.handler = handler
			}
			lastSlash = i
		}
	}

	if paramCount > r.maxParamters {
		r.maxParamters = paramCount
		r.paramPool.Init(paramCount, poolSize, maxPoolSize)
	}

}

// parse a request url and call the right handler
func (r *Router) executeHandler(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Path
	method := methodToInt(req.Method)
	var parameters *ParameterList

	if method == -1 {
		r.notAllowedMethod(w, req, nil)
		log.Println("Method not allowed: " + req.Method)
		return
	}

	// return index page
	if len(url) == 0 || url == "/" {
		r.index(w, req, nil)
		return
	}

	currentNode := r.pathTrees[method]
	// check if there is an handler for the request method
	if currentNode == nil {
		r.notAllowedMethod(w, req, nil)
		log.Println("Method not allowed, no handler set for: " + req.Method)
		return
	}

	lastSlash := 0
	// start from one to skip the first /
	size := len(url)
	var sch string

	for i := 1; i < size; i++ {
		// do something only when a / is found (or end of url is reached)
		if url[i] == '/' || i == size-1 {
			// grab the paramter from the url with a slice
			if i != size-1 {
				sch = url[lastSlash:i]
			} else {
				sch = url[lastSlash:]
			}

			// first search static nodes
			var staticNode *pathNode
			if currentNode.staticRoutes != nil {
				sz := len(sch)
				if sz < len(currentNode.staticRoutes) {
					staticNode = currentNode.staticRoutes[sz].get(sch)
				} else {
					staticNode = nil
				}
			} else {
				staticNode = nil
			}

			if staticNode != nil {
				currentNode = staticNode
			} else if currentNode.parameterHandler != nil { // then check if the value could be a paramter
				currentNode = currentNode.parameterHandler
				if parameters == nil {
					parameters = r.paramPool.Get()
				}
				parameters.Set(currentNode.name, sch[1:])
			} else { // in nothing is found then we can print a not found message
				r.fallback(w, req, nil) // not found any possible match
				return
			}

			lastSlash = i // update last slash pos after operations
		}
	}

	// when we are here we are in the last node of the url so we can execute the action
	if currentNode.handler != nil {
		currentNode.handler(w, req, parameters)
	} else {
		r.fallback(w, req, nil)
	}
	r.paramPool.Push(parameters)
}

//*********************************************************************************************************************

// MakeRouter creates a new router with no middleware and with default index and error pages
func MakeRouter() *Router {
	r := &Router{fallback: defaultFallback, index: defaultIndex, notAllowedMethod: defaultNotAllowedMethod}
	return r
}

// SetFallback sets a custom error page used when a element is not found
func (r *Router) SetFallback(handler RequestHandler) {
	r.fallback = handler
}

// SetNotAllowedHandler sets a custom error page used when a unsupported method is received
func (r *Router) SetNotAllowedHandler(handler RequestHandler) {
	r.notAllowedMethod = handler
}

// Handle adds (or reset) a route handler
func (r *Router) Handle(method, path string, handler RequestHandler) {
	r.setPath(methodToInt(method), path, handler)
}

// GET sets a request handler for the specified url only for GET requests
// this is equivalent to call Handle("GET", ...)
func (r *Router) GET(path string, handler RequestHandler) {
	r.setPath(httpGET, path, handler)
}

// POST sets a request handler for the specified url only for POST requests
// this is equivalent to call Handle("POST", ...)
func (r *Router) POST(path string, handler RequestHandler) {
	r.setPath(httpPOST, path, handler)
}

// PATCH sets a request handler for the specified url only for PATCH requests
// this is equivalent to call Handle("PATCH", ...)
func (r *Router) PATCH(path string, handler RequestHandler) {
	r.setPath(httpPATCH, path, handler)
}

// PUT sets a request handler for the specified url only forPUT requests
// this is equivalent to call Handle("PUT", ...)
func (r *Router) PUT(path string, handler RequestHandler) {
	r.setPath(httpPUT, path, handler)
}

// DELETE sets a request handler for the specified url only for DELETE requests
// this is equivalent to call Handle("DELETE", ...)
func (r *Router) DELETE(path string, handler RequestHandler) {
	r.setPath(httpDELETE, path, handler)
}

// ServeHTTP implements http.handler interface to allow this router to be easly used with std server
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.executeHandler(w, req)
}
