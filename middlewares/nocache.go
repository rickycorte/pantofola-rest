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

// NoCache sets the header Cache-Control to no-cache for the path
func NoCache(w http.ResponseWriter, r *http.Request, p *router.ParameterList, handler router.RequestHandler) router.RequestHandler {

	return func(w http.ResponseWriter, r *http.Request, p *router.ParameterList) {
		w.Header().Set("Cache-Control", "no-cache")
		handler(w, r, p)
	}
}

type noCacheRouter struct {
	inner http.Handler
}

func (lr *noCacheRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	lr.inner.ServeHTTP(w, r)
}

// GlobalNoCache returs a router with Cache-Control header set to no-cache for all the paths
func GlobalNoCache(router http.Handler) http.Handler {

	return &noCacheRouter{inner: router}
}
