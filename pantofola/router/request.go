package router

import (
	"fmt"
	"net/http"
)

// Request is a compat version of default go request/responce combo
// that offers some utility functions to make this look prettier
type Request struct {
	status       int
	reader       *http.Request
	writer       http.ResponseWriter
	relativePath string
	isHandled    bool
}

// Reply writes a responce to the request
func (req *Request) Reply(code int, headers map[string]string, body string) {

	req.status = code
	req.isHandled = true
	req.writer.WriteHeader(code)

	// set headers
	if headers != nil {
		for k, v := range headers {
			req.writer.Header().Set(k, v)
		}
	}

	fmt.Fprint(req.writer, body)
}

// MakeRequest from go internal structs
func MakeRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{status: 0, reader: r, writer: w, relativePath: "", isHandled: false}
}
