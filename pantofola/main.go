package main

import (
	"fmt"
	"net/http"
	"pantofola/pantofola/router"
)

func hello(r *router.Request) {
	r.Reply(200, nil, "Hello!")
}

func rippe(r *router.Request) {
	r.Reply(404, nil, "Rip this is the default handler!")
}

func apiA(r *router.Request) {
	r.Reply(200, nil, "Api A :#")
}

func apiB(r *router.Request) {
	r.Reply(200, nil, "Api B :3")
}

func middlewareA(r *router.Request, next *router.MiddlewareChain) {
	fmt.Println("Middleware A before")
	next.Next(r)
	fmt.Println("Middleware A after")
}

func middlewareB(r *router.Request, next *router.MiddlewareChain) {
	fmt.Println("Middleware B before")
	next.Next(r)
	fmt.Println("Middleware B after")
}

func ripMiddleware(r *router.Request, next *router.MiddlewareChain) {

	r.Reply(404, nil, "MUHAHAHA the evil middleware killed your request!")
}

func main() {

	rout := router.MakeRouter()
	rout.AddHandler("/hello", hello)
	rout.SetFallbackHandler(rippe)

	rout.UseMiddleware(router.WriteReqMetaMiddleware)

	rout.AddHandlerChain("/shortChain", []router.Middleware{middlewareA}, hello)

	rout.AddHandlerChain("/longChain", []router.Middleware{middlewareA, middlewareB}, hello)

	rout.AddHandlerChain("/evilChain", []router.Middleware{middlewareA, ripMiddleware}, hello)

	subrouter := router.MakeRouter()
	subrouter.AddHandler("/a", apiA)
	subrouter.AddHandler("/b", apiB)

	rout.AddSubRouter("/api", subrouter)

	http.ListenAndServe(":8080", rout)

}
