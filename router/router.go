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
	"fmt"
	"log"
	"net/http"
)

const (
	httpGET          = 0
	httpPOST         = 1
	httpPUT          = 2
	httpPATCH        = 3
	httpDELETE       = 4
	httpTotalMethods = 5
)

// RequestHandler2 is a direct function call that handles a http request
type RequestHandler2 func(http.ResponseWriter, *http.Request, map[string]string)

// Router type holds all the data to handle and correctly route requests to direct
// handler or other subrouters
type pathNode struct {
	staticRoutes     map[string]*pathNode
	parameterHandler *pathNode
	handler          RequestHandler2
	parameterName    string
}

// Router is the main block of the api and hold all registered paths
// this shoul be used instead of the default server mux
type Router struct {
	pathTrees    [httpTotalMethods]*pathNode
	index        RequestHandler2
	fallback     RequestHandler2
	maxParamters int
}

//*********************************************************************************************************************

// generate a tree from a path and a method
func (r *Router) setPath(method int, path string, handler RequestHandler2) {

	//log.Println("Parsing " + path)

	currentNode := r.pathTrees[method]

	// create first node if not exist
	if currentNode == nil {
		currentNode = &pathNode{}
		r.pathTrees[method] = currentNode
	}

	if method == -1 {
		panic("Unsupported method for path " + path)
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

				if currentNode.parameterHandler != nil && currentNode.parameterHandler.parameterName != sch[2:] {
					panic("Paramter name mismatch for route " + path)
				}
				//check if there is no handler
				if currentNode.parameterHandler == nil {
					currentNode.parameterHandler = &pathNode{parameterName: sch[2:]}
				}
				currentNode = currentNode.parameterHandler
				paramCount++
			} else {
				// map need to be generated
				if currentNode.staticRoutes == nil {
					currentNode.staticRoutes = make(map[string]*pathNode)
				}
				// static part of path
				existNode := currentNode.staticRoutes[sch]
				// check if exist
				if existNode == nil {
					currentNode.staticRoutes[sch] = &pathNode{}
				}

				currentNode = currentNode.staticRoutes[sch]
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
	}

}

// convert method from string into a mapped int
func methodToInt(method string) int {
	if len(method) < 3 {
		return -1
	}
	c := method[0]
	if c == 'G' {
		return httpGET
	}
	if c == 'P' {
		m := method[1]
		if m == 'O' {
			return httpPOST
		}
		if m == 'U' {
			return httpPUT
		}
		return httpPATCH
	}
	if c == 'D' {
		return httpDELETE
	}
	return -1
}

// parse a request url and call the right handler
func (r *Router) executeHandler(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Path
	method := methodToInt(req.Method)
	var parameters map[string]string

	if method == -1 {
		//TODO: custimuze error function
		w.WriteHeader(400)
		fmt.Fprintf(w, "Unsupported method")
		log.Fatal("Unsupported method: " + req.Method)
		return
	}

	//log.Println("Executing: " + req.Reader.Method + " " + url)

	// return index page
	if len(url) == 0 || url == "/" {
		r.index(w, req, nil)
		return
	}

	currentNode := r.pathTrees[method]
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
				if parameters == nil {
					parameters = make(map[string]string)
				}
				parameters[currentNode.parameterName] = sch[1:]
			} else { // in nothing is found then we can print a not found message
				r.fallback(w, req, nil) // not found any possible match
				//log.Println("Not found: " + sch)
				return
			}

			lastSlash = i // update last slash pos after operations
		}
	}

	// when we are here we are in the last node of the url so we can execute the action
	currentNode.handler(w, req, parameters)
}

// default error page
func defaultFallback(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found")
}

// default index page
func defaultIndex(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Welcome to Pantofola-Rest!")
}

//*********************************************************************************************************************

// MakeRouter creates a new router with no middleware and with default index and error pages
func MakeRouter() *Router {
	r := &Router{fallback: defaultFallback, index: defaultIndex}
	return r
}

// SetFallback sets the default error page used when a element is not found
func (r *Router) SetFallback(handler RequestHandler2) {
	r.fallback = handler
}

// SetIndex sets the default error page
func (r *Router) SetIndex(handler RequestHandler2) {
	r.index = handler
}

// Handle adds (or reset) a route handler
func (r *Router) Handle(method, path string, handler RequestHandler2) {
	r.setPath(methodToInt(method), path, handler)
}

// ServeHTTP implements http.handler interface to allow this router to be easly used with std server
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.executeHandler(w, req)
}
