package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"fmt"
	"time"
)

func main() {
	unauthenticatedRouter := NewRouter()
	unauthenticatedRouter.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	unauthenticatedRouter.GET("/", HandleHome)
	unauthenticatedRouter.GET("/register", HandleUserNew)

	authenticatedRouter := NewRouter()
	authenticatedRouter.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	authenticatedRouter.GET("/images/new", HandleImageNew)

	middleware := Middleware{}
	middleware.Add(unauthenticatedRouter)
	middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	fmt.Printf("Serving requests : %s%s", time.Now(), "\n")

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
