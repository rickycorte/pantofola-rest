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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// write api| in body
func apiMiddleware(r *Request, next *MiddlewareChain) {
	r.writer.Write([]byte("api|"))
	next.Next(r)
}

// this handler writes 404| in body
// and set return value accordigly
func otherhandler(r *Request) {

	r.writer.Write([]byte("404|"))
	r.status = 404
	r.isHandled = true
}

func runRequest(router *Router, path, expected string, t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)

	router.ServeHTTP(recorder, req)

	b, _ := ioutil.ReadAll(recorder.Body)
	body := string(b)

	if body != expected {
		t.Errorf("Mismatch in execution of %s. Expected: %s, got: %s", path, expected, body)
	}

}

func TestRouter(t *testing.T) {

	router := MakeRouter()

	if router.fallbackHandler != nil {
		t.Errorf("Router fallback route must be empty by default")
	}

	//setup all routers
	router.AddHandler("/ah", simplehandler)

	router.AddHandlerChain("/ah2", nil, simplehandler)

	router.AddHandlerChain("/ch", []Middleware{middleware300, middleware400}, simplehandler)

	router.SetFallbackHandler(otherhandler)

	runRequest(router, "/ah", "200|", t)
	runRequest(router, "/ah2", "200|", t)                // same as vanilla handler becase no middleware is passed
	runRequest(router, "/ch", "300|400|200|400|300|", t) // middleware action
	runRequest(router, "/not-found", "404|", t)          // fallback

}

func TestRouterWithSubRoutingAndMiddlwares(t *testing.T) {

	router := MakeRouter()
	sub := MakeRouter()

	//setup all routers
	router.SetFallbackHandler(simplehandler)

	sub.UseMiddleware(apiMiddleware)

	sub.AddHandler("/ah", simplehandler)
	sub.AddHandlerChain("/ch", []Middleware{middleware300, middleware400}, simplehandler)
	sub.SetFallbackHandler(otherhandler)

	router.AddSubRouter("/api", sub)

	runRequest(router, "/ah", "200|", t) // fallback on main

	runRequest(router, "/api/ah", "api|200|", t)                 // same as vanilla handler becase no middleware is passed
	runRequest(router, "/api/ch", "api|300|400|200|400|300|", t) // middleware action
	runRequest(router, "/api/not-found", "api|404|", t)          // fallback

}
