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

// this middleware writes 300| in body
func middleware300(r *Request, next *MiddlewareChain) {

	r.writer.Write([]byte("300|"))
	next.Next(r)
	r.writer.Write([]byte("300|"))
}

// this middleware writes 400| in body
func middleware400(r *Request, next *MiddlewareChain) {

	r.writer.Write([]byte("400|"))
	next.Next(r)
	r.writer.Write([]byte("400|"))
}

// this middleware stops the chains and return 500 code
// and writes 500| in body
func middlwareStopChain(r *Request, next *MiddlewareChain) {
	r.writer.Write([]byte("500|"))
	r.status = 500
}

// handler that sets responce code to 200
// and writes 200| in body
func simplehandler(r *Request) {
	r.writer.Write([]byte("200|"))
	r.status = 200
	r.isHandled = true
}

func TestNextExecutionWithNoKiller(t *testing.T) {

	middlewares := []Middleware{middleware300, middleware400, middleware300}

	chain := MakeMiddlewareChain(middlewares, simplehandler)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health-check", nil)

	r := MakeRequest(recorder, req)

	chain.Next(r) // run complete chain

	body, _ := ioutil.ReadAll(recorder.Body)
	sbody := string(body)
	// to understand this string see what the middleware do
	const expected = "300|400|300|200|300|400|300|"
	if sbody != expected {
		t.Errorf("Mismatched execution order of the chain. Expected: %s, got %s", expected, sbody)
	}

	if r.status != 200 {
		t.Errorf("Mismatched return code. Expected 200, got: %v", recorder.Code)
	}

}

func TestNextExecutionWithKiller(t *testing.T) {
	middlewares := []Middleware{middleware300, middlwareStopChain, middleware300}
	chain := MakeMiddlewareChain(middlewares, simplehandler)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health-check", nil)

	r := MakeRequest(recorder, req)

	chain.Next(r) // run complete chain

	const expected = "300|500|300|"

	body, _ := ioutil.ReadAll(recorder.Body)
	sbody := string(body)
	// to understand this string see what the middleware do
	if sbody != expected {
		t.Errorf("Mismatched execution order of the chain. Expected: %s, got %s", expected, sbody)
	}

	if r.status != 500 {
		t.Errorf("Mismatched return code. Expected 500, got: %v", recorder.Code)
	}

}

func TestLogRequestInfoMiddleware(t *testing.T) {
	middlewares := []Middleware{LogRequestInfoMiddleware}
	chain := MakeMiddlewareChain(middlewares, simplehandler)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health-check", nil)

	r := MakeRequest(recorder, req)

	chain.Next(r) // run complete chain

	const expected = "200|"

	body, _ := ioutil.ReadAll(recorder.Body)
	sbody := string(body)
	// to understand this string see what the middleware do
	if sbody != expected {
		t.Errorf("Mismatched execution order of the chain. Expected: %s, got %s", expected, sbody)
	}

	if r.status != 200 {
		t.Errorf("Mismatched return code. Expected 500, got: %v", recorder.Code)
	}
}
