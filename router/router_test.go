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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func runRequest(router *Router, method, path, expected string, t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)

	router.ServeHTTP(recorder, req)

	b, _ := ioutil.ReadAll(recorder.Body)
	body := string(b)

	if body != expected {
		t.Errorf("Mismatch in execution of %s. Expected: %s, got: %s", path, expected, body)
	}

}

func printHello(w http.ResponseWriter, _ *http.Request, _ *ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "hello")
}

func writeData(w http.ResponseWriter, r *http.Request, p *ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, p.Get("user")+"-"+p.Get("activity")+"-"+p.Get("comment"))
}

func TestStaticRoutes(t *testing.T) {
	router := MakeRouter()
	router.Handle("GET", "/", printHello)
	router.GET("/static/path/to/hello", printHello)
	router.Handle("GET", "/static/path/hello.html", printHello)

	runRequest(router, "GET", "/", "hello", t)
	runRequest(router, "GET", "/static/path/to/hello", "hello", t)
	runRequest(router, "GET", "/static/path/to/hello.html", "Not Found", t)
	runRequest(router, "GET", "/static/path/hello.html", "hello", t)
}

func TestParametricRoutes(t *testing.T) {
	router := MakeRouter()
	router.GET("/activity/:user", writeData)
	router.POST("/activity/:user/:activity", writeData)
	router.GET("/activity/:user/:activity/comments/:comment", writeData)

	runRequest(router, "GET", "/activity/raccoon", "raccoon--", t)
	runRequest(router, "POST", "/activity/raccoon/123", "raccoon-123-", t)
	runRequest(router, "GET", "/activity/raccoon/123/comments/456", "raccoon-123-456", t)
}

func TestNotFound(t *testing.T) {

	router := MakeRouter()
	router.GET("/zello/yes", printHello)
	router.GET("/hello", printHello)
	router.GET("/activity/:user", writeData)

	// completely wrong path
	runRequest(router, "GET", "/notFound", "Not Found", t)

	// partial match
	runRequest(router, "GET", "/zello/random", "Not Found", t)
	runRequest(router, "GET", "/activity", "Not Found", t)
	runRequest(router, "GET", "/activity/", "Not Found", t)

	//test method with no handlers
	runRequest(router, "PUT", "/activity/123", "Method Not Allowed", t)
	runRequest(router, "PATCH", "/activity/123", "Method Not Allowed", t)
	runRequest(router, "DELETE", "/activity/123", "Method Not Allowed", t)
}

func TestUnsupportedMethod(t *testing.T) {

	router := MakeRouter()
	router.GET("/zello/yes", printHello)
	router.GET("/hello", printHello)
	router.GET("/activity/:user", writeData)

	runRequest(router, "HEAD", "/activity/", "Method Not Allowed", t)
}
