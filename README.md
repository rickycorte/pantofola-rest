# Pantofola REST

Simple and fast framework to create awesome REST APIs in seconds!

## Features

- Simple routing
- Subrouters 
- Middleware chains


## Install

First you need to install this library in your project:

`go get -u github.com/rickycorte/pantofola-rest`

## Create your first router

Now that you have installed the library you are ready to create your first router and setup your first handler!

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


## Creating a subrouter

A subrouter is a router that bundles a set of other subrouters or handlers under the same path.

For exaple we want to add `/api/a` and `/api/b`to our application and make sure we can change `/api` prefix at any time withour worring about updateding all the routers. Here subrouters come to shine!

Let's create our api subrouter:

```go

// first we create two new api handlers for our application

func apiA(r *router.Request) {
	r.Reply(200, nil, "Api A :#")
}

func apiB(r *router.Request) {
	r.Reply(200, nil, "Api B :3")
}

func main() {

    //...

    // we create a new router
    apiRouter := router.MakeRouter()

    // we add relative api path to the subrouter
    apiRouter.AddHandler("/a", apiA)
    apiRouter.AddHandler("/b", apiB)
    
    // we add our subrouter to the main one and set its path
    mainRouter.AddSubRouter("/api", apiRouter)

    //...

}

```
Easy right? Now you can visit `/api/a` and `/api/b` and see your handlers being called!

Remember you can also set a default handler for the api router if you want to create a custom reply for `/api/*` paths


## Middlewares

Middlwares are a way to customize the exection of a HTTP request. You can easly create your own middlewares and chain them how you want.

Any middleware can stop the chain execution that will always end with and handler. They can be applied at different levels to archive what you want.

For example you can add middleware to router, subrouter and handler level!

Every middleware chain is evaluated from the top router until the handler is reached.

How this can be useful?

This can be used in various ways but let's try to make a simple example. We could add to the main router a loggin middleware to log every request that reached our server and then add a auth middleware only to the api router to protect only `/api` path without adding the middleware to every single handler.

### Create your first middleware

Middlewares are just go functions, so let's create a middleware that logs a string before and after of handler:

```go

func logMiddleware(r *router.Request, next *router.MiddlewareChain) {
	fmt.Println("Before handler")
	next.Next(r)
	fmt.Println("After handler")
}
```
That's it! Just remember to call `next.Next(...)` to make the chain go forward, otherwise the chain will end and the request will never reach the next middleware (or the handler)

This is an example of a middleware that kills the chain:
```go
func killMiddleware(r *router.Request, next *router.MiddlewareChain) {

    // reply with 404 to the user and terminate the chain
	r.Reply(404, nil, "KillMiddleware killed your request!")
}
```

### Create middleware chain on handlers

You can create chains of middleware that run before an handler by simply supplying a ordered list to the router

For example:
``` go 
mainRouter.AddHandlerChain("/longChain", []router.Middleware{middlewareA, middlewareB}, hello)
```

In this example the call order will be `middlwareA` -> `middlwareB` -> `hello`

### Add a middleware to the router

We can add a new middleware to the end of the execution chain of a router with
```go
mainRouter.UseMiddleware(logMiddleware)
```

Remember, every kind of chain has no limit on the number of middlewares!