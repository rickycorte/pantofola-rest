# Pantofola REST

Simple and fast framework to create awesome REST APIs in seconds!

## Features

- Simple routing
- Subrouters 
- Middleware chains


# Example API

This simple code is just a sample to demostrate how simple and clean is the code to create a REST API with `pantofola-rest`.


```go
import (
	"fmt"
	"net/http"
	"github.com/rickycorte/pantofola-rest/router"
)

// hello is our first handler!
func hello(r *router.Request) {

	r.Reply(200, nil, "Hello!")
}


// notFound is our default handler that will be called if no route match
// remember, routers have no default router set! You have to set it manually
// this choise was made to prevent unexpected behaviours in subrouters
func notFound(r *router.Request) {
	r.Reply(404, nil, "Ohu nou a racoon ate our database!")
}

func main() {

    mainRouter := router.MakeRouter() // create the main router used by our app
    
    mainRouter.AddHandler("/hello", hello) // we add our handler as the main router

    mainRouter.SetFallbackHandler(notFound) // add the default fallback handler

    http.ListenAndServe(":8080", mainRouter) // start the server end enjoy your REST API!
}

```

Visit [wiki pages](https://github.com/rickycorte/pantofola-rest/wiki) to learn how to use all the features of `pantofola-rest`