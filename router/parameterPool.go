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
	"sync"
)

// Parameter is a pair of strings
type Parameter struct {
	key   string
	value string
}

// ParameterList is a list of parameters
type ParameterList struct {
	data []Parameter
	size int
}

// ParametersPool is an allocation pool for paramters to keep allocated and not used paramters blocks
type ParametersPool struct {
	paramStack    []*ParameterList
	currentSize   int
	maxSize       int
	maxParameters int
	mutex         *sync.Mutex
}

//*********************************************************************************************************************
// ParamterList

func (pl *ParameterList) Set(key, value string) {
	pl.data[pl.size].key = key
	pl.data[pl.size].value = value
	pl.size++
}

// Get a parameter for a list by key
func (pl *ParameterList) Get(key string) string {

	for i := 0; i < pl.size; i++ {
		if pl.data[i].key == key {
			return pl.data[i].value
		}
	}

	return ""
}

//*********************************************************************************************************************
// ParamterPool

// Init allocates size
func (pp *ParametersPool) Init(maxParameters, size, maxSize int) {
	pp.maxParameters = maxParameters
	pp.currentSize = size
	pp.maxSize = maxSize
	// allocate queue
	pp.paramStack = make([]*ParameterList, size)
	// allocate all internal data
	for i := 0; i < size; i++ {
		pp.paramStack[i] = &ParameterList{data: make([]Parameter, maxParameters)}
	}

	pp.mutex = &sync.Mutex{}
}

// SetParamterSize sets the size of every paramter list
// please notice that this function  deletes the allocated pool!
func (pp *ParametersPool) SetParamterSize(size int) {
	pp.currentSize = 0
	pp.maxSize = 0
	pp.maxParameters = size
	pp.paramStack = nil
}

// Get a parameter array from the pool (or allocate a new one if none is available)
func (pp *ParametersPool) Get() *ParameterList {
	// first give a pre allocated values if there is any
	if pp.paramStack != nil && pp.currentSize > 0 {
		pp.mutex.Lock()
		r := pp.paramStack[pp.currentSize-1]
		r.size = 0
		pp.currentSize--
		pp.mutex.Unlock()
		return r
	}

	return &ParameterList{data: make([]Parameter, pp.maxParameters)}
}

// Push readds a paramter list to the pool, this will work until maxSize is reached
// then all the pushed values will be deleted by the garbage collector
func (pp *ParametersPool) Push(pl *ParameterList) {

	if pl != nil && pp.currentSize < pp.maxSize {
		pp.mutex.Lock()
		pp.paramStack[pp.currentSize] = pl
		pp.currentSize++
		pp.mutex.Unlock()
	}
}
