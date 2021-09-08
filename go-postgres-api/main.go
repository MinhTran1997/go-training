package main

import (
	"fmt"
	"go-postgres-api/router"
	"log"
	"net/http"
)

// The main.go is our server. It will start a server on xxxx port (declared bellow in func ListenAndServe) and serve all the Router.

func main() {
	// using the Router func in router.go
	r := router.Router()

	// // fs := http.FileServer(http.Dir("build"))
	// // http.Handle("/", fs)

	fmt.Println("Starting server on the port 10001...")
	// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections
	// func ListenAndServe(addr string, handler Handler) error
	// --> https://pkg.go.dev/net/http#ListenAndServe for more details
	log.Fatal(http.ListenAndServe(":10001", r))
}
