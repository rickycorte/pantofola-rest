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

package staticfiles

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/rickycorte/pantofola-rest/router"
)

// Serve wraps default go file serve but adds a custom error handler if file is not found
// it also allows to redirect system files to another path by setting diffentet prefix and systemPath
// prefix is set when adding this handler to a router
func Serve(systemPath string, notFound router.RequestHandler) router.RequestHandler {

	return func(w http.ResponseWriter, r *http.Request, p *router.ParameterList) {

		target := "/index.html"

		if p != nil {
			target = path.Clean(p.Get("*"))
		}

		target = systemPath + target

		_, err := os.Stat(target)
		log.Printf("Seaching: %s\n", target)
		if os.IsNotExist(err) {
			notFound(w, r, nil)
		} else {
			http.ServeFile(w, r, target)
		}
	}

}
