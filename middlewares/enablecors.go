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
	"net/http"

	"github.com/rickycorte/pantofola-rest/router"
)

// Cors adds access-control header to allow cors request to this handler
func Cors(handler router.RequestHandler) router.RequestHandler {

	return func(w http.ResponseWriter, r *http.Request, p *router.ParameterList) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler(w, r, p)
	}
}

type corsRouter struct {
	inner     http.Handler
	preflight bool
	methods   string
	headers   string
	maxAge    int
}

func (lr *corsRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if lr.preflight && r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", lr.methods)
		w.Header().Set("Access-Control-Allow-Headers", lr.headers)
		w.Header().Set("Access-Control-Max-Age", string(lr.maxAge))
		w.WriteHeader(204) // no body
		return
	}

	lr.inner.ServeHTTP(w, r)
}

// GlobalCors create a router with enable cors for every path
func GlobalCors(router http.Handler) http.Handler {
	return &corsRouter{inner: router}
}

// compile the header into a single line string
func compileHeader(data []string) string {
	compiledHeader := ""
	for i := 0; data != nil && i < len(data); i++ {
		compiledHeader += data[i]
		if i != len(data)-1 {
			compiledHeader += ","
		}
	}
	return compiledHeader
}

// GlobalCorsPreflight enables cors requests but also enables preflight OPTIONS requests with custom
// data passed as parameter
func GlobalCorsPreflight(router http.Handler, methods, headers []string) http.Handler {
	return &corsRouter{inner: router, preflight: true, methods: compileHeader(methods), headers: compileHeader(headers), maxAge: 86400}
}
