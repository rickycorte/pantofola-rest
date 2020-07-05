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
	"time"
)

// Middleware is a chain component that definse a path execution mode
type Middleware func(r *Request, next *MiddlewareChain)

// structs

// MiddlewareChain is a linked list that makes possibile parallel executions
// without holding a state for every single connection we have on the sever
// we just move forword in the list and call all the middlewares
type MiddlewareChain struct {
	current Middleware
	next    *MiddlewareChain
	handler RequestHandler
}

//*********************************************************************************************************************

// Next runs the next middleware or the hander of the chain
// use this function to move forword in the chain
// or never call this to simply stop the chain execution
func (mc *MiddlewareChain) Next(req *Request) {
	if mc.current == nil {
		mc.handler(req)
	} else {
		mc.current(req, mc.next)
	}
}

// MakeMiddlewareChain creats a runnable chain of middlewares
// passing a nil middleware list will produce a chain that is the same as directly running the handler
func MakeMiddlewareChain(middlewares []Middleware, handler RequestHandler) *MiddlewareChain {
	if middlewares == nil {
		return &MiddlewareChain{handler: handler}
	}

	var first, last *MiddlewareChain
	for i := 0; i < len(middlewares); i++ {

		mc := &MiddlewareChain{current: middlewares[i], next: nil, handler: nil}

		if i == 0 {
			first = mc // the first must be saved
		} else {
			last.next = mc // update last chain item to poin the new one created
		}

		last = mc
	}

	// append handler to last
	last.next = &MiddlewareChain{handler: handler}

	return first // return head of list
}

//*********************************************************************************************************************

// LogRequestInfoMiddleware prints a compact formateted log about received http requests
func LogRequestInfoMiddleware(r *Request, next *MiddlewareChain) {
	start := time.Now()
	next.Next(r)
	elapsed := time.Since(start)
	log.Printf("HTTP %s %s - %d in %dms\n", r.reader.Method, r.reader.URL, r.status, elapsed.Milliseconds())
}
