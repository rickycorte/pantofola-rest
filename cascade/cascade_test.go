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

package cascade

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rickycorte/pantofola-rest/router"
)

func RunRequest(router http.Handler, method, path string, status int, expected string, t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)

	router.ServeHTTP(recorder, req)

	b, _ := ioutil.ReadAll(recorder.Body)
	body := string(b)

	if recorder.Code != status {
		t.Errorf("Mismatch result code of %s. Expected: %v, got: %v", path, status, recorder.Code)
	}

	if body != expected {
		t.Errorf("Mismatch in execution of %s. Expected: %s, got: %s", path, expected, body)
	}

}

func helloApi(w http.ResponseWriter, r *http.Request, _ *router.ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "hello"+r.RequestURI)
}

func apiA(w http.ResponseWriter, r *http.Request, _ *router.ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "A "+r.URL.Path)
}

func apiB(w http.ResponseWriter, r *http.Request, _ *router.ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "B "+r.URL.Path)
}

func apiNotFound(w http.ResponseWriter, _ *http.Request, _ *router.ParameterList) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "No api")
}

func TestCascadeSearch(t *testing.T) {

	apiRouter := router.MakeRouter()
	apiRouter.GET("/", helloApi)
	apiRouter.GET("/a", apiA)
	apiRouter.GET("/b", apiB)
	apiRouter.SetNotFoundHandler(apiNotFound)

	cascade := MakeCascade()
	cascade.Set("/api", apiRouter)

	// main router that is empy
	RunRequest(cascade, "GET", "/", 405, "Method Not Allowed", t)
	RunRequest(cascade, "GET", "/notApi", 405, "Method Not Allowed", t)

	// api requests
	RunRequest(cascade, "GET", "/api", 200, "hello", t)
	RunRequest(cascade, "GET", "/api/a", 200, "A /api/a", t)
	RunRequest(cascade, "GET", "/api/random", 404, "No api", t)
}

func TestNestedCascade(t *testing.T) {

	apiRouter := router.MakeRouter()
	apiRouter.GET("/", helloApi)
	apiRouter.GET("/a", apiA)
	apiRouter.GET("/b", apiB)
	apiRouter.SetNotFoundHandler(apiNotFound)

	apiRouter2 := router.MakeRouter()
	apiRouter2.GET("/", helloApi)
	apiRouter2.GET("/a", apiA)
	apiRouter2.GET("/b", apiB)
	apiRouter2.SetNotFoundHandler(apiNotFound)

	apiCascade := MakeCascade()
	// this will change path on responces
	apiCascade.Set("/v1", apiRouter)
	apiCascade.Set("/v2", apiRouter2)

	cascade := MakeCascade()
	cascade.Set("", router.MakeRouter())
	cascade.Set("/api", apiCascade)

	// main router that is empy
	RunRequest(cascade, "GET", "/", 405, "Method Not Allowed", t)
	RunRequest(cascade, "GET", "/notApi", 405, "Method Not Allowed", t)

	// api requests
	RunRequest(cascade, "GET", "/api", 405, "Method Not Allowed", t)
	RunRequest(cascade, "GET", "/api/v1/a", 200, "A /api/v1/a", t)
	RunRequest(cascade, "GET", "/api/v2/b", 200, "B /api/v2/b", t)
	RunRequest(cascade, "GET", "/api/v1/alsdhd", 404, "No api", t)

}

func TestPanicOnWrongPrefix(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should panic with not root prefix")
		}
	}()

	cascade := MakeCascade()
	cascade.Set("/api/aldj", router.MakeRouter())

}
