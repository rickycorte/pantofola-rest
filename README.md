![Go Test](https://github.com/rickycorte/pantofola-rest/workflows/Go%20Test/badge.svg)
[![codecov](https://codecov.io/gh/rickycorte/pantofola-rest/branch/master/graph/badge.svg)](https://codecov.io/gh/rickycorte/pantofola-rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickycorte/pantofola-rest)](https://goreportcard.com/report/github.com/rickycorte/pantofola-rest)
[![GoDoc](https://godoc.org/github.comrickycorte/pantofola-rest?status.svg)](http://godoc.org/github.com/rickycorte/pantofola-rest)

# Pantofola REST

Simple and fast framework to create awesome REST APIs in seconds!

## Features

- Incredible fast routing
- Optimized for dynamic path with multiple parameters
- Namad paramters
- Parameter pool for 0 allocations and max speed
- Middlwares (included: cors, no-cache, simple logging)
- Cascade routers for complex API

All this speed comes to a cost, a good amout of used memory due to the pool and a really slow initialization process.

Please notice that adding path after initialization is not a good practice because it leads to temporary performance degradation. Inizialitation is not designed to be fast, all the speed comes after the cost of booting everything!

# Example API

This simple code is just a sample to demostrate how simple and clean is the code to create a REST API with `pantofola-rest`.


```go
import (
	"fmt"
	"net/http"
	pantofola "github.com/rickycorte/pantofola-rest/router"
)

// hello is our first handler!
func hello(w http.ResponseWriter, _ *http.Request, _ *ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hi sen(pi) <3")
}

// generic handler for a parametric route that prints available values by name
func writeData(w http.ResponseWriter, r *http.Request, p *pantofola.ParameterList) {
	w.WriteHeader(200)
	fmt.Fprintf(w, r.RequestURI+"\n\nUser: "+p.Get("user")
		+"\nActivity: "+p.Get("activity")+"\nComment: "+p.Get("comment"))
}

func main() {

    mainRouter := router.MakeRouter() // create the main router used by our app
	
	// we add our handler as the main router
    mainRouter.GET("/hello", hello) 
	
	// let's make some parametric handlers
	mainRouter.GET("/activity/:user", writeData)
	mainRouter.GET("/activity/:user/:activity", writeData)
	mainRouter.GET("/activity/:user/:activity/comments/:comment", writeData)

    http.ListenAndServe(":8080", mainRouter) // start the server end enjoy your REST API!
}

```

Visit [wiki pages](https://github.com/rickycorte/pantofola-rest/wiki) to learn how to use all the features of `pantofola-rest`