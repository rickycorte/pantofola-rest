![Go Test](https://github.com/rickycorte/pantofola-rest/workflows/Go%20Test/badge.svg)
[![codecov](https://codecov.io/gh/rickycorte/pantofola-rest/branch/master/graph/badge.svg)](https://codecov.io/gh/rickycorte/pantofola-rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickycorte/pantofola-rest)](https://goreportcard.com/report/github.com/rickycorte/pantofola-rest)
[![GoDoc](https://godoc.org/github.comrickycorte/pantofola-rest?status.svg)](http://godoc.org/github.com/rickycorte/pantofola-rest)

# Pantofola REST

Simple and fast framework to create awesome REST APIs in seconds!

## Features

- Incredible fast routing
- Optimized for dinamic path with multiple parameters
- Namad paramters
- Parameter pool for 0 allocations and max speed

All this speed comes to a cost, a good amout of used memory due to the pool and a really slow initialization process.

Please notice that adding path after initialization is not a good practice because it leads to temporary performance degradation. Inizialitation is not designed to be fast, all the speed comes after the cost of booting everything!

## Road Map

Feature | Version 
--- | ---  
Fast routing | 0.1.0
Named parameters | 0.1.0 
Parameter pool | 0.1.0 
Custom not found handler | 0.1.0 
Custom not allowed method handler | 0.1.0 
 | 
Custom pool settings | TBD
Panic Handler | TBD
Path Builder | TBD
Serve files with caching | TBD
Built-in middlwares | TBD


# Example API

This simple code is just a sample to demostrate how simple and clean is the code to create a REST API with `pantofola-rest`.


```go
import (
	"fmt"
	"net/http"
	"github.com/rickycorte/pantofola-rest/router"
)

// hello is our first handler!
func hello(w http.ResponseWriter, _ *http.Request, _ *ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "HI sen(pi) <3")
}

// generic handler for a parametric route that prints available values by name
func writeData(w http.ResponseWriter, r *http.Request, p *pantofola.ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, r.RequestURI+"\n\nUser: "+p.Get("user")+"\nActivity: "+p.Get("activity")+"\nComment: "+p.Get("comment"))
}

func main() {

    mainRouter := router.MakeRouter() // create the main router used by our app
	
	// we add our handler as the main router
	mainRouter.GET("/hello", hello) 
	
	// let's make some parametric handlers
	rout.GET("/activity/:user", writeData)
	rout.GET("/activity/:user/:activity", writeData)
	rout.GET("/activity/:user/:activity/comments/:comment", writeData)

    http.ListenAndServe(":8080", mainRouter) // start the server end enjoy your REST API!
}

```

Visit [wiki pages](https://github.com/rickycorte/pantofola-rest/wiki) to learn how to use all the features of `pantofola-rest`