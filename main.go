package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	unauthenticatedRouter := NewRouter()
	unauthenticatedRouter.GET("/", HandleHome)
	unauthenticatedRouter.GET("/register", HandleUserNew)

	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", HandleImageNew)

	mux := http.NewServeMux()
	mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))),
	)

	middleware := Middleware{}
	middleware.Add(unauthenticatedRouter)
	middleware.Add(mux)
	middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	http.ListenAndServe(":3000", middleware)
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()

	// Note this was
//	router.NotFound = func(http.ResponseWriter, *http.Request) {}

	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return router
}
