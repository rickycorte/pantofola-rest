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

// convert method from string into a mapped int
func methodToInt(method string) int {
	if len(method) < 3 {
		return -1
	}
	c := method[0]
	if c == 'G' {
		return httpGET
	}
	if c == 'P' {
		m := method[1]
		if m == 'O' {
			return httpPOST
		}
		if m == 'U' {
			return httpPUT
		}
		return httpPATCH
	}
	if c == 'D' {
		return httpDELETE
	}
	return -1
}
