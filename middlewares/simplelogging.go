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

package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/rickycorte/pantofola-rest/router"
)

// SimpleRequestLogging is a single handler middleware that logs base information about the executed handler
// like the url, request status, process time
func SimpleRequestLogging(handler router.RequestHandler) router.RequestHandler {

	return func(w http.ResponseWriter, r *http.Request, p *router.ParameterList) {
		start := time.Now().UnixNano()
		writer := &captureResponceWriter{w, http.StatusOK}
		handler(writer, r, p)
		delta := time.Now().UnixNano() - start
		log.Printf("HTTP %s %s - %d in %.2fms\n", r.Method, r.URL, writer.status, float64(delta)/1000000)
	}
}

type logginRouter struct {
	inner http.Handler
}

func (lr *logginRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()
	writer := &captureResponceWriter{w, 0}
	lr.inner.ServeHTTP(writer, r)
	delta := time.Now().UnixNano() - start
	log.Printf("HTTP %s %s - %d in %.2fms\n", r.Method, r.URL, writer.status, float64(delta)/1000000)
}

// GlobalSimpleRequestLogging is a router middleware that logs base information about all the requests passed to the router
// like the url, request status, process time
func GlobalSimpleRequestLogging(router http.Handler) http.Handler {

	return &logginRouter{inner: router}
}
