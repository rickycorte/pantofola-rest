package main

import (
	"net/http"
	"pantofola/pantofola/router"
)

func hello(r *router.Request) {
	r.Reply(200, nil, "Hello!")
}

func rippe(r *router.Request) {
	r.Reply(200, nil, "Rip this is the default handler!")
}

func apiA(r *router.Request) {
	r.Reply(200, nil, "Api A :#")
}

func apiB(r *router.Request) {
	r.Reply(200, nil, "Api B :3")
}

func main() {

	rout := router.MakeRouter()
	rout.AddHandler("/hello", hello)
	rout.SetFallbackHandler(rippe)

	subrouter := router.MakeRouter()
	subrouter.AddHandler("/a", apiA)
	subrouter.AddHandler("/b", apiB)

	rout.AddSubRouter("/api", subrouter)

	http.ListenAndServe(":8080", rout)

}
