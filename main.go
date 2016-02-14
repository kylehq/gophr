package main

import (
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := NewRouter()
	router.Handle("GET", "/", HandleHome)
	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)
	// Add user route handler
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)

	middleware := Middleware{}
	middleware.Add(router)

	fmt.Printf("Serving requests : %s%s", time.Now(), "\n")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()

	// Note this was
	//	router.NotFound = func(http.ResponseWriter, *http.Request) {}

	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return router
}
