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
	"net/http"
)

// default error page
func defaultFallback(w http.ResponseWriter, _ *http.Request, _ *ParameterList) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found")
}

// default not allowed method responce
func defaultNotAllowedMethod(w http.ResponseWriter, _ *http.Request, _ *ParameterList) {
	w.WriteHeader(405)
	fmt.Fprintf(w, "Method Not Allowed")
}

// default panic handler
func defaultPanicHandler(w http.ResponseWriter, _ *http.Request, _ interface{}) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "Something went wrong with your request")
}
