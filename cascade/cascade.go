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
	"net/http"
	"regexp"
	"strings"

	"github.com/rickycorte/pantofola-rest/router"
)

// Handler usable from a cascade router
type Handler interface {
	http.Handler
	// UsePrefix is used to set the path prefix
	// that the sub router should ignore
	UsePrefix(string)
}

// CascadeRouter is a way to nest routers with prefixies
// the main use of there is in large application that should have different behaviours
// for example we want to run an api endpoint on "/api" prefix and also serve static files on "/"
// cascade routers are compatible with base global middlewares and can be nested at infinite levels to create really complex
// and heterogeneous behaviours
// note: cascade routers CHANGE the request path passed to lower levels by removing their assigned prefix
type CascadeRouter struct {
	mainRouter Handler // main router has no prefix and will be used if no match is found in sub routers
	subRouters map[string]Handler
	prefix     string
}

//*********************************************************************************************************************

// MakeCascade cretes an empty cascade router
func MakeCascade() *CascadeRouter {
	return &CascadeRouter{}
}

// Set set the router handler for a prefix
// pass "" as prefix to set main router that will be used in case of no match with other prefiex
// if no main router is set an empy default router will be used
// note that prefix must not contain sub bats eg: /path/to/me is no a valid prefix
// only /path is valid, to archive the path above combina a router that handes relative path /to/me
// or use multiple cascade routers nested
func (cr *CascadeRouter) Set(prefix string, rout Handler) {

	if prefix == "" {
		cr.mainRouter = rout
		return
	}

	if !regexp.MustCompile("^/[a-zA-Z0-9]+$").MatchString(prefix) {
		panic("Please use a top level prefix eg: /prefix")
	}

	if cr.subRouters == nil {
		cr.subRouters = make(map[string]Handler)
	}

	rout.UsePrefix(cr.prefix + prefix)

	cr.subRouters[prefix] = rout

	if cr.mainRouter == nil {
		cr.mainRouter = router.MakeRouter()
	}
}

// ServeHTTP searches the most approprate router or cascade to use
func (cr *CascadeRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	url := strings.TrimPrefix(req.URL.Path, cr.prefix)
	secondSlash := -1
	if len(url) > 1 {
		secondSlash = strings.Index(url[1:], "/") // skip first char that is a /
	}

	var prefix string
	// grab prefix and remove it from path
	if secondSlash == -1 {
		prefix = url
	} else {
		prefix = url[:secondSlash+1]
	}
	r := cr.mainRouter

	if cr.subRouters != nil {
		if temp := cr.subRouters[prefix]; temp != nil {
			r = temp
		}
	}

	r.ServeHTTP(w, req)
}

// UsePrefix set a path prefix that should be removed before parsing the request
func (cr *CascadeRouter) UsePrefix(prefix string) {
	cr.prefix = prefix
	// propagate prefix changes to childs
	for k, v := range cr.subRouters {
		v.UsePrefix(cr.prefix + k)
	}
}
